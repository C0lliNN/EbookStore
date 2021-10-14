package http

import "github.com/gin-gonic/gin"

func (h ShopHandler) AuthRoutes(r *gin.RouterGroup) {
	r.GET("/orders", h.getOrders)
	r.GET("/orders/:id", h.getOrder)
	r.POST("/orders", h.createOrder)
}

func (h ShopHandler) UnAuthRoutes(r *gin.RouterGroup) {
	r.POST("/stripe/webhook", h.handleStripeWebhook)
}