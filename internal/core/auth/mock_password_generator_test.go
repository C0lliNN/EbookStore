// Code generated by mockery v2.40.1. DO NOT EDIT.

package auth

import mock "github.com/stretchr/testify/mock"

// MockPasswordGenerator is an autogenerated mock type for the PasswordGenerator type
type MockPasswordGenerator struct {
	mock.Mock
}

// NewPassword provides a mock function with given fields:
func (_m *MockPasswordGenerator) NewPassword() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewPassword")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewMockPasswordGenerator creates a new instance of MockPasswordGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPasswordGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPasswordGenerator {
	mock := &MockPasswordGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}