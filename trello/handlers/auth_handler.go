package handlers

import (
	"net/http"

	"github.com/zayyadi/trello/dto" // Import new dto package
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest // Use dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	user, token, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	userResp := dto.MapUserToResponse(user) // Use dto.MapUserToResponse
	RespondWithSuccess(c, http.StatusCreated, "User registered successfully", dto.AuthResponse{User: userResp, Token: token}) // Use dto.AuthResponse
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest // Use dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	userResp := dto.MapUserToResponse(user) // Use dto.MapUserToResponse
	RespondWithSuccess(c, http.StatusOK, "Login successful", dto.AuthResponse{User: userResp, Token: token}) // Use dto.AuthResponse
}

// MapUserToResponse function is now in dto/auth_dto.go and will be removed from here.
