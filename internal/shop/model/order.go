package model

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
	Total           int64
	PaymentIntentID *string
	ClientSecret    *string
	BookID          string
	UserID          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (o *Order) Complete() {
	o.Status = Paid
}
