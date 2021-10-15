//go:build unit
// +build unit

package dto

import (
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromBook(t *testing.T) {
	book := factory.NewBook()

	expected := BookResponse{
		ID:              book.ID,
		Title:           book.Title,
		Description:     book.Description,
		AuthorName:      book.AuthorName,
		PosterImageLink: book.PosterImageLink,
		Price:           book.Price,
		ReleaseDate:     book.ReleaseDate,
		CreatedAt:       book.CreatedAt,
		UpdatedAt:       book.UpdatedAt,
	}

	actual := FromBook(book)

	assert.Equal(t, expected, actual)
}

func TestFromPaginatedBooks(t *testing.T) {
	book1 := factory.NewBook()
	book2 := factory.NewBook()
	book3 := factory.NewBook()

	paginatedBooks := model.PaginatedBooks{
		Books:      []model.Book{book1, book2, book3},
		Limit:      10,
		Offset:     0,
		TotalBooks: 3,
	}

	expected := PaginatedBooksResponse{
		Results:     []BookResponse{FromBook(book1), FromBook(book2), FromBook(book3)},
		CurrentPage: 1,
		PerPage:     10,
		TotalPages:  1,
		TotalItems:  3,
	}
	actual := FromPaginatedBooks(paginatedBooks)

	assert.Equal(t, expected, actual)
}
