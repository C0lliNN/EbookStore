package shop

import (
	"context"
	"fmt"

	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/log"
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
	GetBookContentURL(ctx context.Context, bookId string) (string, error)
}

type IDGenerator interface {
	NewID() string
}

type Validator interface {
	Validate(i interface{}) error
}

type Config struct {
	Repository     Repository
	PaymentClient  PaymentClient
	CatalogService CatalogService
	IDGenerator    IDGenerator
	Validator      Validator
}

type Shop struct {
	Config
}

func New(c Config) *Shop {
	return &Shop{Config: c}
}

func (s *Shop) FindOrders(ctx context.Context, request SearchOrders) (PaginatedOrdersResponse, error) {
	log.FromContext(ctx).Info("new request for fetching orders")

	query := request.OrderQuery()
	if !isAdmin(ctx) {
		// Non-admin users should only see their orders
		query.UserID = userId(ctx)
	}

	paginatedOrders, err := s.Repository.FindByQuery(ctx, query)
	if err != nil {
		return PaginatedOrdersResponse{}, fmt.Errorf("(FindOrders) failed fetching orders: %w", err)
	}

	return NewPaginatedOrdersResponse(paginatedOrders), nil
}

func (s *Shop) FindOrderByID(ctx context.Context, id string) (OrderResponse, error) {
	log.FromContext(ctx).Infof("new request for fetching order %s", id)

	order, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("(FindOrderByID) failed fetching order: %w", err)
	}

	if !s.isUserAllowedToReadOrder(ctx, order) {
		return OrderResponse{}, fmt.Errorf("(FindOrderByID) failed validating read conditions: %w", ErrForbiddenOrderAccess)
	}

	return NewOrderResponse(order), nil
}

func (s *Shop) CreateOrder(ctx context.Context, request CreateOrder) (OrderResponse, error) {

	if err := s.Validator.Validate(request); err != nil {
		return OrderResponse{}, fmt.Errorf("(CreateOrder) failed validating request: %w", err)
	}

	order := request.Order(s.IDGenerator.NewID(), userId(ctx))
	log.FromContext(ctx).Infof("creating new order with id %s", order.ID)

	book, err := s.CatalogService.FindBookByID(ctx, order.BookID)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("(CreateOrder) failed finding book by id %s: %w", order.BookID, err)
	}
	order.Total = int64(book.Price)

	if err = s.PaymentClient.CreatePaymentIntentForOrder(ctx, &order); err != nil {
		return OrderResponse{}, fmt.Errorf("(CreateOrder) failed creating payment intent: %w", err)
	}

	if err = s.Repository.Create(ctx, &order); err != nil {
		return OrderResponse{}, fmt.Errorf("(CreateOrder) failed creating order: %w", err)
	}

	return NewOrderResponse(order), nil
}

func (s *Shop) CompleteOrder(ctx context.Context, orderID string) error {
	log.FromContext(ctx).Infof("completing the order %s", orderID)

	order, err := s.Repository.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("(CompleteOrder) failed finding order by id %s: %w", orderID, err)
	}

	order.Complete()
	if err = s.Repository.Update(ctx, &order); err != nil {
		return fmt.Errorf("(CompleteOrder) failed updating order %s: %w", orderID, err)
	}

	return nil
}

func (s *Shop) GetOrderDeliverableContent(ctx context.Context, orderID string) (ShopBookResponse, error) {
	log.FromContext(ctx).Infof("getting deliverable content for order %s", orderID)

	order, err := s.Repository.FindByID(ctx, orderID)
	if err != nil {
		return ShopBookResponse{}, fmt.Errorf("(GetOrderDeliverableContent) failed finding order by id %s: %w", orderID, err)
	}

	if order.Status != Paid {
		return ShopBookResponse{}, fmt.Errorf("(GetOrderDeliverableContent) failed validating order conditions %s: %w", orderID, ErrOrderNotPaid)
	}

	if !s.isUserAllowedToReadOrder(ctx, order) {
		return ShopBookResponse{}, fmt.Errorf("(GetOrderDeliverableContent) failed validating user conditions: %w", ErrForbiddenOrderAccess)
	}

	url, err := s.CatalogService.GetBookContentURL(ctx, order.BookID)
	if err != nil {
		return ShopBookResponse{}, fmt.Errorf("(GetOrderDeliverableContent) failed getting book content: %w", err)
	}

	return ShopBookResponse{URL: url}, nil
}

func (s *Shop) isUserAllowedToReadOrder(ctx context.Context, order Order) bool {
	return isAdmin(ctx) || order.UserID == userId(ctx)
}

func isAdmin(ctx context.Context) bool {
	admin, ok := ctx.Value("admin").(bool)
	if !ok {
		return false
	}

	return admin
}

func userId(ctx context.Context) string {
	id, _ := ctx.Value("userId").(string)
	return id
}
