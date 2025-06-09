package dto

import (
	"time"

	"github.com/zayyadi/trello/models"
	// "github.com/zayyadi/trello/dto" // For UserResponse, assumed in same package
)

// Card DTOs
type CreateCardRequest struct {
	Title          string     `json:"title" binding:"required,min=1,max=255"`
	Description    string     `json:"description" binding:"max=1000"`
	Position       *uint      `json:"position"`
	DueDate        *time.Time `json:"dueDate,omitempty"`
	AssignedUserID *uint      `json:"assignedUserID,omitempty"`
	SupervisorID   *uint      `json:"supervisorID,omitempty"`
	Color          *string    `json:"color,omitempty"`
}

type UpdateCardRequest struct {
	Title          *string            `json:"title" binding:"omitempty,min=1,max=255"`
	Description    *string            `json:"description" binding:"omitempty,max=1000"`
	Position       *uint              `json:"position"`
	DueDate        *time.Time         `json:"dueDate,omitempty"`
	AssignedUserID **uint             `json:"assignedUserID,omitempty"`
	SupervisorID   **uint             `json:"supervisorID,omitempty"`
	Status         *models.CardStatus `json:"status,omitempty"`
	Color          *string            `json:"color,omitempty"`
}

type CardResponse struct {
	ID             uint              `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	ListID         uint              `json:"listID"`
	Position       uint              `json:"position"`
	DueDate        *time.Time        `json:"dueDate,omitempty"`
	Status         models.CardStatus `json:"status"`
	AssignedUserID *uint             `json:"assignedUserID,omitempty"`
	AssignedUser   *UserResponse     `json:"assignedUser,omitempty"` // Uses dto.UserResponse
	SupervisorID   *uint             `json:"supervisorID,omitempty"`
	Supervisor     *UserResponse     `json:"supervisor,omitempty"`   // Uses dto.UserResponse
	Color          *string           `json:"color,omitempty"`
	Collaborators  []UserResponse    `json:"collaborators,omitempty"` // Uses dto.UserResponse
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type MoveCardRequest struct {
	TargetListID uint `json:"targetListID" binding:"required"`
	NewPosition  uint `json:"newPosition" binding:"required,min=1"`
}

// Card Collaborator DTOs
type CardAddCollaboratorRequest struct {
	Email  *string `json:"email"`
	UserID *uint   `json:"userID"`
}

// MapCardToResponse maps model.Card to CardResponse
func MapCardToResponse(card *models.Card, includeUserDetails bool) CardResponse {
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
		Status:         card.Status,
		AssignedUserID: card.AssignedUserID,
		SupervisorID:   card.SupervisorID,
		Color:          card.Color,
		CreatedAt:      card.CreatedAt,
		UpdatedAt:      card.UpdatedAt,
	}

	if includeUserDetails {
		if card.AssignedUser != nil && card.AssignedUser.ID != 0 {
			resp.AssignedUser = &UserResponse{}             // Create a new UserResponse pointer
			*resp.AssignedUser = MapUserToResponse(card.AssignedUser) // Map and assign
		}
		if card.Supervisor != nil && card.Supervisor.ID != 0 {
			resp.Supervisor = &UserResponse{}               // Create a new UserResponse pointer
			*resp.Supervisor = MapUserToResponse(card.Supervisor)   // Map and assign
		}
	}

	if card.Collaborators != nil {
		resp.Collaborators = make([]UserResponse, len(card.Collaborators))
		for i, collaborator := range card.Collaborators {
			resp.Collaborators[i] = MapUserToResponse(collaborator) // Assumes MapUserToResponse is in same 'dto' package
		}
	} else {
		resp.Collaborators = []UserResponse{} // Ensure empty slice instead of null
	}

	return resp
}
