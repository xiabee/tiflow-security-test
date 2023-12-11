// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/api/v2/changefeed.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v2 "github.com/pingcap/tiflow/cdc/api/v2"
	v20 "github.com/pingcap/tiflow/pkg/api/v2"
)

// MockChangefeedsGetter is a mock of ChangefeedsGetter interface.
type MockChangefeedsGetter struct {
	ctrl     *gomock.Controller
	recorder *MockChangefeedsGetterMockRecorder
}

// MockChangefeedsGetterMockRecorder is the mock recorder for MockChangefeedsGetter.
type MockChangefeedsGetterMockRecorder struct {
	mock *MockChangefeedsGetter
}

// NewMockChangefeedsGetter creates a new mock instance.
func NewMockChangefeedsGetter(ctrl *gomock.Controller) *MockChangefeedsGetter {
	mock := &MockChangefeedsGetter{ctrl: ctrl}
	mock.recorder = &MockChangefeedsGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChangefeedsGetter) EXPECT() *MockChangefeedsGetterMockRecorder {
	return m.recorder
}

// Changefeeds mocks base method.
func (m *MockChangefeedsGetter) Changefeeds() v20.ChangefeedInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Changefeeds")
	ret0, _ := ret[0].(v20.ChangefeedInterface)
	return ret0
}

// Changefeeds indicates an expected call of Changefeeds.
func (mr *MockChangefeedsGetterMockRecorder) Changefeeds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Changefeeds", reflect.TypeOf((*MockChangefeedsGetter)(nil).Changefeeds))
}

// MockChangefeedInterface is a mock of ChangefeedInterface interface.
type MockChangefeedInterface struct {
	ctrl     *gomock.Controller
	recorder *MockChangefeedInterfaceMockRecorder
}

// MockChangefeedInterfaceMockRecorder is the mock recorder for MockChangefeedInterface.
type MockChangefeedInterfaceMockRecorder struct {
	mock *MockChangefeedInterface
}

// NewMockChangefeedInterface creates a new mock instance.
func NewMockChangefeedInterface(ctrl *gomock.Controller) *MockChangefeedInterface {
	mock := &MockChangefeedInterface{ctrl: ctrl}
	mock.recorder = &MockChangefeedInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChangefeedInterface) EXPECT() *MockChangefeedInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockChangefeedInterface) Create(ctx context.Context, cfg *v2.ChangefeedConfig) (*v2.ChangeFeedInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, cfg)
	ret0, _ := ret[0].(*v2.ChangeFeedInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockChangefeedInterfaceMockRecorder) Create(ctx, cfg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockChangefeedInterface)(nil).Create), ctx, cfg)
}

// Delete mocks base method.
func (m *MockChangefeedInterface) Delete(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockChangefeedInterfaceMockRecorder) Delete(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockChangefeedInterface)(nil).Delete), ctx, name)
}

// Get mocks base method.
func (m *MockChangefeedInterface) Get(ctx context.Context, name string) (*v2.ChangeFeedInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, name)
	ret0, _ := ret[0].(*v2.ChangeFeedInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockChangefeedInterfaceMockRecorder) Get(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockChangefeedInterface)(nil).Get), ctx, name)
}

// List mocks base method.
func (m *MockChangefeedInterface) List(ctx context.Context, state string) ([]v2.ChangefeedCommonInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, state)
	ret0, _ := ret[0].([]v2.ChangefeedCommonInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockChangefeedInterfaceMockRecorder) List(ctx, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockChangefeedInterface)(nil).List), ctx, state)
}

// Pause mocks base method.
func (m *MockChangefeedInterface) Pause(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pause", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Pause indicates an expected call of Pause.
func (mr *MockChangefeedInterfaceMockRecorder) Pause(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pause", reflect.TypeOf((*MockChangefeedInterface)(nil).Pause), ctx, name)
}

// Resume mocks base method.
func (m *MockChangefeedInterface) Resume(ctx context.Context, cfg *v2.ResumeChangefeedConfig, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resume", ctx, cfg, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Resume indicates an expected call of Resume.
func (mr *MockChangefeedInterfaceMockRecorder) Resume(ctx, cfg, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resume", reflect.TypeOf((*MockChangefeedInterface)(nil).Resume), ctx, cfg, name)
}

// Update mocks base method.
func (m *MockChangefeedInterface) Update(ctx context.Context, cfg *v2.ChangefeedConfig, name string) (*v2.ChangeFeedInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, cfg, name)
	ret0, _ := ret[0].(*v2.ChangeFeedInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockChangefeedInterfaceMockRecorder) Update(ctx, cfg, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockChangefeedInterface)(nil).Update), ctx, cfg, name)
}

// VerifyTable mocks base method.
func (m *MockChangefeedInterface) VerifyTable(ctx context.Context, cfg *v2.VerifyTableConfig) (*v2.Tables, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyTable", ctx, cfg)
	ret0, _ := ret[0].(*v2.Tables)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyTable indicates an expected call of VerifyTable.
func (mr *MockChangefeedInterfaceMockRecorder) VerifyTable(ctx, cfg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyTable", reflect.TypeOf((*MockChangefeedInterface)(nil).VerifyTable), ctx, cfg)
}