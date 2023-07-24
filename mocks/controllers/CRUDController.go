// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	httperrors "github.com/ditrit/badaas/httperrors"
	mock "github.com/stretchr/testify/mock"
)

// CRUDController is an autogenerated mock type for the CRUDController type
type CRUDController struct {
	mock.Mock
}

// GetModel provides a mock function with given fields: w, r
func (_m *CRUDController) GetModel(w http.ResponseWriter, r *http.Request) (interface{}, httperrors.HTTPError) {
	ret := _m.Called(w, r)

	var r0 interface{}
	var r1 httperrors.HTTPError
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) (interface{}, httperrors.HTTPError)); ok {
		return rf(w, r)
	}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) interface{}); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) httperrors.HTTPError); ok {
		r1 = rf(w, r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(httperrors.HTTPError)
		}
	}

	return r0, r1
}

// GetModels provides a mock function with given fields: w, r
func (_m *CRUDController) GetModels(w http.ResponseWriter, r *http.Request) (interface{}, httperrors.HTTPError) {
	ret := _m.Called(w, r)

	var r0 interface{}
	var r1 httperrors.HTTPError
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) (interface{}, httperrors.HTTPError)); ok {
		return rf(w, r)
	}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) interface{}); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) httperrors.HTTPError); ok {
		r1 = rf(w, r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(httperrors.HTTPError)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewCRUDController interface {
	mock.TestingT
	Cleanup(func())
}

// NewCRUDController creates a new instance of CRUDController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCRUDController(t mockConstructorTestingTNewCRUDController) *CRUDController {
	mock := &CRUDController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
