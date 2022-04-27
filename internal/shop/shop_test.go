package shop_test

import (
	"context"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"github.com/c0llinn/ebook-store/internal/shop"
	mocks "github.com/c0llinn/ebook-store/mocks/shop"
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
)

type ShopTestSuite struct {
	suite.Suite
	repo           *mocks.Repository
	paymentClient  *mocks.PaymentClient
	catalogService *mocks.CatalogService
	idGenerator    *mocks.IDGenerator
	shop           *shop.Shop
}

func (s *ShopTestSuite) SetupTest() {
	s.repo = new(mocks.Repository)
	s.paymentClient = new(mocks.PaymentClient)
	s.catalogService = new(mocks.CatalogService)
	s.idGenerator = new(mocks.IDGenerator)

	s.shop = shop.New(shop.Config{
		Repository:     s.repo,
		PaymentClient:  s.paymentClient,
		CatalogService: s.catalogService,
		IDGenerator:    s.idGenerator,
	})
}

func TestShop(t *testing.T) {
	suite.Run(t, new(ShopTestSuite))
}

func (s *ShopTestSuite) TestFindOrders_Successfully() {
	request := shop.SearchOrders{}
	query := request.OrderQuery()

	paginatedOrders := shop.PaginatedOrders{
		Orders: []shop.Order{{ID: "some-id"}},
		Limit:  10,
	}
	s.repo.On(findOrdersByQueryMethod, context.TODO(), query).Return(paginatedOrders, nil)

	expected := shop.NewPaginatedOrdersResponse(paginatedOrders)
	actual, err := s.shop.FindOrders(context.TODO(), request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repo.AssertCalled(s.T(), findOrdersByQueryMethod, context.TODO(), query)
}

func (s *ShopTestSuite) TestFindOrders_WithError() {
	request := shop.SearchOrders{}
	query := request.OrderQuery()

	s.repo.On(findOrdersByQueryMethod, context.TODO(), query).Return(shop.PaginatedOrders{}, fmt.Errorf("some error"))

	_, err := s.shop.FindOrders(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrdersByQueryMethod, context.TODO(), query)
}

func (s *ShopTestSuite) TestFindOrderByID_Successfully() {
	order := shop.Order{
		ID:    "order-rid",
		Total: 4000,
	}
	s.repo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)

	expected := shop.NewOrderResponse(order)
	actual, err := s.shop.FindOrderByID(context.TODO(), order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
}

func (s *ShopTestSuite) TestFindOrderByID_WithError() {
	id := "some-id"
	s.repo.On(findOrderByIDMethod, context.TODO(), id).Return(shop.Order{}, fmt.Errorf("some error"))

	_, err := s.shop.FindOrderByID(context.TODO(), id)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), id)
}

func (s *ShopTestSuite) TestCreateOrder_WhenCatalogServiceFails() {
	request := shop.CreateOrder{
		BookID: "some-id",
	}
	orderId := "orderId"
	userId := "userId"
	order := request.Order(orderId, userId)

	s.idGenerator.On(newIdMethod).Return(orderId)
	s.catalogService.On(findBookByIDMethod, context.TODO(), order.BookID).Return(catalog.Book{}, fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, context.TODO(), order.BookID)
	s.paymentClient.AssertNotCalled(s.T(), createPaymentIntentMethod, context.TODO(), &order)
	s.repo.AssertNotCalled(s.T(), createOrderMethod, context.TODO(), &order)
}

func (s *ShopTestSuite) TestCreateOrder_WhenPaymentClientFails() {
	request := shop.CreateOrder{
		BookID: "some-id",
	}
	orderId := "orderId"
	userId := "userId"
	order := request.Order(orderId, userId)
	book := catalog.Book{ID: "some-id", Price: 2000}

	s.idGenerator.On(newIdMethod).Return(orderId)
	s.catalogService.On(findBookByIDMethod, context.TODO(), book.ID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.paymentClient.On(createPaymentIntentMethod, context.TODO(), &updatedOrder).Return(fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, context.TODO(), &updatedOrder)
	s.catalogService.AssertNotCalled(s.T(), findBookByIDMethod, context.TODO(), order.ID)
	s.repo.AssertNotCalled(s.T(), createOrderMethod, context.TODO(), &order)
}

func (s *ShopTestSuite) TestCreateOrder_WhenRepositoryFails() {
	request := shop.CreateOrder{
		BookID: "some-id",
	}
	orderId := "orderId"
	userId := "userId"
	order := request.Order(orderId, userId)
	book := catalog.Book{ID: "some-id", Price: 2000}

	s.idGenerator.On(newIdMethod).Return(orderId)
	s.catalogService.On(findBookByIDMethod, context.TODO(), book.ID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.paymentClient.On(createPaymentIntentMethod, context.TODO(), &updatedOrder).Return(nil)
	s.repo.On(createOrderMethod, context.TODO(), &updatedOrder).Return(fmt.Errorf("some error"))

	_, err := s.shop.CreateOrder(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, context.TODO(), &updatedOrder)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, context.TODO(), order.BookID)
	s.repo.AssertCalled(s.T(), createOrderMethod, context.TODO(), &updatedOrder)
}

func (s *ShopTestSuite) TestCreateOrder_Successfully() {
	request := shop.CreateOrder{
		BookID: "some-id",
	}
	orderId := "orderId"
	userId := "userId"
	order := request.Order(orderId, userId)
	book := catalog.Book{ID: "some-id", Price: 4000}

	s.idGenerator.On(newIdMethod).Return(orderId)
	s.catalogService.On(findBookByIDMethod, context.TODO(), order.BookID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.paymentClient.On(createPaymentIntentMethod, context.TODO(), &updatedOrder).Return(nil)
	s.repo.On(createOrderMethod, context.TODO(), &updatedOrder).Return(nil)

	expected := shop.NewOrderResponse(updatedOrder)
	actual, err := s.shop.CreateOrder(context.TODO(), request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)

	s.idGenerator.AssertCalled(s.T(), newIdMethod)
	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, context.TODO(), &updatedOrder)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, context.TODO(), order.BookID)
	s.repo.AssertCalled(s.T(), createOrderMethod, context.TODO(), &updatedOrder)
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

func (s *ShopTestSuite) TestGetOrderDeliverableContent_Successfully() {
	order := shop.Order{
		ID:     "some-id",
		BookID: "some-book-id",
		Status: shop.Paid,
	}
	s.repo.On(findOrderByIDMethod, context.TODO(), order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, context.TODO(), order.BookID).Return(io.NopCloser(strings.NewReader("test")), nil)

	actual, err := s.shop.GetOrderDeliverableContent(context.TODO(), order.ID)

	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, context.TODO(), order.ID)
	s.catalogService.AssertCalled(s.T(), getBookContent, context.TODO(), order.BookID)
}
