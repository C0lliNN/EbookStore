package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ebookstore/internal/core/auth"
	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/core/shop"
	"github.com/ebookstore/internal/log"
	"github.com/ebookstore/internal/platform/persistence"
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

		log.Warnf(context, "error processing the request: %v", err)

		var (
			response          *ErrorResponse
			bindingErr        *BindingErr
			validationErr     validator.ValidationErrors
			duplicateKeyErr   *persistence.ErrDuplicateKey
			entityNotFoundErr *persistence.ErrEntityNotFound
		)

		switch {
		case errors.As(err, &bindingErr):
			response = newBindingErrorResponse(bindingErr)
		case errors.As(err, &validationErr):
			response = newValidationErrorResponse(validationErr)
		case errors.As(err, &entityNotFoundErr):
			response = newErrorResponse(http.StatusNotFound, entityNotFoundErr)
		case errors.As(err, &duplicateKeyErr):
			response = newErrorResponse(http.StatusConflict, duplicateKeyErr)
		case errors.Is(err, auth.ErrWrongPassword):
			response = newErrorResponse(http.StatusUnauthorized, err)
		case errors.Is(err, catalog.ErrForbiddenCatalogAccess), errors.Is(err, shop.ErrForbiddenOrderAccess):
			response = newErrorResponse(http.StatusForbidden, err)
		case errors.Is(err, shop.ErrOrderNotPaid):
			response = newErrorResponse(http.StatusPaymentRequired, err)
		default:
			response = newGenericErrorResponse(err)
		}

		context.JSON(response.Code, response)
	}
}

func newErrorResponse(code int, err error) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: unwrappedError(err).Error(),
		Details: errorStack(err),
	}
}

func newBindingErrorResponse(err *BindingErr) *ErrorResponse {
	return &ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "invalid request body. check the documentation",
		Details: errorStack(err),
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
	return &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Some unexpected error happened",
		Details: errorStack(err),
	}
}

func unwrappedError(err error) error {
	if errors.Unwrap(err) == nil {
		return err
	}

	return unwrappedError(errors.Unwrap(err))
}

func errorStack(err error) []string {
	stack := strings.Split(err.Error(), ":")
	for i := range stack {
		stack[i] = strings.TrimSpace(stack[i])
	}

	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}

	return stack
}
