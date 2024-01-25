//nolint:staticcheck
package shop_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ebookstore/internal/core/catalog"
	"github.com/stretchr/testify/mock"

	"github.com/ebookstore/internal/core/query"
	"github.com/ebookstore/internal/core/shop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	findOrdersByQueryMethod   = "FindByQuery"
	findOrderByIDMethod       = "FindByID"
	createOrderMethod         = "Create"
	updateOrderMethod         = "Update"
	findBookByID              = "FindBookByID"
	findCartByUserIDMethod    = "FindByUserID"
	deleteCartByUserIDMethod  = "DeleteByUserID"
	createPaymentIntentMethod = "CreatePaymentIntentForOrder"
	getBookContent            = "GetBookContentURL"
	newIdMethod               = "NewID"
)

type ShopTestSuite struct {
	suite.Suite
	orderRepo      *shop.MockOrderRepository
	cartRepo       *shop.MockCartRepository
	paymentClient  *shop.MockPaymentClient
	catalogService *shop.MockCatalogService
	idGenerator    *shop.MockIDGenerator
	validator      *shop.MockValidator
	shop           *shop.Shop
}

func (s *ShopTestSuite) SetupTest() {
	s.orderRepo = new(shop.MockOrderRepository)
	s.cartRepo = new(shop.MockCartRepository)
	s.paymentClient = new(shop.MockPaymentClient)
	s.catalogService = new(shop.MockCatalogService)
	s.idGenerator = new(shop.MockIDGenerator)
	s.validator = new(shop.MockValidator)

	s.shop = shop.New(shop.Config{
		OrderRepository: s.orderRepo,
		CartRepository:  s.cartRepo,
		PaymentClient:   s.paymentClient,
		CatalogService:  s.catalogService,
		IDGenerator:     s.idGenerator,
		Validator:       s.validator,
	})
}

func TestShop(t *testing.T) {
	suite.Run(t, new(ShopTestSuite))
}

func (s *ShopTestSuite) TestFindOrders_Admin_Successfully() {
	request := shop.SearchOrders{}
	query := request.CreateQuery()
	page := request.CreatePage()

	paginatedOrders := shop.PaginatedOrders{
		Orders: []shop.Order{{ID: "some-id"}},
		Limit:  10,
	}
	ctx := context.WithValue(context.Background(), "admin", true)
	s.orderRepo.On(findOrdersByQueryMethod, ctx, query, page).Return(paginatedOrders, nil)

	expected := shop.NewPaginatedOrdersResponse(paginatedOrders)
	actual, err := s.shop.FindOrders(ctx, request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.orderRepo.AssertCalled(s.T(), findOrdersByQueryMethod, ctx, query, page)
}

func (s *ShopTestSuite) TestFindOrders_NonAdmin_Successfully() {
	request := shop.SearchOrders{}

	paginatedOrders := shop.PaginatedOrders{
		Orders: []shop.Order{{ID: "some-id"}},
		Limit:  10,
	}

	q := request.CreateQuery()
	q = *q.And(query.Condition{Field: "user_id", Operator: query.Equal, Value: "some-user-id"})
	page := request.CreatePage()

	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	s.orderRepo.On(findOrdersByQueryMethod, ctx, q, page).Return(paginatedOrders, nil)

	expected := shop.NewPaginatedOrdersResponse(paginatedOrders)
	actual, err := s.shop.FindOrders(ctx, request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.orderRepo.AssertCalled(s.T(), findOrdersByQueryMethod, ctx, q, page)
}

func (s *ShopTestSuite) TestFindOrders_WithError() {
	request := shop.SearchOrders{}
	q := request.CreateQuery()
	q = *q.And(query.Condition{Field: "user_id", Operator: query.Equal, Value: "some-user-id"})
	page := request.CreatePage()

	ctx := context.WithValue(context.Background(), "userId", "some-user-id")

	s.orderRepo.On(findOrdersByQueryMethod, ctx, q, page).Return(shop.PaginatedOrders{}, fmt.Errorf("some error"))

	_, err := s.shop.FindOrders(ctx, request)

	assert.Error(s.T(), err)
	s.orderRepo.AssertCalled(s.T(), findOrdersByQueryMethod, ctx, q, page)
}

func (s *ShopTestSuite) TestFindOrderByID_Admin_Successfully() {
	order := shop.Order{
		ID: "order-rid",
	}
	ctx := context.WithValue(context.Background(), "admin", true)
	s.orderRepo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)

	expected := shop.NewOrderResponse(order)
	actual, err := s.shop.FindOrderByID(ctx, order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
}

func (s *ShopTestSuite) TestFindOrderByID_NonAdmin_Successfully() {
	order := shop.Order{
		ID:     "order-rid",
		UserID: "current-user-id",
	}
	ctx := context.WithValue(context.Background(), "userId", "current-user-id")
	s.orderRepo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)

	expected := shop.NewOrderResponse(order)
	actual, err := s.shop.FindOrderByID(ctx, order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
}

func (s *ShopTestSuite) TestFindOrderByID_NonAdmin_Unauthorized() {
	order := shop.Order{
		ID:     "order-rid",
		UserID: "current-user-id",
	}
	ctx := context.WithValue(context.Background(), "userId", "another-user-id")
	s.orderRepo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)

	_, err := s.shop.FindOrderByID(ctx, order.ID)

	assert.ErrorIs(s.T(), err, shop.ErrForbiddenOrderAccess)
	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
}

func (s *ShopTestSuite) TestFindOrderByID_WithError() {
	id := "some-id"
	s.orderRepo.On(findOrderByIDMethod, context.TODO(), id).Return(shop.Order{}, fmt.Errorf("some error"))

	_, err := s.shop.FindOrderByID(context.TODO(), id)

	assert.Error(s.T(), err)
	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), id)
}

func (s *ShopTestSuite) TestCreateOrder_WhenCartIsNotFound() {
	userId := "some-user-id"
	ctx := context.WithValue(context.Background(), "userId", userId)

	s.cartRepo.On(findCartByUserIDMethod, ctx, userId).Return(&shop.Cart{}, fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(ctx)
	assert.Error(s.T(), err)

	s.cartRepo.AssertCalled(s.T(), findCartByUserIDMethod, ctx, userId)
	s.validator.AssertExpectations(s.T())
	s.idGenerator.AssertExpectations(s.T())
	s.catalogService.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
	s.orderRepo.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestCreateOrder_WhenPaymentClientFails() {
	userId := "some-user-id"
	ctx := context.WithValue(context.Background(), "userId", userId)

	orderId := "some-order-id"
	cart := &shop.Cart{
		ID: "some-cart-id",
	}

	order := cart.CreateOrder(orderId)
	s.idGenerator.On(newIdMethod).Return(orderId)

	s.cartRepo.On(findCartByUserIDMethod, ctx, userId).Return(cart, nil)
	s.paymentClient.On(createPaymentIntentMethod, ctx, &order).Return(fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(ctx)
	assert.Error(s.T(), err)

	s.cartRepo.AssertCalled(s.T(), findCartByUserIDMethod, ctx, userId)
	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, ctx, &order)
	s.validator.AssertExpectations(s.T())
	s.catalogService.AssertExpectations(s.T())
	s.orderRepo.AssertExpectations(s.T())

}

func (s *ShopTestSuite) TestCreateOrder_WhenRepositoryFails() {
	userId := "some-user-id"
	ctx := context.WithValue(context.Background(), "userId", userId)

	orderId := "some-order-id"
	cart := &shop.Cart{
		ID: "some-cart-id",
	}

	order := cart.CreateOrder(orderId)
	s.idGenerator.On(newIdMethod).Return(orderId)

	s.cartRepo.On(findCartByUserIDMethod, ctx, userId).Return(cart, nil)
	s.paymentClient.On(createPaymentIntentMethod, ctx, &order).Return(nil)
	s.orderRepo.On(createOrderMethod, ctx, &order).Return(fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(ctx)
	assert.Error(s.T(), err)

	s.cartRepo.AssertCalled(s.T(), findCartByUserIDMethod, ctx, userId)
	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, ctx, &order)
	s.orderRepo.AssertCalled(s.T(), createOrderMethod, ctx, &order)
	s.validator.AssertExpectations(s.T())
	s.catalogService.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestCreateOrder_Successfully() {
	userId := "some-user-id"
	ctx := context.WithValue(context.Background(), "userId", userId)

	orderId := "some-order-id"
	cart := &shop.Cart{
		ID: "some-cart-id",
	}

	order := cart.CreateOrder(orderId)
	s.idGenerator.On(newIdMethod).Return(orderId)

	s.cartRepo.On(findCartByUserIDMethod, ctx, userId).Return(cart, nil)
	s.paymentClient.On(createPaymentIntentMethod, ctx, &order).Return(nil)
	s.orderRepo.On(createOrderMethod, ctx, &order).Return(nil)
	s.cartRepo.On(deleteCartByUserIDMethod, ctx, userId).Return(nil)

	orderResponse, err := s.shop.CreateOrder(ctx)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), shop.NewOrderResponse(order), orderResponse)

	s.cartRepo.AssertCalled(s.T(), findCartByUserIDMethod, ctx, userId)
	s.cartRepo.AssertCalled(s.T(), deleteCartByUserIDMethod, ctx, userId)
	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, ctx, &order)
	s.orderRepo.AssertCalled(s.T(), createOrderMethod, ctx, &order)
	s.validator.AssertExpectations(s.T())
	s.catalogService.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestCompleteOrder_WhenOrderCouldNotBeFound() {
	id := "some-order-id"
	s.orderRepo.On(findOrderByIDMethod, context.TODO(), id).Return(shop.Order{}, fmt.Errorf("some error"))

	err := s.shop.CompleteOrder(context.TODO(), id)

	assert.Error(s.T(), err)
	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), id)
	s.orderRepo.AssertNumberOfCalls(s.T(), updateOrderMethod, 0)
}

func (s *ShopTestSuite) TestCompleteOrder_WhenUpdateFails() {
	order := shop.Order{
		ID: "some-id",
	}
	s.orderRepo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)

	order.Complete()
	s.orderRepo.On(updateOrderMethod, context.TODO(), &order).Return(fmt.Errorf("some error"))

	err := s.shop.CompleteOrder(context.TODO(), order.ID)

	assert.Error(s.T(), err)
	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.orderRepo.AssertCalled(s.T(), updateOrderMethod, context.TODO(), &order)
}

func (s *ShopTestSuite) TestCompleteOrder_Successfully() {
	order := shop.Order{
		ID: "some-id",
	}
	s.orderRepo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)

	order.Complete()
	s.orderRepo.On(updateOrderMethod, context.TODO(), &order).Return(nil)

	err := s.shop.CompleteOrder(context.TODO(), order.ID)

	assert.Nil(s.T(), err)
	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.orderRepo.AssertCalled(s.T(), updateOrderMethod, context.TODO(), &order)
}

func (s *ShopTestSuite) TestDownloadOrderItemContent_WhenOrderCouldNotBeFound() {
	order := shop.Order{
		ID:    "some-id",
		Items: []shop.Item{{ID: "some-book-id"}},
	}
	s.orderRepo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(shop.Order{}, fmt.Errorf("some error"))

	request := shop.DownloadOrderContentRequest{OrderID: order.ID, ItemID: order.Items[0].ID}

	_, err := s.shop.DownloadOrderItemContent(context.TODO(), request)

	assert.Error(s.T(), err)

	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.catalogService.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestDownloadOrderItemContent_WhenOrderIsNotPaid() {
	order := shop.Order{
		ID:     "some-id",
		Items:  []shop.Item{{ID: "some-book-id"}},
		Status: shop.Pending,
	}
	s.orderRepo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)

	request := shop.DownloadOrderContentRequest{OrderID: order.ID, ItemID: order.Items[0].ID}

	_, err := s.shop.DownloadOrderItemContent(context.TODO(), request)

	assert.ErrorIs(s.T(), err, shop.ErrOrderNotCompleted)

	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.catalogService.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestDownloadOrderItemContent_WhenItemDoesNotBelongToOrder() {
	order := shop.Order{
		ID:     "some-id",
		Items:  []shop.Item{{ID: "some-book-id"}},
		Status: shop.Paid,
	}
	s.orderRepo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)

	request := shop.DownloadOrderContentRequest{OrderID: order.ID, ItemID: "some random id"}

	_, err := s.shop.DownloadOrderItemContent(context.TODO(), request)

	assert.ErrorIs(s.T(), err, shop.ErrItemNotFoundInOrder)

	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.catalogService.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestDownloadOrderItemContent_WithError() {
	order := shop.Order{
		ID:     "some-id",
		Items:  []shop.Item{{ID: "some-book-id"}},
		Status: shop.Paid,
	}
	s.orderRepo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, context.TODO(), "some-book-id").
		Return("", fmt.Errorf("some error"))

	request := shop.DownloadOrderContentRequest{OrderID: order.ID, ItemID: order.Items[0].ID}

	_, err := s.shop.DownloadOrderItemContent(context.TODO(), request)
	assert.Error(s.T(), err)

	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.catalogService.AssertCalled(s.T(), getBookContent, context.TODO(), "some-book-id")
}

func (s *ShopTestSuite) TestDownloadOrderItemContent_NonAdmin_Forbidden() {
	ctx := context.WithValue(context.Background(), "admin", false)
	ctx = context.WithValue(ctx, "userId", "some-user-id2")

	order := shop.Order{
		ID:     "some-id",
		Items:  []shop.Item{{ID: "some-book-id"}},
		UserID: "some-user-id",
		Status: shop.Paid,
	}
	s.orderRepo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)

	request := shop.DownloadOrderContentRequest{OrderID: order.ID, ItemID: order.Items[0].ID}

	_, err := s.shop.DownloadOrderItemContent(ctx, request)

	assert.ErrorIs(s.T(), err, shop.ErrForbiddenOrderAccess)

	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
	s.catalogService.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestDownloadOrderItemContent_NonAdmin_Successfully() {
	ctx := context.WithValue(context.Background(), "admin", false)
	ctx = context.WithValue(ctx, "userId", "some-user-id2")

	order := shop.Order{
		ID:     "some-id",
		Items:  []shop.Item{{ID: "some-book-id"}},
		UserID: "some-user-id2",
		Status: shop.Paid,
	}
	s.orderRepo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, ctx, "some-book-id").
		Return("test", nil)

	request := shop.DownloadOrderContentRequest{OrderID: order.ID, ItemID: order.Items[0].ID}

	response, err := s.shop.DownloadOrderItemContent(ctx, request)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test", response.URL)

	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
	s.catalogService.AssertCalled(s.T(), getBookContent, ctx, "some-book-id")
}

func (s *ShopTestSuite) TestDownloadOrderItemContent_Admin_Successfully() {
	ctx := context.WithValue(context.Background(), "admin", true)
	ctx = context.WithValue(ctx, "userId", "some-user-id2")

	order := shop.Order{
		ID:     "some-id",
		Items:  []shop.Item{{ID: "some-book-id"}},
		UserID: "some-user-id",
		Status: shop.Paid,
	}

	s.orderRepo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, ctx, "some-book-id").
		Return("test", nil)

	request := shop.DownloadOrderContentRequest{OrderID: order.ID, ItemID: order.Items[0].ID}

	response, err := s.shop.DownloadOrderItemContent(ctx, request)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test", response.URL)

	s.orderRepo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
	s.catalogService.AssertCalled(s.T(), getBookContent, ctx, "some-book-id")
}

func (s *ShopTestSuite) TestGetCart_WhenCartCouldNotBeFound() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(&shop.Cart{}, fmt.Errorf("some error"))

	_, err := s.shop.GetCart(ctx)
	assert.Error(s.T(), err)

	s.cartRepo.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestGetCart_Successfully() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	cart := &shop.Cart{
		ID: "some-cart-id",
	}
	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(cart, nil)

	response, err := s.shop.GetCart(ctx)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), shop.NewCartResponse(*cart), response)

	s.cartRepo.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestAddItemToCart_WhenBookCouldNotBeFound() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	itemId := "some-item-id"

	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(&shop.Cart{}, nil)
	s.catalogService.On(findBookByID, ctx, itemId).Return(catalog.BookResponse{}, fmt.Errorf("some error"))

	_, err := s.shop.AddItemToCart(ctx, itemId)
	assert.Error(s.T(), err)

	s.catalogService.AssertCalled(s.T(), "FindBookByID", ctx, itemId)
	s.cartRepo.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestAddItemToCart_WhenCartDoesNotExist() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	itemId := "some-item-id"

	book := catalog.BookResponse{
		ID:    itemId,
		Title: "some-title",
		Price: 100,
	}

	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(nil, fmt.Errorf("some error"))
	s.catalogService.On(findBookByID, ctx, itemId).Return(book, nil)
	s.cartRepo.On("Save", ctx, mock.Anything).Return(nil)
	s.idGenerator.On(newIdMethod).Return("some-cart-id")

	cart, err := s.shop.AddItemToCart(ctx, itemId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), cart.ID, "some-cart-id")
	assert.Equal(s.T(), cart.Items[0].ID, itemId)
	assert.Equal(s.T(), cart.Items[0].Name, book.Title)
	assert.Equal(s.T(), cart.Items[0].Price, int64(book.Price))

	s.catalogService.AssertCalled(s.T(), findBookByID, ctx, itemId)
	s.cartRepo.AssertCalled(s.T(), findCartByUserIDMethod, ctx, "some-user-id")
	s.cartRepo.AssertCalled(s.T(), "Save", ctx, mock.Anything)
	s.idGenerator.AssertCalled(s.T(), newIdMethod)
}

func (s *ShopTestSuite) TestAddItemToCart_WhenCartExists() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	itemId := "some-item-id"

	book := catalog.BookResponse{
		ID:    itemId,
		Title: "some-title",
		Price: 100,
	}

	cart := &shop.Cart{
		ID:     "some-cart-id",
		UserID: "some-user-id",
	}

	s.catalogService.On(findBookByID, ctx, itemId).Return(book, nil)
	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(cart, nil)
	s.cartRepo.On("Save", ctx, cart).Return(nil)

	response, err := s.shop.AddItemToCart(ctx, itemId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), shop.NewCartResponse(*cart), response)

	s.catalogService.AssertCalled(s.T(), findBookByID, ctx, itemId)
	s.cartRepo.AssertCalled(s.T(), findCartByUserIDMethod, ctx, "some-user-id")
	s.cartRepo.AssertCalled(s.T(), "Save", ctx, cart)
}

func (s *ShopTestSuite) TestAddItemToCart_WhenCartCouldNotBeSaved() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	itemId := "some-item-id"

	book := catalog.BookResponse{
		ID:    itemId,
		Title: "some-title",
		Price: 100,
	}

	cart := &shop.Cart{
		ID:     "some-cart-id",
		UserID: "some-user-id",
	}

	s.catalogService.On(findBookByID, ctx, itemId).Return(book, nil)
	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(cart, nil)
	s.cartRepo.On("Save", ctx, cart).Return(fmt.Errorf("some error"))

	_, err := s.shop.AddItemToCart(ctx, itemId)
	assert.Error(s.T(), err)

	s.catalogService.AssertCalled(s.T(), findBookByID, ctx, itemId)
	s.cartRepo.AssertCalled(s.T(), findCartByUserIDMethod, ctx, "some-user-id")
	s.cartRepo.AssertCalled(s.T(), "Save", ctx, cart)
}

func (s *ShopTestSuite) TestRemoveItemFromCart_WhenCartCouldNotBeFound() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	itemId := "some-item-id"

	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(&shop.Cart{}, fmt.Errorf("some error"))

	_, err := s.shop.RemoveItemFromCart(ctx, itemId)
	assert.Error(s.T(), err)

	s.cartRepo.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestRemoveItemFromCart_WhenItemDoesNotExist() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	cart := &shop.Cart{
		ID: "some-cart-id",
	}
	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(cart, nil)

	_, err := s.shop.RemoveItemFromCart(ctx, "some-item-id")

	assert.ErrorIs(s.T(), err, shop.ErrItemNotFoundInCart)
	s.cartRepo.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestRemoveItemFromCart_WhenCartCouldNotBeSaved() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	cart := &shop.Cart{
		ID:    "some-cart-id",
		Items: []shop.Item{{ID: "some-item-id"}},
	}
	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(cart, nil)
	s.cartRepo.On("Save", ctx, cart).Return(fmt.Errorf("some error"))

	_, err := s.shop.RemoveItemFromCart(ctx, "some-item-id")

	assert.Error(s.T(), err)
	s.cartRepo.AssertExpectations(s.T())
}

func (s *ShopTestSuite) TestRemoveItemFromCart_Successfully() {
	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	cart := &shop.Cart{
		ID:    "some-cart-id",
		Items: []shop.Item{{ID: "some-item-id"}},
	}
	s.cartRepo.On(findCartByUserIDMethod, ctx, "some-user-id").Return(cart, nil)
	s.cartRepo.On("Save", ctx, cart).Return(nil)

	response, err := s.shop.RemoveItemFromCart(ctx, "some-item-id")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), shop.NewCartResponse(*cart), response)
	s.cartRepo.AssertExpectations(s.T())
}
