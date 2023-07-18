// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	badorm "github.com/ditrit/badaas/badorm"
	mock "github.com/stretchr/testify/mock"

	reflect "reflect"
)

// iFieldIdentifier is an autogenerated mock type for the iFieldIdentifier type
type iFieldIdentifier struct {
	mock.Mock
}

// ColumnName provides a mock function with given fields: query, table
func (_m *iFieldIdentifier) ColumnName(query *badorm.Query, table badorm.Table) string {
	ret := _m.Called(query, table)

	var r0 string
	if rf, ok := ret.Get(0).(func(*badorm.Query, badorm.Table) string); ok {
		r0 = rf(query, table)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ColumnSQL provides a mock function with given fields: query, table
func (_m *iFieldIdentifier) ColumnSQL(query *badorm.Query, table badorm.Table) string {
	ret := _m.Called(query, table)

	var r0 string
	if rf, ok := ret.Get(0).(func(*badorm.Query, badorm.Table) string); ok {
		r0 = rf(query, table)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetModelType provides a mock function with given fields:
func (_m *iFieldIdentifier) GetModelType() reflect.Type {
	ret := _m.Called()

	var r0 reflect.Type
	if rf, ok := ret.Get(0).(func() reflect.Type); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(reflect.Type)
		}
	}

	return r0
}

type mockConstructorTestingTnewIFieldIdentifier interface {
	mock.TestingT
	Cleanup(func())
}

// newIFieldIdentifier creates a new instance of iFieldIdentifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newIFieldIdentifier(t mockConstructorTestingTnewIFieldIdentifier) *iFieldIdentifier {
	mock := &iFieldIdentifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
