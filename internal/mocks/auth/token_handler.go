// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	auth "github.com/ebookstore/internal/core/auth"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// TokenHandler is an autogenerated mock type for the TokenHandler type
type TokenHandler struct {
	mock.Mock
}

// ExtractUserFromToken provides a mock function with given fields: tokenString
func (_m *TokenHandler) ExtractUserFromToken(tokenString string) (auth.User, error) {
	ret := _m.Called(tokenString)

	var r0 auth.User
	if rf, ok := ret.Get(0).(func(string) auth.User); ok {
		r0 = rf(tokenString)
	} else {
		r0 = ret.Get(0).(auth.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateTokenForUser provides a mock function with given fields: user
func (_m *TokenHandler) GenerateTokenForUser(user auth.User) (string, error) {
	ret := _m.Called(user)

	var r0 string
	if rf, ok := ret.Get(0).(func(auth.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(auth.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTokenHandler creates a new instance of TokenHandler. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTokenHandler(t testing.TB) *TokenHandler {
	mock := &TokenHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
