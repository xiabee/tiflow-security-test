// Code generated by mockery v2.14.1. DO NOT EDIT.

package file

import mock "github.com/stretchr/testify/mock"

// mockFileWriter is an autogenerated mock type for the fileWriter type
type mockFileWriter struct {
	mock.Mock
}

// AdvanceTs provides a mock function with given fields: commitTs
func (_m *mockFileWriter) AdvanceTs(commitTs uint64) {
	_m.Called(commitTs)
}

// Close provides a mock function with given fields:
func (_m *mockFileWriter) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Flush provides a mock function with given fields:
func (_m *mockFileWriter) Flush() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsRunning provides a mock function with given fields:
func (_m *mockFileWriter) IsRunning() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Write provides a mock function with given fields: p
func (_m *mockFileWriter) Write(p []byte) (int, error) {
	ret := _m.Called(p)

	var r0 int
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockFileWriter interface {
	mock.TestingT
	Cleanup(func())
}

// newMockFileWriter creates a new instance of mockFileWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockFileWriter(t mockConstructorTestingTnewMockFileWriter) *mockFileWriter {
	mock := &mockFileWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
