package http

import "github.com/gin-gonic/gin"

func (h ShopHandler) AuthRoutes(r *gin.RouterGroup) {
	r.GET("/orders", h.getOrders)
	r.GET("/orders/:id", h.getOrder)
	r.GET("/orders/:id/download", h.downloadOrder)
	r.POST("/orders", h.createOrder)
}

func (h ShopHandler) Routes(r *gin.Engine) {
	r.POST("/stripe/webhook", h.handleStripeWebhook)
}
