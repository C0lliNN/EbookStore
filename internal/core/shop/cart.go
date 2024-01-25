package shop

import (
	"time"
)

// Cart represents a Shopping Cart. It might be turned into an Order.
type Cart struct {
	ID        string
	Items     []Item
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Cart) TotalPrice() int64 {
	var total int64
	for _, i := range c.Items {
		total += i.Price
	}
	return total
}

func (c *Cart) AddItem(item Item) error {
	for _, i := range c.Items {
		if i.ID == item.ID {
			return ErrItemAlreadyInCart
		}
	}

	c.Items = append(c.Items, item)
	return nil
}

func (c *Cart) RemoveItem(itemID string) error {
	for i, item := range c.Items {
		if item.ID == itemID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return nil
		}
	}

	return ErrItemNotFoundInCart
}

func (c *Cart) CreateOrder(ID string) Order {
	for i := range c.Items {
		c.Items[i].OrderID = ID
	}

	return Order{
		ID:     ID,
		Status: Pending,
		Items:  c.Items,
		UserID: c.UserID,
	}
}
