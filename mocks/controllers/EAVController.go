// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	httperrors "github.com/ditrit/badaas/httperrors"
	mock "github.com/stretchr/testify/mock"
)

// EAVController is an autogenerated mock type for the EAVController type
type EAVController struct {
	mock.Mock
}

// CreateObject provides a mock function with given fields: w, r
func (_m *EAVController) CreateObject(w http.ResponseWriter, r *http.Request) (interface{}, httperrors.HTTPError) {
	ret := _m.Called(w, r)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) interface{}); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 httperrors.HTTPError
	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) httperrors.HTTPError); ok {
		r1 = rf(w, r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(httperrors.HTTPError)
		}
	}

	return r0, r1
}

// DeleteObject provides a mock function with given fields: w, r
func (_m *EAVController) DeleteObject(w http.ResponseWriter, r *http.Request) (interface{}, httperrors.HTTPError) {
	ret := _m.Called(w, r)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) interface{}); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 httperrors.HTTPError
	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) httperrors.HTTPError); ok {
		r1 = rf(w, r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(httperrors.HTTPError)
		}
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: w, r
func (_m *EAVController) GetAll(w http.ResponseWriter, r *http.Request) (interface{}, httperrors.HTTPError) {
	ret := _m.Called(w, r)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) interface{}); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 httperrors.HTTPError
	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) httperrors.HTTPError); ok {
		r1 = rf(w, r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(httperrors.HTTPError)
		}
	}

	return r0, r1
}

// GetObject provides a mock function with given fields: w, r
func (_m *EAVController) GetObject(w http.ResponseWriter, r *http.Request) (interface{}, httperrors.HTTPError) {
	ret := _m.Called(w, r)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) interface{}); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 httperrors.HTTPError
	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) httperrors.HTTPError); ok {
		r1 = rf(w, r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(httperrors.HTTPError)
		}
	}

	return r0, r1
}

// ModifyObject provides a mock function with given fields: w, r
func (_m *EAVController) ModifyObject(w http.ResponseWriter, r *http.Request) (interface{}, httperrors.HTTPError) {
	ret := _m.Called(w, r)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) interface{}); ok {
		r0 = rf(w, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 httperrors.HTTPError
	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) httperrors.HTTPError); ok {
		r1 = rf(w, r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(httperrors.HTTPError)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewEAVController interface {
	mock.TestingT
	Cleanup(func())
}

// NewEAVController creates a new instance of EAVController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEAVController(t mockConstructorTestingTNewEAVController) *EAVController {
	mock := &EAVController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}