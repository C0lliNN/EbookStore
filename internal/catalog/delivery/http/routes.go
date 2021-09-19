package http

import "github.com/gin-gonic/gin"

func (h CatalogHandler) Routes(engine *gin.RouterGroup) {
	engine.GET("/books", h.getBooks)
}
