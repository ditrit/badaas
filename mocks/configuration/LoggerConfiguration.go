// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	zap "go.uber.org/zap"
)

// LoggerConfiguration is an autogenerated mock type for the LoggerConfiguration type
type LoggerConfiguration struct {
	mock.Mock
}

// GetMode provides a mock function with given fields:
func (_m *LoggerConfiguration) GetMode() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetRequestTemplate provides a mock function with given fields:
func (_m *LoggerConfiguration) GetRequestTemplate() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Log provides a mock function with given fields: logger
func (_m *LoggerConfiguration) Log(logger *zap.Logger) {
	_m.Called(logger)
}

// Reload provides a mock function with given fields:
func (_m *LoggerConfiguration) Reload() {
	_m.Called()
}

type mockConstructorTestingTNewLoggerConfiguration interface {
	mock.TestingT
	Cleanup(func())
}

// NewLoggerConfiguration creates a new instance of LoggerConfiguration. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLoggerConfiguration(t mockConstructorTestingTNewLoggerConfiguration) *LoggerConfiguration {
	mock := &LoggerConfiguration{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
