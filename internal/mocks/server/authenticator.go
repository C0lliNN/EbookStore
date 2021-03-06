// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/ebookstore/internal/auth"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Authenticator is an autogenerated mock type for the Authenticator type
type Authenticator struct {
	mock.Mock
}

// Login provides a mock function with given fields: _a0, _a1
func (_m *Authenticator) Login(_a0 context.Context, _a1 auth.LoginRequest) (auth.CredentialsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 auth.CredentialsResponse
	if rf, ok := ret.Get(0).(func(context.Context, auth.LoginRequest) auth.CredentialsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(auth.CredentialsResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, auth.LoginRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: _a0, _a1
func (_m *Authenticator) Register(_a0 context.Context, _a1 auth.RegisterRequest) (auth.CredentialsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 auth.CredentialsResponse
	if rf, ok := ret.Get(0).(func(context.Context, auth.RegisterRequest) auth.CredentialsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(auth.CredentialsResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, auth.RegisterRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResetPassword provides a mock function with given fields: _a0, _a1
func (_m *Authenticator) ResetPassword(_a0 context.Context, _a1 auth.PasswordResetRequest) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.PasswordResetRequest) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAuthenticator creates a new instance of Authenticator. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthenticator(t testing.TB) *Authenticator {
	mock := &Authenticator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
