// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// FilenameGenerator is an autogenerated mock type for the FilenameGenerator type
type FilenameGenerator struct {
	mock.Mock
}

// NewUniqueName provides a mock function with given fields: filename
func (_m *FilenameGenerator) NewUniqueName(filename string) string {
	ret := _m.Called(filename)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(filename)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewFilenameGenerator creates a new instance of FilenameGenerator. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFilenameGenerator(t testing.TB) *FilenameGenerator {
	mock := &FilenameGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
