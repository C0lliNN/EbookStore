package server_test

import (
	"net/http"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/ebookstore/internal/core/auth"
)

func (s *ServerSuiteTest) TestRegister_Successfully() {
	password := "password"

	apitest.New().
		EnableNetworking(http.DefaultClient).
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
		Assert(jsonpath.Present("$.token")).
		End()
}

func (s *ServerSuiteTest) TestRegister_WithInvalidData() {
	password := "password"

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/register").
		JSON(auth.RegisterRequest{
			FirstName:            "",
			LastName:             "Collin",
			Email:                "raphael@test.com",
			Password:             password,
			PasswordConfirmation: password,
		}).
		Expect(s.T()).
		Status(http.StatusBadRequest).
		Assert(jsonpath.Equal("$.message", "the payload is not valid")).
		End()
}

func (s *ServerSuiteTest) TestLogin_Failure() {
	s.createDefaultCustomer()

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/login").
		JSON(auth.LoginRequest{
			Email:    "raphael@test.com",
			Password: "password2",
		}).
		Expect(s.T()).
		Status(http.StatusUnauthorized).
		End()
}

func (s *ServerSuiteTest) TestLogin_Success() {
	s.createDefaultCustomer()

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/login").
		JSON(auth.LoginRequest{
			Email:    "raphael@test.com",
			Password: "password",
		}).
		Expect(s.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.token")).
		End()
}

func (s *ServerSuiteTest) TestResetPassword_Failure() {
	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/password-reset").
		JSON(auth.PasswordResetRequest{
			Email: "raphael@test.com",
		}).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.message", "the provided User was not found")).
		End()
}

func (s *ServerSuiteTest) TestResetPassword_Success() {
	s.createDefaultCustomer()

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/password-reset").
		JSON(auth.PasswordResetRequest{
			Email: "raphael@test.com",
		}).
		Expect(s.T()).
		Status(http.StatusNoContent).
		End()
}
