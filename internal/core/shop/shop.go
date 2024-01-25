package shop

import (
	"context"
	"fmt"

	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/core/query"
	"github.com/ebookstore/internal/log"
)

type OrderRepository interface {
	FindByQuery(ctx context.Context, q query.Query, p query.Page) (PaginatedOrders, error)
	FindByID(ctx context.Context, id string) (Order, error)
	Create(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
}

type CartRepository interface {
	FindByUserID(ctx context.Context, userID string) (*Cart, error)
	Save(ctx context.Context, cart *Cart) error
	DeleteByUserID(ctx context.Context, userID string) error
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
	OrderRepository OrderRepository
	CartRepository  CartRepository
	PaymentClient   PaymentClient
	CatalogService  CatalogService
	IDGenerator     IDGenerator
	Validator       Validator
}

type Shop struct {
	Config
}

func New(c Config) *Shop {
	return &Shop{Config: c}
}

func (s *Shop) FindOrders(ctx context.Context, request SearchOrders) (PaginatedOrdersResponse, error) {
	log.Infof(ctx, "new request for fetching orders")

	q := request.CreateQuery()
	if !isAdmin(ctx) {
		// Non-admin users should only see their orders
		q = *q.And(query.Condition{Field: "user_id", Operator: query.Equal, Value: userId(ctx)})
	}

	paginatedOrders, err := s.OrderRepository.FindByQuery(ctx, q, request.CreatePage())
	if err != nil {
		return PaginatedOrdersResponse{}, fmt.Errorf("(FindOrders) failed fetching orders: %w", err)
	}

	return NewPaginatedOrdersResponse(paginatedOrders), nil
}

func (s *Shop) FindOrderByID(ctx context.Context, id string) (OrderResponse, error) {
	log.Infof(ctx, "new request for fetching order %s", id)

	order, err := s.OrderRepository.FindByID(ctx, id)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("(FindOrderByID) failed fetching order: %w", err)
	}

	if !s.isUserAllowedToReadOrder(ctx, order) {
		return OrderResponse{}, fmt.Errorf("(FindOrderByID) failed validating read conditions: %w", ErrForbiddenOrderAccess)
	}

	return NewOrderResponse(order), nil
}

// CreateOrder creates a new order for the current user.
// Precondition: the user must have a cart.
func (s *Shop) CreateOrder(ctx context.Context) (OrderResponse, error) {
	log.Infof(ctx, "new request for creating order")

	cart, err := s.CartRepository.FindByUserID(ctx, userId(ctx))
	if err != nil {
		return OrderResponse{}, fmt.Errorf("(CreateOrder) failed finding cart: %w", err)
	}

	order := cart.CreateOrder(s.IDGenerator.NewID())

	if err = s.PaymentClient.CreatePaymentIntentForOrder(ctx, &order); err != nil {
		return OrderResponse{}, fmt.Errorf("(CreateOrder) failed creating payment intent: %w", err)
	}

	if err = s.OrderRepository.Create(ctx, &order); err != nil {
		return OrderResponse{}, fmt.Errorf("(CreateOrder) failed creating order: %w", err)
	}

	if err = s.CartRepository.DeleteByUserID(ctx, userId(ctx)); err != nil {
		log.Warnf(ctx, "(CreateOrder) failed deleting cart: %v", err)
	}

	return NewOrderResponse(order), nil
}

func (s *Shop) CompleteOrder(ctx context.Context, orderID string) error {
	log.Infof(ctx, "completing the order %s", orderID)

	order, err := s.OrderRepository.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("(CompleteOrder) failed finding order by id %s: %w", orderID, err)
	}

	order.Complete()
	if err = s.OrderRepository.Update(ctx, &order); err != nil {
		return fmt.Errorf("(CompleteOrder) failed updating order %s: %w", orderID, err)
	}

	return nil
}

func (s *Shop) DownloadOrderItemContent(ctx context.Context, request DownloadOrderContentRequest) (DownloadResponse, error) {
	log.Infof(ctx, "getting content for item %s of order %s", request.ItemID, request.OrderID)

	order, err := s.OrderRepository.FindByID(ctx, request.OrderID)
	if err != nil {
		return DownloadResponse{}, fmt.Errorf("(DownloadOrderItemContent) failed finding order by id %s: %w", request.OrderID, err)
	}

	if !order.Completed() {
		return DownloadResponse{}, fmt.Errorf("(DownloadOrderItemContent) failed validating order conditions %s: %w", request.OrderID, ErrOrderNotCompleted)
	}

	if !order.HasItem(request.ItemID) {
		return DownloadResponse{}, fmt.Errorf("(DownloadOrderItemContent) failed validating order conditions %s: %w", request.OrderID, ErrItemNotFoundInOrder)
	}

	if !s.isUserAllowedToReadOrder(ctx, order) {
		return DownloadResponse{}, fmt.Errorf("(DownloadOrderItemContent) failed validating user conditions: %w", ErrForbiddenOrderAccess)
	}

	url, err := s.CatalogService.GetBookContentURL(ctx, request.ItemID)
	if err != nil {
		return DownloadResponse{}, fmt.Errorf("(DownloadOrderItemContent) failed getting book content: %w", err)
	}

	return DownloadResponse{URL: url}, nil
}

func (s *Shop) GetCart(ctx context.Context) (CartResponse, error) {
	log.Infof(ctx, "getting cart")

	cart, err := s.CartRepository.FindByUserID(ctx, userId(ctx))
	if err != nil {
		return CartResponse{}, fmt.Errorf("(GetCart) failed finding cart: %w", err)
	}

	return NewCartResponse(*cart), nil
}

// AddItemToCart adds an item to the current user's cart.
// Post-condition: the user has a cart in valid state.
func (s *Shop) AddItemToCart(ctx context.Context, itemID string) (CartResponse, error) {
	log.Infof(ctx, "adding item to cart")

	cart, err := s.findOrCreateCart(ctx)
	if err != nil {
		return CartResponse{}, fmt.Errorf("(AddItemToCart) failed finding cart: %w", err)
	}

	book, err := s.CatalogService.FindBookByID(ctx, itemID)
	if err != nil {
		return CartResponse{}, fmt.Errorf("(AddItemToCart) failed finding item by id %s: %w", itemID, err)
	}
	item := Item{ID: itemID, Name: book.Title, Price: int64(book.Price), PreviewImageID: book.MainImageID}

	if err = cart.AddItem(item); err != nil {
		return CartResponse{}, fmt.Errorf("(AddItemToCart) failed adding item to cart: %w", err)
	}

	if err = s.CartRepository.Save(ctx, cart); err != nil {
		return CartResponse{}, fmt.Errorf("(AddItemToCart) failed saving cart: %w", err)
	}

	return NewCartResponse(*cart), nil
}

func (s *Shop) RemoveItemFromCart(ctx context.Context, itemID string) (CartResponse, error) {
	log.Infof(ctx, "removing item from cart")

	cart, err := s.findOrCreateCart(ctx)
	if err != nil {
		return CartResponse{}, fmt.Errorf("(RemoveItemFromCart) failed finding cart: %w", err)
	}

	if err = cart.RemoveItem(itemID); err != nil {
		return CartResponse{}, fmt.Errorf("(RemoveItemFromCart) failed removing item from cart: %w", err)
	}

	if err = s.CartRepository.Save(ctx, cart); err != nil {
		return CartResponse{}, fmt.Errorf("(RemoveItemFromCart) failed saving cart: %w", err)
	}

	return NewCartResponse(*cart), nil
}

func (s *Shop) findOrCreateCart(ctx context.Context) (*Cart, error) {
	cart, err := s.CartRepository.FindByUserID(ctx, userId(ctx))
	if err != nil {
		log.Debugf(ctx, "cart not found for user %s, creating a new one", userId(ctx))
	}

	if cart == nil {
		cart = &Cart{
			ID:     s.IDGenerator.NewID(),
			UserID: userId(ctx),
		}
	}

	return cart, nil
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
