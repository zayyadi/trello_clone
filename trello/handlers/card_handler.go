package handlers

import (
	"net/http"
	"strconv"

	"github.com/zayyadi/trello/models"
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

	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil { /* ... error handling ... */
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
	RespondWithSuccess(c, http.StatusCreated, "Card created successfully", MapCardToResponse(card, true))
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
	RespondWithSuccess(c, http.StatusOK, "Card retrieved successfully", MapCardToResponse(card, true))
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
	var cardResponses []CardResponse
	for _, card := range cards {
		cardResponses = append(cardResponses, MapCardToResponse(&card, true))
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

	var req UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil { /* ... error handling ... */
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
	RespondWithSuccess(c, http.StatusOK, "Card updated successfully", MapCardToResponse(card, true))
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

	var req MoveCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	card, err := h.cardService.MoveCard(uint(cardID), req.TargetListID, req.NewPosition, userID.(uint))
	if err != nil {
		HandleServiceError(c, err)
		return
	}
	RespondWithSuccess(c, http.StatusOK, "Card moved successfully", MapCardToResponse(card, true))
}

// Helper to map model.Card to CardResponse
func MapCardToResponse(card *models.Card, includeAssignedUser bool) CardResponse {
	if card == nil {
		return CardResponse{}
	}
	resp := CardResponse{
		ID:             card.ID,
		Title:          card.Title,
		Description:    card.Description,
		ListID:         card.ListID,
		Position:       card.Position,
		DueDate:        card.DueDate,
		Status:         card.Status, // New field
		AssignedUserID: card.AssignedUserID,
		SupervisorID:   card.SupervisorID, // New field
		Color:          card.Color,        // New field
		CreatedAt:      card.CreatedAt,
		UpdatedAt:      card.UpdatedAt,
	}
	if includeAssignedUser { // Re-purpose this flag or add a new one for supervisor
		if card.AssignedUser != nil && card.AssignedUser.ID != 0 {
			resp.AssignedUser = &UserResponse{
				ID: card.AssignedUser.ID, Username: card.AssignedUser.Username, Email: card.AssignedUser.Email,
			}
		}
		if card.Supervisor != nil && card.Supervisor.ID != 0 { // New field
			resp.Supervisor = &UserResponse{
				ID: card.Supervisor.ID, Username: card.Supervisor.Username, Email: card.Supervisor.Email,
			}
		}
	}
	return resp
}

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

	var req CardAddCollaboratorRequest
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
	RespondWithSuccess(c, http.StatusOK, "Collaborator added successfully", MapUserToResponse(addedUser))
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

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = MapUserToResponse(&user)
	}
	RespondWithSuccess(c, http.StatusOK, "Collaborators retrieved successfully", userResponses)
}
