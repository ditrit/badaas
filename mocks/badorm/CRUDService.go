// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	badorm "github.com/ditrit/badaas/badorm"
	mock "github.com/stretchr/testify/mock"
)

// CRUDService is an autogenerated mock type for the CRUDService type
type CRUDService[T interface{}, ID badorm.BadaasID] struct {
	mock.Mock
}

// CreateEntity provides a mock function with given fields: attributeValues
func (_m *CRUDService[T, ID]) CreateEntity(attributeValues map[string]interface{}) (*T, error) {
	ret := _m.Called(attributeValues)

	var r0 *T
	var r1 error
	if rf, ok := ret.Get(0).(func(map[string]interface{}) (*T, error)); ok {
		return rf(attributeValues)
	}
	if rf, ok := ret.Get(0).(func(map[string]interface{}) *T); ok {
		r0 = rf(attributeValues)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	if rf, ok := ret.Get(1).(func(map[string]interface{}) error); ok {
		r1 = rf(attributeValues)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteEntity provides a mock function with given fields: entityID
func (_m *CRUDService[T, ID]) DeleteEntity(entityID ID) error {
	ret := _m.Called(entityID)

	var r0 error
	if rf, ok := ret.Get(0).(func(ID) error); ok {
		r0 = rf(entityID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEntities provides a mock function with given fields: conditions
func (_m *CRUDService[T, ID]) GetEntities(conditions map[string]interface{}) ([]*T, error) {
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

// GetEntity provides a mock function with given fields: id
func (_m *CRUDService[T, ID]) GetEntity(id ID) (*T, error) {
	ret := _m.Called(id)

	var r0 *T
	var r1 error
	if rf, ok := ret.Get(0).(func(ID) (*T, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(ID) *T); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	if rf, ok := ret.Get(1).(func(ID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEntity provides a mock function with given fields: entityID, newValues
func (_m *CRUDService[T, ID]) UpdateEntity(entityID ID, newValues map[string]interface{}) (*T, error) {
	ret := _m.Called(entityID, newValues)

	var r0 *T
	var r1 error
	if rf, ok := ret.Get(0).(func(ID, map[string]interface{}) (*T, error)); ok {
		return rf(entityID, newValues)
	}
	if rf, ok := ret.Get(0).(func(ID, map[string]interface{}) *T); ok {
		r0 = rf(entityID, newValues)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	if rf, ok := ret.Get(1).(func(ID, map[string]interface{}) error); ok {
		r1 = rf(entityID, newValues)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCRUDService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCRUDService creates a new instance of CRUDService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCRUDService[T interface{}, ID badorm.BadaasID](t mockConstructorTestingTNewCRUDService) *CRUDService[T, ID] {
	mock := &CRUDService[T, ID]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}