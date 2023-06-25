package catalog

import (
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
