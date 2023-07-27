// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	verdeter "github.com/ditrit/verdeter"
	mock "github.com/stretchr/testify/mock"
)

// keySetter is an autogenerated mock type for the keySetter type
type keySetter struct {
	mock.Mock
}

// Set provides a mock function with given fields: cmd, key
func (_m *keySetter) Set(cmd *verdeter.VerdeterCommand, key configuration.keyDefinition) error {
	ret := _m.Called(cmd, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(*verdeter.VerdeterCommand, configuration.keyDefinition) error); ok {
		r0 = rf(cmd, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewKeySetter interface {
	mock.TestingT
	Cleanup(func())
}

// newKeySetter creates a new instance of keySetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newKeySetter(t mockConstructorTestingTnewKeySetter) *keySetter {
	mock := &keySetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
