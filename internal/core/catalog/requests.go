package catalog

import (
	"time"

	"github.com/ebookstore/internal/core/query"
)

type SearchBooks struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	AuthorName  string `form:"authorName"`
	Page        int    `form:"page"`
	PerPage     int    `form:"perPage"`
}

func (s *SearchBooks) CreateQuery() query.Query {
	q := query.New()

	if s.Title != "" {
		q.And(query.Condition{Field: "title", Operator: query.Match, Value: s.Title})
	}

	if s.Description != "" {
		q.And(query.Condition{Field: "description", Operator: query.Match, Value: s.Description})
	}

	if s.AuthorName != "" {
		q.And(query.Condition{Field: "author_name", Operator: query.Match, Value: s.AuthorName})
	}

	return *q
}

func (s *SearchBooks) CreatePage() query.Page {
	p := query.DefaultPage

	if (s.Page > 0) {
		p.Number = s.Page
	}

	if (s.PerPage > 0) {
		p.Size = s.PerPage
	}

	return p
}

type CreateBook struct {
	Title       string    `json:"title" validate:"required,max=100"`
	Description string    `json:"description" validate:"required"`
	AuthorName  string    `json:"authorName" validate:"required,max=100"`
	ContentID   string    `json:"contentId" validate:"required,max=100"`
	Price       int       `json:"price" validate:"required,gt=0"`
	ReleaseDate time.Time `json:"releaseDate" validate:"required"`

	Images []ImageRequest `json:"images"`
}

type ImageRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func (r ImageRequest) Image(bookId string) Image {
	return Image{
		ID:          r.ID,
		Description: r.Description,
		BookID:      bookId,
	}
}

func (c CreateBook) Book(id string) Book {
	images := make([]Image, 0, len(c.Images))
	for _, img := range c.Images {
		images = append(images, img.Image(id))
	}

	return Book{
		ID:          id,
		Title:       c.Title,
		Description: c.Description,
		AuthorName:  c.AuthorName,
		ContentID:   c.ContentID,
		Price:       c.Price,
		ReleaseDate: c.ReleaseDate,
		Images:      images,
	}
}

type UpdateBook struct {
	ID          string
	Title       *string `validate:"omitempty,max=100"`
	Description *string `validate:"omitempty"`
	AuthorName  *string `validate:"omitempty,max=100"`
	Images      []ImageRequest
}

func (u UpdateBook) Update(existing Book) Book {
	updated := existing

	if u.Title != nil && *u.Title != "" {
		updated.Title = *u.Title
	}

	if u.Description != nil && *u.Description != "" {
		updated.Description = *u.Description
	}

	if u.AuthorName != nil && *u.AuthorName != "" {
		updated.AuthorName = *u.AuthorName
	}

	images := make([]Image, 0, len(u.Images))
	for _, img := range u.Images {
		images = append(images, img.Image(existing.ID))
	}

	updated.Images = images

	return updated
}
