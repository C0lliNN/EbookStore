package shop

import "time"

type OrderStatus string

const (
	Pending   OrderStatus = "PENDING"
	Paid      OrderStatus = "PAID"
	Cancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID              string
	Status          OrderStatus
	PaymentIntentID *string
	ClientSecret    *string
	Items           []Item
	UserID          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (o *Order) Complete() {
	o.Status = Paid
}

func (o *Order) Completed() bool {
	return o.Status == Paid
}

func (o *Order) TotalPrice() int64 {
	var total int64
	for _, i := range o.Items {
		total += i.Price
	}
	return total
}

func (o *Order) HasItem(itemID string) bool {
	for _, i := range o.Items {
		if i.ID == itemID {
			return true
		}
	}
	return false
}
