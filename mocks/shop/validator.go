// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Validator is an autogenerated mock type for the Validator type
type Validator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: i
func (_m *Validator) Validate(i interface{}) error {
	ret := _m.Called(i)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewValidator creates a new instance of Validator. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewValidator(t testing.TB) *Validator {
	mock := &Validator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
