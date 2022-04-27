package factory

import (
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/shop"
	"time"
)

func NewOrder() shop.Order {
	paymentIntentId := faker.UUIDHyphenated()
	clientSecret := faker.UUIDHyphenated()

	return shop.Order{
		ID:              faker.UUIDHyphenated(),
		Status:          shop.Pending,
		Total:           1000,
		PaymentIntentID: &paymentIntentId,
		ClientSecret:    &clientSecret,
		BookID:          faker.UUIDHyphenated(),
		UserID:          faker.UUIDHyphenated(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
