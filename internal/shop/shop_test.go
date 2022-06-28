//nolint:staticcheck
package shop_test

import (
	"context"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/catalog"
	mocks2 "github.com/c0llinn/ebook-store/internal/mocks/shop"
	"github.com/c0llinn/ebook-store/internal/shop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"strings"
	"testing"
)

const (
	findOrdersByQueryMethod   = "FindByQuery"
	findOrderByIDMethod       = "FindByID"
	createOrderMethod         = "Create"
	updateOrderMethod         = "Update"
	findBookByIDMethod        = "FindBookByID"
	createPaymentIntentMethod = "CreatePaymentIntentForOrder"
	getBookContent            = "GetBookContent"
	newIdMethod               = "NewID"
	validateMethod            = "Validate"
)

type ShopTestSuite struct {
	suite.Suite
	repo           *mocks2.Repository
	paymentClient  *mocks2.PaymentClient
	catalogService *mocks2.CatalogService
	idGenerator    *mocks2.IDGenerator
	validator      *mocks2.Validator
	shop           *shop.Shop
}

func (s *ShopTestSuite) SetupTest() {
	s.repo = new(mocks2.Repository)
	s.paymentClient = new(mocks2.PaymentClient)
	s.catalogService = new(mocks2.CatalogService)
	s.idGenerator = new(mocks2.IDGenerator)
	s.validator = new(mocks2.Validator)

	s.shop = shop.New(shop.Config{
		Repository:     s.repo,
		PaymentClient:  s.paymentClient,
		CatalogService: s.catalogService,
		IDGenerator:    s.idGenerator,
		Validator:      s.validator,
	})
}

func TestShop(t *testing.T) {
	suite.Run(t, new(ShopTestSuite))
}

func (s *ShopTestSuite) TestFindOrders_Admin_Successfully() {
	request := shop.SearchOrders{}
	query := request.OrderQuery()

	paginatedOrders := shop.PaginatedOrders{
		Orders: []shop.Order{{ID: "some-id"}},
		Limit:  10,
	}
	ctx := context.WithValue(context.Background(), "admin", true)
	s.repo.On(findOrdersByQueryMethod, ctx, query).Return(paginatedOrders, nil)

	expected := shop.NewPaginatedOrdersResponse(paginatedOrders)
	actual, err := s.shop.FindOrders(ctx, request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repo.AssertCalled(s.T(), findOrdersByQueryMethod, ctx, query)
}

func (s *ShopTestSuite) TestFindOrders_NonAdmin_Successfully() {
	request := shop.SearchOrders{}

	paginatedOrders := shop.PaginatedOrders{
		Orders: []shop.Order{{ID: "some-id"}},
		Limit:  10,
	}

	query := request.OrderQuery()
	query.UserID = "some-user-id"

	ctx := context.WithValue(context.Background(), "userId", "some-user-id")
	s.repo.On(findOrdersByQueryMethod, ctx, query).Return(paginatedOrders, nil)

	expected := shop.NewPaginatedOrdersResponse(paginatedOrders)
	actual, err := s.shop.FindOrders(ctx, request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repo.AssertCalled(s.T(), findOrdersByQueryMethod, ctx, query)
}

func (s *ShopTestSuite) TestFindOrders_WithError() {
	request := shop.SearchOrders{}
	query := request.OrderQuery()

	s.repo.On(findOrdersByQueryMethod, context.TODO(), query).Return(shop.PaginatedOrders{}, fmt.Errorf("some error"))

	_, err := s.shop.FindOrders(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrdersByQueryMethod, context.TODO(), query)
}

func (s *ShopTestSuite) TestFindOrderByID_Admin_Successfully() {
	order := shop.Order{
		ID:    "order-rid",
		Total: 4000,
	}
	ctx := context.WithValue(context.Background(), "admin", true)
	s.repo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)

	expected := shop.NewOrderResponse(order)
	actual, err := s.shop.FindOrderByID(ctx, order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
}

func (s *ShopTestSuite) TestFindOrderByID_NonAdmin_Successfully() {
	order := shop.Order{
		ID:     "order-rid",
		UserID: "current-user-id",
		Total:  4000,
	}
	ctx := context.WithValue(context.Background(), "userId", "current-user-id")
	s.repo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)

	expected := shop.NewOrderResponse(order)
	actual, err := s.shop.FindOrderByID(ctx, order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
}

func (s *ShopTestSuite) TestFindOrderByID_NonAdmin_Unauthorized() {
	order := shop.Order{
		ID:     "order-rid",
		UserID: "current-user-id",
		Total:  4000,
	}
	ctx := context.WithValue(context.Background(), "userId", "another-user-id")
	s.repo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)

	_, err := s.shop.FindOrderByID(ctx, order.ID)

	assert.Equal(s.T(), shop.ErrForbiddenOrderAccess, err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
}

func (s *ShopTestSuite) TestFindOrderByID_WithError() {
	id := "some-id"
	s.repo.On(findOrderByIDMethod, context.TODO(), id).Return(shop.Order{}, fmt.Errorf("some error"))

	_, err := s.shop.FindOrderByID(context.TODO(), id)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), id)
}

func (s *ShopTestSuite) TestCreateOrder_ValidationFails() {
	request := shop.CreateOrder{
		BookID: "some-id",
	}
	orderId := "orderId"
	userId := "userId"
	order := request.Order(orderId, userId)

	s.validator.On(validateMethod, request).Return(fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertCalled(s.T(), validateMethod, request)
	s.idGenerator.AssertNotCalled(s.T(), newIdMethod)
	s.catalogService.AssertNotCalled(s.T(), findBookByIDMethod, context.TODO(), order.BookID)
	s.paymentClient.AssertNotCalled(s.T(), createPaymentIntentMethod, context.TODO(), &order)
	s.repo.AssertNotCalled(s.T(), createOrderMethod, context.TODO(), &order)
}

func (s *ShopTestSuite) TestCreateOrder_WhenCatalogServiceFails() {
	request := shop.CreateOrder{
		BookID: "some-id",
	}
	orderId := "orderId"
	userId := "some-user-id"
	order := request.Order(orderId, userId)

	ctx := context.WithValue(context.Background(), "userId", userId)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return(orderId)
	s.catalogService.On(findBookByIDMethod, ctx, order.BookID).Return(catalog.BookResponse{}, fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(ctx, request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertCalled(s.T(), validateMethod, request)
	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, ctx, order.BookID)
	s.paymentClient.AssertNotCalled(s.T(), createPaymentIntentMethod, ctx, &order)
	s.repo.AssertNotCalled(s.T(), createOrderMethod, ctx, &order)
}

func (s *ShopTestSuite) TestCreateOrder_WhenPaymentClientFails() {
	request := shop.CreateOrder{
		BookID: "some-id",
	}
	orderId := "orderId"
	userId := "userId"
	order := request.Order(orderId, userId)
	book := catalog.BookResponse{ID: "some-id", Price: 2000}

	ctx := context.WithValue(context.Background(), "userId", userId)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return(orderId)
	s.catalogService.On(findBookByIDMethod, ctx, book.ID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.paymentClient.On(createPaymentIntentMethod, ctx, &updatedOrder).Return(fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(ctx, request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertCalled(s.T(), validateMethod, request)
	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, ctx, &updatedOrder)
	s.catalogService.AssertNotCalled(s.T(), findBookByIDMethod, ctx, order.ID)
	s.repo.AssertNotCalled(s.T(), createOrderMethod, ctx, &order)
}

func (s *ShopTestSuite) TestCreateOrder_WhenRepositoryFails() {
	request := shop.CreateOrder{
		BookID: "some-id",
	}
	orderId := "orderId"
	userId := "some-user-id"
	order := request.Order(orderId, userId)
	book := catalog.BookResponse{ID: "some-id", Price: 2000}

	ctx := context.WithValue(context.Background(), "userId", userId)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return(orderId)
	s.catalogService.On(findBookByIDMethod, ctx, book.ID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.paymentClient.On(createPaymentIntentMethod, ctx, &updatedOrder).Return(nil)
	s.repo.On(createOrderMethod, ctx, &updatedOrder).Return(fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(ctx, request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertCalled(s.T(), validateMethod, request)
	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, ctx, &updatedOrder)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, ctx, order.BookID)
	s.repo.AssertCalled(s.T(), createOrderMethod, ctx, &updatedOrder)
}

func (s *ShopTestSuite) TestCreateOrder_Successfully() {
	request := shop.CreateOrder{
		BookID: "some-id",
	}
	orderId := "orderId"
	userId := "some-user-id"
	order := request.Order(orderId, userId)
	book := catalog.BookResponse{ID: "some-id", Price: 4000}

	ctx := context.WithValue(context.Background(), "userId", userId)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return(orderId)
	s.catalogService.On(findBookByIDMethod, ctx, order.BookID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.paymentClient.On(createPaymentIntentMethod, ctx, &updatedOrder).Return(nil)
	s.repo.On(createOrderMethod, ctx, &updatedOrder).Return(nil)

	expected := shop.NewOrderResponse(updatedOrder)
	actual, err := s.shop.CreateOrder(ctx, request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)

	s.validator.AssertCalled(s.T(), validateMethod, request)
	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, ctx, &updatedOrder)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, ctx, order.BookID)
	s.repo.AssertCalled(s.T(), createOrderMethod, ctx, &updatedOrder)
}

func (s *ShopTestSuite) TestCompleteOrder_WhenOrderCouldNotBeFound() {
	id := "some-order-id"
	s.repo.On(findOrderByIDMethod, context.TODO(), id).Return(shop.Order{}, fmt.Errorf("some error"))

	err := s.shop.CompleteOrder(context.TODO(), id)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), id)
	s.repo.AssertNumberOfCalls(s.T(), updateOrderMethod, 0)
}

func (s *ShopTestSuite) TestCompleteOrder_WhenUpdateFails() {
	order := shop.Order{
		ID: "some-id",
	}
	s.repo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)

	order.Complete()
	s.repo.On(updateOrderMethod, context.TODO(), &order).Return(fmt.Errorf("some error"))

	err := s.shop.CompleteOrder(context.TODO(), order.ID)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.repo.AssertCalled(s.T(), updateOrderMethod, context.TODO(), &order)
}

func (s *ShopTestSuite) TestCompleteOrder_Successfully() {
	order := shop.Order{
		ID: "some-id",
	}
	s.repo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)

	order.Complete()
	s.repo.On(updateOrderMethod, context.TODO(), &order).Return(nil)

	err := s.shop.CompleteOrder(context.TODO(), order.ID)

	assert.Nil(s.T(), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.repo.AssertCalled(s.T(), updateOrderMethod, context.TODO(), &order)
}

func (s *ShopTestSuite) TestGetOrderDeliverableContent_WhenOrderCouldNotBeFound() {
	order := shop.Order{
		ID:     "some-id",
		BookID: "some-book-id",
	}
	s.repo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(shop.Order{}, fmt.Errorf("some error"))

	_, err := s.shop.GetOrderDeliverableContent(context.TODO(), order.ID)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.catalogService.AssertNotCalled(s.T(), getBookContent, context.TODO(), order.BookID)
}

func (s *ShopTestSuite) TestGetOrderDeliverableContent_WhenOrderIsNotPaid() {
	order := shop.Order{
		ID:     "some-id",
		BookID: "some-book-id",
		Status: shop.Pending,
	}
	s.repo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)

	_, err := s.shop.GetOrderDeliverableContent(context.TODO(), order.ID)

	assert.Equal(s.T(), shop.ErrOrderNotPaid, err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.catalogService.AssertNotCalled(s.T(), getBookContent, context.TODO(), order.BookID)
}

func (s *ShopTestSuite) TestGetOrderDeliverableContent_WithError() {
	order := shop.Order{
		ID:     "some-id",
		BookID: "some-book-id",
		Status: shop.Paid,
	}
	s.repo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, context.TODO(), order.BookID).Return(nil, fmt.Errorf("some error"))

	_, err := s.shop.GetOrderDeliverableContent(context.TODO(), order.ID)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.catalogService.AssertCalled(s.T(), getBookContent, context.TODO(), order.BookID)
}

func (s *ShopTestSuite) TestGetOrderDeliverableContent_NonAdmin_Forbidden() {
	order := shop.Order{
		ID:     "some-id",
		BookID: "some-book-id",
		UserID: "some-user-id",
		Status: shop.Paid,
	}

	ctx := context.WithValue(context.Background(), "admin", false)
	ctx = context.WithValue(ctx, "userId", "some-user-id2")

	s.repo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, ctx, order.BookID).Return(io.NopCloser(strings.NewReader("test")), nil)

	_, err := s.shop.GetOrderDeliverableContent(ctx, order.ID)

	assert.Equal(s.T(), shop.ErrForbiddenOrderAccess, err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
	s.catalogService.AssertNotCalled(s.T(), getBookContent, ctx, order.BookID)
}

func (s *ShopTestSuite) TestGetOrderDeliverableContent_NonAdmin_Successfully() {
	order := shop.Order{
		ID:     "some-id",
		BookID: "some-book-id",
		UserID: "some-user-id",
		Status: shop.Paid,
	}

	ctx := context.WithValue(context.Background(), "admin", false)
	ctx = context.WithValue(ctx, "userId", "some-user-id")

	s.repo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, ctx, order.BookID).Return(io.NopCloser(strings.NewReader("test")), nil)

	actual, err := s.shop.GetOrderDeliverableContent(ctx, order.ID)

	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
	s.catalogService.AssertCalled(s.T(), getBookContent, ctx, order.BookID)
}

func (s *ShopTestSuite) TestGetOrderDeliverableContent_Admin_Successfully() {
	order := shop.Order{
		ID:     "some-id",
		BookID: "some-book-id",
		UserID: "some-user-id",
		Status: shop.Paid,
	}

	ctx := context.WithValue(context.Background(), "admin", true)

	s.repo.On(findOrderByIDMethod, ctx, order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, ctx, order.BookID).Return(io.NopCloser(strings.NewReader("test")), nil)

	actual, err := s.shop.GetOrderDeliverableContent(ctx, order.ID)

	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, ctx, order.ID)
	s.catalogService.AssertCalled(s.T(), getBookContent, ctx, order.BookID)
}
