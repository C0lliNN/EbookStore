package catalog

import (
	"math"
	"time"
)

type BookResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	AuthorName  string          `json:"authorName"`
	Images      []ImageResponse `json:"images"`
	Price       int             `json:"price"`
	ReleaseDate time.Time       `json:"releaseDate"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func NewBookResponse(book Book, links []string) BookResponse {
	images := make([]ImageResponse, 0, len(book.Images))
	for i := range book.Images {
		images = append(images, NewImageResponse(book.Images[i], links[i]))
	}

	return BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		AuthorName:  book.AuthorName,
		Images:      images,
		Price:       book.Price,
		ReleaseDate: book.ReleaseDate,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}
}

type ImageResponse struct {
	ID          string `json:"id"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

func NewImageResponse(image Image, link string) ImageResponse {
	return ImageResponse{
		ID:          image.ID,
		Link:        link,
		Description: image.Description,
	}
}

type PaginatedBooksResponse struct {
	Results     []BookResponse `json:"results"`
	CurrentPage int            `json:"currentPage"`
	PerPage     int            `json:"perPage"`
	TotalPages  int            `json:"totalPages"`
	TotalItems  int64          `json:"totalItems"`
}

func NewPaginatedBooksResponse(paginatedBooks PaginatedBooks, imageLinks map[string][]string) PaginatedBooksResponse {
	books := make([]BookResponse, 0, len(paginatedBooks.Books))
	for _, b := range paginatedBooks.Books {
		books = append(books, NewBookResponse(b, imageLinks[b.ID]))
	}

	return PaginatedBooksResponse{
		Results:     books,
		CurrentPage: (paginatedBooks.Offset / paginatedBooks.Limit) + 1,
		PerPage:     paginatedBooks.Limit,
		TotalPages:  int(math.Ceil(float64(paginatedBooks.TotalBooks) / float64(paginatedBooks.Limit))),
		TotalItems:  paginatedBooks.TotalBooks,
	}
}

type PresignURLResponse struct {
	URL string `json:"url"`
}
