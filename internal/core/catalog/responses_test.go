package catalog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBookResponse(t *testing.T) {
	book := Book{
		ID:          "some-id",
		Title:       "Clean Code",
		Description: "Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		ContentID:   "some-key2",
		Images: []Image{
			{
				ID:          "some-id",
				Description: "some-description",
				BookID:      "some-id",
			},
		},
		Price:       40000,
		ReleaseDate: time.Date(2020, time.September, 23, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2022, time.September, 23, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2022, time.September, 23, 20, 0, 0, 0, time.UTC),
	}

	expected := BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		AuthorName:  book.AuthorName,
		Images: []ImageResponse{
			{
				ID:          "some-id",
				Description: "some-description",
				Link:        "link",
			},
		},
		Price:       book.Price,
		ReleaseDate: book.ReleaseDate,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}

	actual := NewBookResponse(book, []string{"link"})

	assert.Equal(t, expected, actual)
}

func TestNewImageResponse(t *testing.T) {
	image := Image{
		ID:          "some-id",
		Description: "some-description",
		BookID:      "some-book-id",
	}

	expected := ImageResponse{
		ID:          image.ID,
		Description: image.Description,
		Link:        "link",
	}

	actual := NewImageResponse(image, "link")

	assert.Equal(t, expected, actual)
}

func TestNewPaginatedBooksResponse(t *testing.T) {
	book1 := Book{
		ID:          "some-id",
		Title:       "Clean Code",
		Description: "Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		ContentID:   "some-key2",
		Price:       40000,
		ReleaseDate: time.Date(2020, time.September, 23, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2022, time.September, 23, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2022, time.September, 23, 20, 0, 0, 0, time.UTC),
	}

	book2 := Book{
		ID:          "some-id2",
		Title:       "Clean Coder",
		Description: "Professional Guide",
		AuthorName:  "Robert C. Martin",
		Images: []Image{
			{
				ID:          "some-id",
				Description: "some-description",
				BookID:      "some-id2",
			},
		},
		ContentID:   "some-key22",
		Price:       45000,
		ReleaseDate: time.Date(2020, time.September, 22, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2022, time.September, 22, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2022, time.September, 22, 20, 0, 0, 0, time.UTC),
	}

	book3 := Book{
		ID:          "some-id4",
		Title:       "Clean Architecture",
		Description: "Architecture Guide",
		AuthorName:  "Robert C. Martin",
		ContentID:   "some-key23",
		Price:       60000,
		ReleaseDate: time.Date(2021, time.September, 23, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2023, time.September, 23, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, time.September, 23, 20, 0, 0, 0, time.UTC),
	}

	paginatedBooks := PaginatedBooks{
		Books:      []Book{book1, book2, book3},
		Limit:      10,
		Offset:     0,
		TotalBooks: 3,
	}

	expected := PaginatedBooksResponse{
		Results:     []BookResponse{NewBookResponse(book1, nil), NewBookResponse(book2, []string{"link"}), NewBookResponse(book3, nil)},
		CurrentPage: 1,
		PerPage:     10,
		TotalPages:  1,
		TotalItems:  3,
	}
	actual := NewPaginatedBooksResponse(paginatedBooks, map[string][]string{"some-id2": {"link"}})

	assert.Equal(t, expected, actual)
}
