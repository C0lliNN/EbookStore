package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ebookstore/internal/core/shop"
	"github.com/ebookstore/internal/log"
	"github.com/gin-gonic/gin"
)

type Shop interface {
	FindOrders(context.Context, shop.SearchOrders) (shop.PaginatedOrdersResponse, error)
	FindOrderByID(context.Context, string) (shop.OrderResponse, error)
	CreateOrder(context.Context, shop.CreateOrder) (shop.OrderResponse, error)
	CompleteOrder(context.Context, string) error
	GetOrderDeliverableContent(context.Context, string) (shop.ShopBookResponse, error)
}

type ShopHandler struct {
	shop Shop
}

func NewShopHandler(shop Shop) *ShopHandler {
	return &ShopHandler{
		shop: shop,
	}
}

func (h *ShopHandler) Routes() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/orders", Handler: h.getOrders, Public: false},
		{Method: http.MethodGet, Path: "/orders/:id", Handler: h.getOrder, Public: false},
		{Method: http.MethodGet, Path: "/orders/:id/download", Handler: h.downloadOrder, Public: false},
		{Method: http.MethodPost, Path: "/orders", Handler: h.createOrder, Public: false},
		{Method: http.MethodPost, Path: "/stripe/webhook", Handler: h.handleStripeWebhook, Public: true},
	}
}

// getOrders godoc
// @Summary Fetch Orders
// @Tags Shop
// @Produce  json
// @Param params query shop.SearchOrders true "Filters"
// @Success 200 {object} shop.PaginatedOrdersResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/orders [get]
func (h *ShopHandler) getOrders(c *gin.Context) {
	var request shop.SearchOrders
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(&BindingErr{Err: fmt.Errorf("(getOrders) failed binding query: %w", err)})
		return
	}

	response, err := h.shop.FindOrders(c, request)
	if err != nil {
		_ = c.Error(fmt.Errorf("(getOrders) failed handling find request: %w", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// getOrder godoc
// @Summary Fetch Order by ID
// @Tags Shop
// @Produce  json
// @Param id path string true "orderId ID"
// @Success 200 {object} shop.OrderResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/orders/{id} [get]
func (h *ShopHandler) getOrder(c *gin.Context) {
	response, err := h.shop.FindOrderByID(c, c.Param("id"))
	if err != nil {
		_ = c.Error(fmt.Errorf("(getOrder) failed handling find request: %w", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// createOrder godoc
// @Summary Create a new Order
// @Tags Shop
// @Accept json
// @Produce  json
// @Param payload body shop.CreateOrder true "Order Payload"
// @Success 201 {object} shop.OrderResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/orders [post]
func (h *ShopHandler) createOrder(c *gin.Context) {
	var request shop.CreateOrder
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(&BindingErr{Err: fmt.Errorf("(createOrder) failed binding request: %w", err)})
		return
	}

	response, err := h.shop.CreateOrder(c, request)
	if err != nil {
		_ = c.Error(fmt.Errorf("(createOrder) failed handling create request: %w", err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// downloadOrder godoc
// @Summary Download the book for the given Order
// @Tags Shop
// @Produce  json
// @Param id path string true "orderId ID"
// @Success 200 {object} shop.ShopBookResponse
// @Failure 402 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/orders/{id}/download [get]
func (h *ShopHandler) downloadOrder(c *gin.Context) {
	response, err := h.shop.GetOrderDeliverableContent(c, c.Param("id"))
	if err != nil {
		_ = c.Error(fmt.Errorf("(downloadOrder) failed handling get book content: %w", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// handleStripeWebhook godoc
// @Summary Handle stripe webhooks
// @Tags Shop
// @Accept json
// @Produce  json
// @Success 200 "Success"
// @Success 500 {object} ErrorResponse
// @Router /api/v1/stripe/webhook [post]
func (h *ShopHandler) handleStripeWebhook(c *gin.Context) {
	log.FromContext(c).Infof("processing new stripe webhook request")

	var request shop.HandleStripeWebhook
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(&BindingErr{Err: fmt.Errorf("(handleStripeWebhook) failed binding request: %w", err)})
		return
	}

	if request.Type == "payment_intent.succeeded" {
		orderID := request.Data["object"].(map[string]interface{})["metadata"].(map[string]interface{})["orderID"].(string)
		if err := h.shop.CompleteOrder(c, orderID); err != nil {
			_ = c.Error(fmt.Errorf("(handleStripeWebhook) failed handling complete order request: %w", err))
			return
		}
	}

	c.Status(http.StatusOK)
}
