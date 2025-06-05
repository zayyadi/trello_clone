package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zayyadi/trello/services"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"` // Optional additional error detail
}

func RespondWithSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

// HandleServiceError translates service errors to HTTP responses
func HandleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidCredentials):
		RespondWithError(c, http.StatusUnauthorized, "Invalid email or password")
	case errors.Is(err, services.ErrEmailExists):
		RespondWithError(c, http.StatusConflict, "Email already exists")
	case errors.Is(err, services.ErrUsernameExists):
		RespondWithError(c, http.StatusConflict, "Username already exists")
	case errors.Is(err, services.ErrUserNotFound):
		RespondWithError(c, http.StatusNotFound, "User not found")
	case errors.Is(err, services.ErrBoardNotFound):
		RespondWithError(c, http.StatusNotFound, "Board not found")
	case errors.Is(err, services.ErrListNotFound):
		RespondWithError(c, http.StatusNotFound, "List not found")
	case errors.Is(err, services.ErrCardNotFound):
		RespondWithError(c, http.StatusNotFound, "Card not found")
	case errors.Is(err, services.ErrUnauthorized):
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized access")
	case errors.Is(err, services.ErrForbidden):
		RespondWithError(c, http.StatusForbidden, "Forbidden: You don't have permission to perform this action.")
	case errors.Is(err, services.ErrUserAlreadyMember):
		RespondWithError(c, http.StatusConflict, "User is already a member of this board.")
	case errors.Is(err, services.ErrBoardMemberNotFound):
		RespondWithError(c, http.StatusNotFound, "Board member not found.")
	case errors.Is(err, services.ErrCannotRemoveOwner):
		RespondWithError(c, http.StatusBadRequest, "Board owner cannot be removed this way.")
	case errors.Is(err, services.ErrInvalidInput):
		RespondWithError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrPositionOutOfBound):
		RespondWithError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrSameListMove):
		RespondWithError(c, http.StatusBadRequest, err.Error())
	default:
		// Log the error for server-side inspection
		c.Error(err) // Gin's internal error logging
		RespondWithError(c, http.StatusInternalServerError, "An unexpected error occurred.")
	}
}
