package handlers

import (
	"net/http"
	"strconv"

	"github.com/zayyadi/trello/dto" // Import new dto package
	"github.com/zayyadi/trello/services"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cardService services.CardServiceInterface // Use interface
}

func NewCardHandler(cardService services.CardServiceInterface) *CardHandler { // Use interface
	return &CardHandler{cardService: cardService}
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	listIDStr := c.Param("listID")
	listID, err := strconv.ParseUint(listIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid list ID")
		return
	}

	var req dto.CreateCardRequest                  // Use dto type
	if err := c.ShouldBindJSON(&req); err != nil { /* ... error handling ... */
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error()) // Ensure proper error response
		return
	}

	card, err := h.cardService.CreateCard(
		uint(listID),
		req.Title,
		req.Description,
		req.Position,
		req.DueDate,
		req.AssignedUserID,
		req.SupervisorID, // Pass new field
		req.Color,        // Pass new field
		userID.(uint),
	)
	if err != nil { // Error handling for service call
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusCreated, "Card created successfully", dto.MapCardToResponse(card, true)) // Use dto mapper
}

func (h *CardHandler) GetCardByID(c *gin.Context) {
	userID, _ := c.Get("userID")
	cardIDStr := c.Param("cardID")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid card ID")
		return
	}

	card, err := h.cardService.GetCardByID(uint(cardID), userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Card retrieved successfully", dto.MapCardToResponse(card, true)) // Use dto mapper
}

func (h *CardHandler) GetCardsByListID(c *gin.Context) {
	userID, _ := c.Get("userID")
	listIDStr := c.Param("listID")
	listID, err := strconv.ParseUint(listIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid list ID")
		return
	}

	cards, err := h.cardService.GetCardsByListID(uint(listID), userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	var cardResponses []dto.CardResponse // Use dto type
	for _, card := range cards {
		cardResponses = append(cardResponses, dto.MapCardToResponse(&card, true)) // Use dto mapper
	}
	RespondWithSuccess(c, http.StatusOK, "Cards retrieved successfully", cardResponses)
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	cardIDStr := c.Param("cardID")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid card ID")
		return
	}

	var req dto.UpdateCardRequest                  // Use dto type
	if err := c.ShouldBindJSON(&req); err != nil { /* ... error handling ... */
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error()) // Ensure proper error response
		return
	}

	card, err := h.cardService.UpdateCard(
		uint(cardID),
		req.Title,
		req.Description,
		req.Position,
		req.DueDate,
		req.AssignedUserID,
		req.SupervisorID, // Pass new field
		req.Status,       // Pass new field
		req.Color,        // Pass new field
		userID.(uint),
	)
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Card updated successfully", dto.MapCardToResponse(card, true)) // Use dto mapper
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	cardIDStr := c.Param("cardID")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid card ID")
		return
	}

	err = h.cardService.DeleteCard(uint(cardID), userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Card deleted successfully", nil)
}

func (h *CardHandler) MoveCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	cardIDStr := c.Param("cardID")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid card ID")
		return
	}

	var req dto.MoveCardRequest // Use dto type
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	card, err := h.cardService.MoveCard(uint(cardID), req.TargetListID, req.NewPosition, userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Card moved successfully", dto.MapCardToResponse(card, true)) // Use dto mapper
}

// MapCardToResponse function is now in dto/card_dto.go
// MapUserToResponse is in dto/auth_dto.go

func (h *CardHandler) AddCollaborator(c *gin.Context) {
	currentUserID, _ := c.Get("userID")
	uCurrentUserID, ok := currentUserID.(uint)
	if !ok {
		RespondWithError(c, http.StatusInternalServerError, "Invalid user ID type in token")
		return
	}

	cardIDStr := c.Param("cardID")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid card ID format")
		return
	}

	var req dto.CardAddCollaboratorRequest // Use dto type
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if (req.Email == nil || *req.Email == "") && (req.UserID == nil || *req.UserID == 0) {
		RespondWithError(c, http.StatusBadRequest, "Either email or userID must be provided to add a collaborator")
		return
	}
	if req.Email != nil && *req.Email != "" && req.UserID != nil && *req.UserID != 0 {
		RespondWithError(c, http.StatusBadRequest, "Provide either email or userID, not both")
		return
	}

	var targetUserEmail string
	if req.Email != nil {
		targetUserEmail = *req.Email
	}

	addedUser, err := h.cardService.AddCollaboratorToCard(uint(cardID), uCurrentUserID, targetUserEmail, req.UserID)
	if err != nil {
		HandleServiceError(c, err) // HandleServiceError should map service errors to HTTP errors
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Collaborator added successfully", dto.MapUserToResponse(addedUser)) // Use dto mapper
}

func (h *CardHandler) RemoveCollaborator(c *gin.Context) {
	currentUserID, _ := c.Get("userID")
	uCurrentUserID, ok := currentUserID.(uint)
	if !ok {
		RespondWithError(c, http.StatusInternalServerError, "Invalid user ID type in token")
		return
	}

	cardIDStr := c.Param("cardID")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid card ID format")
		return
	}

	targetUserIDStr := c.Param("userID") // This is the UserID of the collaborator to remove
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid target user ID format")
		return
	}

	err = h.cardService.RemoveCollaboratorFromCard(uint(cardID), uCurrentUserID, uint(targetUserID))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Collaborator removed successfully", nil)
}

func (h *CardHandler) GetCollaborators(c *gin.Context) {
	currentUserID, _ := c.Get("userID")
	uCurrentUserID, ok := currentUserID.(uint)
	if !ok {
		RespondWithError(c, http.StatusInternalServerError, "Invalid user ID type in token")
		return
	}

	cardIDStr := c.Param("cardID")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid card ID format")
		return
	}

	users, err := h.cardService.GetCardCollaborators(uint(cardID), uCurrentUserID)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	var userResponses []dto.UserResponse // Use dto type
	for i, user := range users {
		userResponses[i] = dto.MapUserToResponse(&user) // Use dto mapper
	}
	RespondWithSuccess(c, http.StatusOK, "Collaborators retrieved successfully", userResponses)
}
