// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	badorm "github.com/ditrit/badaas/badorm"
	mock "github.com/stretchr/testify/mock"
)

// CRUDUnsafeService is an autogenerated mock type for the CRUDUnsafeService type
type CRUDUnsafeService[T interface{}, ID badorm.BadaasID] struct {
	mock.Mock
}

// GetEntities provides a mock function with given fields: conditions
func (_m *CRUDUnsafeService[T, ID]) GetEntities(conditions map[string]interface{}) ([]*T, error) {
	ret := _m.Called(conditions)

	var r0 []*T
	var r1 error
	if rf, ok := ret.Get(0).(func(map[string]interface{}) ([]*T, error)); ok {
		return rf(conditions)
	}
	if rf, ok := ret.Get(0).(func(map[string]interface{}) []*T); ok {
		r0 = rf(conditions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*T)
		}
	}

	if rf, ok := ret.Get(1).(func(map[string]interface{}) error); ok {
		r1 = rf(conditions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCRUDUnsafeService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCRUDUnsafeService creates a new instance of CRUDUnsafeService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCRUDUnsafeService[T interface{}, ID badorm.BadaasID](t mockConstructorTestingTNewCRUDUnsafeService) *CRUDUnsafeService[T, ID] {
	mock := &CRUDUnsafeService[T, ID]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
