// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	query "github.com/ditrit/badaas/orm/query"
)

// Operator is an autogenerated mock type for the Operator type
type Operator[T interface{}] struct {
	mock.Mock
}

// InterfaceVerificationMethod provides a mock function with given fields: _a0
func (_m *Operator[T]) InterfaceVerificationMethod(_a0 T) {
	_m.Called(_a0)
}

// ToSQL provides a mock function with given fields: _a0, columnName
func (_m *Operator[T]) ToSQL(_a0 *query.GormQuery, columnName string) (string, []interface{}, error) {
	ret := _m.Called(_a0, columnName)

	var r0 string
	var r1 []interface{}
	var r2 error
	if rf, ok := ret.Get(0).(func(*query.GormQuery, string) (string, []interface{}, error)); ok {
		return rf(_a0, columnName)
	}
	if rf, ok := ret.Get(0).(func(*query.GormQuery, string) string); ok {
		r0 = rf(_a0, columnName)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*query.GormQuery, string) []interface{}); ok {
		r1 = rf(_a0, columnName)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]interface{})
		}
	}

	if rf, ok := ret.Get(2).(func(*query.GormQuery, string) error); ok {
		r2 = rf(_a0, columnName)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewOperator interface {
	mock.TestingT
	Cleanup(func())
}

// NewOperator creates a new instance of Operator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOperator[T interface{}](t mockConstructorTestingTNewOperator) *Operator[T] {
	mock := &Operator[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
