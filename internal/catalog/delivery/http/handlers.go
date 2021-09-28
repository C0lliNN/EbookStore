package http

import (
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Service interface {
	FindBooks(query model.BookQuery) (paginatedBooks model.PaginatedBooks, err error)
	FindBookByID(id string) (book model.Book, err error)
	CreateBook(book *model.Book, posterImage io.ReadSeeker, bookContent io.ReadSeeker) error
	UpdateBook(book *model.Book) error
	DeleteBook(id string) error
}

type CatalogHandler struct {
	service Service
}

func NewCatalogHandler(service Service) CatalogHandler {
	return CatalogHandler{service: service}
}

func (h CatalogHandler) getBooks(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}
