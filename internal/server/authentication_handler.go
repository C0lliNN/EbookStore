package server

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Authenticator interface {
	Register(context.Context, auth.RegisterRequest) (auth.CredentialsResponse, error)
	Login(context.Context, auth.LoginRequest) (auth.CredentialsResponse, error)
	ResetPassword(context.Context, auth.PasswordResetRequest) error
}

type AuthenticatorHandler struct {
	engine        *gin.Engine
	authenticator Authenticator
}

func NewAuthenticatorHandler(engine *gin.Engine, authenticator Authenticator) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		engine: engine,
		authenticator: authenticator,
	}
}

func (h *AuthenticatorHandler) RegisterRoutes(engine *gin.Engine) {
	engine.POST("/register", h.register)
	engine.POST("/login", h.login)
	engine.POST("/password-reset", h.resetPassword)
}

// register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce  json
// @Param payload body auth.RegisterRequest true "Register Payload"
// @Success 201 {object} auth.CredentialsResponse
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /register [post]
func (h *AuthenticatorHandler) register(c *gin.Context) {
	var request auth.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "RegisterRequest", Err: err})
		return
	}

	response, err := h.authenticator.Register(c.Request.Context(), request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// login godoc
// @Summary Login using email and password
// @Tags Auth
// @Accept json
// @Produce  json
// @Param payload body auth.LoginRequest true "Login Payload"
// @Success 201 {object} auth.CredentialsResponse
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /login [post]
func (h *AuthenticatorHandler) login(c *gin.Context) {
	var request auth.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "LoginRequest", Err: err})
		return
	}

	response, err := h.authenticator.Login(c.Request.Context(), request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// resetPassword godoc
// @Summary Reset the password for the given email
// @Tags Auth
// @Accept json
// @Produce  json
// @Param payload body auth.PasswordResetRequest true "Register Payload"
// @Success 204 "success"
// @Failure 400 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /password-reset [post]
func (h *AuthenticatorHandler) resetPassword(c *gin.Context) {
	var request auth.PasswordResetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "PasswordResetRequest", Err: err})
		return
	}

	if err := h.authenticator.ResetPassword(c.Request.Context(), request); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
