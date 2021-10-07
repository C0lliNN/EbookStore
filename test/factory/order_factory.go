package factory

import (
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/shop/model"
	"time"
)

func NewOrder() model.Order {
	paymentIntentId := faker.UUIDHyphenated()
	clientSecret := faker.UUIDHyphenated()

	return model.Order{
		ID:              faker.UUIDHyphenated(),
		Status:          model.Pending,
		Total:           1000,
		PaymentIntentID: &paymentIntentId,
		ClientSecret:    &clientSecret,
		BookID:          faker.UUIDHyphenated(),
		UserID:          faker.UUIDHyphenated(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func NewPaginatedOrders() model.PaginatedOrders {
	return model.PaginatedOrders{
		Orders:      []model.Order{NewOrder(), NewOrder(), NewOrder()},
		Limit:       10,
		Offset:      0,
		TotalOrders: 3,
	}
}
