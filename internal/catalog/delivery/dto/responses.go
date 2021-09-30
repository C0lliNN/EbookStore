package dto

import (
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"math"
	"time"
)

type BookResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	AuthorName      string    `json:"authorName"`
	PosterImageLink string    `json:"posterImageLink"`
	Price           int       `json:"price"`
	ReleaseDate     time.Time `json:"releaseDate"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func FromBook(book model.Book) BookResponse {
	return BookResponse{
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
}

type PaginatedBooksResponse struct {
	Results     []BookResponse `json:"results"`
	CurrentPage int            `json:"currentPage"`
	PerPage     int            `json:"perPage"`
	TotalPages  int            `json:"totalPages"`
	TotalItems  int64          `json:"totalItems"`
}

func FromPaginatedBooks(paginatedBooks model.PaginatedBooks) PaginatedBooksResponse {
	books := make([]BookResponse, 0, len(paginatedBooks.Books))
	for _, b := range paginatedBooks.Books {
		books = append(books, FromBook(b))
	}

	return PaginatedBooksResponse{
		Results:     books,
		CurrentPage: (paginatedBooks.Offset / paginatedBooks.Limit) + 1,
		PerPage:     paginatedBooks.Limit,
		TotalPages:  int(math.Ceil(float64(paginatedBooks.TotalBooks) / float64(paginatedBooks.Limit))),
		TotalItems:  paginatedBooks.TotalBooks,
	}
}
