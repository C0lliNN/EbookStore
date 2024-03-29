// Code generated by mockery v2.40.1. DO NOT EDIT.

package shop

import (
	context "context"

	catalog "github.com/ebookstore/internal/core/catalog"

	mock "github.com/stretchr/testify/mock"
)

// MockCatalogService is an autogenerated mock type for the CatalogService type
type MockCatalogService struct {
	mock.Mock
}

// FindBookByID provides a mock function with given fields: ctx, bookId
func (_m *MockCatalogService) FindBookByID(ctx context.Context, bookId string) (catalog.BookResponse, error) {
	ret := _m.Called(ctx, bookId)

	if len(ret) == 0 {
		panic("no return value specified for FindBookByID")
	}

	var r0 catalog.BookResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (catalog.BookResponse, error)); ok {
		return rf(ctx, bookId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) catalog.BookResponse); ok {
		r0 = rf(ctx, bookId)
	} else {
		r0 = ret.Get(0).(catalog.BookResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, bookId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBookContentURL provides a mock function with given fields: ctx, bookId
func (_m *MockCatalogService) GetBookContentURL(ctx context.Context, bookId string) (string, error) {
	ret := _m.Called(ctx, bookId)

	if len(ret) == 0 {
		panic("no return value specified for GetBookContentURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, bookId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, bookId)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, bookId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockCatalogService creates a new instance of MockCatalogService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCatalogService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCatalogService {
	mock := &MockCatalogService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
