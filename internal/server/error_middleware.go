package server

import (
	"fmt"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/gin-gonic/gin"
)

type ErrorMiddleware struct{}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

type ErrorResponse struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func fromNotValid(err *common.ErrNotValid) *ErrorResponse {
	return &ErrorResponse{
		Code:    400,
		Message: "The provided payload is not valid",
		Details: err.Error(),
	}
}

func fromDuplicateKey(err *common.ErrDuplicateKey) *ErrorResponse {
	return &ErrorResponse{
		Code:    409,
		Message: fmt.Sprintf("this %s is already being used", err.Key),
		Details: err.Error(),
	}
}

func fromEntityNotFound(err *common.ErrEntityNotFound) *ErrorResponse {
	return &ErrorResponse{
		Code:    404,
		Message: fmt.Sprintf("%s with the provided parameters could not be found", err.Entity),
		Details: err.Error(),
	}
}

func fromWrongPassword(err *common.ErrWrongPassword) *ErrorResponse {
	return &ErrorResponse{
		Code:    401,
		Message: "the provided password is invalid",
		Details: err.Error(),
	}
}

func fromOrderNotPaid(err *common.ErrOrderNotPaid) *ErrorResponse {
	return &ErrorResponse{
		Code:    402,
		Message: "you are allowed to download unpaid orders",
		Details: err.Error(),
	}
}

func fromGeneric(err error) *ErrorResponse {
	return &ErrorResponse{
		Code:    500,
		Message: "Some unexpected error happened",
		Details: err.Error(),
	}
}

func (e *ErrorMiddleware) Handler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		if len(context.Errors) <= 0 {
			return
		}
		log.FromContext(context).Warn("error found: ", context.Errors.Last().Err)

		var response *ErrorResponse

		switch err := context.Errors.Last().Err.(type) {
		case *common.ErrNotValid:
			response = fromNotValid(err)
		case *common.ErrDuplicateKey:
			response = fromDuplicateKey(err)
		case *common.ErrEntityNotFound:
			response = fromEntityNotFound(err)
		case *common.ErrWrongPassword:
			response = fromWrongPassword(err)
		case *common.ErrOrderNotPaid:
			response = fromOrderNotPaid(err)
		default:
			response = fromGeneric(err)
		}

		context.AbortWithStatusJSON(response.Code, response)
	}
}
