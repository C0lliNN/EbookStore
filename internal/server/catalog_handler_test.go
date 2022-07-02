//nolint:unused
package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/ebookstore/internal/auth"
	"github.com/ebookstore/internal/catalog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CatalogHandlerTestSuite struct {
	ServerSuiteTest
}

func (s *CatalogHandlerTestSuite) TearDownTest() {
	s.db.Delete(&auth.User{}, "1 = 1")
	s.db.Delete(&catalog.Book{}, "1 = 1")
}

func TestCatalogHandler(t *testing.T) {
	//if testing.Short() {
	//	t.Skip("skipping test in short mode.")
	//}

	t.Skip()
	suite.Run(t, new(CatalogHandlerTestSuite))
}

func (s *CatalogHandlerTestSuite) TestCreateBook_Failure() {
	token := s.createCustomer()

	payload := new(bytes.Buffer)
	mp := multipart.NewWriter(payload)

	_ = mp.SetBoundary("---WebKitFormBoundary7MA4YWxkTrZu0gW")
	_ = mp.WriteField("title", "Domain Driver Design")
	_ = mp.WriteField("description", "Complexity")
	_ = mp.WriteField("authorName", "Eric Evans")
	_ = mp.WriteField("price", "40000")
	_ = mp.WriteField("releaseDate", time.Date(2002, time.October, 28, 0, 0, 0, 0, time.UTC).String())
	_, _ = mp.CreateFormFile("posterImage", "poster.png")
	_, _ = mp.CreateFormFile("bookContent", "content.pdf")

	request, err := http.NewRequest(http.MethodPost, s.baseURL+"/books", payload)
	require.Nil(s.T(), err)
	request.Header.Add("Content-Type", "multipart/form-data;boundary=---WebKitFormBoundary7MA4YWxkTrZu0gW")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	client := &http.Client{}
	response, err := client.Do(request)

	require.Nil(s.T(), err)

	require.Equal(s.T(), http.StatusForbidden, response.Status)
}

func (s *CatalogHandlerTestSuite) createCustomer() string {
	password := "password"
	payload := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	response, err := http.Post(s.baseURL+"/register", "application/json", bytes.NewReader(data))

	require.Nil(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, response.StatusCode)

	credentials := &auth.CredentialsResponse{}
	err = json.NewDecoder(response.Body).Decode(credentials)
	require.Nil(s.T(), err)

	return credentials.Token
}

func (s *CatalogHandlerTestSuite) createAdmin() {
	password := "password"
	payload := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael2@test.com",
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	response, err := http.Post(s.baseURL+"/register", "application/json", bytes.NewReader(data))
	require.Nil(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, response.StatusCode)

	result := s.db.Model(&auth.User{}).Where("email = 'raphael2@test.com'").Update("role", auth.Admin)
	require.Nil(s.T(), result.Error)
}
