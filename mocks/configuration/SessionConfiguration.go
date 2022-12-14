// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"

	zap "go.uber.org/zap"
)

// SessionConfiguration is an autogenerated mock type for the SessionConfiguration type
type SessionConfiguration struct {
	mock.Mock
}

// GetPullInterval provides a mock function with given fields:
func (_m *SessionConfiguration) GetPullInterval() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// GetRollDuration provides a mock function with given fields:
func (_m *SessionConfiguration) GetRollDuration() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// GetSessionDuration provides a mock function with given fields:
func (_m *SessionConfiguration) GetSessionDuration() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// Log provides a mock function with given fields: logger
func (_m *SessionConfiguration) Log(logger *zap.Logger) {
	_m.Called(logger)
}

// Reload provides a mock function with given fields:
func (_m *SessionConfiguration) Reload() {
	_m.Called()
}

type mockConstructorTestingTNewSessionConfiguration interface {
	mock.TestingT
	Cleanup(func())
}

// NewSessionConfiguration creates a new instance of SessionConfiguration. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSessionConfiguration(t mockConstructorTestingTNewSessionConfiguration) *SessionConfiguration {
	mock := &SessionConfiguration{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
