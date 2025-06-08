package realtime

// WebSocketMessage represents a message sent over WebSocket.
type WebSocketMessage struct {
	Type    string      `json:"type"`              // Type of the message (e.g., "BOARD_UPDATED")
	Payload interface{} `json:"payload,omitempty"` // Actual data payload
	BoardID uint        `json:"-"`                 // Target BoardID, used by Hub for routing, not sent to client
	UserID  uint        `json:"-"`                 // Optional: Originating UserID, not sent to client, for potential exclusion
}

// Message Types Constants
const (
	MessageTypeBoardCreated      = "BOARD_CREATED"
	MessageTypeBoardUpdated      = "BOARD_UPDATED"
	MessageTypeBoardDeleted      = "BOARD_DELETED"
	MessageTypeBoardMemberAdded  = "BOARD_MEMBER_ADDED"
	MessageTypeBoardMemberRemoved = "BOARD_MEMBER_REMOVED"

	MessageTypeListCreated = "LIST_CREATED"
	MessageTypeListUpdated = "LIST_UPDATED"
	MessageTypeListDeleted = "LIST_DELETED"
	// MessageTypeListMoved   = "LIST_MOVED" // If list reordering is implemented as a distinct event

	MessageTypeCardCreated             = "CARD_CREATED"
	MessageTypeCardUpdated             = "CARD_UPDATED"
	MessageTypeCardDeleted             = "CARD_DELETED"
	MessageTypeCardMoved               = "CARD_MOVED"
	MessageTypeCardAssigned            = "CARD_ASSIGNED" // When a user is assigned
	MessageTypeCardUnassigned          = "CARD_UNASSIGNED"
	MessageTypeCardCollaboratorAdded   = "CARD_COLLABORATOR_ADDED"
	MessageTypeCardCollaboratorRemoved = "CARD_COLLABORATOR_REMOVED"
	// Add more as needed, e.g., CARD_COMMENT_ADDED
)

// Example Payloads (can also use DTOs from handlers package directly if suitable)

// BoardBasicInfo might be used for delete operations or simple updates
type BoardBasicInfo struct {
	ID uint `json:"id"`
}

// ListBasicInfo for delete operations
type ListBasicInfo struct {
	ID      uint `json:"id"`
	BoardID uint `json:"boardId"` // Include BoardID for client-side context
}

// CardBasicInfo for delete operations
type CardBasicInfo struct {
	ID      uint `json:"id"`
	ListID  uint `json:"listId"`
	BoardID uint `json:"boardId"` // Include BoardID for client-side context
}

// CardMovedPayload details the specifics of a card move operation
type CardMovedPayload struct {
	CardID       uint `json:"cardId"`
	OldListID    uint `json:"oldListId"`
	NewListID    uint `json:"newListId"`
	OldPosition  uint `json:"oldPosition"`
	NewPosition  uint `json:"newPosition"`
	BoardID      uint `json:"boardId"`      // For client-side context if card is moved between boards (not current model)
	UpdatedCards []struct { // Optional: if positions of other cards in affected lists are sent
		ID       uint `json:"id"`
		Position uint `json:"position"`
		ListID   uint `json:"listId"`
	} `json:"updatedCards,omitempty"`
}

// BoardMemberPayload for member changes
type BoardMemberPayload struct {
	BoardID  uint `json:"boardId"`
	UserID   uint `json:"userId"`
	UserName string `json:"userName,omitempty"` // Or full User DTO
}

// CardCollaboratorPayload for collaborator changes
type CardCollaboratorPayload struct {
	CardID   uint `json:"cardId"`
	UserID   uint `json:"userId"`
	BoardID  uint `json:"boardId"` // For client-side context
	UserName string `json:"userName,omitempty"`
}
