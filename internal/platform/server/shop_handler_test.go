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

func (s *ServerSuiteTest) TestCreateOrder_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	s.createCart(token, book)
	order := s.createOrder(token)

	s.Equal(string(shop.Pending), order.Status)
	s.Equal(book.ID, order.Items[0].ID)
	s.Equal(book.Price, int(order.TotalPrice))
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
	s.createCart(token, book)
	order := s.createOrder(token)

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
	s.createCart(token, book)
	expected := s.createOrder(token)
	var actual shop.OrderResponse

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/"+expected.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&actual)

	s.Equal(expected.Items[0].ID, actual.Items[0].ID)
	s.Equal(expected.Status, actual.Status)
	s.Equal(expected.TotalPrice, actual.TotalPrice)
}

func (s *ServerSuiteTest) TestDownloadOrder_Unknown() {
	token := s.createDefaultCustomer()

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/id1/items/id2/download").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.message", "the provided order was not found")).
		End()
}

func (s *ServerSuiteTest) TestDownloadOrder_NotPaid() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	s.createCart(token, book)
	order := s.createOrder(token)

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/"+order.ID+"/items/"+book.ID+"/download").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusPaymentRequired).
		Assert(jsonpath.Equal("$.message", "only books from completed orders can be downloaded")).
		End()
}

func (s *ServerSuiteTest) TestDownloadOrder_Unauthorized() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	s.createCart(token, book)
	order := s.createOrder(token)

	result := s.container.DB().Model(&shop.Order{}).Where("id = ?", order.ID).Update("status", shop.Paid)
	require.NoError(s.T(), result.Error)

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/"+order.ID+"/items/"+book.ID+"/download").
		Header("Authorization", fmt.Sprintf("Bearer %v", s.createRandomCustomer())).
		Expect(s.T()).
		Status(http.StatusForbidden).
		Assert(jsonpath.Equal("$.message", "the access to this order is restricted to allowed users")).
		End()
}

func (s *ServerSuiteTest) TestDownloadOrder_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())

	s.createCart(token, book)
	order := s.createOrder(token)
	result := s.container.DB().Model(&shop.Order{}).Where("id = ?", order.ID).Update("status", shop.Paid)
	require.NoError(s.T(), result.Error)

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/orders/"+order.ID+"/items/"+book.ID+"/download").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.url")).
		End()
}

func (s *ServerSuiteTest) TestGetOrders_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	s.createCart(token, book)
	order := s.createOrder(token)

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
	s.Equal(order.Items, response.Results[0].Items)
	s.Equal(order.Status, response.Results[0].Status)
	s.Equal(order.TotalPrice, response.Results[0].TotalPrice)
}

func (s *ServerSuiteTest) TestRemoveItemFromCart_NoExistent() {
	token := s.createDefaultCustomer()

	apitest.New().
		EnableNetworking().
		Delete(s.baseURL+"/api/v1/cart/items/id1").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.message", "item not found in cart")).
		End()
}

func (s *ServerSuiteTest) TestRemoveItemFromCart_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	s.createCart(token, book)

	var response shop.CartResponse

	apitest.New().
		EnableNetworking().
		Delete(s.baseURL+"/api/v1/cart/items/"+book.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&response)

	s.Equal(0, len(response.Items))
}

func (s *ServerSuiteTest) TestGetCart_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())
	s.createCart(token, book)

	var response shop.CartResponse

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/active-cart").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&response)

	s.Equal(book.ID, response.Items[0].ID)
	s.Equal(book.Price, int(response.TotalPrice))
}

func (s *ServerSuiteTest) createCart(customerToken string, book catalog.BookResponse) shop.CartResponse {
	var response shop.CartResponse

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/cart/items/"+book.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", customerToken)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&response)

	return response
}

// createOrder creates a new order using the active cart. Must be called after createCart.
func (s *ServerSuiteTest) createOrder(customerToken string) shop.OrderResponse {
	var response shop.OrderResponse

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/orders").
		Header("Authorization", fmt.Sprintf("Bearer %v", customerToken)).
		Expect(s.T()).
		Status(http.StatusCreated).
		End().
		JSON(&response)

	return response
}
