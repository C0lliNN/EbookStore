package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/config"
	"github.com/c0llinn/ebook-store/internal/email"
	"github.com/c0llinn/ebook-store/internal/generator"
	"github.com/c0llinn/ebook-store/internal/hash"
	"github.com/c0llinn/ebook-store/internal/persistence"
	"github.com/c0llinn/ebook-store/internal/token"
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

type AuthenticatorHandlerTestSuite struct {
	suite.Suite
	baseURL  string
	context  *gin.Context
	recorder *httptest.ResponseRecorder
	db       *gorm.DB
	handler  *AuthenticationHandler
}

func (s *AuthenticatorHandlerTestSuite) SetupTest() {
	test.SetEnvironmentVariables()
	s.T().Log(viper.GetString("DATABASE_URL"))
	config.LoadMigrations("file:../../migrations")

	s.db = config.NewConnection()
	s.baseURL = fmt.Sprintf("http://localhost:%s", viper.GetString("PORT"))

	userRepository := persistence.NewUserRepository(s.db)
	jwtWrapper := token.NewJWTWrapper(token.NewHMACSecret())
	ses := config.NewSNSService()
	emailClient := email.NewEmailClient(ses)
	passwordGenerator := generator.NewPasswordGenerator()
	bcryptWrapper := hash.NewBcryptWrapper()
	idGenerator := generator.NewUUIDGenerator()

	authenticator := auth.New(auth.Config{
		Repository:        userRepository,
		Tokener:           jwtWrapper,
		Hasher:            bcryptWrapper,
		EmailClient:       emailClient,
		PasswordGenerator: passwordGenerator,
		IDGenerator:       idGenerator,
	})

	s.handler = NewAuthenticatorHandler(gin.New(), authenticator)

	s.recorder = httptest.NewRecorder()
	s.context, _ = gin.CreateTestContext(s.recorder)
}

func (s *AuthenticatorHandlerTestSuite) TearDownTest() {
	s.db.Delete(&auth.User{}, "1 = 1")
}

func TestAuthenticatorHandler(t *testing.T) {
	suite.Run(t, new(AuthenticatorHandlerTestSuite))
}

func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidFirstName() {
	password := faker.PASSWORD
	payload := auth.RegisterRequest{
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

	var user auth.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidLastName() {
	password := faker.PASSWORD
	payload := auth.RegisterRequest{
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

	var user auth.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidEmail() {
	password := faker.PASSWORD
	payload := auth.RegisterRequest{
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

	var user auth.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidPassword() {
	payload := auth.RegisterRequest{
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

	var user auth.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidPasswordConfirmation() {
	password := faker.PASSWORD
	payload := auth.RegisterRequest{
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

	var user auth.User
	err = s.db.First(&user).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthenticatorHandlerTestSuite) TestRegister_WithDuplicateEmail() {
	persistedUser := factory.NewUser()
	err := s.db.Create(persistedUser).Error
	require.Nil(s.T(), err)

	password := faker.PASSWORD
	payload := auth.RegisterRequest{
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

	var user auth.User
	err = s.db.First(&user, "first_name = ? AND last_name = ?", payload.FirstName, payload.LastName).Error
	assert.NotNil(s.T(), err)
}

func (s *AuthenticatorHandlerTestSuite) TestRegister_Successfully() {
	password := faker.PASSWORD
	payload := auth.RegisterRequest{
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

	var credentials auth.CredentialsResponse
	err = json.NewDecoder(s.recorder.Result().Body).Decode(&credentials)
	require.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusCreated, s.recorder.Code)
	assert.NotEmpty(s.T(), credentials.Token)

	var user auth.User
	err = s.db.First(&user).Error
	require.Nil(s.T(), err)

	assert.Len(s.T(), user.ID, 36)
	assert.Equal(s.T(), payload.FirstName, user.FirstName)
	assert.Equal(s.T(), payload.LastName, user.LastName)
	assert.Equal(s.T(), payload.Email, user.Email)
	assert.NotEmpty(s.T(), user.Password)
	assert.Equal(s.T(), auth.Customer, user.Role)
	assert.NotZero(s.T(), user.CreatedAt)
}

func (s *AuthenticatorHandlerTestSuite) TestLogin_WithMalFormedEmail() {
	payload := auth.LoginRequest{
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

func (s *AuthenticatorHandlerTestSuite) TestLogin_WithMalFormedPassword() {
	payload := auth.LoginRequest{
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

func (s *AuthenticatorHandlerTestSuite) TestLogin_UnknownEmail() {
	payload := auth.LoginRequest{
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

func (s *AuthenticatorHandlerTestSuite) TestLogin_WithInvalidPassword() {
	user := factory.NewUser()
	err := s.db.Create(user).Error
	require.Nil(s.T(), err)

	payload := auth.LoginRequest{
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

func (s *AuthenticatorHandlerTestSuite) TestLogin_Successfully() {
	password := faker.PASSWORD
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user := factory.NewUser()
	user.Password = string(hashedPassword)

	err := s.db.Create(user).Error
	require.Nil(s.T(), err)

	payload := auth.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
	s.context.Request = request

	s.handler.login(s.context)

	var credentials auth.CredentialsResponse
	err = json.NewDecoder(s.recorder.Result().Body).Decode(&credentials)
	require.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, s.recorder.Code)
	assert.NotEmpty(s.T(), credentials.Token)
}

func (s *AuthenticatorHandlerTestSuite) TestResetPassword_WithMalformedEmail() {
	payload := auth.PasswordResetRequest{
		Email: "invalid-email",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/password-reset", bytes.NewReader(data))
	s.context.Request = request

	s.handler.resetPassword(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())
}

func (s *AuthenticatorHandlerTestSuite) TestResetPassword_WithUnknownEmail() {
	user := factory.NewUser()

	payload := auth.PasswordResetRequest{
		Email: user.Email,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/password-reset", bytes.NewReader(data))
	s.context.Request = request

	s.handler.resetPassword(s.context)

	assert.NotEmpty(s.T(), s.context.Errors.Errors())
}

func (s *AuthenticatorHandlerTestSuite) TestResetPassword_Successfully() {
	user := factory.NewUser()
	err := s.db.Create(user).Error

	payload := auth.PasswordResetRequest{
		Email: user.Email,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	request := httptest.NewRequest("POST", s.baseURL+"/password-reset", bytes.NewReader(data))
	s.context.Request = request

	s.handler.resetPassword(s.context)

	assert.Empty(s.T(), s.context.Errors.Errors())

	var updated auth.User
	err = s.db.First(&updated, "id = ?", user.ID).Error
	require.Nil(s.T(), err)

	assert.NotEqual(s.T(), updated.Password, user.Password)
}
