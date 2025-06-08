package services

import (
	"log"

	"github.com/zayyadi/trello/handlers" // For DTOs as payload
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
	// Send the structured message to the hub's broadcast channel
	select {
	case hub.Broadcast <- msg:
		log.Printf("Message of type %s for board %d queued for broadcast.", messageType, boardID)
	default:
		// This case should ideally not be hit if the broadcast channel is buffered
		// or if the hub is processing messages quickly enough.
		// If hit, it means the hub's broadcast channel is full, which is a bottleneck.
		log.Printf("Error: Hub broadcast channel full. Failed to queue message type %s for board %d.", messageType, boardID)
	}
}

// --- Payload Mapping Helpers (moved from individual services for reusability if needed) ---
// These functions can remain in their respective handler (or DTO) files if preferred,
// but are placed here if services need to construct DTOs for WebSocket payloads
// without directly calling handler mapping functions in some scenarios.
// For now, services will mostly use models or existing DTOs.

// Example: MapBoardToBoardResponse is in handlers/dto.go.
// If a specific WebSocket payload is needed that's different from a handler DTO,
// it can be defined and mapped here.

// MapBoardToPayload converts a models.Board to handlers.BoardResponse for WebSocket payload
func MapBoardToPayload(board *models.Board) handlers.BoardResponse {
	return handlers.MapBoardToBoardResponse(board) // Reuse existing mapper
}

// MapListToPayload converts a models.List to handlers.ListResponse
func MapListToPayload(list *models.List) handlers.ListResponse {
	return handlers.MapListToListResponse(list) // Reuse existing mapper
}

// MapCardToPayload converts a models.Card to handlers.CardResponse
func MapCardToPayload(card *models.Card) handlers.CardResponse {
	return handlers.MapCardToCardResponse(card) // Reuse existing mapper
}
