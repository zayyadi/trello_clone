package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/services"
)

// CommentHandler handles HTTP requests related to comments.
type CommentHandler struct {
	commentService services.CommentServiceInterface
}

// NewCommentHandler creates a new CommentHandler.
func NewCommentHandler(commentService services.CommentServiceInterface) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateCommentRequest defines the request body for creating a comment.
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// CommentResponse defines the structure for comment responses.
type CommentResponse struct {
	ID        uint         `json:"id"`
	Content   string       `json:"content"`
	CardID    uint         `json:"cardID"`
	UserID    uint         `json:"userID"`
	User      UserResponse `json:"user"` // Reusing UserResponse from dto.go
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

// MapCommentToResponse maps a models.Comment to CommentResponse.
func MapCommentToResponse(comment *models.Comment) CommentResponse {
	userResp := UserResponse{} // Default empty if user not preloaded
	if comment.User.ID != 0 {  // Check if User struct is populated (not just zero values)
		userResp = MapUserToResponse(&comment.User)
	}
	return CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CardID:    comment.CardID,
		UserID:    comment.UserID,
		User:      userResp,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

// MapCommentsToResponse maps a slice of models.Comment to a slice of CommentResponse.
func MapCommentsToResponse(comments []models.Comment) []CommentResponse {
	commentResponses := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = MapCommentToResponse(&comment)
	}
	return commentResponses
}

// CreateComment handles POST /cards/:cardID/comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	currentUserID, ok := userID.(uint)
	if !ok {
		RespondWithError(c, http.StatusInternalServerError, "User ID is of invalid type")
		return
	}

	cardIDStr := c.Param("cardID")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid card ID format")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	comment, err := h.commentService.CreateComment(uint(cardID), currentUserID, req.Content)
	if err != nil {
		// Error mapping based on service errors
		if errors.Is(err, services.ErrCardNotFound) || errors.Is(err, services.ErrListNotFound) || errors.Is(err, services.ErrBoardNotFound) {
			RespondWithError(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, services.ErrForbidden) {
			RespondWithError(c, http.StatusForbidden, err.Error())
		} else if err.Error() == "comment content cannot be empty" { // Specific error from service
			RespondWithError(c, http.StatusBadRequest, err.Error())
		} else {
			RespondWithError(c, http.StatusInternalServerError, "Failed to create comment: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, MapCommentToResponse(comment))
}

// GetCommentsByCardID handles GET /cards/:cardID/comments
func (h *CommentHandler) GetCommentsByCardID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	currentUserID, ok := userID.(uint)
	if !ok {
		RespondWithError(c, http.StatusInternalServerError, "User ID is of invalid type")
		return
	}

	cardIDStr := c.Param("cardID")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid card ID format")
		return
	}

	comments, err := h.commentService.GetCommentsByCardID(uint(cardID), currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrCardNotFound) || errors.Is(err, services.ErrListNotFound) || errors.Is(err, services.ErrBoardNotFound) {
			RespondWithError(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, services.ErrForbidden) {
			RespondWithError(c, http.StatusForbidden, err.Error())
		} else {
			RespondWithError(c, http.StatusInternalServerError, "Failed to get comments: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, MapCommentsToResponse(comments))
}
