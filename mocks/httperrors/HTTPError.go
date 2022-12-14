// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	zap "go.uber.org/zap"
)

// HTTPError is an autogenerated mock type for the HTTPError type
type HTTPError struct {
	mock.Mock
}

// Error provides a mock function with given fields:
func (_m *HTTPError) Error() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Log provides a mock function with given fields:
func (_m *HTTPError) Log() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ToJSON provides a mock function with given fields:
func (_m *HTTPError) ToJSON() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Write provides a mock function with given fields: httpResponse, logger
func (_m *HTTPError) Write(httpResponse http.ResponseWriter, logger *zap.Logger) {
	_m.Called(httpResponse, logger)
}

type mockConstructorTestingTNewHTTPError interface {
	mock.TestingT
	Cleanup(func())
}

// NewHTTPError creates a new instance of HTTPError. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHTTPError(t mockConstructorTestingTNewHTTPError) *HTTPError {
	mock := &HTTPError{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
