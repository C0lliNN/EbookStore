package server_test

import (
	"fmt"
	"net/http"

	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/core/shop"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
)

func (s *ServerSuiteTest) TestCreateOrder_InvalidPayload() {
	token := s.createDefaultCustomer()

	req := shop.CreateOrder{BookID: ""}

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/orders").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		JSON(req).
		Expect(s.T()).
		Status(http.StatusBadRequest).
		Assert(jsonpath.Equal("$.message", "the payload is not valid")).
		End()
}

func (s *ServerSuiteTest) TestCreateOrder_UnknownBook() {
	token := s.createDefaultCustomer()

	req := shop.CreateOrder{BookID: "id"}

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/orders").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		JSON(req).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.message", "the provided book was not found")).
		End()
}

func (s *ServerSuiteTest) TestCreateOrder_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	order := s.createOrder(book, token)

	s.Equal(string(shop.Pending), order.Status)
	s.Equal(book.ID, order.BookID)
	s.Equal(book.Price, int(order.Total))
}

func (s *ServerSuiteTest) TestGetOrder_Unknown() {
	token := s.createDefaultCustomer()

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/id1").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.message", "the provided order was not found")).
		End()
}

func (s *ServerSuiteTest) TestGetOrder_Unauthorized() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	order := s.createOrder(book, token)

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/"+order.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", s.createRandomCustomer())).
		Expect(s.T()).
		Status(http.StatusForbidden).
		Assert(jsonpath.Equal("$.message", "the access to this order is restricted to allowed users")).
		End()
}

func (s *ServerSuiteTest) TestGetOrder_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	expected := s.createOrder(book, token)
	var actual shop.OrderResponse

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/"+expected.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&actual)

	s.Equal(expected.BookID, actual.BookID)
	s.Equal(expected.Status, actual.Status)
	s.Equal(expected.Total, actual.Total)
}

func (s *ServerSuiteTest) TestDownloadOrder_Unknown() {
	token := s.createDefaultCustomer()

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/id1/download").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.message", "the provided order was not found")).
		End()
}

func (s *ServerSuiteTest) TestDownloadOrder_NotPaid() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	order := s.createOrder(book, token)

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/"+order.ID+"/download").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusPaymentRequired).
		Assert(jsonpath.Equal("$.message", "only books from paid orders can be downloaded")).
		End()
}

func (s *ServerSuiteTest) TestDownloadOrder_Unauthorized() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	order := s.createOrder(book, token)
	
	result := s.container.DB().Model(&shop.Order{}).Where("id = ?", order.ID).Update("status", shop.Paid)
	require.NoError(s.T(), result.Error)

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/"+order.ID+"/download").
		Header("Authorization", fmt.Sprintf("Bearer %v", s.createRandomCustomer())).
		Expect(s.T()).
		Status(http.StatusForbidden).
		Assert(jsonpath.Equal("$.message", "the access to this order is restricted to allowed users")).
		End()
}

func (s *ServerSuiteTest) TestDownloadOrder_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	
	order := s.createOrder(book, token)
	result := s.container.DB().Model(&shop.Order{}).Where("id = ?", order.ID).Update("status", shop.Paid)
	require.NoError(s.T(), result.Error)

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/"+order.ID+"/download").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.url")).
		End()
}

func (s *ServerSuiteTest) TestGetOrders_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	order := s.createOrder(book, token)

	var response shop.PaginatedOrdersResponse

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&response)

	s.Equal(1, response.CurrentPage)
	s.Equal(15, response.PerPage)
	s.Equal(int64(1), response.TotalItems)
	s.Equal(1, response.TotalPages)
	s.Equal(order.BookID, response.Results[0].BookID)
	s.Equal(order.Status, response.Results[0].Status)
	s.Equal(order.Total, response.Results[0].Total)
}

func (s *ServerSuiteTest) createOrder(book catalog.BookResponse, customerToken string) shop.OrderResponse {
	req := shop.CreateOrder{BookID: book.ID}

	var response shop.OrderResponse

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/orders").
		Header("Authorization", fmt.Sprintf("Bearer %v", customerToken)).
		JSON(req).
		Expect(s.T()).
		Status(http.StatusCreated).
		End().
		JSON(&response)

	return response
}
