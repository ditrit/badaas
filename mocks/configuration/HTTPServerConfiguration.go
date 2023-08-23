// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"

	zap "go.uber.org/zap"
)

// HTTPServerConfiguration is an autogenerated mock type for the HTTPServerConfiguration type
type HTTPServerConfiguration struct {
	mock.Mock
}

// GetAddr provides a mock function with given fields:
func (_m *HTTPServerConfiguration) GetAddr() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetHost provides a mock function with given fields:
func (_m *HTTPServerConfiguration) GetHost() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMaxTimeout provides a mock function with given fields:
func (_m *HTTPServerConfiguration) GetMaxTimeout() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// GetPort provides a mock function with given fields:
func (_m *HTTPServerConfiguration) GetPort() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Log provides a mock function with given fields: logger
func (_m *HTTPServerConfiguration) Log(logger *zap.Logger) {
	_m.Called(logger)
}

// Reload provides a mock function with given fields:
func (_m *HTTPServerConfiguration) Reload() {
	_m.Called()
}

type mockConstructorTestingTNewHTTPServerConfiguration interface {
	mock.TestingT
	Cleanup(func())
}

// NewHTTPServerConfiguration creates a new instance of HTTPServerConfiguration. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHTTPServerConfiguration(t mockConstructorTestingTNewHTTPServerConfiguration) *HTTPServerConfiguration {
	mock := &HTTPServerConfiguration{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
