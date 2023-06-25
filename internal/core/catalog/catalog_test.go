//nolint:staticcheck
package catalog_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/ebookstore/internal/core/catalog"
	mocks2 "github.com/ebookstore/internal/mocks/core/catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	findByQueryMethod          = "FindByQuery"
	findByIdMethod             = "FindByID"
	generatePreSignedUrlMethod = "GenerateGetPreSignedUrl"
	newUniqueNameMethod        = "NewUniqueName"
	newIdMethod                = "NewID"
	saveFileMethod             = "SaveFile"
	retrieveFileMethod         = "RetrieveFile"
	createBookMethod           = "Create"
	updateBookMethod           = "Update"
	deleteBookMethod           = "Delete"
	validateMethod             = "Validate"
)

type CatalogTestSuite struct {
	suite.Suite
	repo          *mocks2.Repository
	storageClient *mocks2.StorageClient
	idGenerator   *mocks2.IDGenerator
	validator     *mocks2.Validator
	catalog       *catalog.Catalog
}

func (s *CatalogTestSuite) SetupTest() {
	s.repo = new(mocks2.Repository)
	s.storageClient = new(mocks2.StorageClient)
	s.idGenerator = new(mocks2.IDGenerator)
	s.validator = new(mocks2.Validator)

	config := catalog.Config{
		Repository:    s.repo,
		StorageClient: s.storageClient,
		IDGenerator:   s.idGenerator,
		Validator:     s.validator,
	}

	s.catalog = catalog.New(config)
}

func TestCatalog(t *testing.T) {
	suite.Run(t, new(CatalogTestSuite))
}

func (s *CatalogTestSuite) TestFindByQuery_WhenRepositoryFails() {
	request := catalog.SearchBooks{}
	query := request.BookQuery()

	s.repo.On(findByQueryMethod, context.TODO(), query).Return(catalog.PaginatedBooks{}, fmt.Errorf("some error"))

	_, err := s.catalog.FindBooks(context.TODO(), request)

	assert.Error(s.T(), err)

	s.repo.AssertCalled(s.T(), findByQueryMethod, context.TODO(), query)
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod)
}

func (s *CatalogTestSuite) TestFindByQuery_WhenStorageClientFails() {
	request := catalog.SearchBooks{}
	query := request.BookQuery()

	paginatedBooks := catalog.PaginatedBooks{
		Books: []catalog.Book{
			{Images: []catalog.Image{{ID: "some-key"}}},
		},
	}

	s.repo.On(findByQueryMethod, context.TODO(), query).Return(paginatedBooks, nil)
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), "some-key").Return("", fmt.Errorf("some error"))

	_, err := s.catalog.FindBooks(context.TODO(), request)

	assert.Error(s.T(), err)

	s.repo.AssertCalled(s.T(), findByQueryMethod, context.TODO(), query)
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 1)
}

func (s *CatalogTestSuite) TestFindByQuery_Successfully() {
	request := catalog.SearchBooks{}
	query := request.BookQuery()

	paginatedBooks := catalog.PaginatedBooks{
		Books: []catalog.Book{
			{ID: "book-id-1", Images: []catalog.Image{{ID: "some-key"}, {ID: "some-key2"}}},
			{ID: "book-id-2", Images: []catalog.Image{{ID: "some-key3"}}},
		},
		Limit: 10,
	}

	s.repo.On(findByQueryMethod, context.TODO(), query).Return(paginatedBooks, nil)
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), "some-key").Return("some-link-1", nil).Once()
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), "some-key2").Return("some-link-2", nil).Once()
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), "some-key3").Return("some-link-3", nil).Once()

	expected := catalog.NewPaginatedBooksResponse(
		paginatedBooks,
		map[string][]string{"book-id-1": {"some-link-1", "some-link-2"}, "book-id-2": {"some-link-3"}},
	)

	actual, err := s.catalog.FindBooks(context.TODO(), request)

	assert.Equal(s.T(), expected, actual)
	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), findByQueryMethod, context.TODO(), query)
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 3)
}

func (s *CatalogTestSuite) TestFindBookByID_WhenRepositoryFails() {
	s.repo.On(findByIdMethod, context.TODO(), "some-id").Return(catalog.Book{}, fmt.Errorf("some error"))

	_, err := s.catalog.FindBookByID(context.TODO(), "some-id")

	assert.Error(s.T(), err)

	s.repo.AssertCalled(s.T(), findByIdMethod, context.TODO(), "some-id")
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 0)
}

func (s *CatalogTestSuite) TestFindBookByID_WhenStorageClientFails() {
	book := catalog.Book{
		ID:     "some-id",
		Images: []catalog.Image{{ID: "some-key"}},
	}

	s.repo.On(findByIdMethod, context.TODO(), book.ID).Return(book, nil)
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), "some-key").Return("", fmt.Errorf("some error"))

	_, err := s.catalog.FindBookByID(context.TODO(), book.ID)

	assert.Error(s.T(), err)

	s.repo.AssertCalled(s.T(), findByIdMethod, context.TODO(), book.ID)
	s.storageClient.AssertCalled(s.T(), generatePreSignedUrlMethod, context.TODO(), "some-key")
}

func (s *CatalogTestSuite) TestFindBookByID_Successfully() {
	book := catalog.Book{
		ID:     "some-id",
		Images: []catalog.Image{{ID: "some-key"}},
	}

	s.repo.On(findByIdMethod, context.TODO(), book.ID).Return(book, nil)
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), "some-key").Return("some-link", nil)

	actual, err := s.catalog.FindBookByID(context.TODO(), book.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), catalog.NewBookResponse(book, []string{"some-link"}), actual)

	s.repo.AssertCalled(s.T(), findByIdMethod, context.TODO(), book.ID)
	s.storageClient.AssertCalled(s.T(), generatePreSignedUrlMethod, context.TODO(), "some-key")
}

func (s *CatalogTestSuite) TestGetBookContent_WhenBookCouldNotBeFound() {
	id := "some-id"
	s.repo.On(findByIdMethod, context.TODO(), id).Return(catalog.Book{}, fmt.Errorf("some error"))

	_, err := s.catalog.GetBookContent(context.TODO(), id)

	assert.Error(s.T(), err)

	s.repo.AssertCalled(s.T(), findByIdMethod, context.TODO(), id)
	s.storageClient.AssertNumberOfCalls(s.T(), retrieveFileMethod, 0)
}

func (s *CatalogTestSuite) TestGetBookContent_WithError() {
	book := catalog.Book{
		ID:               "some-id",
		ContentBucketKey: "some-key",
	}
	s.repo.On(findByIdMethod, context.TODO(), book.ID).Return(book, nil)
	s.storageClient.On(retrieveFileMethod, context.TODO(), book.ContentBucketKey).Return(nil, fmt.Errorf("some error"))

	_, err := s.catalog.GetBookContent(context.TODO(), book.ID)

	assert.Error(s.T(), err)

	s.repo.AssertCalled(s.T(), findByIdMethod, context.TODO(), book.ID)
	s.storageClient.AssertCalled(s.T(), retrieveFileMethod, context.TODO(), book.ContentBucketKey)
}

func (s *CatalogTestSuite) TestGetBookContent_Successfully() {
	book := catalog.Book{
		ID:               "some-id",
		ContentBucketKey: "some-key",
	}
	s.repo.On(findByIdMethod, context.TODO(), book.ID).Return(book, nil)
	s.storageClient.On(retrieveFileMethod, context.TODO(), book.ContentBucketKey).Return(io.NopCloser(strings.NewReader("test")), nil)

	reader, err := s.catalog.GetBookContent(context.TODO(), book.ID)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), reader)

	s.repo.AssertCalled(s.T(), findByIdMethod, context.TODO(), book.ID)
	s.storageClient.AssertCalled(s.T(), retrieveFileMethod, context.TODO(), book.ContentBucketKey)
}

func (s *CatalogTestSuite) TestCreateBook_WithNonAdminUser() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
	}

	ctx := context.WithValue(context.Background(), "admin", false)
	_, err := s.catalog.CreateBook(ctx, request)

	assert.Equal(s.T(), catalog.ErrForbiddenCatalogAccess, errors.Unwrap(err))

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 0)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 0)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 0)
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertNumberOfCalls(s.T(), createBookMethod, 0)
}

func (s *CatalogTestSuite) TestCreateBook_ValidationFails() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
	}
	s.validator.On(validateMethod, request).Return(fmt.Errorf("some error"))

	ctx := context.WithValue(context.Background(), "admin", true)
	_, err := s.catalog.CreateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 0)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 0)
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertNumberOfCalls(s.T(), createBookMethod, 0)
}

func (s *CatalogTestSuite) TestCreateBook_WhenRepositoryFails() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
	}

	book := request.Book("some-id")
	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return("some-id")

	s.repo.On(createBookMethod, ctx, &book).Return(fmt.Errorf("some error"))

	_, err := s.catalog.CreateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), createBookMethod, 1)
}

func (s *CatalogTestSuite) TestCreateBook_WhenPreSigningFails() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
		Images: []catalog.ImageRequest{
			{
				ID: "some-key",
			},
		},
	}

	book := request.Book("some-id")
	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return("some-id")
	s.repo.On(createBookMethod, ctx, &book).Return(nil)
	s.repo.On(findByIdMethod, ctx, book.ID).Return(book, nil)
	s.storageClient.On(generatePreSignedUrlMethod, ctx, "some-key").Return("", fmt.Errorf("some error"))

	_, err := s.catalog.CreateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByIdMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), createBookMethod, 1)
}

func (s *CatalogTestSuite) TestCreateBook_Successfully() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
		Images: []catalog.ImageRequest{
			{
				ID: "some-key",
			},
		},
	}

	book := request.Book("some-id")
	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return("some-id")
	s.repo.On(findByIdMethod, ctx, book.ID).Return(book, nil)
	s.repo.On(createBookMethod, ctx, &book).Return(nil)
	s.storageClient.On(generatePreSignedUrlMethod, ctx, "some-key").Return("link", nil)

	_, err := s.catalog.CreateBook(ctx, request)

	assert.NoError(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByIdMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), createBookMethod, 1)
}

func (s *CatalogTestSuite) TestUpdateBook_WithNonAdminUser() {
	newTitle := "new title"
	request := catalog.UpdateBook{
		ID:    "some-id",
		Title: &newTitle,
	}

	ctx := context.Background()
	err := s.catalog.UpdateBook(ctx, request)

	assert.Equal(s.T(), catalog.ErrForbiddenCatalogAccess, errors.Unwrap(err))

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 0)
	s.repo.AssertNotCalled(s.T(), findByIdMethod, ctx, request.ID)
}

func (s *CatalogTestSuite) TestUpdateBook_ValidationFails() {
	newTitle := "new title"
	request := catalog.UpdateBook{
		ID:    "some-id",
		Title: &newTitle,
	}
	s.validator.On(validateMethod, request).Return(fmt.Errorf("some error"))

	ctx := context.WithValue(context.Background(), "admin", true)
	err := s.catalog.UpdateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNotCalled(s.T(), findByIdMethod, ctx, request.ID)
}

func (s *CatalogTestSuite) TestUpdateBook_WhenBookIsNotFound() {
	newTitle := "new title"
	request := catalog.UpdateBook{
		ID:    "some-id",
		Title: &newTitle,
	}

	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByIdMethod, ctx, request.ID).Return(catalog.Book{}, fmt.Errorf("some error"))

	err := s.catalog.UpdateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertCalled(s.T(), findByIdMethod, ctx, request.ID)
}

func (s *CatalogTestSuite) TestUpdateBook_WhenUpdateFails() {
	newTitle := "new title"
	request := catalog.UpdateBook{
		ID:    "some-id",
		Title: &newTitle,
	}
	book := catalog.Book{
		ID:    "some-id",
		Title: "old-title",
	}

	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByIdMethod, ctx, request.ID).Return(book, nil)

	updated := request.Update(book)
	s.repo.On(updateBookMethod, ctx, &updated).Return(fmt.Errorf("some error"))

	err := s.catalog.UpdateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertCalled(s.T(), findByIdMethod, ctx, request.ID)
	s.repo.AssertCalled(s.T(), updateBookMethod, ctx, &updated)
}

func (s *CatalogTestSuite) TestUpdateBook_Successfully() {
	newTitle := "new title"
	request := catalog.UpdateBook{
		ID:    "some-id",
		Title: &newTitle,
	}
	book := catalog.Book{
		ID:    "some-id",
		Title: "old-title",
	}

	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByIdMethod, ctx, request.ID).Return(book, nil)

	updated := request.Update(book)
	s.repo.On(updateBookMethod, ctx, &updated).Return(nil)

	err := s.catalog.UpdateBook(ctx, request)

	assert.Nil(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertCalled(s.T(), findByIdMethod, ctx, request.ID)
	s.repo.AssertCalled(s.T(), updateBookMethod, ctx, &updated)
}

func (s *CatalogTestSuite) TestDeleteBook_WithNonAdminUser() {
	id := "some-id"

	ctx := context.Background()
	err := s.catalog.DeleteBook(ctx, id)

	assert.Equal(s.T(), catalog.ErrForbiddenCatalogAccess, errors.Unwrap(err))

	s.repo.AssertNotCalled(s.T(), deleteBookMethod, ctx, id)
}

func (s *CatalogTestSuite) TestDeleteBook_WhenRepositoryFails() {
	id := "some-id"
	ctx := context.WithValue(context.Background(), "admin", true)
	s.repo.On(deleteBookMethod, ctx, id).Return(fmt.Errorf("some error"))

	err := s.catalog.DeleteBook(ctx, id)

	assert.Error(s.T(), err)

	s.repo.AssertCalled(s.T(), deleteBookMethod, ctx, id)
}

func (s *CatalogTestSuite) TestDeleteBook_Successfully() {
	id := "some-id"
	ctx := context.WithValue(context.Background(), "admin", true)
	s.repo.On(deleteBookMethod, ctx, id).Return(nil)

	err := s.catalog.DeleteBook(ctx, id)
	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), deleteBookMethod, ctx, id)
}
