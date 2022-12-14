// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	zap "go.uber.org/zap"
)

// ConfigurationHolder is an autogenerated mock type for the ConfigurationHolder type
type ConfigurationHolder struct {
	mock.Mock
}

// Log provides a mock function with given fields: logger
func (_m *ConfigurationHolder) Log(logger *zap.Logger) {
	_m.Called(logger)
}

// Reload provides a mock function with given fields:
func (_m *ConfigurationHolder) Reload() {
	_m.Called()
}

type mockConstructorTestingTNewConfigurationHolder interface {
	mock.TestingT
	Cleanup(func())
}

// NewConfigurationHolder creates a new instance of ConfigurationHolder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfigurationHolder(t mockConstructorTestingTNewConfigurationHolder) *ConfigurationHolder {
	mock := &ConfigurationHolder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
