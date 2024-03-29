// Code generated by MockGen. DO NOT EDIT.
// Source: changefeed.go

// Package mock_v2 is a generated GoMock package.
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

// GetInfo mocks base method.
func (m *MockChangefeedInterface) GetInfo(ctx context.Context, name string) (*v2.ChangeFeedInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInfo", ctx, name)
	ret0, _ := ret[0].(*v2.ChangeFeedInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInfo indicates an expected call of GetInfo.
func (mr *MockChangefeedInterfaceMockRecorder) GetInfo(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInfo", reflect.TypeOf((*MockChangefeedInterface)(nil).GetInfo), ctx, name)
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
