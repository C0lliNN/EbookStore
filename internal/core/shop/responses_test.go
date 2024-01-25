package shop

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderResponse(t *testing.T) {
	order := Order{
		ID:     "order-id",
		Status: Paid,
		Items: []Item{
			{Price: 10},
			{Price: 20},
			{Price: 30},
		},
		PaymentIntentID: nil,
		ClientSecret:    nil,
		UserID:          "user-id",
		CreatedAt:       time.Date(2022, time.September, 24, 18, 30, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2022, time.September, 24, 18, 40, 0, 0, time.UTC),
	}

	expected := OrderResponse{
		ID:              order.ID,
		Status:          string(order.Status),
		TotalPrice:      int64(60),
		PaymentIntentID: order.PaymentIntentID,
		ClientSecret:    order.ClientSecret,
		Items:           []ItemResponse{{Price: 10}, {Price: 20}, {Price: 30}},
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
		PaymentIntentID: nil,
		ClientSecret:    nil,
		Items:           []Item{{Price: 10}, {Price: 20}, {Price: 30}},
		UserID:          "user-id",
		CreatedAt:       time.Date(2022, time.September, 24, 18, 30, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2022, time.September, 24, 18, 40, 0, 0, time.UTC),
	}
	order2 := Order{
		ID:              "order-id2",
		Status:          Pending,
		PaymentIntentID: nil,
		ClientSecret:    nil,
		Items:           []Item{{Price: 10}, {Price: 30}, {Price: 30}},
		UserID:          "user-id2",
		CreatedAt:       time.Date(2022, time.September, 25, 18, 30, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2022, time.September, 25, 18, 40, 0, 0, time.UTC),
	}
	order3 := Order{
		ID:              "order-id3",
		Status:          Cancelled,
		PaymentIntentID: nil,
		ClientSecret:    nil,
		Items:           []Item{{Price: 10}, {Price: 30}, {Price: 40}},
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

func TestNewItemResponse(t *testing.T) {
	item := Item{
		ID:             "item-id",
		Name:           "item-name",
		Price:          10,
		PreviewImageID: "preview-image-id",
		OrderID:        "order-id",
	}

	expected := ItemResponse{
		ID:             item.ID,
		Name:           item.Name,
		Price:          item.Price,
		PreviewImageID: item.PreviewImageID,
	}

	actual := newItemResponse(item)

	assert.Equal(t, expected, actual)
}

func TestNewCartResponse(t *testing.T) {
	cart := Cart{
		ID:     "cart-id",
		Items:  []Item{{Price: 10}, {Price: 20}, {Price: 30}},
		UserID: "user-id",
	}

	expected := CartResponse{
		ID:         cart.ID,
		Items:      []ItemResponse{{Price: 10}, {Price: 20}, {Price: 30}},
		UserID:     cart.UserID,
		TotalPrice: cart.TotalPrice(),
	}

	actual := NewCartResponse(cart)

	assert.Equal(t, expected, actual)
}
