package shop

import (
	"math"
	"time"
)

type OrderResponse struct {
	ID              string         `json:"id"`
	Status          string         `json:"status"`
	TotalPrice      int64          `json:"total"`
	PaymentIntentID *string        `json:"paymentIntentId"`
	ClientSecret    *string        `json:"clientSecret"`
	Items           []ItemResponse `json:"bookId"`
	UserID          string         `json:"userId"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

func NewOrderResponse(order Order) OrderResponse {
	items := make([]ItemResponse, 0, len(order.Items))
	for _, i := range order.Items {
		items = append(items, newItemResponse(i))

	}
	return OrderResponse{
		ID:              order.ID,
		Status:          string(order.Status),
		TotalPrice:      order.TotalPrice(),
		PaymentIntentID: order.PaymentIntentID,
		ClientSecret:    order.ClientSecret,
		Items:           items,
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

type DownloadResponse struct {
	URL string `json:"url"`
}

type CartResponse struct {
	ID         string         `json:"id"`
	Items      []ItemResponse `json:"items"`
	UserID     string         `json:"userId"`
	TotalPrice int64          `json:"total"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
}

func NewCartResponse(cart Cart) CartResponse {
	items := make([]ItemResponse, 0, len(cart.Items))
	for _, i := range cart.Items {
		items = append(items, newItemResponse(i))
	}
	return CartResponse{
		ID:         cart.ID,
		Items:      items,
		UserID:     cart.UserID,
		TotalPrice: cart.TotalPrice(),
		CreatedAt:  cart.CreatedAt,
		UpdatedAt:  cart.UpdatedAt,
	}
}

type ItemResponse struct {
	ID             string `json:"id"`
	Name           string
	Price          int64  `json:"price"`
	PreviewImageID string `json:"previewImageId"`
}

func newItemResponse(item Item) ItemResponse {
	return ItemResponse{
		ID:             item.ID,
		Name:           item.Name,
		Price:          item.Price,
		PreviewImageID: item.PreviewImageID,
	}
}
