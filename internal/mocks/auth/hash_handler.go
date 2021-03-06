// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// HashHandler is an autogenerated mock type for the HashHandler type
type HashHandler struct {
	mock.Mock
}

// CompareHashAndPassword provides a mock function with given fields: hashedPassword, password
func (_m *HashHandler) CompareHashAndPassword(hashedPassword string, password string) error {
	ret := _m.Called(hashedPassword, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(hashedPassword, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HashPassword provides a mock function with given fields: password
func (_m *HashHandler) HashPassword(password string) (string, error) {
	ret := _m.Called(password)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewHashHandler creates a new instance of HashHandler. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewHashHandler(t testing.TB) *HashHandler {
	mock := &HashHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
