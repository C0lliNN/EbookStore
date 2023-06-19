package catalog

import (
	"io"
	"time"
)

type SearchBooks struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	AuthorName  string `form:"authorName"`
	Page        int    `form:"page"`
	PerPage     int    `form:"perPage"`
}

func (s *SearchBooks) BookQuery() BookQuery {
	if s.Page == 0 {
		s.Page = 1
	}

	if s.PerPage == 0 {
		s.PerPage = 10
	}

	return BookQuery{
		Title:       s.Title,
		Description: s.Description,
		AuthorName:  s.AuthorName,
		Limit:       s.PerPage,
		Offset:      (s.Page - 1) * s.PerPage,
	}
}

type CreateBook struct {
	Title       string    `form:"title" validate:"required,max=100"`
	Description string    `form:"description" validate:"required"`
	AuthorName  string    `form:"authorName" validate:"required,max=100"`
	Price       int       `form:"price" validate:"required,gt=0"`
	ReleaseDate time.Time `form:"releaseDate" validate:"required"`

	PosterImage io.ReadSeeker `form:"-" swaggerignore:"true"`
	BookContent io.ReadSeeker `form:"-" swaggerignore:"true"`
}

func (c CreateBook) Book(id string) Book {
	return Book{
		ID:          id,
		Title:       c.Title,
		Description: c.Description,
		AuthorName:  c.AuthorName,
		Price:       c.Price,
		ReleaseDate: c.ReleaseDate,
	}
}

type UpdateBook struct {
	ID          string
	Title       *string `form:"title" validate:"omitempty,max=100"`
	Description *string `form:"description" validate:"omitempty"`
	AuthorName  *string `form:"authorName" validate:"omitempty,max=100"`
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

	return updated
}
