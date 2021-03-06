package catalog

import (
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

func (s *SearchBooksTestSuite) TestBookQuery_WithNoFields() {
	dto := SearchBooks{}

	expected := BookQuery{
		Limit: 10,
	}
	actual := dto.BookQuery()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestBookQuery_WithTitle() {
	dto := SearchBooks{Title: "some-title"}

	expected := BookQuery{
		Title: "some-title",
		Limit: 10,
	}
	actual := dto.BookQuery()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestBookQuery_WithDescription() {
	dto := SearchBooks{Description: "some-description"}

	expected := BookQuery{
		Description: "some-description",
		Limit:       10,
	}
	actual := dto.BookQuery()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestBookQuery_WithAuthorName() {
	dto := SearchBooks{AuthorName: "some-name"}

	expected := BookQuery{
		AuthorName: "some-name",
		Limit:      10,
	}
	actual := dto.BookQuery()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestBookQuery_WithPage() {
	dto := SearchBooks{Page: 4}

	expected := BookQuery{
		Limit:  10,
		Offset: 30,
	}
	actual := dto.BookQuery()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestBookQuery_WithPerPage() {
	dto := SearchBooks{PerPage: 20}

	expected := BookQuery{
		Limit: 20,
	}
	actual := dto.BookQuery()

	assert.Equal(s.T(), expected, actual)
}

func (s *SearchBooksTestSuite) TestBookQuery_WithMultipleFields() {
	dto := SearchBooks{Page: 3, PerPage: 20, Title: "some-title"}

	expected := BookQuery{
		Limit:  20,
		Offset: 40,
		Title:  "some-title",
	}
	actual := dto.BookQuery()

	assert.Equal(s.T(), expected, actual)
}

type CreateBookTestSuite struct {
	suite.Suite
}

func TestCreateBookRun(t *testing.T) {
	suite.Run(t, new(CreateBookTestSuite))
}

func (s *CreateBookTestSuite) TestBook() {
	id := uuid.NewString()
	dto := CreateBook{
		Title:       "some-title",
		Description: "description",
		AuthorName:  "fake name",
		Price:       10000,
		ReleaseDate: time.Date(2021, time.September, 29, 17, 28, 0, 0, time.UTC),
	}

	expected := Book{
		ID:          id,
		Title:       "some-title",
		Description: "description",
		AuthorName:  "fake name",
		Price:       10000,
		ReleaseDate: time.Date(2021, time.September, 29, 17, 28, 0, 0, time.UTC),
	}

	actual := dto.Book(id)

	assert.Equal(s.T(), expected, actual)
}

type UpdateBookTestSuite struct {
	suite.Suite
}

func TestUpdateBookRun(t *testing.T) {
	suite.Run(t, new(UpdateBookTestSuite))
}

func (s *UpdateBookTestSuite) TestUpdate_WithNoFields() {
	book := Book{ID: "some-id"}
	dto := UpdateBook{}

	expected := book
	actual := dto.Update(book)

	assert.Equal(s.T(), expected, actual)
}

func (s *UpdateBookTestSuite) TestUpdate_WithTitle() {
	book := Book{
		ID: "some-id",
		Title: "other-title",
	}

	title := "some-title"
	dto := UpdateBook{Title: &title}

	expected := book
	expected.Title = title
	actual := dto.Update(book)

	assert.Equal(s.T(), expected, actual)
}

func (s *UpdateBookTestSuite) TestUpdate_WithDescription() {
	book := Book{
		ID: "some-id",
		Description: "other-description",
	}

	description := "some-description"
	dto := UpdateBook{Description: &description}

	expected := book
	expected.Description = description
	actual := dto.Update(book)

	assert.Equal(s.T(), expected, actual)
}

func (s *UpdateBookTestSuite) TestUpdate_WithAuthorName() {
	book := Book{
		ID: "some-id",
		AuthorName: "other-name",
	}

	authorName := "new name"
	dto := UpdateBook{AuthorName: &authorName}

	expected := book
	expected.AuthorName = authorName
	actual := dto.Update(book)

	assert.Equal(s.T(), expected, actual)
}

func (s *UpdateBookTestSuite) TestUpdate_WithMultipleFields() {
	book := Book{
		ID: "some-id",
		Title: "other-title",
		AuthorName: "other-name",
	}

	title := "title"
	authorName := "new name"
	dto := UpdateBook{Title: &title, AuthorName: &authorName}

	expected := book
	expected.Title = title
	expected.AuthorName = authorName
	actual := dto.Update(book)

	assert.Equal(s.T(), expected, actual)
}
