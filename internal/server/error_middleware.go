package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/c0llinn/ebook-store/internal/persistence"
	"github.com/c0llinn/ebook-store/internal/shop"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorMiddleware struct{}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

type ErrorResponse struct {
	Code    int      `json:"-"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func (e *ErrorMiddleware) Handler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		if len(context.Errors) <= 0 {
			return
		}
		err := context.Errors.Last().Err

		log.FromContext(context).Warnf("error processing the request: %v", err)

		var (
			response          *ErrorResponse
			validationErr     validator.ValidationErrors
			duplicateKeyErr   *persistence.ErrDuplicateKey
			entityNotFoundErr *persistence.ErrEntityNotFound
		)

		switch {
		case errors.Is(err, auth.ErrWrongPassword):
			response = newErrorResponse(http.StatusUnauthorized, err.Error())
		case errors.Is(err, catalog.ErrForbiddenCatalogAccess), errors.Is(err, shop.ErrForbiddenOrderAccess):
			response = newErrorResponse(http.StatusForbidden, err.Error())
		case errors.Is(err, fmt.Errorf("invalid request")), errors.Is(err, fmt.Errorf("failed binding request body")):
			response = newErrorResponse(http.StatusBadRequest, "invalid request body. check the documentation")
		case errors.Is(err, shop.ErrOrderNotPaid):
			response = newErrorResponse(http.StatusPaymentRequired, err.Error())
		case errors.As(err, &validationErr):
			response = newValidationErrorResponse(validationErr)
		case errors.As(err, &duplicateKeyErr):
			response = newErrorResponse(http.StatusConflict, duplicateKeyErr.Error())
		case errors.As(err, &entityNotFoundErr):
			response = newErrorResponse(http.StatusNotFound, entityNotFoundErr.Error())
		default:
			response = newGenericErrorResponse(err)
		}

		context.JSON(response.Code, response)
	}
}

func newErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

func newValidationErrorResponse(errors validator.ValidationErrors) *ErrorResponse {
	details := make([]string, 0, len(errors))

	for _, err := range errors {
		details = append(details, fmt.Sprintf("the validation of the field %s for tag %s failed", err.Field(), err.Tag()))
	}

	return &ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "the payload is not valid",
		Details: details,
	}
}

func newGenericErrorResponse(err error) *ErrorResponse {
	details := strings.Split(err.Error(), ":")
	for i := range details {
		details[i] = strings.TrimSpace(details[i])
	}

	return &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Some unexpected error happened",
		Details: details,
	}
}
