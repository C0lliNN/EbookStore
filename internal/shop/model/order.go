package model

import "time"

type Order struct {
	ID            string
	Status        string
	PaymentMethod *string
	PaymentIntent *string
	BookID        string
	UserID        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
