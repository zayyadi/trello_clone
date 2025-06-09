package dto

import (
	"time"

	"github.com/zayyadi/trello/models"
	// "github.com/zayyadi/trello/dto" // This would be a circular dependency if ListResponse is here.
	// ListResponse will be in its own file. We assume it's available in the package 'dto'.
)

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
	Owner       UserResponse          `json:"owner,omitempty"` // Uses dto.UserResponse
	Lists       []ListResponse        `json:"lists,omitempty"` // Uses dto.ListResponse (to be created)
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
	User      UserResponse `json:"user"` // Uses dto.UserResponse
	CreatedAt time.Time    `json:"createdAt"`
}

// MapBoardToResponse maps model.Board to BoardResponse
// includeOwner and includeMembersAndLists flags control depth of mapping.
func MapBoardToResponse(board *models.Board, includeOwner, includeMembersAndLists bool) BoardResponse {
	if board == nil {
		return BoardResponse{}
	}
	resp := BoardResponse{
		ID:          board.Model.ID,
		Name:        board.Name,
		Description: board.Description,
		OwnerID:     board.OwnerID,
		CreatedAt:   board.Model.CreatedAt,
		UpdatedAt:   board.Model.UpdatedAt,
	}
	if includeOwner && board.Owner.ID != 0 {
		resp.Owner = MapUserToResponse(&board.Owner) // Assumes MapUserToResponse is in the same 'dto' package
	}
	if includeMembersAndLists {
		if len(board.Members) > 0 {
			resp.Members = []BoardMemberResponse{}
			for _, m := range board.Members {
				resp.Members = append(resp.Members, MapBoardMemberToResponse(&m)) // Assumes MapBoardMemberToResponse is in this file
			}
		}
		if len(board.Lists) > 0 {
			resp.Lists = []ListResponse{}
			for _, l := range board.Lists {
				// For lists within board response, don't include cards by default to avoid huge payloads
				// Assumes MapListToResponse is in the same 'dto' package and handles *models.List to dto.ListResponse
				resp.Lists = append(resp.Lists, MapListToResponse(&l, false))
			}
		}
	}
	return resp
}

// MapBoardMemberToResponse maps model.BoardMember to BoardMemberResponse
func MapBoardMemberToResponse(member *models.BoardMember) BoardMemberResponse {
	if member == nil {
		return BoardMemberResponse{}
	}
	return BoardMemberResponse{
		BoardID:   member.BoardID,
		UserID:    member.UserID,
		User:      MapUserToResponse(&member.User), // Assumes MapUserToResponse is in the same 'dto' package
		CreatedAt: member.CreatedAt,
	}
}
