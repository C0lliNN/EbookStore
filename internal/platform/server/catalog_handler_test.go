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
	token := s.createDefaultCustomer()

	req := catalog.CreateBook{
		Title:       "Domain Design Driven",
		Description: "Complexity",
		AuthorName:  "Eric Evans",
		ContentID:   "123",
		Images: []catalog.ImageRequest{
			{ID: "image-1", Description: "The first image"},
		},
		Price:       40000,
		ReleaseDate: time.Date(2002, time.October, 28, 0, 0, 0, 0, time.UTC),
	}

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/books").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		JSON(req).
		Expect(s.T()).
		Status(http.StatusForbidden).
		Assert(jsonpath.Equal("$.message", "the access to this action is restricted to allowed users")).
		End()
}

func (s *ServerSuiteTest) TestCreateBook_InvalidPayload() {
	token := s.createDefaultAdmin()

	req := catalog.CreateBook{
		Title:       "",
		Description: "Complexity",
		AuthorName:  "Eric Evans",
		ContentID:   "123",
		Images: []catalog.ImageRequest{
			{ID: "image-1", Description: "The first image"},
		},
		Price:       40000,
		ReleaseDate: time.Date(2002, time.October, 28, 0, 0, 0, 0, time.UTC),
	}

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/books").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		JSON(req).
		Expect(s.T()).
		Status(http.StatusBadRequest).
		Assert(jsonpath.Equal("$.message", "the payload is not valid")).
		End()
}

func (s *ServerSuiteTest) TestCreateBook_Success() {
	token := s.createDefaultAdmin()
	response := s.createBook(token)
	s.T().Log(response)
}

func (s *ServerSuiteTest) TestGetBooks_Success() {
	book := s.createBook(s.createDefaultAdmin())
	token := s.createDefaultCustomer()

	expected := catalog.PaginatedBooksResponse{
		Results:     []catalog.BookResponse{book},
		CurrentPage: 1,
		PerPage:     15,
		TotalPages:  1,
		TotalItems:  1,
	}
	var actual catalog.PaginatedBooksResponse

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/books").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&actual)

	s.Equal(expected.CurrentPage, actual.CurrentPage)
	s.Equal(expected.PerPage, actual.PerPage)
	s.Equal(expected.TotalPages, actual.TotalPages)
	s.Equal(expected.TotalItems, actual.TotalItems)
	s.Equal(expected.Results[0].ID, actual.Results[0].ID)
}

func (s *ServerSuiteTest) TestGetBook_NotFound() {
	token := s.createDefaultCustomer()

	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/books/"+token).
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.message", "the provided book was not found")).
		End()
}

func (s *ServerSuiteTest) TestGetBook_Success() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())

	var actual catalog.BookResponse
	apitest.New().
		EnableNetworking().
		Get(s.baseURL+"/api/v1/books/"+book.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&actual)

	s.Equal(book.ID, actual.ID)
	s.Equal(book.Title, actual.Title)
	s.Equal(book.Description, actual.Description)
	s.Equal(book.AuthorName, actual.AuthorName)
}

func (s *ServerSuiteTest) TestUpdateBook_Unauthorized() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())

	req := catalog.UpdateBook{
		ID:          book.ID,
		Title:       &book.Title,
		Description: &book.Description,
		AuthorName:  &book.AuthorName,
	}

	apitest.New().
		EnableNetworking().
		Patch(s.baseURL+"/api/v1/books/"+book.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		JSON(req).
		Expect(s.T()).
		Status(http.StatusForbidden).
		Assert(jsonpath.Equal("$.message", "the access to this action is restricted to allowed users")).
		End()
}

func (s *ServerSuiteTest) TestUpdateBook_Success() {
	token := s.createDefaultAdmin()
	book := s.createBook(token)

	req := catalog.UpdateBook{
		ID:          book.ID,
		Title:       &book.Title,
		Description: &book.Description,
		AuthorName:  &book.AuthorName,
	}

	apitest.New().
		EnableNetworking().
		Patch(s.baseURL+"/api/v1/books/"+book.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		JSON(req).
		Expect(s.T()).
		Status(http.StatusNoContent).
		End()
}

func (s *ServerSuiteTest) TestDeleteBook_Unauthorized() {
	token := s.createDefaultCustomer()
	book := s.createBook(s.createDefaultAdmin())

	apitest.New().
		EnableNetworking().
		Delete(s.baseURL+"/api/v1/books/"+book.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusForbidden).
		Assert(jsonpath.Equal("$.message", "the access to this action is restricted to allowed users")).
		End()
}

func (s *ServerSuiteTest) TestDeleteBook_NotFound() {
	token := s.createDefaultAdmin()

	apitest.New().
		EnableNetworking().
		Delete(s.baseURL+"/api/v1/books/id-1").
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.message", "the provided book was not found")).
		End()
}

func (s *ServerSuiteTest) TestDeleteBook_Success() {
	token := s.createDefaultAdmin()
	book := s.createBook(token)

	apitest.New().
		EnableNetworking().
		Delete(s.baseURL+"/api/v1/books/"+book.ID).
		Header("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect(s.T()).
		Status(http.StatusNoContent).
		End()
}

func (s *ServerSuiteTest) createBook(authToken string) catalog.BookResponse {
	var presignResponse catalog.PresignURLResponse

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/presign-url").
		Header("Authorization", fmt.Sprintf("Bearer %v", authToken)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&presignResponse)

	s.putFile(presignResponse.URL, "testdata/book1_image.jpg")
	imageID := presignResponse.ID

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/presign-url").
		Header("Authorization", fmt.Sprintf("Bearer %v", authToken)).
		Expect(s.T()).
		Status(http.StatusOK).
		End().
		JSON(&presignResponse)

	s.putFile(presignResponse.URL, "testdata/book1_content.pdf")
	contentID := presignResponse.ID

	var bookResponse catalog.BookResponse

	apitest.New().
		EnableNetworking().
		Post(s.baseURL+"/api/v1/books").
		Header("Authorization", fmt.Sprintf("Bearer %v", authToken)).
		JSON(catalog.CreateBook{
			Title:       "Domain Driven Design",
			Description: "Complexity",
			AuthorName:  "Eric Evans",
			ContentID:   contentID,
			Images: []catalog.ImageRequest{
				{ID: imageID, Description: "The first image"},
			},
			Price:       40000,
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
