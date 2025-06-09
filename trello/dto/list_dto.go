package dto

import (
	"time"

	"github.com/zayyadi/trello/models"
	// "github.com/zayyadi/trello/dto" // For CardResponse, assumed in same package
)

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
	Cards     []CardResponse `json:"cards,omitempty"` // Uses dto.CardResponse (to be created)
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// MapListToResponse maps model.List to ListResponse
// includeCards flag controls depth of mapping.
func MapListToResponse(list *models.List, includeCards bool) ListResponse {
	if list == nil {
		return ListResponse{}
	}
	resp := ListResponse{
		ID:        list.ID,
		Name:      list.Name,
		BoardID:   list.BoardID,
		Position:  list.Position,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
	if includeCards && len(list.Cards) > 0 {
		resp.Cards = []CardResponse{}
		for _, c := range list.Cards {
			// Assumes MapCardToResponse is in the same 'dto' package and handles *models.Card to dto.CardResponse
			resp.Cards = append(resp.Cards, MapCardToResponse(&c, true)) // Include assigned user for cards by default here
		}
	}
	return resp
}
