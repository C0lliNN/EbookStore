package http

import (
	"github.com/c0llinn/ebook-store/internal/auth/delivery/dto"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IDGenerator interface {
	NewID() string
}

type UseCase interface {
	Register(user model.User) (model.Credentials, error)
	Login(email, password string) (model.Credentials, error)
	ResetPassword(email string) error
}

type AuthHandler struct {
	useCase UseCase
	idGenerator IDGenerator
}

func NewAuthHandler(useCase UseCase, idGenerator IDGenerator) AuthHandler {
	return AuthHandler{useCase: useCase, idGenerator: idGenerator}
}

// register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce  json
// @Param payload body dto.RegisterRequest true "Register Payload"
// @Success 201 {object} dto.CredentialsResponse
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /register [post]
func (h AuthHandler) register(context *gin.Context) {
	var r dto.RegisterRequest
	if err := context.ShouldBindJSON(&r); err != nil {
		context.Error(&common.ErrNotValid{Input: "RegisterRequest", Err: err})
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

// login godoc
// @Summary Login using email and password
// @Tags Auth
// @Accept json
// @Produce  json
// @Param payload body dto.LoginRequest true "Register Payload"
// @Success 201 {object} dto.CredentialsResponse
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /login [post]
func (h AuthHandler) login(context *gin.Context) {
	var r dto.LoginRequest
	if err := context.ShouldBindJSON(&r); err != nil {
		context.Error(&common.ErrNotValid{Input: "LoginRequest", Err: err})
		return
	}

	credentials, err := h.useCase.Login(r.Email, r.Password)
	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, dto.FromCredentials(credentials))
}

// resetPassword godoc
// @Summary Reset the password for the given email
// @Tags Auth
// @Accept json
// @Produce  json
// @Param payload body dto.PasswordResetRequest true "Register Payload"
// @Success 204 "success"
// @Failure 400 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /password-reset [post]
func (h AuthHandler) resetPassword(context *gin.Context) {
	var r dto.PasswordResetRequest

	if err := context.ShouldBindJSON(&r); err != nil {
		context.Error(&common.ErrNotValid{Input: "PasswordResetRequest", Err: err})
		return
	}

	if err := h.useCase.ResetPassword(r.Email); err != nil {
		context.Error(err)
		return
	}

	context.Status(http.StatusNoContent)

}