package server

import (
	"bytes"
	"context"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/shop"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Shop interface {
	FindOrders(context.Context, shop.SearchOrders) (shop.PaginatedOrdersResponse, error)
	FindOrderByID(context.Context, string) (shop.OrderResponse, error)
	CreateOrder(context.Context, shop.CreateOrder) (shop.OrderResponse, error)
	CompleteOrder(context.Context, string) error
	GetOrderDeliverableContent(context.Context, string) (io.ReadCloser, error)
}

type ShopHandler struct {
	engine *gin.Engine
	shop   Shop
}

func NewShopHandler(engine *gin.Engine, shop Shop) *ShopHandler {
	return &ShopHandler{
		engine: engine,
		shop:   shop,
	}
}

func (h *ShopHandler) Routes() {
	h.engine.GET("/orders", h.getOrders)
	h.engine.GET("/orders/:id", h.getOrder)
	h.engine.GET("/orders/:id/download", h.downloadOrder)
	h.engine.POST("/orders", h.createOrder)
	h.engine.POST("/stripe/webhook", h.handleStripeWebhook)
}

// getOrders godoc
// @Summary Fetch Orders
// @Tags Shop
// @Produce  json
// @Param payload body shop.SearchOrders true "Filters"
// @Success 200 {object} shop.PaginatedOrdersResponse
// @Failure 500 {object} api.Error
// @Router /orders [get]
func (h *ShopHandler) getOrders(c *gin.Context) {
	var request shop.SearchOrders
	if err := c.ShouldBindQuery(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "SearchOrders", Err: err})
		return
	}

	response, err := h.shop.FindOrders(c, request)
	if err != nil {
		c.Error(err)
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
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /orders/{id} [get]
func (h *ShopHandler) getOrder(c *gin.Context) {
	response, err := h.shop.FindOrderByID(c, c.Param("id"))
	if err != nil {
		c.Error(err)
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
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /orders [post]
func (h *ShopHandler) createOrder(c *gin.Context) {
	var request shop.CreateOrder
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "CreateOrder", Err: err})
		return
	}

	response, err := h.shop.CreateOrder(c, request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// downloadOrder godoc
// @Summary Download the book for the given Order
// @Tags Shop
// @Produce  application/pdf
// @Param payload body dto.CreateOrder true "Order Payload"
// @Success 200 "Success"
// @Failure 402 {object} api.Error
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /orders/{id}/download [get]
func (h *ShopHandler) downloadOrder(c *gin.Context) {
	content, err := h.shop.GetOrderDeliverableContent(c, c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	buffer := new(bytes.Buffer)
	if _, err = buffer.ReadFrom(content); err != nil {
		c.Error(err)
		return
	}

	c.DataFromReader(http.StatusOK, int64(buffer.Len()), "application/pdf", buffer, nil)
}

// handleStripeWebhook godoc
// @Summary Handle stripe webhooks
// @Tags Shop
// @Accept json
// @Produce  json
// @Success 200 "Success"
// @Success 500 {object} api.Error
// @Router /stripe/webhook [post]
func (h *ShopHandler) handleStripeWebhook(c *gin.Context) {
	var request shop.HandleStripeWebhook
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	if request.Type == "payment_intent.succeeded" {
		orderID := request.Data["object"].(map[string]interface{})["metadata"].(map[string]interface{})["orderID"].(string)
		if err := h.shop.CompleteOrder(c.Request.Context(), orderID); err != nil {
			c.Error(err)
			return
		}
	}

	c.Status(http.StatusOK)
}
