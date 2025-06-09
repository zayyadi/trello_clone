package services

import (
	"log"

	"github.com/zayyadi/trello/dto" // Use new dto package for mappers
	"github.com/zayyadi/trello/models" // Keep models import
	"github.com/zayyadi/trello/realtime"
)

// broadcastMessage constructs a WebSocketMessage and sends it to the hub.
// userID is optional and can be used to prevent sending messages back to the originating user.
func broadcastMessage(
	hub *realtime.Hub,
	boardID uint,
	messageType string,
	payload interface{},
	originatingUserID ...uint, // Optional: for excluding the sender
) {
	if hub == nil {
		log.Printf("Warning: Hub is nil. Cannot broadcast message type %s for board %d.", messageType, boardID)
		return
	}

	msg := &realtime.WebSocketMessage{
		BoardID: boardID,
		Type:    messageType,
		Payload: payload,
	}

	if len(originatingUserID) > 0 {
		msg.UserID = originatingUserID[0]
	}

	// The hub will handle JSON marshalling of the client-facing parts (Type & Payload)
	// Send the structured message to the hub via its exported method
	hub.Submit(msg)
	log.Printf("Message of type %s for board %d submitted to hub.", messageType, boardID)
}

// --- Payload Mapping Helpers (moved from individual services for reusability if needed) ---
// These functions can remain in their respective handler (or DTO) files if preferred,
// but are placed here if services need to construct DTOs for WebSocket payloads
// without directly calling handler mapping functions in some scenarios.
// For now, services will mostly use models or existing DTOs.

// Example: MapBoardToBoardResponse is now in dto/board_dto.go.
// If a specific WebSocket payload is needed that's different from a handler DTO,
// it can be defined and mapped here.

// MapBoardToPayload converts a models.Board to dto.BoardResponse for WebSocket payload
func MapBoardToPayload(board *models.Board) dto.BoardResponse {
	// The boolean flags for includeOwner and includeMembersAndLists might need to be decided here.
	// For WebSocket, usually, we want to send comprehensive data.
	return dto.MapBoardToResponse(board, true, true)
}

// MapListToPayload converts a models.List to dto.ListResponse
func MapListToPayload(list *models.List) dto.ListResponse {
	// Similar to board, decide if cards should be included.
	return dto.MapListToResponse(list, true)
}

// MapCardToPayload converts a models.Card to dto.CardResponse
func MapCardToPayload(card *models.Card) dto.CardResponse {
	// Decide if user details should be included.
	return dto.MapCardToResponse(card, true)
}
