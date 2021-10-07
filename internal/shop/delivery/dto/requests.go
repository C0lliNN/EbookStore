package dto

import (
	"github.com/c0llinn/ebook-store/internal/shop/model"
)

type SearchOrders struct {
	Status  string `form:"status"`
	Page    int    `form:"page"`
	PerPage int    `form:"perPage"`
}

func (s *SearchOrders) ToDomain() model.OrderQuery {
	if s.Page == 0 {
		s.Page = 1
	}

	if s.PerPage == 0 {
		s.PerPage = 10
	}

	return model.OrderQuery{
		Status: model.OrderStatus(s.Status),
		Limit:  s.PerPage,
		Offset: (s.Page - 1) * s.PerPage,
	}
}

type CreateOrder struct {
	BookID string `json:"bookId" binding:"required,max=36"`
}

func (c CreateOrder) ToDomain(orderId, userId string) model.Order {
	return model.Order{
		ID:     orderId,
		BookID: c.BookID,
		UserID: userId,
	}
}
