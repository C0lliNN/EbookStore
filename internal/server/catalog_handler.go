package server

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Catalog interface {
	FindBooks(context.Context, catalog.SearchBooks) (catalog.PaginatedBooksResponse, error)
	FindBookByID(context.Context, string) (catalog.BookResponse, error)
	CreateBook(context.Context, catalog.CreateBook) (catalog.BookResponse, error)
	UpdateBook(context.Context, catalog.UpdateBook) error
	DeleteBook(context.Context, string) error
}

type CatalogHandler struct {
	engine  *gin.Engine
	catalog Catalog
}

func NewCatalogHandler(engine *gin.Engine, catalog Catalog) *CatalogHandler {
	return &CatalogHandler{
		engine:  engine,
		catalog: catalog,
	}
}

func (h *CatalogHandler) Routes() {
	h.engine.GET("/books", h.getBooks)
	h.engine.GET("/books/:id", h.getBook)

	// Admin Routes
	h.engine.POST("/books", h.createBook)
	h.engine.PATCH("/books/:id", h.updateBook)
	h.engine.DELETE("/books/:id", h.deleteBook)
}

// getBooks godoc
// @Summary Fetch Books
// @Tags Catalog
// @Produce  json
// @Param payload body catalog.SearchBooks true "Filters"
// @Success 200 {object} catalog.PaginatedBooksResponse
// @Failure 500 {object} api.Error
// @Router /books [get]
func (h *CatalogHandler) getBooks(c *gin.Context) {
	var request catalog.SearchBooks
	if err := c.ShouldBindQuery(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "SearchBooks", Err: err})
		return
	}

	response, err := h.catalog.FindBooks(c.Request.Context(), request)
	if err != nil {
		c.Error(err)
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
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /books/{id} [get]
func (h *CatalogHandler) getBook(c *gin.Context) {
	response, err := h.catalog.FindBookByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Error(err)
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
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /books [post]
func (h *CatalogHandler) createBook(c *gin.Context) {
	var request catalog.CreateBook
	if err := c.ShouldBind(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "CreateBook", Err: err})
		return
	}

	poster, err := c.FormFile("poster")
	if err != nil {
		c.Error(err)
		return
	}

	posterFile, err := poster.Open()
	if err != nil {
		c.Error(err)
		return
	}

	content, err := c.FormFile("content")
	if err != nil {
		c.Error(err)
	}

	contentFile, err := content.Open()
	if err != nil {
		c.Error(err)
		return
	}

	request.PosterImage = posterFile
	request.BookContent = contentFile

	response, err := h.catalog.CreateBook(c.Request.Context(), request)
	if err != nil {
		c.Error(err)
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
// @Failure 400 {object} api.Error
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /books/{id} [patch]
func (h *CatalogHandler) updateBook(c *gin.Context) {
	var request catalog.UpdateBook
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "UpdateBook", Err: err})
		return
	}

	request.ID = c.Param("id")
	if err := h.catalog.UpdateBook(c.Request.Context(), request); err != nil {
		c.Error(err)
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
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /books/{id} [delete]
func (h *CatalogHandler) deleteBook(c *gin.Context) {
	if err := h.catalog.DeleteBook(c.Request.Context(), c.Param("id")); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}