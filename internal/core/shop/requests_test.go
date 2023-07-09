package shop

import (
	"testing"

	"github.com/ebookstore/internal/core/query"
	"github.com/stretchr/testify/assert"
)

func TestSearchOrders_CreateQuery_WithEmptyData(t *testing.T) {
	dto := SearchOrders{}

	expected := *query.New()
	actual := dto.CreateQuery()

	assert.Equal(t, expected, actual)
}

func TestSearchOrders_CreateQuery_WitStatus(t *testing.T) {
	dto := SearchOrders{
		Status:  "PAID",
	}

	expected := *query.New().And(query.Condition{Field: "status", Operator: query.Equal, Value: "PAID"})
	actual := dto.CreateQuery()

	assert.Equal(t, expected, actual)
}

func TestSearchOrders_CreatePage_WithPage(t *testing.T) {
	dto := SearchOrders{Page: 4}

	expected := query.Page{
		Size:  15,
		Number: 4,
	}
	actual := dto.CreatePage()

	assert.Equal(t, expected, actual)
}

func TestSearchOrders_CreatePage_WithPerPage(t *testing.T) {
	dto := SearchOrders{PerPage: 20}

	expected := query.Page{
		Number: 1,
		Size: 20,
	}
	actual := dto.CreatePage()

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
