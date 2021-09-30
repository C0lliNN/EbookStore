package http

import (
	"github.com/gin-gonic/gin"
)

func (h CatalogHandler) AuthRoutes(engine *gin.RouterGroup) {
	engine.GET("/books", h.getBooks)
	engine.GET("/books/:id", h.getBook)
}

func (h CatalogHandler) AdminRoutes(engine *gin.RouterGroup) {
	engine.POST("/books", h.createBook)
	engine.PATCH("/books/:id", h.updateBook)
	engine.DELETE("/books/:id", h.deleteBook)
}
