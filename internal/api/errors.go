package api

import (
	"fmt"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func fromNotValid(err *model.ErrNotValid) *Error {
	return &Error{
		Code:    400,
		Message: "The provided payload is not valid",
		Details: err.Error(),
	}
}

func fromDuplicateKey(err *model.ErrDuplicateKey) *Error {
	return &Error{
		Code:    409,
		Message: fmt.Sprintf("this %s is already being used", err.Key),
		Details: err.Error(),
	}
}

func fromEntityNotFound(err *model.ErrEntityNotFound) *Error {
	return &Error{
		Code:    404,
		Message: fmt.Sprintf("%s with the provided parameters could not be found", err.Entity),
		Details: err.Error(),
	}
}

func fromWrongPassword(err *model.ErrWrongPassword) *Error {
	return &Error{
		Code:    401,
		Message: "the provided password is invalid",
		Details: err.Error(),
	}
}

func fromGeneric(err error) *Error {
	return &Error{
		Code:    500,
		Message: "Some unexpected error happened",
		Details: err.Error(),
	}
}

func Errors() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		detectedErrors := context.Errors.ByType(gin.ErrorTypeAny)
		var apiError *Error

		if len(detectedErrors) > 0 {
			log.Logger.Warn("error detected: " + detectedErrors.Errors()[0])
			err := detectedErrors[0].Err

			switch parsed := err.(type) {
			case *model.ErrNotValid:
				apiError = fromNotValid(parsed)
			case *model.ErrDuplicateKey:
				apiError = fromDuplicateKey(parsed)
			case *model.ErrEntityNotFound:
				apiError = fromEntityNotFound(parsed)
			case *model.ErrWrongPassword:
				apiError = fromWrongPassword(parsed)
			default:
				apiError = fromGeneric(parsed)
			}

			context.AbortWithStatusJSON(apiError.Code, apiError)
		}
	}
}