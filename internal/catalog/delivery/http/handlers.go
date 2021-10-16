package http

import (
	"github.com/c0llinn/ebook-store/internal/catalog/delivery/dto"
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/internal/common"
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

type IDGenerator interface {
	NewID() string
}

type CatalogHandler struct {
	service     Service
	idGenerator IDGenerator
}

func NewCatalogHandler(service Service, idGenerator IDGenerator) CatalogHandler {
	return CatalogHandler{
		service:     service,
		idGenerator: idGenerator,
	}
}

// getBooks godoc
// @Summary Fetch Books
// @Tags Catalog
// @Produce  json
// @Param payload body dto.SearchBooks true "Filters"
// @Success 200 {object} dto.PaginatedBooksResponse
// @Failure 500 {object} api.Error
// @Router /books [get]
func (h CatalogHandler) getBooks(context *gin.Context) {
	var s dto.SearchBooks
	if err := context.ShouldBindQuery(&s); err != nil {
		context.Error(&common.ErrNotValid{Input: "SearchBooks", Err: err})
		return
	}

	paginatedBooks, err := h.service.FindBooks(s.ToDomain())
	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, dto.FromPaginatedBooks(paginatedBooks))
}

// getBook godoc
// @Summary Fetch Book by ID
// @Tags Catalog
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} dto.BookResponse
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /books/{id} [get]
func (h CatalogHandler) getBook(context *gin.Context) {
	book, err := h.service.FindBookByID(context.Param("id"))
	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, dto.FromBook(book))
}

// createBook godoc
// @Summary Create a new Book
// @Tags Catalog
// @Accept mpfd
// @Produce  json
// @Param payload formData dto.CreateBook true "Book Payload"
// @Param poster formData file true "Book Poster"
// @Param content formData file true "Book Content in PDF"
// @Success 201 {object} dto.BookResponse
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /books [post]
func (h CatalogHandler) createBook(context *gin.Context) {
	var c dto.CreateBook
	if err := context.ShouldBind(&c); err != nil {
		context.Error(&common.ErrNotValid{Input: "CreateBook", Err: err})
		return
	}

	poster, err := context.FormFile("poster")
	if err != nil {
		context.Error(err)
		return
	}

	posterFile, err := poster.Open()
	if err != nil {
		context.Error(err)
		return
	}

	content, err := context.FormFile("content")
	if err != nil {
		context.Error(err)
	}

	contentFile, err := content.Open()
	if err != nil {
		context.Error(err)
		return
	}

	book := c.ToDomain(h.idGenerator.NewID())
	if err = h.service.CreateBook(&book, posterFile, contentFile); err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusCreated, dto.FromBook(book))
}

// updateBook godoc
// @Summary Update the provided Book
// @Tags Catalog
// @Accept json
// @Produce  json
// @Param payload body dto.UpdateBook true "Book Payload"
// @Param id path string true "Book ID"
// @Success 204 "Success"
// @Failure 400 {object} api.Error
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /books/{id} [patch]
func (h CatalogHandler) updateBook(context *gin.Context) {
	var u dto.UpdateBook
	if err := context.ShouldBindJSON(&u); err != nil {
		context.Error(&common.ErrNotValid{Input: "UpdateBook", Err: err})
		return
	}

	existingBook, err := h.service.FindBookByID(context.Param("id"))
	if err != nil {
		context.Error(err)
		return
	}

	newBook := u.ToDomain(existingBook)
	if err = h.service.UpdateBook(&newBook); err != nil {
		context.Error(err)
		return
	}

	context.Status(http.StatusNoContent)
}

// deleteBook godoc
// @Summary Delete a Book
// @Tags Catalog
// @Produce  json
// @Param id path string true "Book ID"
// @Success 204 "Success"
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /books/{id} [delete]
func (h CatalogHandler) deleteBook(context *gin.Context) {
	if err := h.service.DeleteBook(context.Param("id")); err != nil {
		context.Error(err)
		return
	}

	context.Status(http.StatusNoContent)
}
