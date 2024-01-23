// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package sinkmanager

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/processor/sourcemanager/engine"
	"github.com/pingcap/tiflow/cdc/processor/tablepb"
	"github.com/pingcap/tiflow/cdc/sinkv2/tablesink"
	cerrors "github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/retry"
	"github.com/tikv/client-go/v2/oracle"
	pd "github.com/tikv/pd/client"
	"go.uber.org/zap"
)

var tableSinkWrapperVersion uint64 = 0

// tableSinkWrapper is a wrapper of TableSink, it is used in SinkManager to manage TableSink.
// Because in the SinkManager, we write data to TableSink and RedoManager concurrently,
// so current sink node can not be reused.
type tableSinkWrapper struct {
	version uint64

	// changefeed used for logging.
	changefeed model.ChangeFeedID
	tableID    model.TableID

	tableSinkCreater func() (tablesink.TableSink, uint64)

	// tableSink is the underlying sink.
	tableSink struct {
		sync.RWMutex
		s       tablesink.TableSink
		version uint64 // it's generated by `tableSinkCreater`.

		innerMu      sync.Mutex
		advanced     time.Time
		resolvedTs   model.ResolvedTs
		checkpointTs model.ResolvedTs
		lastSyncedTs model.Ts
	}

	// state used to control the lifecycle of the table.
	state *tablepb.TableState

	// startTs is the start ts of the table.
	startTs model.Ts
	// targetTs is the upper bound of the table sink.
	targetTs model.Ts

	// barrierTs is the barrier bound of the table sink.
	barrierTs atomic.Uint64
	// receivedSorterResolvedTs is the resolved ts received from the sorter.
	// We use this to advance the redo log.
	receivedSorterResolvedTs atomic.Uint64

	// replicateTs is the ts that the table sink has started to replicate.
	replicateTs    model.Ts
	genReplicateTs func(ctx context.Context) (model.Ts, error)

	// lastCleanTime indicates the last time the table has been cleaned.
	lastCleanTime time.Time

	// rangeEventCounts is for clean the table engine.
	// If rangeEventCounts[i].events is greater than 0, it means there must be
	// events in the range (rangeEventCounts[i-1].lastPos, rangeEventCounts[i].lastPos].
	rangeEventCounts   []rangeEventCount
	rangeEventCountsMu sync.Mutex
}

type rangeEventCount struct {
	// firstPos and lastPos are used to merge many rangeEventCount into one.
	firstPos engine.Position
	lastPos  engine.Position
	events   int
}

func newRangeEventCount(pos engine.Position, events int) rangeEventCount {
	return rangeEventCount{
		firstPos: pos,
		lastPos:  pos,
		events:   events,
	}
}

func newTableSinkWrapper(
	changefeed model.ChangeFeedID,
	tableID model.TableID,
	tableSinkCreater func() (tablesink.TableSink, uint64),
	state tablepb.TableState,
	startTs model.Ts,
	targetTs model.Ts,
	genReplicateTs func(ctx context.Context) (model.Ts, error),
) *tableSinkWrapper {
	res := &tableSinkWrapper{
		version:          atomic.AddUint64(&tableSinkWrapperVersion, 1),
		changefeed:       changefeed,
		tableID:          tableID,
		tableSinkCreater: tableSinkCreater,
		state:            &state,
		startTs:          startTs,
		targetTs:         targetTs,
		genReplicateTs:   genReplicateTs,
	}

	res.tableSink.version = 0
	res.tableSink.checkpointTs = model.NewResolvedTs(startTs)
	res.tableSink.resolvedTs = model.NewResolvedTs(startTs)
	res.tableSink.advanced = time.Now()

	res.receivedSorterResolvedTs.Store(startTs)
	res.barrierTs.Store(startTs)
	return res
}

func (t *tableSinkWrapper) start(ctx context.Context, startTs model.Ts) (err error) {
	if t.replicateTs != 0 {
		log.Panic("The table sink has already started",
			zap.String("namespace", t.changefeed.Namespace),
			zap.String("changefeed", t.changefeed.ID),
			zap.Int64("tableID", t.tableID),
			zap.Uint64("startTs", startTs),
			zap.Uint64("oldReplicateTs", t.replicateTs),
		)
	}

	// FIXME(qupeng): it can be re-fetched later instead of fails.
	if t.replicateTs, err = t.genReplicateTs(ctx); err != nil {
		return errors.Trace(err)
	}

	log.Info("Sink is started",
		zap.String("namespace", t.changefeed.Namespace),
		zap.String("changefeed", t.changefeed.ID),
		zap.Int64("tableID", t.tableID),
		zap.Uint64("startTs", startTs),
		zap.Uint64("replicateTs", t.replicateTs),
	)

	// This start ts maybe greater than the initial start ts of the table sink.
	// Because in two phase scheduling, the table sink may be advanced to a later ts.
	// And we can just continue to replicate the table sink from the new start ts.
	for {
		old := t.receivedSorterResolvedTs.Load()
		if startTs <= old || t.receivedSorterResolvedTs.CompareAndSwap(old, startTs) {
			break
		}
	}
	if model.NewResolvedTs(startTs).Greater(t.tableSink.checkpointTs) {
		t.tableSink.checkpointTs = model.NewResolvedTs(startTs)
		t.tableSink.resolvedTs = model.NewResolvedTs(startTs)
		t.tableSink.advanced = time.Now()
	}
	t.state.Store(tablepb.TableStateReplicating)
	return nil
}

func (t *tableSinkWrapper) appendRowChangedEvents(events ...*model.RowChangedEvent) error {
	t.tableSink.RLock()
	defer t.tableSink.RUnlock()
	if t.tableSink.s == nil {
		// If it's nil it means it's closed.
		return tablesink.NewSinkInternalError(errors.New("table sink cleared"))
	}
	t.tableSink.s.AppendRowChangedEvents(events...)
	return nil
}

func (t *tableSinkWrapper) updateBarrierTs(ts model.Ts) {
	for {
		old := t.barrierTs.Load()
		if ts <= old || t.barrierTs.CompareAndSwap(old, ts) {
			break
		}
	}
}

func (t *tableSinkWrapper) updateReceivedSorterResolvedTs(ts model.Ts) {
	for {
		old := t.receivedSorterResolvedTs.Load()
		if ts <= old {
			return
		}
		if t.receivedSorterResolvedTs.CompareAndSwap(old, ts) {
			if t.state.Load() == tablepb.TableStatePreparing {
				t.state.Store(tablepb.TableStatePrepared)
			}
			return
		}
	}
}

func (t *tableSinkWrapper) updateResolvedTs(ts model.ResolvedTs) error {
	t.tableSink.RLock()
	defer t.tableSink.RUnlock()
	if t.tableSink.s == nil {
		// If it's nil it means it's closed.
		return tablesink.NewSinkInternalError(errors.New("table sink cleared"))
	}
	t.tableSink.innerMu.Lock()
	defer t.tableSink.innerMu.Unlock()
	t.tableSink.resolvedTs = ts
	return t.tableSink.s.UpdateResolvedTs(ts)
}

func (t *tableSinkWrapper) getLastSyncedTs() uint64 {
	t.tableSink.RLock()
	defer t.tableSink.RUnlock()
	if t.tableSink.s != nil {
		return t.tableSink.s.GetLastSyncedTs()
	}
	return t.tableSink.lastSyncedTs
}

func (t *tableSinkWrapper) getCheckpointTs() model.ResolvedTs {
	t.tableSink.RLock()
	defer t.tableSink.RUnlock()
	t.tableSink.innerMu.Lock()
	defer t.tableSink.innerMu.Unlock()

	if t.tableSink.s != nil {
		checkpointTs := t.tableSink.s.GetCheckpointTs()
		if t.tableSink.checkpointTs.Less(checkpointTs) {
			t.tableSink.checkpointTs = checkpointTs
			t.tableSink.advanced = time.Now()
		} else if !checkpointTs.Less(t.tableSink.resolvedTs) {
			t.tableSink.advanced = time.Now()
		}
	}
	return t.tableSink.checkpointTs
}

func (t *tableSinkWrapper) getReceivedSorterResolvedTs() model.Ts {
	return t.receivedSorterResolvedTs.Load()
}

func (t *tableSinkWrapper) getState() tablepb.TableState {
	return t.state.Load()
}

// getUpperBoundTs returns the upperbound of the table sink.
// It is used by sinkManager to generate sink task.
// upperBoundTs should be the minimum of the following two values:
// 1. the resolved ts of the sorter
// 2. the barrier ts of the table
func (t *tableSinkWrapper) getUpperBoundTs() model.Ts {
	resolvedTs := t.getReceivedSorterResolvedTs()
	barrierTs := t.barrierTs.Load()
	if resolvedTs > barrierTs {
		resolvedTs = barrierTs
	}
	return resolvedTs
}

func (t *tableSinkWrapper) markAsClosing() {
	for {
		curr := t.state.Load()
		if curr == tablepb.TableStateStopping || curr == tablepb.TableStateStopped {
			break
		}
		if t.state.CompareAndSwap(curr, tablepb.TableStateStopping) {
			log.Info("Sink is closing",
				zap.String("namespace", t.changefeed.Namespace),
				zap.String("changefeed", t.changefeed.ID),
				zap.Int64("tableID", t.tableID))
			break
		}
	}
}

func (t *tableSinkWrapper) markAsClosed() {
	for {
		curr := t.state.Load()
		if curr == tablepb.TableStateStopped {
			return
		}
		if t.state.CompareAndSwap(curr, tablepb.TableStateStopped) {
			log.Info("Sink is closed",
				zap.String("namespace", t.changefeed.Namespace),
				zap.String("changefeed", t.changefeed.ID),
				zap.Int64("tableID", t.tableID))
			return
		}
	}
}

func (t *tableSinkWrapper) asyncStop() bool {
	t.markAsClosing()
	if t.asyncCloseAndClearTableSink() {
		t.markAsClosed()
		return true
	}
	return false
}

// Return true means the internal table sink has been initialized.
func (t *tableSinkWrapper) initTableSink() bool {
	t.tableSink.Lock()
	defer t.tableSink.Unlock()
	if t.tableSink.s == nil {
		t.tableSink.s, t.tableSink.version = t.tableSinkCreater()
		if t.tableSink.s != nil {
			t.tableSink.advanced = time.Now()
			return true
		}
		return false
	}
	return true
}

func (t *tableSinkWrapper) asyncCloseTableSink() bool {
	t.tableSink.RLock()
	defer t.tableSink.RUnlock()
	if t.tableSink.s == nil {
		return true
	}
	return t.tableSink.s.AsyncClose()
}

func (t *tableSinkWrapper) closeTableSink() {
	t.tableSink.RLock()
	defer t.tableSink.RUnlock()
	if t.tableSink.s == nil {
		return
	}
	t.tableSink.s.Close()
}

func (t *tableSinkWrapper) asyncCloseAndClearTableSink() bool {
	closed := t.asyncCloseTableSink()
	if closed {
		t.doTableSinkClear()
	}
	return closed
}

func (t *tableSinkWrapper) closeAndClearTableSink() {
	t.closeTableSink()
	t.doTableSinkClear()
}

func (t *tableSinkWrapper) doTableSinkClear() {
	t.tableSink.Lock()
	defer t.tableSink.Unlock()
	if t.tableSink.s == nil {
		return
	}
	checkpointTs := t.tableSink.s.GetCheckpointTs()
	t.tableSink.innerMu.Lock()
	if t.tableSink.checkpointTs.Less(checkpointTs) {
		t.tableSink.checkpointTs = checkpointTs
	}
	t.tableSink.resolvedTs = checkpointTs
	t.tableSink.lastSyncedTs = t.tableSink.s.GetLastSyncedTs()
	t.tableSink.advanced = time.Now()
	t.tableSink.innerMu.Unlock()
	t.tableSink.s = nil
	t.tableSink.version = 0
}

func (t *tableSinkWrapper) checkTableSinkHealth() (err error) {
	t.tableSink.RLock()
	defer t.tableSink.RUnlock()
	if t.tableSink.s != nil {
		err = t.tableSink.s.CheckHealth()
	}
	return
}

// When the attached sink fail, there can be some events that have already been
// committed at downstream but we don't know. So we need to update `replicateTs`
// of the table so that we can re-send those events later.
func (t *tableSinkWrapper) restart(ctx context.Context) (err error) {
	if t.replicateTs, err = t.genReplicateTs(ctx); err != nil {
		return errors.Trace(err)
	}
	log.Info("Sink is restarted",
		zap.String("namespace", t.changefeed.Namespace),
		zap.String("changefeed", t.changefeed.ID),
		zap.Int64("tableID", t.tableID),
		zap.Uint64("replicateTs", t.replicateTs))
	return nil
}

func (t *tableSinkWrapper) updateRangeEventCounts(eventCount rangeEventCount) {
	t.rangeEventCountsMu.Lock()
	defer t.rangeEventCountsMu.Unlock()

	countsLen := len(t.rangeEventCounts)
	if countsLen == 0 {
		t.rangeEventCounts = append(t.rangeEventCounts, eventCount)
		return
	}
	if t.rangeEventCounts[countsLen-1].lastPos.Compare(eventCount.lastPos) < 0 {
		// If two rangeEventCounts are close enough, we can merge them into one record
		// to save memory usage. When merging B into A, A.lastPos will be updated but
		// A.firstPos will be kept so that we can determine whether to continue to merge
		// more events or not based on timeDiff(C.lastPos, A.firstPos).
		lastPhy := oracle.ExtractPhysical(t.rangeEventCounts[countsLen-1].firstPos.CommitTs)
		currPhy := oracle.ExtractPhysical(eventCount.lastPos.CommitTs)
		if (currPhy - lastPhy) >= 1000 { // 1000 means 1000ms.
			t.rangeEventCounts = append(t.rangeEventCounts, eventCount)
		} else {
			t.rangeEventCounts[countsLen-1].lastPos = eventCount.lastPos
			t.rangeEventCounts[countsLen-1].events += eventCount.events
		}
	}
}

func (t *tableSinkWrapper) cleanRangeEventCounts(upperBound engine.Position, minEvents int) bool {
	t.rangeEventCountsMu.Lock()
	defer t.rangeEventCountsMu.Unlock()

	idx := sort.Search(len(t.rangeEventCounts), func(i int) bool {
		return t.rangeEventCounts[i].lastPos.Compare(upperBound) > 0
	})
	if len(t.rangeEventCounts) == 0 || idx == 0 {
		return false
	}

	count := 0
	for _, events := range t.rangeEventCounts[0:idx] {
		count += events.events
	}
	shouldClean := count >= minEvents

	if !shouldClean {
		// To reduce engine.CleanByTable calls.
		t.rangeEventCounts[idx-1].events = count
		t.rangeEventCounts = t.rangeEventCounts[idx-1:]
	} else {
		t.rangeEventCounts = t.rangeEventCounts[idx:]
	}
	return shouldClean
}

func (t *tableSinkWrapper) sinkMaybeStuck(stuckCheck time.Duration) (bool, uint64) {
	t.getCheckpointTs()

	t.tableSink.RLock()
	defer t.tableSink.RUnlock()
	t.tableSink.innerMu.Lock()
	defer t.tableSink.innerMu.Unlock()

	// What these conditions mean:
	// 1. the table sink has been associated with a valid sink;
	// 2. its checkpoint hasn't been advanced for a while;
	version := t.tableSink.version
	advanced := t.tableSink.advanced
	if version > 0 && time.Since(advanced) > stuckCheck {
		return true, version
	}
	return false, uint64(0)
}

// handleRowChangedEvents uses to convert RowChangedEvents to TableSinkRowChangedEvents.
// It will deal with the old value compatibility.
func handleRowChangedEvents(
	changefeed model.ChangeFeedID, tableID model.TableID, events ...*model.PolymorphicEvent,
) ([]*model.RowChangedEvent, uint64) {
	size := 0
	rowChangedEvents := make([]*model.RowChangedEvent, 0, len(events))
	for _, e := range events {
		if e == nil || e.Row == nil {
			log.Warn("skip emit nil event",
				zap.String("namespace", changefeed.Namespace),
				zap.String("changefeed", changefeed.ID),
				zap.Int64("tableID", tableID),
				zap.Any("event", e))
			continue
		}

		colLen := len(e.Row.Columns)
		preColLen := len(e.Row.PreColumns)
		// Some transactions could generate empty row change event, such as
		// begin; insert into t (id) values (1); delete from t where id=1; commit;
		// Just ignore these row changed events.
		if colLen == 0 && preColLen == 0 {
			log.Warn("skip emit empty row event",
				zap.Int64("tableID", tableID),
				zap.String("namespace", changefeed.Namespace),
				zap.String("changefeed", changefeed.ID),
				zap.Any("event", e))
			continue
		}

		size += e.Row.ApproximateBytes()
		rowChangedEvents = append(rowChangedEvents, e.Row)
	}
	return rowChangedEvents, uint64(size)
}

func genReplicateTs(ctx context.Context, pdClient pd.Client) (model.Ts, error) {
	backoffBaseDelayInMs := int64(100)
	totalRetryDuration := 10 * time.Second
	var replicateTs model.Ts
	err := retry.Do(ctx, func() error {
		phy, logic, err := pdClient.GetTS(ctx)
		if err != nil {
			return errors.Trace(err)
		}
		replicateTs = oracle.ComposeTS(phy, logic)
		return nil
	}, retry.WithBackoffBaseDelay(backoffBaseDelayInMs),
		retry.WithTotalRetryDuratoin(totalRetryDuration),
		retry.WithIsRetryableErr(cerrors.IsRetryableError))
	if err != nil {
		return model.Ts(0), errors.Trace(err)
	}
	return replicateTs, nil
}
