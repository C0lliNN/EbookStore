// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// PasswordGenerator is an autogenerated mock type for the PasswordGenerator type
type PasswordGenerator struct {
	mock.Mock
}

// NewPassword provides a mock function with given fields:
func (_m *PasswordGenerator) NewPassword() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}