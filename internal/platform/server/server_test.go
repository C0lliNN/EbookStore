package server_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ebookstore/internal/container"
	"github.com/ebookstore/internal/core/auth"
	"github.com/ebookstore/internal/platform/config"
	"github.com/ebookstore/test"
	"github.com/spf13/viper"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ServerSuiteTest struct {
	suite.Suite
	container 		 	*container.Container
	baseURL string
	postgresContainer   *test.PostgresContainer
	localstackContainer *test.LocalstackContainer
}

func (s *ServerSuiteTest) SetupSuite() {
	config.LoadConfiguration()

	var err error
	ctx := context.TODO()

	s.postgresContainer, err = test.NewPostgresContainer(ctx)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.localstackContainer, err = test.NewLocalstackContainer(ctx)
	if err != nil {
		s.FailNow(err.Error())
	}

	viper.Set("DATABASE_URI", s.postgresContainer.URI)
	viper.Set("AWS_SES_ENDPOINT", fmt.Sprintf("http://localhost:%v", s.localstackContainer.Port))
	viper.Set("AWS_S3_ENDPOINT", fmt.Sprintf("http://s3.localhost.localstack.cloud:%v", s.localstackContainer.Port))

	s.baseURL = fmt.Sprintf("http://%v", viper.GetString("SERVER_ADDR"))
	s.container = container.New()
	
	go func () {
		s.container.Start(context.TODO())
	}()

	require.Eventually(s.T(), func() bool {
		response, err := http.Get(s.baseURL + "/api/v1/healthcheck")
		if err != nil {
			return false
		}

		return response.StatusCode == http.StatusOK
	}, time.Second*15, time.Millisecond*500)
}

func (s *ServerSuiteTest) TearDownTest() {
	s.container.DB().Exec("DELETE FROM users")
	s.container.DB().Exec("DELETE FROM books")
	s.container.DB().Exec("DELETE FROM orders")
}

func (s *ServerSuiteTest) TearDownSuite() {
	ctx := context.TODO()

	_ = s.postgresContainer.Terminate(ctx)
	_ = s.localstackContainer.Terminate(ctx)
}

func TestServer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	suite.Run(t, new(ServerSuiteTest))
}

func (s *ServerSuiteTest) createCustomer() string {
	password := "password"
	
	var response auth.CredentialsResponse

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/register").
		JSON(auth.RegisterRequest{
			FirstName:            "Raphael",
			LastName:             "Collin",
			Email:                "raphael@test.com",
			Password:             password,
			PasswordConfirmation: password,
		}).
		Expect(s.T()).
		Status(http.StatusCreated).
		End().
		JSON(&response)

	return response.Token
}

func (s *ServerSuiteTest) createAdmin() string {
	password := "password"
	
	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/register").
		JSON(auth.RegisterRequest{
			FirstName:            "Raphael",
			LastName:             "Collin",
			Email:                "raphael2@test.com",
			Password:             password,
			PasswordConfirmation: password,
		}).
		Expect(s.T()).
		Status(http.StatusCreated).
		End()

	result := s.container.DB().Model(&auth.User{}).Where("email = 'raphael2@test.com'").Update("role", auth.Admin)
	require.NoError(s.T(), result.Error)

	var response auth.CredentialsResponse

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/login").
		JSON(auth.LoginRequest{
			Email:                "raphael2@test.com",
			Password:             password,
		}).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&response)

	return response.Token
}
