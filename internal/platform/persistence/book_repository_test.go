package persistence_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/core/query"
	"github.com/ebookstore/internal/platform/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BookRepositoryTestSuite struct {
	PostgresRepositoryTestSuite
	repo *persistence.BookRepository
}

func (s *BookRepositoryTestSuite) SetupSuite() {
	s.PostgresRepositoryTestSuite.SetupSuite()

	s.repo = persistence.NewBookRepository(s.db)
}

func TestBookRepositoryRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	suite.Run(t, new(BookRepositoryTestSuite))
}

func (s *BookRepositoryTestSuite) TearDownTest() {
	s.db.Delete(&catalog.Book{}, "1 = 1")
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithEmptyQuery() {
	ctx := context.TODO()

	book1 := catalog.Book{
		ID:          "some-id1",
		Title:       "Clean Code",
		AuthorName:  "Robert c. Martin",
		Price:       5500,
		ReleaseDate: time.Date(2017, time.January, 20, 0, 0, 0, 0, time.UTC),
	}
	book2 := catalog.Book{
		ID:         "some-id2",
		Title:      "Clean Coder",
		AuthorName: "Robert c. Martin",
		Images: []catalog.Image{
			{
				ID:          "some-id2",
				Description: "poster",
				BookID:      "some-id1",
			},
		},
		Price:       7000,
		ReleaseDate: time.Date(2017, time.February, 12, 0, 0, 0, 0, time.UTC),
	}
	book3 := catalog.Book{
		ID:          "some-id3",
		Title:       "DOmain Driver Design",
		AuthorName:  "Eric Evans",
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.UTC),
	}

	err := s.repo.Create(ctx, &book1)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book2)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book3)
	require.Nil(s.T(), err)

	paginatedBooks, err := s.repo.FindByQuery(ctx, *query.New(), query.DefaultPage)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(paginatedBooks.Books))
	assert.Equal(s.T(), int64(3), paginatedBooks.TotalBooks)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithTitle() {
	ctx := context.TODO()

	book1 := catalog.Book{
		ID:          "some-id1",
		Title:       "Clean Code",
		AuthorName:  "Robert c. Martin",
		Price:       5500,
		ReleaseDate: time.Date(2017, time.January, 20, 0, 0, 0, 0, time.UTC),
	}
	book2 := catalog.Book{
		ID:          "some-id2",
		Title:       "Clean Coder",
		AuthorName:  "Robert c. Martin",
		Price:       7000,
		ReleaseDate: time.Date(2017, time.February, 12, 0, 0, 0, 0, time.UTC),
	}
	book3 := catalog.Book{
		ID:    "some-id3",
		Title: "Domain Driver Design",
		Images: []catalog.Image{
			{
				ID:          "some-id2",
				Description: "poster",
				BookID:      "some-id3",
			},
		},
		AuthorName:  "Eric Evans",
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.UTC),
	}

	err := s.repo.Create(ctx, &book1)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book2)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book3)
	require.Nil(s.T(), err)

	q := *query.New().And(query.Condition{Field: "title", Operator: query.Match, Value: "Domain Driver Design"})
	paginatedBooks, err := s.repo.FindByQuery(ctx, q, query.DefaultPage)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(paginatedBooks.Books))
	assert.Equal(s.T(), int64(1), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book3.ID, paginatedBooks.Books[0].ID)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithDescription() {
	ctx := context.TODO()

	book1 := catalog.Book{
		ID:          "some-id1",
		Title:       "Clean Code",
		Description: "Craftsman",
		Images: []catalog.Image{
			{
				ID:          "some-id2",
				Description: "poster",
				BookID:      "some-id1",
			},
		},
		AuthorName:  "Robert c. Martin",
		Price:       5500,
		ReleaseDate: time.Date(2017, time.January, 20, 0, 0, 0, 0, time.UTC),
	}
	book2 := catalog.Book{
		ID:          "some-id2",
		Title:       "Clean Coder",
		AuthorName:  "Robert c. Martin",
		Price:       7000,
		ReleaseDate: time.Date(2017, time.February, 12, 0, 0, 0, 0, time.UTC),
	}
	book3 := catalog.Book{
		ID:          "some-id3",
		Title:       "Domain Driver Design",
		AuthorName:  "Eric Evans",
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.UTC),
	}

	err := s.repo.Create(ctx, &book1)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book2)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book3)
	require.Nil(s.T(), err)

	q := *query.New().And(query.Condition{Field: "description", Operator: query.Match, Value: "Craftsman"})
	paginatedBooks, err := s.repo.FindByQuery(ctx, q, query.DefaultPage)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(paginatedBooks.Books))
	assert.Equal(s.T(), int64(1), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book1.ID, paginatedBooks.Books[0].ID)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithAuthorName() {
	ctx := context.TODO()

	book1 := catalog.Book{
		ID:    "some-id1",
		Title: "Clean Code",
		Images: []catalog.Image{
			{
				ID:          "some-id2",
				Description: "poster",
				BookID:      "some-id1",
			},
		},
		Description: "Craftsman",
		AuthorName:  "Robert c. Martin",
		Price:       5500,
		ReleaseDate: time.Date(2017, time.January, 20, 0, 0, 0, 0, time.UTC),
	}
	book2 := catalog.Book{
		ID:          "some-id2",
		Title:       "Clean Coder",
		AuthorName:  "Robert c. Martin",
		Price:       7000,
		ReleaseDate: time.Date(2017, time.February, 12, 0, 0, 0, 0, time.UTC),
	}
	book3 := catalog.Book{
		ID:          "some-id3",
		Title:       "Domain Driver Design",
		AuthorName:  "Eric Evans",
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.UTC),
	}

	err := s.repo.Create(ctx, &book1)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book2)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book3)
	require.Nil(s.T(), err)

	q := *query.New().And(query.Condition{Field: "author_name", Operator: query.Match, Value: "Eric Evans"})
	paginatedBooks, err := s.repo.FindByQuery(ctx, q, query.DefaultPage)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(paginatedBooks.Books))
	assert.Equal(s.T(), int64(1), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book3.ID, paginatedBooks.Books[0].ID)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithPageSize() {
	ctx := context.TODO()

	book1 := catalog.Book{
		ID:          "some-id1",
		Title:       "Clean Code",
		Description: "Craftsman",
		AuthorName:  "Robert c. Martin",
		Price:       5500,
		ReleaseDate: time.Date(2017, time.January, 20, 0, 0, 0, 0, time.UTC),
	}
	book2 := catalog.Book{
		ID:          "some-id2",
		Title:       "Clean Coder",
		AuthorName:  "Robert c. Martin",
		Price:       7000,
		ReleaseDate: time.Date(2017, time.February, 12, 0, 0, 0, 0, time.UTC),
	}
	book3 := catalog.Book{
		ID:          "some-id3",
		Title:       "Domain Driver Design",
		AuthorName:  "Eric Evans",
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.UTC),
	}

	err := s.repo.Create(ctx, &book1)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book2)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book3)
	require.Nil(s.T(), err)

	page := query.DefaultPage
	page.Size = 2

	paginatedBooks, err := s.repo.FindByQuery(ctx, *query.New(), page)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(paginatedBooks.Books))
	assert.Equal(s.T(), 2, paginatedBooks.Limit)
	assert.Equal(s.T(), int64(3), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book1.ID, paginatedBooks.Books[0].ID)
	assert.Equal(s.T(), book2.ID, paginatedBooks.Books[1].ID)
}

func (s *BookRepositoryTestSuite) TestFindByQuery_WithPageNumber() {
	ctx := context.TODO()

	book1 := catalog.Book{
		ID:          "some-id1",
		Title:       "Clean Code",
		Description: "Craftsman",
		AuthorName:  "Robert c. Martin",
		Price:       5500,
		ReleaseDate: time.Date(2017, time.January, 20, 0, 0, 0, 0, time.UTC),
	}
	book2 := catalog.Book{
		ID:          "some-id2",
		Title:       "Clean Coder",
		AuthorName:  "Robert c. Martin",
		Price:       7000,
		ReleaseDate: time.Date(2017, time.February, 12, 0, 0, 0, 0, time.UTC),
	}
	book3 := catalog.Book{
		ID:          "some-id3",
		Title:       "Domain Driver Design",
		AuthorName:  "Eric Evans",
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.UTC),
	}

	err := s.repo.Create(ctx, &book1)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book2)
	require.Nil(s.T(), err)
	err = s.repo.Create(ctx, &book3)
	require.Nil(s.T(), err)

	page := query.DefaultPage
	page.Number = 2
	page.Size = 1
	paginatedBooks, err := s.repo.FindByQuery(ctx, *query.New(), page)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(paginatedBooks.Books))
	assert.Equal(s.T(), 1, paginatedBooks.Offset)
	assert.Equal(s.T(), int64(3), paginatedBooks.TotalBooks)
	assert.Equal(s.T(), book2.ID, paginatedBooks.Books[0].ID)
}

func (s *BookRepositoryTestSuite) TestFindByID_WithInvalidID() {
	ctx := context.TODO()

	id := "some-id"

	_, err := s.repo.FindByID(ctx, id)

	assert.IsType(s.T(), &persistence.ErrEntityNotFound{}, errors.Unwrap(err))
}

func (s *BookRepositoryTestSuite) TestFindByID_WithValidID() {
	ctx := context.TODO()

	expected := catalog.Book{
		ID:    "some-id",
		Title: "Domain Driver Design",
		Images: []catalog.Image{
			{
				ID:          "some-id2",
				Description: "poster",
				BookID:      "some-id",
			},
		},
		AuthorName:  "Eric Evans",
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.Local),
	}

	err := s.repo.Create(ctx, &expected)
	assert.Nil(s.T(), err)

	actual, err := s.repo.FindByID(ctx, expected.ID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected.ID, actual.ID)
	assert.Equal(s.T(), expected.Title, actual.Title)
	assert.Equal(s.T(), expected.Description, actual.Description)
	assert.Equal(s.T(), expected.Price, actual.Price)
	assert.Equal(s.T(), expected.Images, actual.Images)
}

func (s *BookRepositoryTestSuite) TestCreate_Successfully() {
	ctx := context.TODO()

	book := catalog.Book{
		ID:    "some-id",
		Title: "Domain Driver Design",
		Images: []catalog.Image{
			{
				ID:          "some-id2",
				Description: "poster",
				BookID:      "some-id",
			},
		},
		AuthorName:  "Eric Evans",
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.Local),
	}

	err := s.repo.Create(ctx, &book)
	assert.Nil(s.T(), err)

	persisted, err := s.repo.FindByID(ctx, book.ID)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), book.ID, persisted.ID)
}

func (s *BookRepositoryTestSuite) TestUpdate_Successfully() {
	ctx := context.TODO()

	book := catalog.Book{
		ID:         "some-id",
		Title:      "Domain Driver Design",
		AuthorName: "Eric Evans",
		Images: []catalog.Image{
			{
				ID:          "some-id2",
				Description: "poster",
				BookID:      "some-id",
			},
		},
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.Local),
	}

	err := s.repo.Create(ctx, &book)
	require.Nil(s.T(), err)

	book.Title = "new title"
	book.Images = []catalog.Image{}
	err = s.repo.Update(ctx, &book)
	require.Nil(s.T(), err)

	persisted, err := s.repo.FindByID(ctx, book.ID)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), book.ID, persisted.ID)
	assert.Equal(s.T(), book.Title, persisted.Title)
	assert.Empty(s.T(), book.Images)
}

func (s *BookRepositoryTestSuite) TestDelete_WithInvalidID() {
	ctx := context.TODO()

	id := "some-id"

	err := s.repo.Delete(ctx, id)

	assert.IsType(s.T(), &persistence.ErrEntityNotFound{}, errors.Unwrap(err))
}

func (s *BookRepositoryTestSuite) TestDelete_WithValidID() {
	ctx := context.TODO()

	book := catalog.Book{
		ID:          "some-id",
		Title:       "Domain Driver Design",
		AuthorName:  "Eric Evans",
		Price:       8000,
		ReleaseDate: time.Date(2008, time.December, 12, 0, 0, 0, 0, time.Local),
	}

	err := s.repo.Create(ctx, &book)
	require.Nil(s.T(), err)

	err = s.repo.Delete(ctx, book.ID)

	assert.Nil(s.T(), err)
}
