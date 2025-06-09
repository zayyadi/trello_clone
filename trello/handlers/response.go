package handlers

import (
	"errors"
	"log" // Import log package
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
	path := c.Request.URL.Path
	method := c.Request.Method

	switch {
	case errors.Is(err, services.ErrInvalidCredentials):
		log.Printf("INFO [ServiceError]: InvalidCredentials: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusUnauthorized, "Invalid email or password")
	case errors.Is(err, services.ErrEmailExists):
		log.Printf("INFO [ServiceError]: EmailExists: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusConflict, "Email already exists")
	case errors.Is(err, services.ErrUsernameExists):
		log.Printf("INFO [ServiceError]: UsernameExists: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusConflict, "Username already exists")
	case errors.Is(err, services.ErrUserNotFound):
		log.Printf("INFO [ServiceError]: UserNotFound: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusNotFound, "User not found")
	case errors.Is(err, services.ErrBoardNotFound):
		log.Printf("INFO [ServiceError]: BoardNotFound: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusNotFound, "Board not found")
	case errors.Is(err, services.ErrListNotFound):
		log.Printf("INFO [ServiceError]: ListNotFound: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusNotFound, "List not found")
	case errors.Is(err, services.ErrCardNotFound):
		log.Printf("INFO [ServiceError]: CardNotFound: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusNotFound, "Card not found")
	case errors.Is(err, services.ErrUnauthorized):
		log.Printf("INFO [ServiceError]: Unauthorized: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized access")
	case errors.Is(err, services.ErrForbidden):
		log.Printf("INFO [ServiceError]: Forbidden: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusForbidden, "Forbidden: You don't have permission to perform this action.")
	case errors.Is(err, services.ErrUserAlreadyMember):
		log.Printf("INFO [ServiceError]: UserAlreadyMember: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusConflict, "User is already a member of this board.")
	case errors.Is(err, services.ErrBoardMemberNotFound):
		log.Printf("INFO [ServiceError]: BoardMemberNotFound: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusNotFound, "Board member not found.")
	case errors.Is(err, services.ErrCannotRemoveOwner):
		log.Printf("INFO [ServiceError]: CannotRemoveOwner: %v (Request: %s %s)", err, method, path)
		RespondWithError(c, http.StatusBadRequest, "Board owner cannot be removed this way.")
	case errors.Is(err, services.ErrInvalidInput):
		log.Printf("WARN [ServiceError]: InvalidInput: %v (Request: %s %s)", err, method, path) // Log as WARN
		RespondWithError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrPositionOutOfBound):
		log.Printf("WARN [ServiceError]: PositionOutOfBound: %v (Request: %s %s)", err, method, path) // Log as WARN
		RespondWithError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrSameListMove):
		log.Printf("WARN [ServiceError]: SameListMove: %v (Request: %s %s)", err, method, path) // Log as WARN
		RespondWithError(c, http.StatusBadRequest, err.Error())
	default:
		log.Printf("ERROR [ServiceError]: InternalServerError: %v (Request: %s %s)", err, method, path)
		_ = c.Error(err) // Gin's internal error logging, keep it if it's useful for other middlewares
		RespondWithError(c, http.StatusInternalServerError, "An unexpected error occurred.")
	}
}
