//nolint:unused
package server_test

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ebookstore/internal/core/catalog"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/require"
)

func (s *ServerSuiteTest) TestCreateBook_Unauthorized() {
	token := s.createCustomer()

	req := catalog.CreateBook{
		Title:       "Domain Design Driven",
		Description: "Complexity",
		AuthorName: "Eric Evans",
		ContentID: "123",
		Images: []catalog.ImageRequest{
			{ID: "image-1", Description: "The first image"},
		},
		Price: 40000,
		ReleaseDate: time.Date(2002, time.October, 28, 0, 0, 0, 0, time.UTC),
	}

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/books").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		JSON(req).
		Expect(s.T()).
		Status(http.StatusForbidden).
		Assert(jsonpath.Equal("$.message", "the access to this action is restricted to allowed users")).
		End()
}

func (s *ServerSuiteTest) TestCreateBook_InvalidPayload() {
	token := s.createAdmin()

	req := catalog.CreateBook{
		Title:       "",
		Description: "Complexity",
		AuthorName: "Eric Evans",
		ContentID: "123",
		Images: []catalog.ImageRequest{
			{ID: "image-1", Description: "The first image"},
		},
		Price: 40000,
		ReleaseDate: time.Date(2002, time.October, 28, 0, 0, 0, 0, time.UTC),
	}

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/books").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		JSON(req).
		Expect(s.T()).
		Status(http.StatusBadRequest).
		Assert(jsonpath.Equal("$.message", "the payload is not valid")).
		End()
}

func (s *ServerSuiteTest) TestCreateBook_Success() {
	response := s.createBook()
	s.T().Log(response)
}

func (s *ServerSuiteTest) createBook() catalog.BookResponse {
	token := s.createAdmin()

	var presignResponse catalog.PresignURLResponse

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/presign-url").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&presignResponse)

	s.putFile(presignResponse.URL, "testdata/book1_image.jpg")
	imageID := presignResponse.ID

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/presign-url").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&presignResponse)

	s.putFile(presignResponse.URL, "testdata/book1_content.pdf")
	contentID := presignResponse.ID

	var bookResponse catalog.BookResponse

	apitest.New().
		EnableNetworking().
		Post(s.baseURL + "/api/v1/books").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		JSON(catalog.CreateBook{
			Title:       "Domain Driven Design",
			Description: "Complexity",
			AuthorName: "Eric Evans",
			ContentID: contentID,
			Images: []catalog.ImageRequest{
				{ID: imageID, Description: "The first image"},
			},
			Price: 40000,
			ReleaseDate: time.Date(2002, time.October, 28, 0, 0, 0, 0, time.UTC),
		}).
		Expect(s.T()).
		Status(http.StatusCreated).
		End().
		JSON(&bookResponse)

	return bookResponse
}

func (s *ServerSuiteTest) putFile(url string, filepath string) {
	file, err := os.Open(filepath)
	require.NoError(s.T(), err)

	defer file.Close()

	request, err := http.NewRequest(http.MethodPut, url, file)
	require.NoError(s.T(), err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(s.T(), err)

	defer response.Body.Close()
	require.Equal(s.T(), http.StatusOK, response.StatusCode)
}
