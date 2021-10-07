package dto

import (
	"github.com/c0llinn/ebook-store/internal/shop/model"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromOrder(t *testing.T) {
	order := factory.NewOrder()

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

	actual := FromOrder(order)

	assert.Equal(t, expected, actual)
}

func TestFromPaginatedOrders(t *testing.T) {
	order1, order2, order3 := factory.NewOrder(), factory.NewOrder(), factory.NewOrder()

	paginatedOrders := model.PaginatedOrders{
		Orders:      []model.Order{order1, order2, order3},
		Limit:       10,
		Offset:      0,
		TotalOrders: 3,
	}

	expected := PaginatedOrdersResponse{
		Results:     []OrderResponse{FromOrder(order1), FromOrder(order2), FromOrder(order3)},
		CurrentPage: 1,
		PerPage:     10,
		TotalPages:  1,
		TotalItems:  3,
	}

	actual := FromPaginatedOrders(paginatedOrders)

	assert.Equal(t, expected, actual)
}
