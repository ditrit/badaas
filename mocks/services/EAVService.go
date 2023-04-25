// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	models "github.com/ditrit/badaas/persistence/models"
	mock "github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

// EAVService is an autogenerated mock type for the EAVService type
type EAVService struct {
	mock.Mock
}

// CreateEntity provides a mock function with given fields: entityTypeName, attributeValues
func (_m *EAVService) CreateEntity(entityTypeName string, attributeValues map[string]interface{}) (*models.Entity, error) {
	ret := _m.Called(entityTypeName, attributeValues)

	var r0 *models.Entity
	var r1 error
	if rf, ok := ret.Get(0).(func(string, map[string]interface{}) (*models.Entity, error)); ok {
		return rf(entityTypeName, attributeValues)
	}
	if rf, ok := ret.Get(0).(func(string, map[string]interface{}) *models.Entity); ok {
		r0 = rf(entityTypeName, attributeValues)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Entity)
		}
	}

	if rf, ok := ret.Get(1).(func(string, map[string]interface{}) error); ok {
		r1 = rf(entityTypeName, attributeValues)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteEntity provides a mock function with given fields: entityTypeName, entityID
func (_m *EAVService) DeleteEntity(entityTypeName string, entityID uuid.UUID) error {
	ret := _m.Called(entityTypeName, entityID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, uuid.UUID) error); ok {
		r0 = rf(entityTypeName, entityID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEntities provides a mock function with given fields: entityTypeName, conditions
func (_m *EAVService) GetEntities(entityTypeName string, conditions map[string]string) ([]*models.Entity, error) {
	ret := _m.Called(entityTypeName, conditions)

	var r0 []*models.Entity
	var r1 error
	if rf, ok := ret.Get(0).(func(string, map[string]string) ([]*models.Entity, error)); ok {
		return rf(entityTypeName, conditions)
	}
	if rf, ok := ret.Get(0).(func(string, map[string]string) []*models.Entity); ok {
		r0 = rf(entityTypeName, conditions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Entity)
		}
	}

	if rf, ok := ret.Get(1).(func(string, map[string]string) error); ok {
		r1 = rf(entityTypeName, conditions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEntity provides a mock function with given fields: entityTypeName, id
func (_m *EAVService) GetEntity(entityTypeName string, id uuid.UUID) (*models.Entity, error) {
	ret := _m.Called(entityTypeName, id)

	var r0 *models.Entity
	var r1 error
	if rf, ok := ret.Get(0).(func(string, uuid.UUID) (*models.Entity, error)); ok {
		return rf(entityTypeName, id)
	}
	if rf, ok := ret.Get(0).(func(string, uuid.UUID) *models.Entity); ok {
		r0 = rf(entityTypeName, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Entity)
		}
	}

	if rf, ok := ret.Get(1).(func(string, uuid.UUID) error); ok {
		r1 = rf(entityTypeName, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEntity provides a mock function with given fields: entityTypeName, entityID, newValues
func (_m *EAVService) UpdateEntity(entityTypeName string, entityID uuid.UUID, newValues map[string]interface{}) (*models.Entity, error) {
	ret := _m.Called(entityTypeName, entityID, newValues)

	var r0 *models.Entity
	var r1 error
	if rf, ok := ret.Get(0).(func(string, uuid.UUID, map[string]interface{}) (*models.Entity, error)); ok {
		return rf(entityTypeName, entityID, newValues)
	}
	if rf, ok := ret.Get(0).(func(string, uuid.UUID, map[string]interface{}) *models.Entity); ok {
		r0 = rf(entityTypeName, entityID, newValues)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Entity)
		}
	}

	if rf, ok := ret.Get(1).(func(string, uuid.UUID, map[string]interface{}) error); ok {
		r1 = rf(entityTypeName, entityID, newValues)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEAVService interface {
	mock.TestingT
	Cleanup(func())
}

// NewEAVService creates a new instance of EAVService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEAVService(t mockConstructorTestingTNewEAVService) *EAVService {
	mock := &EAVService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
