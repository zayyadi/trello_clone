package handlers

import (
	"time"

	"github.com/zayyadi/trello/models"
)

// Auth DTOs
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Board DTOs
type CreateBoardRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=255"`
}

type UpdateBoardRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}

type BoardResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	OwnerID     uint                  `json:"ownerID"`
	Owner       UserResponse          `json:"owner,omitempty"` // Include owner details
	Lists       []ListResponse        `json:"lists,omitempty"` // Optionally include lists
	Members     []BoardMemberResponse `json:"members,omitempty"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}

// Board Member DTOs
type AddMemberRequest struct {
	Email  *string `json:"email"`  // Add by email
	UserID *uint   `json:"userID"` // Or by UserID
}

type BoardMemberResponse struct {
	BoardID   uint         `json:"boardID"`
	UserID    uint         `json:"userID"`
	User      UserResponse `json:"user"` // Include member's user details
	CreatedAt time.Time    `json:"createdAt"`
}

// List DTOs
type CreateListRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Position *uint  `json:"position"` // Optional: desired position
}

type UpdateListRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=1,max=100"`
	Position *uint   `json:"position"` // Optional: new position
}

type ListResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	BoardID   uint           `json:"boardID"`
	Position  uint           `json:"position"`
	Cards     []CardResponse `json:"cards,omitempty"` // Optionally include cards
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// Card DTOs
type CreateCardRequest struct {
	Title          string     `json:"title" binding:"required,min=1,max=255"`
	Description    string     `json:"description" binding:"max=1000"`
	Position       *uint      `json:"position"`
	DueDate        *time.Time `json:"dueDate,omitempty"`
	AssignedUserID *uint      `json:"assignedUserID,omitempty"`
	SupervisorID   *uint      `json:"supervisorID,omitempty"` // New field
	Color          *string    `json:"color,omitempty"`        // New field
}

type UpdateCardRequest struct {
	Title          *string            `json:"title" binding:"omitempty,min=1,max=255"`
	Description    *string            `json:"description" binding:"omitempty,max=1000"`
	Position       *uint              `json:"position"`
	DueDate        *time.Time         `json:"dueDate,omitempty"` // Can be null to clear
	AssignedUserID **uint             `json:"assignedUserID,omitempty"`
	SupervisorID   **uint             `json:"supervisorID,omitempty"` // New field, ptr to ptr
	Status         *models.CardStatus `json:"status,omitempty"`       // New field
	Color          *string            `json:"color,omitempty"`        // New field
}

type CardResponse struct {
	ID             uint              `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	ListID         uint              `json:"listID"`
	Position       uint              `json:"position"`
	DueDate        *time.Time        `json:"dueDate,omitempty"`
	Status         models.CardStatus `json:"status"` // New field
	AssignedUserID *uint             `json:"assignedUserID,omitempty"`
	AssignedUser   *UserResponse     `json:"assignedUser,omitempty"`
	SupervisorID   *uint             `json:"supervisorID,omitempty"` // New field
	Supervisor     *UserResponse     `json:"supervisor,omitempty"`   // New field
	Color          *string           `json:"color,omitempty"`        // New field
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type MoveCardRequest struct {
	TargetListID uint `json:"targetListID" binding:"required"`
	NewPosition  uint `json:"newPosition" binding:"required,min=1"` // Min 1 because position is 1-based
}
