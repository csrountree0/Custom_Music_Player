package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"musicapp/backend/internal/services"
	"musicapp/backend/pkg/apierror"
)

type AuthHandler struct {
	authSvc *services.AuthService
}

func NewAuthHandler(authSvc *services.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	resp, err := h.authSvc.Register(input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	resp, err := h.authSvc.Login(input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.New(400, "VALIDATION_ERROR", "refresh_token is required").Respond(c)
		return
	}

	resp, err := h.authSvc.RefreshTokens(body.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	if err := h.authSvc.Logout(userID); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	user, err := h.authSvc.GetMe(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input services.UpdateMeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apierror.New(400, "VALIDATION_ERROR", err.Error()).Respond(c)
		return
	}

	user, err := h.authSvc.UpdateMe(userID, input)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
