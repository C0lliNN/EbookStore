package usecase

import (
	"fmt"
	catalog "github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/internal/shop/mock"
	"github.com/c0llinn/ebook-store/internal/shop/model"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	findOrdersByQueryMethod = "FindByQuery"
	findOrderByIDMethod     = "FindByID"
	createOrderMethod         = "Create"
	updateOrderMethod         = "Update"
	findBookByIDMethod        = "FindBookByID"
	createPaymentIntentMethod = "CreatePaymentIntentForOrder"
)

type ShopUseCaseTestSuite struct {
	suite.Suite
	repo           *mock.OrderRepository
	paymentClient  *mock.PaymentClient
	catalogService *mock.CatalogService
	useCase        ShopUseCase
}

func (s *ShopUseCaseTestSuite) SetupTest() {
	s.repo = new(mock.OrderRepository)
	s.paymentClient = new(mock.PaymentClient)
	s.catalogService = new(mock.CatalogService)

	s.useCase = NewShopUseCase(s.repo, s.paymentClient, s.catalogService)
}

func TestShopUseCaseRun(t *testing.T) {
	suite.Run(t, new(ShopUseCaseTestSuite))
}

func (s *ShopUseCaseTestSuite) TestFindOrders_Successfully() {
	paginated := factory.NewPaginatedOrders()
	s.repo.On(findOrdersByQueryMethod, model.OrderQuery{}).Return(paginated, nil)

	actual, err := s.useCase.FindOrders(model.OrderQuery{})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), paginated, actual)
	s.repo.AssertCalled(s.T(), findOrdersByQueryMethod, model.OrderQuery{})
}

func (s *ShopUseCaseTestSuite) TestFindOrders_WithError() {
	s.repo.On(findOrdersByQueryMethod, model.OrderQuery{}).Return(model.PaginatedOrders{}, fmt.Errorf("some error"))

	_, err := s.useCase.FindOrders(model.OrderQuery{})

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrdersByQueryMethod, model.OrderQuery{})
}

func (s *ShopUseCaseTestSuite) TestFindOrderByID_Successfully() {
	order := factory.NewOrder()
	s.repo.On(findOrderByIDMethod, order.ID).Return(order, nil)

	actual, err := s.useCase.FindOrderByID(order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), order, actual)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, order.ID)
}

func (s *ShopUseCaseTestSuite) TestFindOrderByID_WithError() {
	id := uuid.NewString()
	s.repo.On(findOrderByIDMethod, id).Return(model.Order{}, fmt.Errorf("some error"))

	_, err := s.useCase.FindOrderByID(id)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, id)
}

func (s *ShopUseCaseTestSuite) TestCreateOrder_WhenPaymentClientFails() {
	order := factory.NewOrder()
	s.paymentClient.On(createPaymentIntentMethod, &order).Return(fmt.Errorf("some error"))

	err := s.useCase.CreateOrder(&order)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, &order)
	s.catalogService.AssertNotCalled(s.T(), findBookByIDMethod, order.ID)
	s.repo.AssertNotCalled(s.T(), createOrderMethod, &order)
}

func (s *ShopUseCaseTestSuite) TestCreateOrder_WhenCatalogServiceFails() {
	order := factory.NewOrder()
	s.paymentClient.On(createPaymentIntentMethod, &order).Return(nil)
	s.catalogService.On(findBookByIDMethod, order.BookID).Return(catalog.Book{}, fmt.Errorf("some error"))

	err := s.useCase.CreateOrder(&order)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, &order)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, order.BookID)
	s.repo.AssertNotCalled(s.T(), createOrderMethod, &order)
}

func (s *ShopUseCaseTestSuite) TestCreateOrder_WhenRepositoryFails() {
	order := factory.NewOrder()
	book := factory.NewBook()
	s.paymentClient.On(createPaymentIntentMethod, &order).Return(nil)
	s.catalogService.On(findBookByIDMethod, order.BookID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.repo.On(createOrderMethod, &updatedOrder).Return(fmt.Errorf("some error"))

	err := s.useCase.CreateOrder(&order)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, &order)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, order.BookID)
	s.repo.AssertCalled(s.T(), createOrderMethod, &updatedOrder)
}

func (s *ShopUseCaseTestSuite) TestCreateOrder_Successfully() {
	order := factory.NewOrder()
	book := factory.NewBook()
	s.paymentClient.On(createPaymentIntentMethod, &order).Return(nil)
	s.catalogService.On(findBookByIDMethod, order.BookID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.repo.On(createOrderMethod, &updatedOrder).Return(nil)

	err := s.useCase.CreateOrder(&order)

	assert.Nil(s.T(), err)

	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, &order)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, order.BookID)
	s.repo.AssertCalled(s.T(), createOrderMethod, &updatedOrder)
}

func (s *ShopUseCaseTestSuite) TestUpdateOrder_Successfully() {
	order := factory.NewOrder()
	s.repo.On(updateOrderMethod, &order).Return(nil)

	err := s.useCase.UpdateOrder(&order)

	assert.Nil(s.T(), err)
	s.repo.AssertCalled(s.T(), updateOrderMethod, &order)
}

func (s *ShopUseCaseTestSuite) TestUpdateOrder_WithError() {
	order := factory.NewOrder()
	s.repo.On(updateOrderMethod, &order).Return(fmt.Errorf("some error"))

	err := s.useCase.UpdateOrder(&order)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), updateOrderMethod, &order)
}