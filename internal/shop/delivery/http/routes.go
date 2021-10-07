package http

import "github.com/gin-gonic/gin"

func (h ShopHandler) Routes(r *gin.RouterGroup) {
	r.GET("/orders", h.getOrders)
	r.GET("/orders/:id", h.getOrder)
	r.POST("/orders", h.createOrder)
}
