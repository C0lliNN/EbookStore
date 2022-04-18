package api

import (
	"fmt"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/config"
	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func fromNotValid(err *common.ErrNotValid) *Error {
	return &Error{
		Code:    400,
		Message: "The provided payload is not valid",
		Details: err.Error(),
	}
}

func fromDuplicateKey(err *common.ErrDuplicateKey) *Error {
	return &Error{
		Code:    409,
		Message: fmt.Sprintf("this %s is already being used", err.Key),
		Details: err.Error(),
	}
}

func fromEntityNotFound(err *common.ErrEntityNotFound) *Error {
	return &Error{
		Code:    404,
		Message: fmt.Sprintf("%s with the provided parameters could not be found", err.Entity),
		Details: err.Error(),
	}
}

func fromWrongPassword(err *common.ErrWrongPassword) *Error {
	return &Error{
		Code:    401,
		Message: "the provided password is invalid",
		Details: err.Error(),
	}
}

func fromOrderNotPaid(err *common.ErrOrderNotPaid) *Error {
	return &Error{
		Code:    402,
		Message: "you are allowed to download unpaid orders",
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
			config.Logger.Warn("error detected: " + detectedErrors.Errors()[0])
			err := detectedErrors[0].Err

			switch parsed := err.(type) {
			case *common.ErrNotValid:
				apiError = fromNotValid(parsed)
			case *common.ErrDuplicateKey:
				apiError = fromDuplicateKey(parsed)
			case *common.ErrEntityNotFound:
				apiError = fromEntityNotFound(parsed)
			case *common.ErrWrongPassword:
				apiError = fromWrongPassword(parsed)
			case *common.ErrOrderNotPaid:
				apiError = fromOrderNotPaid(parsed)
			default:
				apiError = fromGeneric(parsed)
			}

			context.AbortWithStatusJSON(apiError.Code, apiError)
		}
	}
}
