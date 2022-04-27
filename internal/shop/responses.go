package shop

import (
	"math"
	"time"
)

type OrderResponse struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	Total           int64     `json:"total"`
	PaymentIntentID *string   `json:"paymentIntentId"`
	ClientSecret    *string   `json:"clientSecret"`
	BookID          string    `json:"bookId"`
	UserID          string    `json:"userId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func NewOrderResponse(order Order) OrderResponse {
	return OrderResponse{
		ID:              order.ID,
		Status:          string(order.Status),
		Total:           order.Total,
		PaymentIntentID: order.PaymentIntentID,
		ClientSecret:    order.ClientSecret,
		BookID:          order.BookID,
		UserID:          order.UserID,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

type PaginatedOrdersResponse struct {
	Results     []OrderResponse `json:"results"`
	CurrentPage int             `json:"currentPage"`
	PerPage     int             `json:"perPage"`
	TotalPages  int             `json:"totalPages"`
	TotalItems  int64           `json:"totalItems"`
}

func NewPaginatedOrdersResponse(paginatedOrders PaginatedOrders) PaginatedOrdersResponse {
	orders := make([]OrderResponse, 0, len(paginatedOrders.Orders))
	for _, o := range paginatedOrders.Orders {
		orders = append(orders, NewOrderResponse(o))
	}

	return PaginatedOrdersResponse{
		Results:     orders,
		CurrentPage: (paginatedOrders.Offset / paginatedOrders.Limit) + 1,
		PerPage:     paginatedOrders.Limit,
		TotalPages:  int(math.Ceil(float64(paginatedOrders.TotalOrders) / float64(paginatedOrders.Limit))),
		TotalItems:  paginatedOrders.TotalOrders,
	}
}
