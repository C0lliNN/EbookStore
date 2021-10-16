// +build unit

package usecase

import (
	"bytes"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/catalog/mock"
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"strings"
	"testing"
)

const (
	findByQueryMethod          = "FindByQuery"
	findByIdMethod             = "FindByID"
	generatePreSignedUrlMethod = "GeneratePreSignedUrl"
	newUniqueNameMethod        = "NewUniqueName"
	saveFileMethod             = "SaveFile"
	retrieveFileMethod         = "RetrieveFile"
	createBookMethod           = "Create"
	updateBookMethod           = "Update"
	deleteBookMethod           = "Delete"
)

type CatalogUseCaseTestSuite struct {
	suite.Suite
	repo              *mock.BookRepository
	storageClient     *mock.StorageClient
	filenameGenerator *mock.FilenameGenerator
	useCase           CatalogUseCase
}

func (s *CatalogUseCaseTestSuite) SetupTest() {
	s.repo = new(mock.BookRepository)
	s.storageClient = new(mock.StorageClient)
	s.filenameGenerator = new(mock.FilenameGenerator)
	s.useCase = NewCatalogUseCase(s.repo, s.storageClient, s.filenameGenerator)
}

func TestCatalogUseCaseRun(t *testing.T) {
	suite.Run(t, new(CatalogUseCaseTestSuite))
}

func (s *CatalogUseCaseTestSuite) TestFindByQuery_WhenRepositoryFails() {
	s.repo.On(findByQueryMethod, model.BookQuery{}).Return(model.PaginatedBooks{}, fmt.Errorf("some error"))

	_, err := s.useCase.FindBooks(model.BookQuery{})

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findByQueryMethod, model.BookQuery{})
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod)
}

func (s *CatalogUseCaseTestSuite) TestFindByQuery_WhenStorageClientFails() {
	paginatedBooks := factory.NewPaginatedBooks()
	s.repo.On(findByQueryMethod, model.BookQuery{}).Return(paginatedBooks, nil)
	s.storageClient.On(generatePreSignedUrlMethod, paginatedBooks.Books[0].PosterImageBucketKey).Return("", fmt.Errorf("some error"))

	_, err := s.useCase.FindBooks(model.BookQuery{})

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findByQueryMethod, model.BookQuery{})
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 1)
}

func (s *CatalogUseCaseTestSuite) TestFindByQuery_Successfully() {
	paginatedBooks := factory.NewPaginatedBooks()
	s.repo.On(findByQueryMethod, model.BookQuery{}).Return(paginatedBooks, nil)
	s.storageClient.On(generatePreSignedUrlMethod, paginatedBooks.Books[0].PosterImageBucketKey).Return("some-link-1", nil).Once()
	s.storageClient.On(generatePreSignedUrlMethod, paginatedBooks.Books[1].PosterImageBucketKey).Return("some-link-2", nil).Once()

	expected := paginatedBooks
	expected.Books[0].PosterImageLink = "some-link-1"
	expected.Books[1].PosterImageLink = "some-link-2"

	actual, err := s.useCase.FindBooks(model.BookQuery{})

	assert.Equal(s.T(), expected, actual)
	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), findByQueryMethod, model.BookQuery{})
	s.storageClient.AssertNumberOfCalls(s.T(), generatePreSignedUrlMethod, 2)
}

func (s *CatalogUseCaseTestSuite) TestFindBookByID_WhenRepositoryFails() {
	s.repo.On(findByIdMethod, "some-id").Return(model.Book{}, fmt.Errorf("some error"))

	_, err := s.useCase.FindBookByID("some-id")

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findByIdMethod, "some-id")
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod)
}

func (s *CatalogUseCaseTestSuite) TestFindBookByID_WhenStorageClientFails() {
	book := factory.NewBook()
	s.repo.On(findByIdMethod, book.ID).Return(book, nil)
	s.storageClient.On(generatePreSignedUrlMethod, book.PosterImageBucketKey).Return("", fmt.Errorf("some error"))

	_, err := s.useCase.FindBookByID(book.ID)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findByIdMethod, book.ID)
	s.storageClient.AssertCalled(s.T(), generatePreSignedUrlMethod, book.PosterImageBucketKey)
}

func (s *CatalogUseCaseTestSuite) TestFindBookByID_Successfully() {
	book := factory.NewBook()
	s.repo.On(findByIdMethod, book.ID).Return(book, nil)
	s.storageClient.On(generatePreSignedUrlMethod, book.PosterImageBucketKey).Return("some-link", nil)

	expected := book
	expected.PosterImageLink = "some-link"
	actual, err := s.useCase.FindBookByID(book.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)

	s.repo.AssertCalled(s.T(), findByIdMethod, book.ID)
	s.storageClient.AssertCalled(s.T(), generatePreSignedUrlMethod, book.PosterImageBucketKey)
}

func (s *CatalogUseCaseTestSuite) TestGetBookContent_WhenBookCouldNotBeFound() {
	id := uuid.NewString()
	s.repo.On(findByIdMethod, id).Return(model.Book{}, fmt.Errorf("some error"))

	_, err := s.useCase.GetBookContent(id)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findByIdMethod, id)
	s.storageClient.AssertNumberOfCalls(s.T(), retrieveFileMethod, 0)
}

func (s *CatalogUseCaseTestSuite) TestGetBookContent_WithError() {
	book := factory.NewBook()
	s.repo.On(findByIdMethod, book.ID).Return(book, nil)
	s.storageClient.On(retrieveFileMethod, book.ContentBucketKey).Return(nil, fmt.Errorf("some error"))

	_, err := s.useCase.GetBookContent(book.ID)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findByIdMethod, book.ID)
	s.storageClient.AssertCalled(s.T(), retrieveFileMethod, book.ContentBucketKey)
}

func (s *CatalogUseCaseTestSuite) TestGetBookContent_Successfully() {
	book := factory.NewBook()
	s.repo.On(findByIdMethod, book.ID).Return(book, nil)
	s.storageClient.On(retrieveFileMethod, book.ContentBucketKey).Return(io.NopCloser(strings.NewReader("test")), nil)

	reader, err := s.useCase.GetBookContent(book.ID)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), reader)

	s.repo.AssertCalled(s.T(), findByIdMethod, book.ID)
	s.storageClient.AssertCalled(s.T(), retrieveFileMethod, book.ContentBucketKey)
}

func (s *CatalogUseCaseTestSuite) TestCreateBook_WhenPosterStorageFails() {
	book := factory.NewBook()
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+book.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+book.Title).Return("content_name").Once()

	posterImage := bytes.NewReader([]byte("some content"))

	s.storageClient.On(saveFileMethod, "poster_name", "image/jpeg", posterImage).Return(fmt.Errorf("some error"))

	err := s.useCase.CreateBook(&book, posterImage, nil)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 1)
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertNotCalled(s.T(), createBookMethod, &book)
}

func (s *CatalogUseCaseTestSuite) TestCreateBook_WhenContentStorageFails() {
	book := factory.NewBook()
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+book.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+book.Title).Return("content_name").Once()

	posterImage := bytes.NewReader([]byte("some content"))
	bookContent := bytes.NewReader([]byte("some book content"))

	s.storageClient.On(saveFileMethod, "poster_name", "image/jpeg", posterImage).Return(nil)
	s.storageClient.On(saveFileMethod, "content_name", "application/pdf", bookContent).Return(fmt.Errorf("some error"))

	err := s.useCase.CreateBook(&book, posterImage, bookContent)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 2)
	s.storageClient.AssertNotCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertNotCalled(s.T(), createBookMethod, &book)
}

func (s *CatalogUseCaseTestSuite) TestCreateBook_WhenPreSigningFails() {
	book := factory.NewBook()
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+book.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+book.Title).Return("content_name").Once()

	posterImage := bytes.NewReader([]byte("some content"))
	bookContent := bytes.NewReader([]byte("some book content"))

	s.storageClient.On(saveFileMethod, "poster_name", "image/jpeg", posterImage).Return(nil)
	s.storageClient.On(saveFileMethod, "content_name", "application/pdf", bookContent).Return(nil)
	s.storageClient.On(generatePreSignedUrlMethod, "poster_name").Return("", fmt.Errorf("some error"))

	err := s.useCase.CreateBook(&book, posterImage, bookContent)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 2)
	s.storageClient.AssertCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertNotCalled(s.T(), createBookMethod, &book)
}

func (s *CatalogUseCaseTestSuite) TestCreateBook_WhenRepositoryFails() {
	book := factory.NewBook()
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+book.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+book.Title).Return("content_name").Once()

	posterImage := bytes.NewReader([]byte("some content"))
	bookContent := bytes.NewReader([]byte("some book content"))

	s.storageClient.On(saveFileMethod, "poster_name", "image/jpeg", posterImage).Return(nil)
	s.storageClient.On(saveFileMethod, "content_name", "application/pdf", bookContent).Return(nil)
	s.storageClient.On(generatePreSignedUrlMethod, "poster_name").Return("some-link", nil)

	newBook := book
	newBook.PosterImageBucketKey = "poster_name"
	newBook.ContentBucketKey = "content_name"
	newBook.PosterImageLink = "some-link"

	s.repo.On(createBookMethod, &newBook).Return(fmt.Errorf("some error"))

	err := s.useCase.CreateBook(&book, posterImage, bookContent)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 2)
	s.storageClient.AssertCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertCalled(s.T(), createBookMethod, &newBook)
}

func (s *CatalogUseCaseTestSuite) TestCreateBook_Successfully() {
	book := factory.NewBook()
	s.filenameGenerator.On(newUniqueNameMethod, "poster_"+book.Title).Return("poster_name").Once()
	s.filenameGenerator.On(newUniqueNameMethod, "content_"+book.Title).Return("content_name").Once()

	posterImage := bytes.NewReader([]byte("some content"))
	bookContent := bytes.NewReader([]byte("some book content"))

	s.storageClient.On(saveFileMethod, "poster_name", "image/jpeg", posterImage).Return(nil)
	s.storageClient.On(saveFileMethod, "content_name", "application/pdf", bookContent).Return(nil)
	s.storageClient.On(generatePreSignedUrlMethod, "poster_name").Return("some-link", nil)

	newBook := book
	newBook.PosterImageBucketKey = "poster_name"
	newBook.ContentBucketKey = "content_name"
	newBook.PosterImageLink = "some-link"

	s.repo.On(createBookMethod, &newBook).Return(nil)

	err := s.useCase.CreateBook(&book, posterImage, bookContent)

	assert.Nil(s.T(), err)

	s.filenameGenerator.AssertNumberOfCalls(s.T(), newUniqueNameMethod, 2)
	s.storageClient.AssertNumberOfCalls(s.T(), saveFileMethod, 2)
	s.storageClient.AssertCalled(s.T(), generatePreSignedUrlMethod, "poster_name")
	s.repo.AssertCalled(s.T(), createBookMethod, &newBook)
}

func (s *CatalogUseCaseTestSuite) TestUpdateBook_WhenRepositoryFails() {
	book := factory.NewBook()
	s.repo.On(updateBookMethod, &book).Return(fmt.Errorf("some error"))

	err := s.useCase.UpdateBook(&book)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), updateBookMethod, &book)
}

func (s *CatalogUseCaseTestSuite) TestUpdateBook_Successfully() {
	book := factory.NewBook()
	s.repo.On(updateBookMethod, &book).Return(nil)

	err := s.useCase.UpdateBook(&book)

	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), updateBookMethod, &book)
}

func (s *CatalogUseCaseTestSuite) TestDeleteBook_WhenRepositoryFails() {
	id := uuid.NewString()
	s.repo.On(deleteBookMethod, id).Return(fmt.Errorf("some error"))

	err := s.useCase.DeleteBook(id)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), deleteBookMethod, id)
}

func (s *CatalogUseCaseTestSuite) TestDeleteBook_Successfully() {
	id := uuid.NewString()
	s.repo.On(deleteBookMethod, id).Return(nil)

	err := s.useCase.DeleteBook(id)

	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), deleteBookMethod, id)
}
