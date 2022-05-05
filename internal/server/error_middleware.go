package server

import (
	"fmt"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/c0llinn/ebook-store/internal/persistence"
	"github.com/c0llinn/ebook-store/internal/shop"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
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
		log.FromContext(context).Warn("an error occurred when processing the request: ", context.Errors.Last().Err)

		var response *ErrorResponse

		err := context.Errors.Last().Err
		switch err {
		case auth.ErrWrongPassword:
			response = newErrorResponse(http.StatusUnauthorized, err.Error())
		case catalog.ErrForbiddenCatalogAccess, shop.ErrForbiddenOrderAccess:
			response = newErrorResponse(http.StatusForbidden, err.Error())
		case fmt.Errorf("invalid request"):
			response = newErrorResponse(http.StatusBadRequest, "invalid request body. check the documentation")
		case shop.ErrOrderNotPaid:
			response = newErrorResponse(http.StatusPaymentRequired, err.Error())
		}

		if response != nil {
			context.AbortWithStatusJSON(response.Code, response)
			return
		}

		switch err := err.(type) {
		case *validator.ValidationErrors:
			response = newValidationErrorResponse(err)
		case *persistence.ErrDuplicateKey:
			response = newErrorResponse(http.StatusConflict, err.Error())
		case *persistence.ErrEntityNotFound:
			response = newErrorResponse(http.StatusNotFound, err.Error())
		default:
			response = newGenericErrorResponse(err)
		}

		context.AbortWithStatusJSON(response.Code, response)
	}
}

func newErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

func newValidationErrorResponse(errors *validator.ValidationErrors) *ErrorResponse {
	details := make([]string, 0, len(*errors))

	for _, err := range *errors {
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
		Details: []string{err.Error()},
	}
}
