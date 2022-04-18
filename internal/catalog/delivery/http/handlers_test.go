//go:build integration
// +build integration

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/catalog/delivery/dto"
	"github.com/c0llinn/ebook-store/internal/catalog/helper"
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/internal/catalog/repository"
	"github.com/c0llinn/ebook-store/internal/catalog/storage"
	"github.com/c0llinn/ebook-store/internal/catalog/usecase"
	"github.com/c0llinn/ebook-store/internal/common"
	config2 "github.com/c0llinn/ebook-store/internal/config"
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
)

type CatalogHandlerTestSuite struct {
	suite.Suite
	baseURL  string
	context  *gin.Context
	recorder *httptest.ResponseRecorder
	db       *gorm.DB
	handler  CatalogHandler
}

func (s *CatalogHandlerTestSuite) SetupTest() {
	test.SetEnvironmentVariables()
	config2.InitLogger()
	config2.LoadMigrations("file:../../../../migration")

	s.db = config2.NewConnection()
	s.baseURL = fmt.Sprintf("http://localhost:%s", viper.GetString("PORT"))

	s.db = config2.NewConnection()
	bookRepository := repository.NewBookRepository(s.db)
	s3Client := storage.NewS3Client(config2.NewS3Service(), config2.NewBucket())
	filenameGenerator := helper.NewFilenameGenerator()
	catalogUseCase := usecase.NewCatalogUseCase(bookRepository, s3Client, filenameGenerator)
	idGenerator := helper.NewIDGenerator()
	s.handler = NewCatalogHandler(catalogUseCase, idGenerator)

	s.recorder = httptest.NewRecorder()
	s.context, _ = gin.CreateTestContext(s.recorder)
}

func (s *CatalogHandlerTestSuite) TearDownTest() {
	s.db.Delete(&model.Book{}, "1 = 1")
}

func TestCatalogHandlerRun(t *testing.T) {
	suite.Run(t, new(CatalogHandlerTestSuite))
}

func (s *CatalogHandlerTestSuite) TestGetBooks() {
	book1 := factory.NewBook()
	book2 := factory.NewBook()
	book3 := factory.NewBook()

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
	book := factory.NewBook()

	err := s.db.Create(book).Error
	require.Nil(s.T(), err)

	request := httptest.NewRequest("GET", s.baseURL+"/books/"+book.ID, nil)
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: book.ID}}

	s.handler.getBook(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())
}

func (s *CatalogHandlerTestSuite) TestGetBook_NotFound() {
	id := uuid.NewString()

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
	payload := dto.CreateBook{
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
	book := factory.NewBook()

	err := s.db.Create(book).Error
	require.Nil(s.T(), err)

	title := uuid.NewString() + uuid.NewString() + uuid.NewString()
	payload := dto.UpdateBook{
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
	id := uuid.NewString()

	title := faker.TitleMale()
	payload := dto.UpdateBook{
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
	book := factory.NewBook()

	err := s.db.Create(&book).Error
	require.Nil(s.T(), err)

	title := faker.TitleMale()
	payload := dto.UpdateBook{
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

	actual := model.Book{}

	err = s.db.First(&actual, "id = ?", book.ID).Error
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), expected.ID, actual.ID)
	assert.Equal(s.T(), expected.Title, actual.Title)
}

func (s *CatalogHandlerTestSuite) TestDeleteBook_Successfully() {
	book := factory.NewBook()

	err := s.db.Create(book).Error
	require.Nil(s.T(), err)

	request := httptest.NewRequest("GET", s.baseURL+"/books/"+book.ID, nil)
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: book.ID}}

	s.handler.deleteBook(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())
}

func (s *CatalogHandlerTestSuite) TestDeleteBook_WithError() {
	id := uuid.NewString()

	request := httptest.NewRequest("GET", s.baseURL+"/books/"+id, nil)
	s.context.Request = request
	s.context.Params = gin.Params{{Key: "id", Value: id}}

	s.handler.deleteBook(s.context)

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, s.context.Errors.Last().Err)
}
