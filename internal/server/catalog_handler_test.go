package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/config"
	"github.com/c0llinn/ebook-store/internal/generator"
	"github.com/c0llinn/ebook-store/internal/persistence"
	"github.com/c0llinn/ebook-store/internal/storage"
	"github.com/c0llinn/ebook-store/test"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http/httptest"
	"testing"
	"time"
)

type CatalogHandlerTestSuite struct {
	suite.Suite
	baseURL  string
	context  *gin.Context
	recorder *httptest.ResponseRecorder
	db       *gorm.DB
	handler  *CatalogHandler
}

func (s *CatalogHandlerTestSuite) SetupTest() {
	test.SetEnvironmentVariables()
	config.LoadMigrations("file:../../migrations")

	s.db = config.NewConnection()
	s.baseURL = fmt.Sprintf("http://localhost:%s", viper.GetString("PORT"))

	s.db = config.NewConnection()
	bookRepository := persistence.NewBookRepository(s.db)
	s3Client := storage.NewS3Client(config.NewS3Service(), config.NewBucket())
	filenameGenerator := generator.NewFilenameGenerator()
	idGenerator := generator.NewUUIDGenerator()

	c := catalog.Config{
		Repository:        bookRepository,
		StorageClient:     s3Client,
		FilenameGenerator: filenameGenerator,
		IDGenerator:       idGenerator,
	}

	catalog := catalog.New(c)

	s.handler = NewCatalogHandler(gin.New(), catalog)

	s.recorder = httptest.NewRecorder()
	s.context, _ = gin.CreateTestContext(s.recorder)
}

func (s *CatalogHandlerTestSuite) TearDownTest() {
	s.db.Delete(&catalog.Book{}, "1 = 1")
}

func TestCatalogHandler(t *testing.T) {
	suite.Run(t, new(CatalogHandlerTestSuite))
}

func (s *CatalogHandlerTestSuite) TestGetBooks() {
	book1 := catalog.Book{
		ID:                   "some-id1",
		Title:                "Clean Code",
		Description:          "Some Description",
		AuthorName:           "Robert C. Martin",
		PosterImageBucketKey: "key1",
		PosterImageLink:      "http://localhost",
		ContentBucketKey:     "key2",
		Price:                4000,
		ReleaseDate:          time.Date(2021, time.September, 25, 0, 0, 0, 0, time.UTC),
		CreatedAt:            time.Date(2022, time.September, 26, 0, 0, 0, 0, time.UTC),
		UpdatedAt:            time.Date(2022, time.September, 27, 0, 0, 0, 0, time.UTC),
	}

	book2 := catalog.Book{
		ID:                   "some-id2",
		Title:                "Clean Coder",
		Description:          "Some Description2",
		AuthorName:           "Robert C. Martin",
		PosterImageBucketKey: "key4",
		PosterImageLink:      "http://localhost/test",
		ContentBucketKey:     "key6",
		Price:                4500,
		ReleaseDate:          time.Date(2021, time.September, 28, 0, 0, 0, 0, time.UTC),
		CreatedAt:            time.Date(2022, time.September, 29, 0, 0, 0, 0, time.UTC),
		UpdatedAt:            time.Date(2022, time.September, 30, 0, 0, 0, 0, time.UTC),
	}

	book3 := catalog.Book{
		ID:                   "some-id4",
		Title:                "Clean Architecture",
		Description:          "Some Description3",
		AuthorName:           "Robert C. Martin",
		PosterImageBucketKey: "key12",
		PosterImageLink:      "http://localhost/test",
		ContentBucketKey:     "key14",
		Price:                6000,
		ReleaseDate:          time.Date(2021, time.October, 28, 0, 0, 0, 0, time.UTC),
		CreatedAt:            time.Date(2022, time.October, 29, 0, 0, 0, 0, time.UTC),
		UpdatedAt:            time.Date(2022, time.October, 30, 0, 0, 0, 0, time.UTC),
	}

	err := s.db.Create(book1).Error
	require.Nil(s.T(), err)

	err = s.db.Create(book2).Error
	require.Nil(s.T(), err)

	err = s.db.Create(book3).Error
	require.Nil(s.T(), err)

	request := httptest.NewRequest("GET", s.baseURL+"/books", nil)
	s.context.Request = request

	s.handler.getBooks(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())
}

func (s *CatalogHandlerTestSuite) TestGetBook_Successfully() {
	book := catalog.Book{
		ID:                   "some-id1",
		Title:                "Clean Code",
		Description:          "Some Description",
		AuthorName:           "Robert C. Martin",
		PosterImageBucketKey: "key1",
		PosterImageLink:      "http://localhost",
		ContentBucketKey:     "key2",
		Price:                4000,
		ReleaseDate:          time.Date(2021, time.September, 25, 0, 0, 0, 0, time.UTC),
		CreatedAt:            time.Date(2022, time.September, 26, 0, 0, 0, 0, time.UTC),
		UpdatedAt:            time.Date(2022, time.September, 27, 0, 0, 0, 0, time.UTC),
	}

	err := s.db.Create(book).Error
	require.Nil(s.T(), err)

	request := httptest.NewRequest("GET", s.baseURL+"/books/"+book.ID, nil)
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: book.ID}}

	s.handler.getBook(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())
}

func (s *CatalogHandlerTestSuite) TestGetBook_NotFound() {
	id := "some-id"

	request := httptest.NewRequest("GET", s.baseURL+"/books/"+id, nil)
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: id}}

	s.handler.getBook(s.context)

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, s.context.Errors.Last().Err)
}

func (s *CatalogHandlerTestSuite) TestCreateBook_WithMalformedPayload() {
	payload := map[string]interface{}{
		"title": 34,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/books", bytes.NewReader(data))
	s.context.Request = request

	s.handler.createBook(s.context)

	assert.IsType(s.T(), &common.ErrNotValid{}, s.context.Errors.Last().Err)
}

func (s *CatalogHandlerTestSuite) TestCreateBook_WithInvalidPayload() {
	payload := catalog.CreateBook{
		Title: "title",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/books", bytes.NewReader(data))
	s.context.Request = request

	s.handler.createBook(s.context)

	assert.IsType(s.T(), &common.ErrNotValid{}, s.context.Errors.Last().Err)
}

func (s *CatalogHandlerTestSuite) TestUpdateBook_WithMalformedPayload() {
	book := factory.NewBook()

	err := s.db.Create(book).Error
	require.Nil(s.T(), err)

	payload := map[string]interface{}{
		"title": 34,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("PATCH", s.baseURL+"/books/"+book.ID, bytes.NewReader(data))
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: book.ID}}

	s.handler.updateBook(s.context)

	assert.IsType(s.T(), &common.ErrNotValid{}, s.context.Errors.Last().Err)
}

func (s *CatalogHandlerTestSuite) TestUpdateBook_WithInvalidPayload() {
	book := catalog.Book{
		ID:                   "some-id1",
		Title:                "Clean Code",
		Description:          "Some Description",
		AuthorName:           "Robert C. Martin",
		PosterImageBucketKey: "key1",
		PosterImageLink:      "http://localhost",
		ContentBucketKey:     "key2",
		Price:                4000,
		ReleaseDate:          time.Date(2021, time.September, 25, 0, 0, 0, 0, time.UTC),
		CreatedAt:            time.Date(2022, time.September, 26, 0, 0, 0, 0, time.UTC),
		UpdatedAt:            time.Date(2022, time.September, 27, 0, 0, 0, 0, time.UTC),
	}

	err := s.db.Create(book).Error
	require.Nil(s.T(), err)

	title := uuid.NewString() + uuid.NewString() + uuid.NewString()
	payload := catalog.UpdateBook{
		Title: &title,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("PATCH", s.baseURL+"/books/"+book.ID, bytes.NewReader(data))
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: book.ID}}

	s.handler.updateBook(s.context)

	assert.IsType(s.T(), &common.ErrNotValid{}, s.context.Errors.Last().Err)
}

func (s *CatalogHandlerTestSuite) TestUpdateBook_WithUnknownID() {
	id := "some-id"

	title := faker.TitleMale()
	payload := catalog.UpdateBook{
		Title: &title,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("PATCH", s.baseURL+"/books/"+id, bytes.NewReader(data))
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: id}}

	s.handler.updateBook(s.context)

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, s.context.Errors.Last().Err)
}

func (s *CatalogHandlerTestSuite) TestUpdateBook_Successfully() {
	book := catalog.Book{
		ID:                   "some-id1",
		Title:                "Clean Code",
		Description:          "Some Description",
		AuthorName:           "Robert C. Martin",
		PosterImageBucketKey: "key1",
		PosterImageLink:      "http://localhost",
		ContentBucketKey:     "key2",
		Price:                4000,
		ReleaseDate:          time.Date(2021, time.September, 25, 0, 0, 0, 0, time.UTC),
		CreatedAt:            time.Date(2022, time.September, 26, 0, 0, 0, 0, time.UTC),
		UpdatedAt:            time.Date(2022, time.September, 27, 0, 0, 0, 0, time.UTC),
	}

	err := s.db.Create(&book).Error
	require.Nil(s.T(), err)

	title := "new title"
	payload := catalog.UpdateBook{
		Title: &title,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("PATCH", s.baseURL+"/books/"+book.ID, bytes.NewReader(data))
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: book.ID}}

	s.handler.updateBook(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())

	expected := book
	expected.Title = title

	actual := catalog.Book{}

	err = s.db.First(&actual, "id = ?", book.ID).Error
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), expected.ID, actual.ID)
	assert.Equal(s.T(), expected.Title, actual.Title)
}

func (s *CatalogHandlerTestSuite) TestDeleteBook_Successfully() {
	book := catalog.Book{
		ID:                   "some-id1",
		Title:                "Clean Code",
		Description:          "Some Description",
		AuthorName:           "Robert C. Martin",
		PosterImageBucketKey: "key1",
		PosterImageLink:      "http://localhost",
		ContentBucketKey:     "key2",
		Price:                4000,
		ReleaseDate:          time.Date(2021, time.September, 25, 0, 0, 0, 0, time.UTC),
		CreatedAt:            time.Date(2022, time.September, 26, 0, 0, 0, 0, time.UTC),
		UpdatedAt:            time.Date(2022, time.September, 27, 0, 0, 0, 0, time.UTC),
	}

	err := s.db.Create(book).Error
	require.Nil(s.T(), err)

	request := httptest.NewRequest("GET", s.baseURL+"/books/"+book.ID, nil)
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: book.ID}}

	s.handler.deleteBook(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())
}

func (s *CatalogHandlerTestSuite) TestDeleteBook_WithError() {
	id := "some-id"

	request := httptest.NewRequest("GET", s.baseURL+"/books/"+id, nil)
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: id}}

	s.handler.deleteBook(s.context)

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, s.context.Errors.Last().Err)
}
