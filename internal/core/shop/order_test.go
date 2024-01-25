package shop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_Complete(t *testing.T) {
	order := Order{}

	order.Complete()

	assert.Equal(t, Paid, order.Status)
}

func TestOrder_TotalPrice(t *testing.T) {
	order := Order{
		Items: []Item{
			{Price: 10},
			{Price: 20},
			{Price: 30},
		},
	}

	total := order.TotalPrice()

	assert.Equal(t, int64(60), total)
}

func TestOrder_HasItem(t *testing.T) {
	order := Order{
		Items: []Item{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
		},
	}

	assert.True(t, order.HasItem("1"))
	assert.True(t, order.HasItem("2"))
	assert.True(t, order.HasItem("3"))
	assert.False(t, order.HasItem("4"))
}

func TestOrder_Completed(t *testing.T) {
	order := Order{}

	assert.False(t, order.Completed())

	order.Complete()
	assert.True(t, order.Completed())
}
