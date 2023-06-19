//nolint:staticcheck
package catalog_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/ebookstore/internal/core/catalog"
	mocks2 "github.com/ebookstore/internal/mocks/catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	findByQueryMethod          = "FindByQuery"
	findByIdMethod             = "FindByID"
	generatePreSignedUrlMethod = "GeneratePreSignedUrl"
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
	repo              *mocks2.Repository
	storageClient     *mocks2.StorageClient
	filenameGenerator *mocks2.FilenameGenerator
	idGenerator       *mocks2.IDGenerator
	validator         *mocks2.Validator
	catalog           *catalog.Catalog
}

func (s *CatalogTestSuite) SetupTest() {
	s.repo = new(mocks2.Repository)
	s.storageClient = new(mocks2.StorageClient)
	s.filenameGenerator = new(mocks2.FilenameGenerator)
	s.idGenerator = new(mocks2.IDGenerator)
	s.validator = new(mocks2.Validator)

	config := catalog.Config{
		Repository:        s.repo,
		StorageClient:     s.storageClient,
		FilenameGenerator: s.filenameGenerator,
		IDGenerator:       s.idGenerator,
		Validator:         s.validator,
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
			{PosterImageBucketKey: "some-key"},
		},
	}

	s.repo.On(findByQueryMethod, context.TODO(), query).Return(paginatedBooks, nil)
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), paginatedBooks.Books[0].PosterImageBucketKey).Return("", fmt.Errorf("some error"))

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
			{PosterImageBucketKey: "some-key"},
			{PosterImageBucketKey: "some-key2"},
		},
		Limit: 10,
	}

	s.repo.On(findByQueryMethod, context.TODO(), query).Return(paginatedBooks, nil)
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), paginatedBooks.Books[0].PosterImageBucketKey).Return("some-link-1", nil).Once()
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), paginatedBooks.Books[1].PosterImageBucketKey).Return("some-link-2", nil).Once()

	expected := catalog.NewPaginatedBooksResponse(paginatedBooks)
	expected.Results[0].PosterImageLink = "some-link-1"
	expected.Results[1].PosterImageLink = "some-link-2"

	actual, err := s.catalog.FindBooks(context.TODO(), request)

	assert.Equal(s.T(), expected, actual)
	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), findByQueryMethod, context.TODO(), query)
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 2)
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
		ID:                   "some-id",
		PosterImageBucketKey: "some-key",
	}

	s.repo.On(findByIdMethod, context.TODO(), book.ID).Return(book, nil)
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), book.PosterImageBucketKey).Return("", fmt.Errorf("some error"))

	_, err := s.catalog.FindBookByID(context.TODO(), book.ID)

	assert.Error(s.T(), err)

	s.repo.AssertCalled(s.T(), findByIdMethod, context.TODO(), book.ID)
	s.storageClient.AssertCalled(s.T(), generatePreSignedUrlMethod, context.TODO(), book.PosterImageBucketKey)
}

func (s *CatalogTestSuite) TestFindBookByID_Successfully() {
	book := catalog.Book{
		ID:                   "some-id",
		PosterImageBucketKey: "some-key",
	}

	s.repo.On(findByIdMethod, context.TODO(), book.ID).Return(book, nil)
	s.storageClient.On(generatePreSignedUrlMethod, context.TODO(), book.PosterImageBucketKey).Return("some-link", nil)

	expected := book
	expected.PosterImageLink = "some-link"
	actual, err := s.catalog.FindBookByID(context.TODO(), book.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), catalog.NewBookResponse(expected), actual)

	s.repo.AssertCalled(s.T(), findByIdMethod, context.TODO(), book.ID)
	s.storageClient.AssertCalled(s.T(), generatePreSignedUrlMethod, context.TODO(), book.PosterImageBucketKey)
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
		PosterImage: bytes.NewReader([]byte("some poster image")),
		BookContent: bytes.NewReader([]byte("some book content content")),
	}

	ctx := context.WithValue(context.Background(), "admin", false)
	_, err := s.catalog.CreateBook(ctx, request)

	assert.Equal(s.T(), catalog.ErrForbiddenCatalogAccess, errors.Unwrap(err))

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 0)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 0)
	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 0)
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
		PosterImage: bytes.NewReader([]byte("some poster image")),
		BookContent: bytes.NewReader([]byte("some book content content")),
	}
	s.validator.On(validateMethod, request).Return(fmt.Errorf("some error"))

	ctx := context.WithValue(context.Background(), "admin", true)
	_, err := s.catalog.CreateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 0)
	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 0)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 0)
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertNumberOfCalls(s.T(), createBookMethod, 0)
}

func (s *CatalogTestSuite) TestCreateBook_WhenPosterStorageFails() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
		PosterImage: bytes.NewReader([]byte("some poster image")),
		BookContent: bytes.NewReader([]byte("some book content content")),
	}

	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return("some-id")
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+request.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+request.Title).Return("content_name").Once()
	s.storageClient.On(saveFileMethod, ctx, "poster_name", "image/jpeg", request.PosterImage).Return(fmt.Errorf("some error"))

	book := request.Book("some-id")

	_, err := s.catalog.CreateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 1)
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertNotCalled(s.T(), createBookMethod, &book)
}

func (s *CatalogTestSuite) TestCreateBook_WhenContentStorageFails() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
		PosterImage: bytes.NewReader([]byte("some poster image")),
		BookContent: bytes.NewReader([]byte("some book content content")),
	}

	book := request.Book("some-id")

	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return("some-id")
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+request.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+request.Title).Return("content_name").Once()
	s.storageClient.On(saveFileMethod, ctx, "poster_name", "image/jpeg", request.PosterImage).Return(nil)
	s.storageClient.On(saveFileMethod, ctx, "content_name", "application/pdf", request.BookContent).Return(fmt.Errorf("some error"))

	_, err := s.catalog.CreateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 2)
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertNotCalled(s.T(), createBookMethod, &book)
}

func (s *CatalogTestSuite) TestCreateBook_WhenPreSigningFails() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
		PosterImage: bytes.NewReader([]byte("some poster image")),
		BookContent: bytes.NewReader([]byte("some book content content")),
	}

	book := request.Book("some-id")
	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return("some-id")
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+book.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+book.Title).Return("content_name").Once()
	s.storageClient.On(saveFileMethod, ctx, "poster_name", "image/jpeg", request.PosterImage).Return(nil)
	s.storageClient.On(saveFileMethod, ctx, "content_name", "application/pdf", request.BookContent).Return(nil)
	s.storageClient.On(generatePreSignedUrlMethod, ctx, "poster_name").Return("", fmt.Errorf("some error"))

	_, err := s.catalog.CreateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), createBookMethod, 0)
}

func (s *CatalogTestSuite) TestCreateBook_WhenRepositoryFails() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
		PosterImage: bytes.NewReader([]byte("some poster image")),
		BookContent: bytes.NewReader([]byte("some book content content")),
	}

	book := request.Book("some-id")
	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return("some-id")
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+book.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+book.Title).Return("content_name").Once()
	s.storageClient.On(saveFileMethod, ctx, "poster_name", "image/jpeg", request.PosterImage).Return(nil)
	s.storageClient.On(saveFileMethod, ctx, "content_name", "application/pdf", request.BookContent).Return(nil)
	s.storageClient.On(generatePreSignedUrlMethod, ctx, "poster_name").Return("some-link", nil)

	newBook := book
	newBook.PosterImageBucketKey = "poster_name"
	newBook.ContentBucketKey = "content_name"
	newBook.PosterImageLink = "some-link"

	s.repo.On(createBookMethod, ctx, &newBook).Return(fmt.Errorf("some error"))

	_, err := s.catalog.CreateBook(ctx, request)

	assert.Error(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), createBookMethod, 1)
}

func (s *CatalogTestSuite) TestCreateBook_Successfully() {
	request := catalog.CreateBook{
		Title:       "Clean Code",
		Description: "A Craftsman Guide",
		AuthorName:  "Robert C. Martin",
		Price:       4000,
		ReleaseDate: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
		PosterImage: bytes.NewReader([]byte("some poster image")),
		BookContent: bytes.NewReader([]byte("some book content content")),
	}

	book := request.Book("some-id")
	ctx := context.WithValue(context.Background(), "admin", true)

	s.validator.On(validateMethod, request).Return(nil)
	s.idGenerator.On(newIdMethod).Return("some-id")
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+book.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+book.Title).Return("content_name").Once()
	s.storageClient.On(saveFileMethod, ctx, "poster_name", "image/jpeg", request.PosterImage).Return(nil)
	s.storageClient.On(saveFileMethod, ctx, "content_name", "application/pdf", request.BookContent).Return(nil)
	s.storageClient.On(generatePreSignedUrlMethod, ctx, "poster_name").Return("some-link", nil)

	newBook := book
	newBook.PosterImageBucketKey = "poster_name"
	newBook.ContentBucketKey = "content_name"
	newBook.PosterImageLink = "some-link"

	s.repo.On(createBookMethod, ctx, &newBook).Return(nil)

	expected := catalog.NewBookResponse(newBook)
	actual, err := s.catalog.CreateBook(ctx, request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 1)
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
