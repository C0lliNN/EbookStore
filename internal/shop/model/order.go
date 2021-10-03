package model

import "time"

type OrderStatus string

const (
	Pending   OrderStatus = "PENDING"
	Paid      OrderStatus = "PAID"
	Cancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID            string
	Status        OrderStatus
	Total         int64
	PaymentMethod *string
	PaymentIntent *string
	ClientSecret  *string
	BookID        string
	UserID        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
