// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Authorizer is an autogenerated mock type for the Authorizer type
type Authorizer struct {
	mock.Mock
}

// Authorize provides a mock function with given fields: ctx, object, action
func (_m *Authorizer) Authorize(ctx context.Context, object interface{}, action string) error {
	ret := _m.Called(ctx, object, action)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, string) error); ok {
		r0 = rf(ctx, object, action)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAuthorizer creates a new instance of Authorizer. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthorizer(t testing.TB) *Authorizer {
	mock := &Authorizer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
