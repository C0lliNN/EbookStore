//go:build unit
// +build unit

package dto

import (
	"github.com/c0llinn/ebook-store/internal/shop/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchOrders_ToDomain_WithEmptyData(t *testing.T) {
	dto := SearchOrders{}

	expected := model.OrderQuery{
		Limit: 10,
	}
	actual := dto.ToDomain()

	assert.Equal(t, expected, actual)
}

func TestSearchOrders_ToDomain_WithAllFields(t *testing.T) {
	dto := SearchOrders{
		Status:  "PAID",
		Page:    4,
		PerPage: 20,
	}

	expected := model.OrderQuery{
		Status: "PAID",
		Limit:  20,
		Offset: 60,
	}

	actual := dto.ToDomain()

	assert.Equal(t, expected, actual)
}

func TestCreateOrder_ToDomain(t *testing.T) {
	dto := CreateOrder{BookID: "some-book-id"}

	expected := model.Order{
		ID:     "some-order-id",
		BookID: "some-book-id",
		UserID: "some-user-id",
	}

	actual := dto.ToDomain("some-order-id", "some-user-id")

	assert.Equal(t, expected, actual)
}
