package shop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchOrders_OrderQuery_WithEmptyData(t *testing.T) {
	dto := SearchOrders{}

	expected := OrderQuery{
		Limit: 10,
	}
	actual := dto.OrderQuery()

	assert.Equal(t, expected, actual)
}

func TestSearchOrders_OrderQuery_WithAllFields(t *testing.T) {
	dto := SearchOrders{
		Status:  "PAID",
		Page:    4,
		PerPage: 20,
	}

	expected := OrderQuery{
		Status: "PAID",
		Limit:  20,
		Offset: 60,
	}

	actual := dto.OrderQuery()

	assert.Equal(t, expected, actual)
}

func TestCreateOrder_Order(t *testing.T) {
	dto := CreateOrder{BookID: "some-book-id"}

	expected := Order{
		ID:     "some-order-id",
		BookID: "some-book-id",
		UserID: "some-user-id",
	}

	actual := dto.Order("some-order-id", "some-user-id")

	assert.Equal(t, expected, actual)
}
