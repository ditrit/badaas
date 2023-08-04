// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	zap "go.uber.org/zap"
)

// Holder is an autogenerated mock type for the Holder type
type Holder struct {
	mock.Mock
}

// Log provides a mock function with given fields: logger
func (_m *Holder) Log(logger *zap.Logger) {
	_m.Called(logger)
}

// Reload provides a mock function with given fields:
func (_m *Holder) Reload() {
	_m.Called()
}

type mockConstructorTestingTNewHolder interface {
	mock.TestingT
	Cleanup(func())
}

// NewHolder creates a new instance of Holder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHolder(t mockConstructorTestingTNewHolder) *Holder {
	mock := &Holder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}