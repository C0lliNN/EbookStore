package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CatalogHandler struct {
}

func NewCatalogHandler() CatalogHandler {
	return CatalogHandler{}
}

func (h CatalogHandler) getBooks(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}
