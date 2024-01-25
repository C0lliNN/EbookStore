package shop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCart_TotalPrice(t *testing.T) {
	cart := Cart{
		Items: []Item{
			{Price: 10},
			{Price: 20},
			{Price: 30},
		},
	}

	expected := int64(60)
	actual := cart.TotalPrice()

	assert.Equal(t, expected, actual)
}

func TestCart_AddItem(t *testing.T) {
	cart := Cart{
		Items: []Item{
			{ID: "item-id"},
		},
	}

	err := cart.AddItem(Item{ID: "item-id"})
	assert.Equal(t, ErrItemAlreadyInCart, err)

	err = cart.AddItem(Item{ID: "item-id2"})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cart.Items))
}

func TestCart_RemoveItem(t *testing.T) {
	cart := Cart{
		Items: []Item{
			{ID: "item-id"},
			{ID: "item-id2"},
		},
	}

	err := cart.RemoveItem("item-id")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(cart.Items))

	err = cart.RemoveItem("item-id")
	assert.Equal(t, ErrItemNotFoundInCart, err)
}

func TestCart_CreateOrder(t *testing.T) {
	cart := Cart{
		Items: []Item{
			{ID: "item-id"},
			{ID: "item-id2"},
		},
		UserID: "user-id",
	}

	order := cart.CreateOrder("order-id")

	assert.Equal(t, "order-id", order.ID)
	assert.Equal(t, Pending, order.Status)
	assert.Equal(t, 2, len(order.Items))
	assert.Equal(t, "order-id", order.Items[0].OrderID)
	assert.Equal(t, "order-id", order.Items[1].OrderID)
	assert.Equal(t, "user-id", order.UserID)
}
