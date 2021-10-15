//go:build integration
// +build integration

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c0llinn/ebook-store/config/aws"
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/config/log"
	auth "github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/internal/catalog/helper"
	catalog "github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/internal/catalog/repository"
	"github.com/c0llinn/ebook-store/internal/catalog/storage"
	"github.com/c0llinn/ebook-store/internal/catalog/usecase"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/shop/client"
	"github.com/c0llinn/ebook-store/internal/shop/delivery/dto"
	helper3 "github.com/c0llinn/ebook-store/internal/shop/helper"
	"github.com/c0llinn/ebook-store/internal/shop/model"
	repository3 "github.com/c0llinn/ebook-store/internal/shop/repository"
	usecase3 "github.com/c0llinn/ebook-store/internal/shop/usecase"
	"github.com/c0llinn/ebook-store/test"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http/httptest"
	"strings"
	"testing"
)

type ShopHandlerTestSuite struct {
	suite.Suite
	baseURL  string
	context  *gin.Context
	recorder *httptest.ResponseRecorder
	db       *gorm.DB
	s3Client storage.S3Client
	handler  ShopHandler
}

func (s *ShopHandlerTestSuite) SetupTest() {
	test.SetEnvironmentVariables()
	log.InitLogger()
	db.LoadMigrations("file:../../../../migration")

	s.db = db.NewConnection()
	s.baseURL = fmt.Sprintf("http://localhost:%s", viper.GetString("PORT"))

	orderRepository := repository3.NewOrderRepository(s.db)
	stripeClient := client.NewStripeClient()
	bucket := aws.NewBucket()
	bookRepository := repository.NewBookRepository(s.db)
	s3 := aws.NewS3Service()
	s.s3Client = storage.NewS3Client(s3, bucket)
	filenameGenerator := helper.NewFilenameGenerator()
	catalogUseCase := usecase.NewCatalogUseCase(bookRepository, s.s3Client, filenameGenerator)
	shopUseCase := usecase3.NewShopUseCase(orderRepository, stripeClient, catalogUseCase)
	helperIDGenerator := helper3.NewIDGenerator()
	s.handler = NewShopHandler(shopUseCase, helperIDGenerator)

	s.recorder = httptest.NewRecorder()
	s.context, _ = gin.CreateTestContext(s.recorder)
}

func TestShopHandlerRun(t *testing.T) {
	suite.Run(t, new(ShopHandlerTestSuite))
}

func (s *ShopHandlerTestSuite) TearDownTest() {
	s.db.Delete(&model.Order{}, "1 = 1")
	s.db.Delete(&catalog.Book{}, "1 = 1")
	s.db.Delete(&auth.User{}, "1 = 1")
}

func (s *ShopHandlerTestSuite) TestGetOrders_Customer() {
	user := factory.NewUser()
	order1, order2 := factory.NewOrder(), factory.NewOrder()

	order1.UserID = user.ID

	err := s.db.Create(user).Error
	require.Nil(s.T(), err)

	err = s.db.Create(order1).Error
	require.Nil(s.T(), err)

	err = s.db.Create(order2).Error
	require.Nil(s.T(), err)

	s.context.Set("user", user)
	s.context.Request = httptest.NewRequest("GET", s.baseURL+"/orders", nil)

	s.handler.getOrders(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())

	var response dto.PaginatedOrdersResponse
	err = json.Unmarshal(s.recorder.Body.Bytes(), &response)

	assert.Equal(s.T(), order1.ID, response.Results[0].ID)
	assert.Equal(s.T(), 10, response.PerPage)
	assert.Equal(s.T(), 1, response.CurrentPage)
	assert.Equal(s.T(), 1, response.TotalPages)
	assert.Equal(s.T(), int64(1), response.TotalItems)
}

func (s *ShopHandlerTestSuite) TestGetOrders_Admin() {
	user := factory.NewUser()
	user.Role = auth.Admin

	order1, order2 := factory.NewOrder(), factory.NewOrder()

	err := s.db.Create(user).Error
	require.Nil(s.T(), err)

	err = s.db.Create(order1).Error
	require.Nil(s.T(), err)

	err = s.db.Create(order2).Error
	require.Nil(s.T(), err)

	s.context.Set("user", user)
	s.context.Request = httptest.NewRequest("GET", s.baseURL+"/orders", nil)

	s.handler.getOrders(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())

	var response dto.PaginatedOrdersResponse
	err = json.Unmarshal(s.recorder.Body.Bytes(), &response)

	assert.Equal(s.T(), order2.ID, response.Results[0].ID)
	assert.Equal(s.T(), order1.ID, response.Results[1].ID)
	assert.Equal(s.T(), 10, response.PerPage)
	assert.Equal(s.T(), 1, response.CurrentPage)
	assert.Equal(s.T(), 1, response.TotalPages)
	assert.Equal(s.T(), int64(2), response.TotalItems)
}

func (s *ShopHandlerTestSuite) TestGetOrder_CustomerSuccessfully() {
	user := factory.NewUser()

	order1 := factory.NewOrder()
	order1.UserID = user.ID

	err := s.db.Create(user).Error
	require.Nil(s.T(), err)

	err = s.db.Create(order1).Error
	require.Nil(s.T(), err)

	s.context.Set("user", user)
	s.context.Request = httptest.NewRequest("GET", s.baseURL+"/orders/"+order1.ID, nil)
	s.context.Params = []gin.Param{{Key: "id", Value: order1.ID}}

	s.handler.getOrder(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())

	var response dto.OrderResponse
	err = json.Unmarshal(s.recorder.Body.Bytes(), &response)

	assert.Equal(s.T(), order1.ID, response.ID)
}

func (s *ShopHandlerTestSuite) TestGetOrder_CustomerWithError() {
	user := factory.NewUser()

	order1 := factory.NewOrder()

	err := s.db.Create(user).Error
	require.Nil(s.T(), err)

	err = s.db.Create(order1).Error
	require.Nil(s.T(), err)

	s.context.Set("user", user)
	s.context.Request = httptest.NewRequest("GET", s.baseURL+"/orders/"+order1.ID, nil)
	s.context.Params = []gin.Param{{Key: "id", Value: order1.ID}}

	s.handler.getOrder(s.context)

	assert.Equal(s.T(), fmt.Errorf("you don't have permission to see this order"), s.context.Errors.Last().Err)
}

func (s *ShopHandlerTestSuite) TestGetOrder_Admin() {
	user := factory.NewUser()
	user.Role = auth.Admin
	order1 := factory.NewOrder()

	err := s.db.Create(user).Error
	require.Nil(s.T(), err)

	err = s.db.Create(order1).Error
	require.Nil(s.T(), err)

	s.context.Set("user", user)
	s.context.Request = httptest.NewRequest("GET", s.baseURL+"/orders/"+order1.ID, nil)
	s.context.Params = []gin.Param{{Key: "id", Value: order1.ID}}

	s.handler.getOrder(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())

	var response dto.OrderResponse
	err = json.Unmarshal(s.recorder.Body.Bytes(), &response)

	assert.Equal(s.T(), order1.ID, response.ID)
}

func (s *ShopHandlerTestSuite) TestCreateOrder_WithInvalidData() {
	payload := dto.CreateOrder{BookID: ""}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	user := factory.NewUser()
	err = s.db.Create(user).Error
	require.Nil(s.T(), err)

	s.context.Set("user", user)
	s.context.Request = httptest.NewRequest("POST", s.baseURL+"/orders", bytes.NewReader(data))

	s.handler.createOrder(s.context)

	assert.IsType(s.T(), &common.ErrNotValid{}, s.context.Errors.Last().Err)
}

func (s *ShopHandlerTestSuite) TestCreateOrder_Successfully() {
	book := factory.NewBook()
	err := s.db.Create(book).Error
	require.Nil(s.T(), err)

	payload := dto.CreateOrder{BookID: book.ID}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	user := factory.NewUser()
	err = s.db.Create(user).Error
	require.Nil(s.T(), err)

	s.context.Set("user", user)
	s.context.Request = httptest.NewRequest("POST", s.baseURL+"/orders", bytes.NewReader(data))

	s.handler.createOrder(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())

	var response dto.OrderResponse
	err = json.Unmarshal(s.recorder.Body.Bytes(), &response)
	require.Nil(s.T(), err)

	assert.Equal(s.T(), book.ID, response.BookID)
	assert.Equal(s.T(), user.ID, response.UserID)
	assert.Equal(s.T(), int64(book.Price), response.Total)
}

func (s *ShopHandlerTestSuite) TestCompleteOrder_WithInvalidPayload() {
	payload := map[string]interface{}{
		"data": map[string]interface{}{},
		"type": "payment_intent.succeeded",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	s.context.Request = httptest.NewRequest("POST", s.baseURL+"/stripe/webhook", bytes.NewReader(data))

	assert.Panics(s.T(), func() {
		s.handler.handleStripeWebhook(s.context)
	})
}

func (s *ShopHandlerTestSuite) TestCompleteOrder_WithUnknownOrderId() {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"object": map[string]interface{}{
				"metadata": map[string]string{
					"orderID": "some-id",
				},
			},
		},
		"type": "payment_intent.succeeded",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	s.context.Request = httptest.NewRequest("POST", s.baseURL+"/stripe/webhook", bytes.NewReader(data))

	s.handler.handleStripeWebhook(s.context)

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, s.context.Errors.Last().Err)
}

func (s *ShopHandlerTestSuite) TestCompleteOrder_Successfully() {
	order := factory.NewOrder()

	err := s.db.Create(order).Error
	require.Nil(s.T(), err)

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"object": map[string]interface{}{
				"metadata": map[string]string{
					"orderID": order.ID,
				},
			},
		},
		"type": "payment_intent.succeeded",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	s.context.Request = httptest.NewRequest("POST", s.baseURL+"/stripe/webhook", bytes.NewReader(data))

	s.handler.handleStripeWebhook(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())

	var updated model.Order
	err = s.db.First(&updated, "id = ?", order.ID).Error
	require.Nil(s.T(), err)

	assert.Equal(s.T(), model.Paid, updated.Status)
}

func (s *ShopHandlerTestSuite) TestDownloadOrder_WhenOrderIsNotFound() {
	s.context.Params = []gin.Param{{Key: "id", Value: "invalid-order-id"}}

	s.handler.downloadOrder(s.context)

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, s.context.Errors.Last().Err)
}

func (s *ShopHandlerTestSuite) TestDownloadOrder_WhenOrderIsNotPaid() {
	order := factory.NewOrder()
	order.Status = model.Pending

	err := s.db.Create(order).Error
	require.Nil(s.T(), err)

	s.context.Params = []gin.Param{{Key: "id", Value: order.ID}}

	s.handler.downloadOrder(s.context)

	assert.IsType(s.T(), &common.ErrOrderNotPaid{}, s.context.Errors.Last().Err)
}

func (s *ShopHandlerTestSuite) TestDownloadOrder_Successfully() {
	book := factory.NewBook()
	err := s.db.Create(book).Error
	require.Nil(s.T(), err)

	err = s.s3Client.SaveFile(book.ContentBucketKey, strings.NewReader("something"))
	require.Nil(s.T(), err)

	order := factory.NewOrder()
	order.Status = model.Paid
	order.BookID = book.ID

	err = s.db.Create(order).Error
	require.Nil(s.T(), err)

	s.context.Params = []gin.Param{{Key: "id", Value: order.ID}}

	s.handler.downloadOrder(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())
}
