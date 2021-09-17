// +build integration

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/config/aws"
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/internal/auth/delivery/dto"
	"github.com/c0llinn/ebook-store/internal/auth/email"
	"github.com/c0llinn/ebook-store/internal/auth/helper"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/internal/auth/repository"
	"github.com/c0llinn/ebook-store/internal/auth/token"
	"github.com/c0llinn/ebook-store/internal/auth/usecase"
	"github.com/c0llinn/ebook-store/test"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
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
	ses := aws.NewSNSService()
	client := email.NewEmailClient(ses)
	passwordGenerator := helper.NewPasswordGenerator()
	authUseCase := usecase.NewAuthUseCase(userRepository, jwtWrapper, client, passwordGenerator)

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

func (s *AuthHandlerTestSuite) TestRegister_WithInvalidFirstName() {
	password := faker.PASSWORD
	payload := dto.RegisterRequest{
		FirstName:            "",
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

	assert.NotEmpty(s.T(), s.context.Errors.Errors())

	var user model.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthHandlerTestSuite) TestRegister_WithInvalidLastName() {
	password := faker.PASSWORD
	payload := dto.RegisterRequest{
		FirstName:            faker.FirstName(),
		LastName:             "",
		Email:                faker.Email(),
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
	s.context.Request = request

	s.handler.register(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())

	var user model.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthHandlerTestSuite) TestRegister_WithInvalidEmail() {
	password := faker.PASSWORD
	payload := dto.RegisterRequest{
		FirstName:            faker.FirstName(),
		LastName:             faker.LastName(),
		Email:                "invalid-email",
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
	s.context.Request = request

	s.handler.register(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())

	var user model.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthHandlerTestSuite) TestRegister_WithInvalidPassword() {
	payload := dto.RegisterRequest{
		FirstName:            faker.FirstName(),
		LastName:             faker.LastName(),
		Email:                faker.Email(),
		Password:             "1234",
		PasswordConfirmation: "1234",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
	s.context.Request = request

	s.handler.register(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())

	var user model.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthHandlerTestSuite) TestRegister_WithInvalidPasswordConfirmation() {
	password := faker.PASSWORD
	payload := dto.RegisterRequest{
		FirstName:            faker.FirstName(),
		LastName:             faker.LastName(),
		Email:                faker.Email(),
		Password:             password,
		PasswordConfirmation: "12341234",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
	s.context.Request = request

	s.handler.register(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())

	var user model.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthHandlerTestSuite) TestRegister_WithDuplicateEmail() {
	persistedUser := factory.NewUser()
	err := s.db.Create(persistedUser).Error
	require.Nil(s.T(), err)

	password := faker.PASSWORD
	payload := dto.RegisterRequest{
		FirstName:            faker.FirstName(),
		LastName:             faker.LastName(),
		Email:                persistedUser.Email,
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
	s.context.Request = request

	s.handler.register(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())

	var user model.User
	err = s.db.First(&user, "first_name = ? AND last_name = ?", payload.FirstName, payload.LastName).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthHandlerTestSuite) TestRegister_Successfully() {
	password := faker.PASSWORD
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

func (s *AuthHandlerTestSuite) TestLogin_WithMalFormedEmail() {
	payload := dto.LoginRequest{
		Email:    "invalid-email",
		Password: "password",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
	s.context.Request = request

	s.handler.login(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())
}

func (s *AuthHandlerTestSuite) TestLogin_WithMalFormedPassword() {
	payload := dto.LoginRequest{
		Email:    faker.Email(),
		Password: "1234",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
	s.context.Request = request

	s.handler.login(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())
}

func (s *AuthHandlerTestSuite) TestLogin_UnknownEmail() {
	payload := dto.LoginRequest{
		Email:    faker.Email(),
		Password: "password",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
	s.context.Request = request

	s.handler.login(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())
}

func (s *AuthHandlerTestSuite) TestLogin_WithInvalidPassword() {
	user := factory.NewUser()
	err := s.db.Create(user).Error
	require.Nil(s.T(), err)

	payload := dto.LoginRequest{
		Email:    user.Email,
		Password: "wrong-password",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
	s.context.Request = request

	s.handler.login(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())
}

func (s *AuthHandlerTestSuite) TestLogin_Successfully() {
	password := faker.PASSWORD
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user := factory.NewUser()
	user.Password = string(hashedPassword)

	err := s.db.Create(user).Error
	require.Nil(s.T(), err)

	payload := dto.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
	s.context.Request = request

	s.handler.login(s.context)

	var credentials dto.CredentialsResponse
	err = json.NewDecoder(s.recorder.Result().Body).Decode(&credentials)
	require.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, s.recorder.Code)
	assert.NotEmpty(s.T(), credentials.Token)
}

func (s *AuthHandlerTestSuite) TestResetPassword_WithMalformedEmail() {
	payload := dto.PasswordResetRequest{
		Email: "invalid-email",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/password-reset", bytes.NewReader(data))
	s.context.Request = request

	s.handler.resetPassword(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())
}

func (s *AuthHandlerTestSuite) TestResetPassword_WithUnknownEmail() {
	user := factory.NewUser()

	payload := dto.PasswordResetRequest{
		Email: user.Email,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/password-reset", bytes.NewReader(data))
	s.context.Request = request

	s.handler.resetPassword(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())
}

func (s *AuthHandlerTestSuite) TestResetPassword_Successfully() {
	user := factory.NewUser()
	err := s.db.Create(user).Error

	payload := dto.PasswordResetRequest{
		Email: user.Email,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/password-reset", bytes.NewReader(data))
	s.context.Request = request

	s.handler.resetPassword(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())

	var updated model.User
	err = s.db.First(&updated, "id = ?", user.ID).Error
	require.Nil(s.T(), err)

	assert.NotEqual(s.T(), updated.Password, user.Password)
}
