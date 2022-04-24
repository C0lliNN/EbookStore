//go:build integration
// +build integration

package persistence

import (
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/internal/common"
	config2 "github.com/c0llinn/ebook-store/internal/config"
	"github.com/c0llinn/ebook-store/test"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BookRepositoryTestSuite struct {
	suite.Suite
	repo BookRepository
}

func (s *BookRepositoryTestSuite) SetupTest() {
	test.SetEnvironmentVariables()
	config2.InitLogger()
	config2.LoadMigrations("file:../../../migrations")

	conn := config2.NewConnection()
	s.repo = BookRepository{conn}
}

func TestBookRepositoryRun(t *testing.T) {
	suite.Run(t, new(BookRepositoryTestSuite))
}

func (s *BookRepositoryTestSuite) TearDownTest() {
	s.repo.db.Delete(&model.Book{}, "1 = 1")
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithEmptyQuery() {
	book1 := factory.NewBook()
	book2 := factory.NewBook()
	book3 := factory.NewBook()

	err := s.repo.Create(&book1)
	err = s.repo.Create(&book2)
	err = s.repo.Create(&book3)
	require.Nil(s.T(), err)

	paginatedBooks, err := s.repo.FindByQuery(model.BookQuery{})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(paginatedBooks.Books))
	assert.Equal(s.T(), int64(3), paginatedBooks.TotalBooks)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithTitle() {
	book1 := factory.NewBook()
	book1.Title = "some title"
	book2 := factory.NewBook()
	book3 := factory.NewBook()

	err := s.repo.Create(&book1)
	err = s.repo.Create(&book2)
	err = s.repo.Create(&book3)
	require.Nil(s.T(), err)

	paginatedBooks, err := s.repo.FindByQuery(model.BookQuery{Title: "title"})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(paginatedBooks.Books))
	assert.Equal(s.T(), int64(1), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book1.ID, paginatedBooks.Books[0].ID)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithDescription() {
	book1 := factory.NewBook()
	book1.Description = "some description"
	book2 := factory.NewBook()
	book3 := factory.NewBook()

	err := s.repo.Create(&book1)
	err = s.repo.Create(&book2)
	err = s.repo.Create(&book3)
	require.Nil(s.T(), err)

	paginatedBooks, err := s.repo.FindByQuery(model.BookQuery{Description: "some"})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(paginatedBooks.Books))
	assert.Equal(s.T(), int64(1), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book1.ID, paginatedBooks.Books[0].ID)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithAuthorName() {
	book1 := factory.NewBook()
	book1.AuthorName = "Raphael Collin"
	book2 := factory.NewBook()
	book3 := factory.NewBook()

	err := s.repo.Create(&book1)
	err = s.repo.Create(&book2)
	err = s.repo.Create(&book3)
	require.Nil(s.T(), err)

	paginatedBooks, err := s.repo.FindByQuery(model.BookQuery{AuthorName: "Raphael Collin"})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(paginatedBooks.Books))
	assert.Equal(s.T(), int64(1), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book1.ID, paginatedBooks.Books[0].ID)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithLimit() {
	book1 := factory.NewBook()
	book2 := factory.NewBook()
	book3 := factory.NewBook()

	err := s.repo.Create(&book1)
	err = s.repo.Create(&book2)
	err = s.repo.Create(&book3)
	require.Nil(s.T(), err)

	paginatedBooks, err := s.repo.FindByQuery(model.BookQuery{Limit: 2})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(paginatedBooks.Books))
	assert.Equal(s.T(), 2, paginatedBooks.Limit)
	assert.Equal(s.T(), int64(3), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book1.ID, paginatedBooks.Books[0].ID)
	assert.Equal(s.T(), book2.ID, paginatedBooks.Books[1].ID)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithOffset() {
	book1 := factory.NewBook()
	book2 := factory.NewBook()
	book3 := factory.NewBook()

	err := s.repo.Create(&book1)
	err = s.repo.Create(&book2)
	err = s.repo.Create(&book3)
	require.Nil(s.T(), err)

	paginatedBooks, err := s.repo.FindByQuery(model.BookQuery{Offset: 1})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(paginatedBooks.Books))
	assert.Equal(s.T(), 1, paginatedBooks.Offset)
	assert.Equal(s.T(), int64(3), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book2.ID, paginatedBooks.Books[0].ID)
	assert.Equal(s.T(), book3.ID, paginatedBooks.Books[1].ID)
}

func (s *BookRepositoryTestSuite) TestFindByID_WithInvalidID() {
	id := uuid.NewString()

	_, err := s.repo.FindByID(id)

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, err)
}

func (s *BookRepositoryTestSuite) TestFindByID_WithValidID() {
	id := uuid.NewString()

	_, err := s.repo.FindByID(id)

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, err)
}

func (s *BookRepositoryTestSuite) TestCreate_Successfully() {
	book := factory.NewBook()

	err := s.repo.Create(&book)
	assert.Nil(s.T(), err)

	persisted, err := s.repo.FindByID(book.ID)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), book.ID, persisted.ID)
}

func (s *BookRepositoryTestSuite) TestUpdate_Successfully() {
	book := factory.NewBook()

	err := s.repo.Create(&book)
	assert.Nil(s.T(), err)

	book.Title = "new title"
	err = s.repo.Update(&book)

	persisted, err := s.repo.FindByID(book.ID)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), book.ID, persisted.ID)
	assert.Equal(s.T(), book.Title, persisted.Title)
}

func (s *BookRepositoryTestSuite) TestDelete_WithInvalidID() {
	id := uuid.NewString()

	err := s.repo.Delete(id)

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, err)
}

func (s *BookRepositoryTestSuite) TestDelete_WithValidID() {
	book := factory.NewBook()

	err := s.repo.Create(&book)
	require.Nil(s.T(), err)

	err = s.repo.Delete(book.ID)

	assert.Nil(s.T(), err)
}
