package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/internal/auth/delivery/dto"
	"github.com/c0llinn/ebook-store/internal/auth/helper"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/internal/auth/repository"
	"github.com/c0llinn/ebook-store/internal/auth/token"
	"github.com/c0llinn/ebook-store/internal/auth/usecase"
	"github.com/c0llinn/ebook-store/test"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	baseURL  string
	context  *gin.Context
	recorder *httptest.ResponseRecorder
	db       *gorm.DB
	handler  AuthHandler
}

func (s *AuthHandlerTestSuite) SetupTest() {
	test.SetEnvironmentVariables()
	log.InitLogger()
	db.LoadMigrations("file:../../../../migration")

	s.db = db.NewConnection()
	s.baseURL = fmt.Sprintf("http://localhost:%s", viper.GetString("PORT"))

	userRepository := repository.NewUserRepository(s.db)
	jwtWrapper := token.NewJWTWrapper(token.NewHMACSecret())
	authUseCase := usecase.NewAuthUseCase(userRepository, jwtWrapper)

	s.handler = NewAuthHandler(authUseCase, helper.NewUUIDGenerator())

	s.recorder = httptest.NewRecorder()
	s.context, _ = gin.CreateTestContext(s.recorder)
}

func (s *AuthHandlerTestSuite) TearDownTest() {
	s.db.Delete(&model.User{}, "1 = 1")
}

func TestAuthHandlerRun(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

func (s *AuthHandlerTestSuite) TestRegisterSuccessfully() {
	password := faker.Password()
	payload := dto.RegisterRequest{
		FirstName:            faker.FirstName(),
		LastName:             faker.LastName(),
		Email:                faker.Email(),
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
	s.context.Request = request

	s.handler.register(s.context)

	var credentials dto.CredentialsResponse
	err = json.NewDecoder(s.recorder.Result().Body).Decode(&credentials)
	require.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusCreated, s.recorder.Code)
	assert.NotEmpty(s.T(), credentials.Token)

	var user model.User
	err = s.db.First(&user).Error
	require.Nil(s.T(), err)

	assert.Len(s.T(), user.ID, 36)
	assert.Equal(s.T(), payload.FirstName, user.FirstName)
	assert.Equal(s.T(), payload.LastName, user.LastName)
	assert.Equal(s.T(), payload.Email, user.Email)
	assert.NotEmpty(s.T(), user.Password)
	assert.Equal(s.T(), model.Customer, user.Role)
	assert.NotZero(s.T(), user.CreatedAt)
}
