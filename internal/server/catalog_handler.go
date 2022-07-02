package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ebookstore/internal/catalog"
	"github.com/gin-gonic/gin"
)

type Catalog interface {
	FindBooks(context.Context, catalog.SearchBooks) (catalog.PaginatedBooksResponse, error)
	FindBookByID(context.Context, string) (catalog.BookResponse, error)
	CreateBook(context.Context, catalog.CreateBook) (catalog.BookResponse, error)
	UpdateBook(context.Context, catalog.UpdateBook) error
	DeleteBook(context.Context, string) error
}

type CatalogHandler struct {
	catalog Catalog
}

func NewCatalogHandler(catalog Catalog) *CatalogHandler {
	return &CatalogHandler{
		catalog: catalog,
	}
}

func (h *CatalogHandler) Routes() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/books", Handler: h.getBooks, Public: true},
		{Method: http.MethodGet, Path: "/books/:id", Handler: h.getBook, Public: true},
		{Method: http.MethodPost, Path: "/books", Handler: h.createBook, Public: false},
		{Method: http.MethodPatch, Path: "/books/:id", Handler: h.updateBook, Public: false},
		{Method: http.MethodDelete, Path: "/books/:id", Handler: h.deleteBook, Public: false},
	}
}

// getBooks godoc
// @Summary Fetch Books
// @Tags Catalog
// @Produce  json
// @Param params query catalog.SearchBooks true "Filters"
// @Success 200 {object} catalog.PaginatedBooksResponse
// @Failure 500 {object} ErrorResponse
// @Router /books [get]
func (h *CatalogHandler) getBooks(c *gin.Context) {
	var request catalog.SearchBooks
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(&BindingErr{Err: fmt.Errorf("(getBooks) failed binding query: %w", err)})
		return
	}

	response, err := h.catalog.FindBooks(c, request)
	if err != nil {
		_ = c.Error(fmt.Errorf("(getBooks) failed handling find request: %w ", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// getBook godoc
// @Summary Fetch Book by ID
// @Tags Catalog
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} catalog.BookResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [get]
func (h *CatalogHandler) getBook(c *gin.Context) {
	response, err := h.catalog.FindBookByID(c, c.Param("id"))
	if err != nil {
		_ = c.Error(fmt.Errorf("(getBook) failed handling get request: %w ", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// createBook godoc
// @Summary Create a new Book
// @Tags Catalog
// @Accept mpfd
// @Produce  json
// @Param payload formData catalog.CreateBook true "Book Payload"
// @Param poster formData file true "Book Poster"
// @Param content formData file true "Book Content in PDF"
// @Success 201 {object} catalog.BookResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books [post]
func (h *CatalogHandler) createBook(c *gin.Context) {
	var request catalog.CreateBook
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(&BindingErr{Err: fmt.Errorf("(createBook) failed binding request body: %w", err)})
		return
	}

	poster, err := c.FormFile("poster")
	if err != nil {
		_ = c.Error(fmt.Errorf("(createBook) failed getting poster file: %w", err))
		return
	}

	posterFile, err := poster.Open()
	if err != nil {
		_ = c.Error(fmt.Errorf("(createBook) failed openning poster file: %w", err))
		return
	}

	content, err := c.FormFile("content")
	if err != nil {
		_ = c.Error(fmt.Errorf("(createBook) failed getting content file: %w", err))
	}

	contentFile, err := content.Open()
	if err != nil {
		_ = c.Error(fmt.Errorf("(createBook) failed openning poster file: %w", err))
		return
	}

	request.PosterImage = posterFile
	request.BookContent = contentFile

	response, err := h.catalog.CreateBook(c, request)
	if err != nil {
		_ = c.Error(fmt.Errorf("(createBook) failed handling create request: %w ", err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// updateBook godoc
// @Summary Update the provided Book
// @Tags Catalog
// @Accept json
// @Produce  json
// @Param payload body catalog.UpdateBook true "Book Payload"
// @Param id path string true "Book ID"
// @Success 204 "Success"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [patch]
func (h *CatalogHandler) updateBook(c *gin.Context) {
	var request catalog.UpdateBook
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(&BindingErr{Err: fmt.Errorf("(updateBook) failed binding request body: %w", err)})
		return
	}

	request.ID = c.Param("id")
	if err := h.catalog.UpdateBook(c, request); err != nil {
		_ = c.Error(fmt.Errorf("(updateBook) failed handling update request: %w ", err))
		return
	}

	c.Status(http.StatusNoContent)
}

// deleteBook godoc
// @Summary Delete a Book
// @Tags Catalog
// @Produce  json
// @Param id path string true "Book ID"
// @Success 204 "Success"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [delete]
func (h *CatalogHandler) deleteBook(c *gin.Context) {
	if err := h.catalog.DeleteBook(c, c.Param("id")); err != nil {
		_ = c.Error(fmt.Errorf("(deleteBook) failed handling delete request: %w ", err))
		return
	}

	c.Status(http.StatusNoContent)
}
