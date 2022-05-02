package shop

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"io"
)

type Repository interface {
	FindByQuery(ctx context.Context, query OrderQuery) (PaginatedOrders, error)
	FindByID(ctx context.Context, id string) (Order, error)
	Create(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
}

type PaymentClient interface {
	CreatePaymentIntentForOrder(ctx context.Context, order *Order) error
}

type CatalogService interface {
	FindBookByID(ctx context.Context, bookId string) (catalog.BookResponse, error)
	GetBookContent(ctx context.Context, bookId string) (io.ReadCloser, error)
}

type IDGenerator interface {
	NewID() string
}

type Validator interface {
	Validate(i interface{}) error
}

type Authorizer interface {
	Authorize(ctx context.Context, object interface{}, action string) error
}

type Config struct {
	Repository     Repository
	PaymentClient  PaymentClient
	CatalogService CatalogService
	IDGenerator    IDGenerator
	Validator      Validator
	Authorizer     Authorizer
}

type Shop struct {
	Config
}

func New(c Config) *Shop {
	return &Shop{Config: c}
}

func (s *Shop) FindOrders(ctx context.Context, request SearchOrders) (PaginatedOrdersResponse, error) {
	query := request.OrderQuery()
	paginatedOrders, err := s.Repository.FindByQuery(ctx, query)
	if err != nil {
		return PaginatedOrdersResponse{}, err
	}

	return NewPaginatedOrdersResponse(paginatedOrders), nil
}

func (s *Shop) FindOrderByID(ctx context.Context, id string) (OrderResponse, error) {
	order, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		return OrderResponse{}, err
	}

	return NewOrderResponse(order), nil
}

func (s *Shop) CreateOrder(ctx context.Context, request CreateOrder) (OrderResponse, error) {
	if err := s.Validator.Validate(request); err != nil {
		return OrderResponse{}, err
	}

	order := request.Order(s.IDGenerator.NewID(), "userId")
	book, err := s.CatalogService.FindBookByID(ctx, order.BookID)
	if err != nil {
		return OrderResponse{}, err
	}
	order.Total = int64(book.Price)

	if err = s.PaymentClient.CreatePaymentIntentForOrder(ctx, &order); err != nil {
		return OrderResponse{}, err
	}

	if err = s.Repository.Create(ctx, &order); err != nil {
		return OrderResponse{}, err
	}

	return NewOrderResponse(order), nil
}

func (s *Shop) CompleteOrder(ctx context.Context, orderID string) error {
	order, err := s.Repository.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.Complete()
	return s.Repository.Update(ctx, &order)
}

func (s *Shop) GetOrderDeliverableContent(ctx context.Context, orderID string) (io.ReadCloser, error) {
	order, err := s.Repository.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != Paid {
		return nil, ErrOrderNotPaid
	}

	return s.CatalogService.GetBookContent(ctx, order.BookID)
}
