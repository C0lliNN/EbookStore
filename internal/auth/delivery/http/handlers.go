package http

import (
	"github.com/c0llinn/ebook-store/internal/auth/delivery/dto"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IDGenerator interface {
	NewID() string
}

type UseCase interface {
	Register(user model.User) (model.Credentials, error)
}

type AuthHandler struct {
	useCase UseCase
	idGenerator IDGenerator
}

func NewAuthHandler(useCase UseCase, idGenerator IDGenerator) AuthHandler {
	return AuthHandler{useCase: useCase, idGenerator: idGenerator}
}

func (h AuthHandler) register(context *gin.Context) {
	var r dto.RegisterRequest
	if err := context.ShouldBindJSON(&r); err != nil {
		context.Error(err)
		return
	}

	user := r.ToDomain(h.idGenerator.NewID())
	credentials, err := h.useCase.Register(user)
	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusCreated, dto.FromCredentials(credentials))
}