//go:build unit
// +build unit

package dto

import (
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SearchBooksTestSuite struct {
	suite.Suite
}

func TestSearchBookRun(t *testing.T) {
	suite.Run(t, new(SearchBooksTestSuite))
}

func (s *SearchBooksTestSuite) TestToDomain_WithNoFields() {
	dto := SearchBooks{}

	expected := model.BookQuery{
		Limit: 10,
	}
	actual := dto.ToDomain()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestToDomain_WithTitle() {
	dto := SearchBooks{Title: "some-title"}

	expected := model.BookQuery{
		Title: "some-title",
		Limit: 10,
	}
	actual := dto.ToDomain()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestToDomain_WithDescription() {
	dto := SearchBooks{Description: "some-description"}

	expected := model.BookQuery{
		Description: "some-description",
		Limit:       10,
	}
	actual := dto.ToDomain()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestToDomain_WithAuthorName() {
	dto := SearchBooks{AuthorName: "some-name"}

	expected := model.BookQuery{
		AuthorName: "some-name",
		Limit:      10,
	}
	actual := dto.ToDomain()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestToDomain_WithPage() {
	dto := SearchBooks{Page: 4}

	expected := model.BookQuery{
		Limit:  10,
		Offset: 30,
	}
	actual := dto.ToDomain()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestToDomain_WithPerPage() {
	dto := SearchBooks{PerPage: 20}

	expected := model.BookQuery{
		Limit: 20,
	}
	actual := dto.ToDomain()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestToDomain_WithMultipleFields() {
	dto := SearchBooks{Page: 3, PerPage: 20, Title: "some-title"}

	expected := model.BookQuery{
		Limit:  20,
		Offset: 40,
		Title:  "some-title",
	}
	actual := dto.ToDomain()

	assert.Equal(s.T(), expected, actual)
}

type CreateBookTestSuite struct {
	suite.Suite
}

func TestCreateBookRun(t *testing.T) {
	suite.Run(t, new(CreateBookTestSuite))
}

func (s *CreateBookTestSuite) TestToDomain() {
	id := uuid.NewString()
	dto := CreateBook{
		Title:       "some-title",
		Description: "description",
		AuthorName:  "fake name",
		Price:       10000,
		ReleaseDate: time.Date(2021, time.September, 29, 17, 28, 0, 0, time.UTC),
	}

	expected := model.Book{
		ID:          id,
		Title:       "some-title",
		Description: "description",
		AuthorName:  "fake name",
		Price:       10000,
		ReleaseDate: time.Date(2021, time.September, 29, 17, 28, 0, 0, time.UTC),
	}

	actual := dto.ToDomain(id)

	assert.Equal(s.T(), expected, actual)
}

type UpdateBookTestSuite struct {
	suite.Suite
}

func TestUpdateBookRun(t *testing.T) {
	suite.Run(t, new(UpdateBookTestSuite))
}

func (s *UpdateBookTestSuite) TestToDomain_WithNoFields() {
	book := factory.NewBook()
	dto := UpdateBook{}

	expected := book
	actual := dto.ToDomain(book)

	assert.Equal(s.T(), expected, actual)
}

func (s *UpdateBookTestSuite) TestToDomain_WithTitle() {
	book := factory.NewBook()

	title := "some-title"
	dto := UpdateBook{Title: &title}

	expected := book
	expected.Title = title
	actual := dto.ToDomain(book)

	assert.Equal(s.T(), expected, actual)
}

func (s *UpdateBookTestSuite) TestToDomain_WithDescription() {
	book := factory.NewBook()

	description := "some-description"
	dto := UpdateBook{Description: &description}

	expected := book
	expected.Description = description
	actual := dto.ToDomain(book)

	assert.Equal(s.T(), expected, actual)
}

func (s *UpdateBookTestSuite) TestToDomain_WithAuthorName() {
	book := factory.NewBook()

	authorName := "new name"
	dto := UpdateBook{AuthorName: &authorName}

	expected := book
	expected.AuthorName = authorName
	actual := dto.ToDomain(book)

	assert.Equal(s.T(), expected, actual)
}

func (s *UpdateBookTestSuite) TestToDomain_WithMultipleFields() {
	book := factory.NewBook()

	title := "title"
	authorName := "new name"
	dto := UpdateBook{Title: &title, AuthorName: &authorName}

	expected := book
	expected.Title = title
	expected.AuthorName = authorName
	actual := dto.ToDomain(book)

	assert.Equal(s.T(), expected, actual)
}
