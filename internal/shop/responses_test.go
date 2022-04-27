package shop

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewOrderResponse(t *testing.T) {
	order := Order{
		ID:              "order-id",
		Status:          Paid,
		Total:           3000,
		PaymentIntentID: nil,
		ClientSecret:    nil,
		BookID:          "some-book-id",
		UserID:          "user-id",
		CreatedAt:       time.Date(2022, time.September, 24, 18, 30, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2022, time.September, 24, 18, 40, 0, 0, time.UTC),
	}

	expected := OrderResponse{
		ID:              order.ID,
		Status:          string(order.Status),
		Total:           order.Total,
		PaymentIntentID: order.PaymentIntentID,
		ClientSecret:    order.ClientSecret,
		BookID:          order.BookID,
		UserID:          order.UserID,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}

	actual := NewOrderResponse(order)

	assert.Equal(t, expected, actual)
}

func TestNewPaginatedOrdersResponse(t *testing.T) {
	order1 := Order{
		ID:              "order-id",
		Status:          Paid,
		Total:           3000,
		PaymentIntentID: nil,
		ClientSecret:    nil,
		BookID:          "some-book-id",
		UserID:          "user-id",
		CreatedAt:       time.Date(2022, time.September, 24, 18, 30, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2022, time.September, 24, 18, 40, 0, 0, time.UTC),
	}
	order2 := Order{
		ID:              "order-id2",
		Status:          Pending,
		Total:           4000,
		PaymentIntentID: nil,
		ClientSecret:    nil,
		BookID:          "some-book-id2",
		UserID:          "user-id2",
		CreatedAt:       time.Date(2022, time.September, 25, 18, 30, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2022, time.September, 25, 18, 40, 0, 0, time.UTC),
	}
	order3 := Order{
		ID:              "order-id3",
		Status:          Cancelled,
		Total:           2800,
		PaymentIntentID: nil,
		ClientSecret:    nil,
		BookID:          "some-book-id3",
		UserID:          "user-id3",
		CreatedAt:       time.Date(2022, time.September, 26, 18, 30, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2022, time.September, 26, 18, 40, 0, 0, time.UTC),
	}

	paginatedOrders := PaginatedOrders{
		Orders:      []Order{order1, order2, order3},
		Limit:       10,
		Offset:      0,
		TotalOrders: 3,
	}

	expected := PaginatedOrdersResponse{
		Results:     []OrderResponse{NewOrderResponse(order1), NewOrderResponse(order2), NewOrderResponse(order3)},
		CurrentPage: 1,
		PerPage:     10,
		TotalPages:  1,
		TotalItems:  3,
	}

	actual := NewPaginatedOrdersResponse(paginatedOrders)

	assert.Equal(t, expected, actual)
}
