// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	configuration "github.com/ditrit/badaas/configuration"
	mock "github.com/stretchr/testify/mock"

	verdeter "github.com/ditrit/verdeter"
)

// KeySetter is an autogenerated mock type for the KeySetter type
type KeySetter struct {
	mock.Mock
}

// Set provides a mock function with given fields: command, key
func (_m *KeySetter) Set(command *verdeter.VerdeterCommand, key configuration.KeyDefinition) error {
	ret := _m.Called(command, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(*verdeter.VerdeterCommand, configuration.KeyDefinition) error); ok {
		r0 = rf(command, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewKeySetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewKeySetter creates a new instance of KeySetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewKeySetter(t mockConstructorTestingTNewKeySetter) *KeySetter {
	mock := &KeySetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
