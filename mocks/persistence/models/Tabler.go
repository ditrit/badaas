// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Tabler is an autogenerated mock type for the Tabler type
type Tabler struct {
	mock.Mock
}

// TableName provides a mock function with given fields:
func (_m *Tabler) TableName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewTabler interface {
	mock.TestingT
	Cleanup(func())
}

// NewTabler creates a new instance of Tabler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTabler(t mockConstructorTestingTNewTabler) *Tabler {
	mock := &Tabler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}