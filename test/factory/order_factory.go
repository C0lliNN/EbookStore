package factory

import (
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/shop/model"
	"time"
)

func NewOrder() model.Order {
	paymentMethod := faker.UUIDHyphenated()
	paymentIntent := faker.UUIDHyphenated()

	return model.Order{
		ID:            faker.UUIDHyphenated(),
		Status:        model.Pending,
		PaymentMethod: &paymentMethod,
		PaymentIntent: &paymentIntent,
		BookID:        faker.UUIDHyphenated(),
		UserID:        faker.UUIDHyphenated(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
