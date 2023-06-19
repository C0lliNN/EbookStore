// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"

	shop "github.com/ebookstore/internal/core/shop"

	testing "testing"
)

// Shop is an autogenerated mock type for the Shop type
type Shop struct {
	mock.Mock
}

// CompleteOrder provides a mock function with given fields: _a0, _a1
func (_m *Shop) CompleteOrder(_a0 context.Context, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateOrder provides a mock function with given fields: _a0, _a1
func (_m *Shop) CreateOrder(_a0 context.Context, _a1 shop.CreateOrder) (shop.OrderResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 shop.OrderResponse
	if rf, ok := ret.Get(0).(func(context.Context, shop.CreateOrder) shop.OrderResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(shop.OrderResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, shop.CreateOrder) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOrderByID provides a mock function with given fields: _a0, _a1
func (_m *Shop) FindOrderByID(_a0 context.Context, _a1 string) (shop.OrderResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 shop.OrderResponse
	if rf, ok := ret.Get(0).(func(context.Context, string) shop.OrderResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(shop.OrderResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOrders provides a mock function with given fields: _a0, _a1
func (_m *Shop) FindOrders(_a0 context.Context, _a1 shop.SearchOrders) (shop.PaginatedOrdersResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 shop.PaginatedOrdersResponse
	if rf, ok := ret.Get(0).(func(context.Context, shop.SearchOrders) shop.PaginatedOrdersResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(shop.PaginatedOrdersResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, shop.SearchOrders) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderDeliverableContent provides a mock function with given fields: _a0, _a1
func (_m *Shop) GetOrderDeliverableContent(_a0 context.Context, _a1 string) (io.ReadCloser, error) {
	ret := _m.Called(_a0, _a1)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(context.Context, string) io.ReadCloser); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewShop creates a new instance of Shop. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewShop(t testing.TB) *Shop {
	mock := &Shop{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
