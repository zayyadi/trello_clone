package handlers

import (
	"net/http"

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
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	user, token, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	userResp := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	RespondWithSuccess(c, http.StatusCreated, "User registered successfully", AuthResponse{User: userResp, Token: token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	userResp := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	RespondWithSuccess(c, http.StatusOK, "Login successful", AuthResponse{User: userResp, Token: token})
}

// Helper to map model.User to UserResponse
func MapUserToResponse(user *models.User) UserResponse {
	if user == nil {
		return UserResponse{}
	}
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
