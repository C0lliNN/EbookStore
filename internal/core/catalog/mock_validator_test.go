// Code generated by mockery v2.40.1. DO NOT EDIT.

package catalog

import mock "github.com/stretchr/testify/mock"

// MockValidator is an autogenerated mock type for the Validator type
type MockValidator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: i
func (_m *MockValidator) Validate(i interface{}) error {
	ret := _m.Called(i)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockValidator creates a new instance of MockValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockValidator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockValidator {
	mock := &MockValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
