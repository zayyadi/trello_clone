package dto

import (
	"time"

	"github.com/zayyadi/trello/models" // Ensure this import is correct
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

// MapUserToResponse maps a models.User to a UserResponse DTO.
// It's moved here from auth_handler.go.
// Note: This function will need to import "github.com/zayyadi/trello/models"
// when used in a context where models.User is passed.
// For now, defining its signature and structure.
// The actual models.User import will be handled when updating auth_handler.go
// or if other DTOs need it directly.
// For now, to make this file self-contained for creation, I will assume
// the MapUserToResponse will be adjusted later if it needs to take models.User.
// Let's define it to take what it needs for now, assuming the model is available.
// No, I must ensure this file is complete. I will add the models import.

// import "github.com/zayyadi/trello/models"
// This comment is just for planning. The actual import will be added by the tool if I use models.User.
// For now, I will write it as if models.User is available in this scope.
// This is tricky. The tool will create the file as is.
// I should ensure the mapping function is complete here.

// To make this file complete, I need to define what User it maps from.
// It maps from models.User. So, I need that import.

// Corrected approach: Define it with the models.User import.
// The tool should handle creating this file with the import.
// However, I cannot use the tool to add imports to an existing file easily
// without replace_with_git_merge_diff.
// So, when I create the file, I MUST include all necessary imports.

/*
import (
	"time"
	"github.com/zayyadi/trello/models" // This is what I need
)
// ... then the MapUserToResponse would be:
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
*/
// The above is the plan. Let's create the file with this structure.
// The tool will take the content below.

// // Auth DTOs
// type RegisterRequest struct {
// 	Username string `json:"username" binding:"required,min=3,max=50"`
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required,min=6,max=100"`
// }

// type LoginRequest struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required"`
// }

// type AuthResponse struct {
// 	User  UserResponse `json:"user"`
// 	Token string       `json:"token"`
// }

// // UserResponse is the DTO for user information.
// type UserResponse struct {
// 	ID        uint      `json:"id"`
// 	Username  string    `json:"username"`
// 	Email     string    `json:"email"`
// 	CreatedAt time.Time `json:"createdAt"`
// 	UpdatedAt time.Time `json:"updatedAt"`
// }

// MapUserToResponse maps a models.User to a UserResponse DTO.
func MapUserToResponse(user *models.User) UserResponse {
	if user == nil {
		return UserResponse{}
	}
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt, // Assumes models.User has CreatedAt directly (it's in models.BaseModel)
		UpdatedAt: user.UpdatedAt, // Assumes models.User has UpdatedAt directly
	}
}
