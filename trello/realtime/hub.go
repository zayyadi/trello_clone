package realtime

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients. The outer map's key is boardID.
	clients map[uint]map[*Client]bool

	// Inbound messages from the services to be broadcast to clients.
	broadcast chan *WebSocketMessage

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *WebSocketMessage), // Use WebSocketMessage
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uint]map[*Client]bool),
	}
}

// Submit queues a message for broadcast by the hub.
func (h *Hub) Submit(msg *WebSocketMessage) {
	// This send is blocking if the broadcast channel is full.
	// The hub's Run method should process messages from this channel quickly.
	// If the broadcast channel is buffered, it can absorb some burst.
	// If this becomes a bottleneck (e.g., services submitting faster than hub can process/broadcast),
	// consider increasing buffer size of h.broadcast or making this Submit non-blocking
	// with a select-default, though that might mean dropping messages under heavy load.
	h.broadcast <- msg
}

// Run starts the hub's event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.clients[client.BoardID]; !ok { // Use client.BoardID
				h.clients[client.BoardID] = make(map[*Client]bool)
			}
			h.clients[client.BoardID][client] = true
			log.Printf("Client registered for board %d. Total clients on board: %d", client.BoardID, len(h.clients[client.BoardID]))
		case client := <-h.unregister:
			if boardClients, ok := h.clients[client.BoardID]; ok {
				if _, ok := boardClients[client]; ok {
					delete(boardClients, client)
					// Do not close client.Send here, it's closed by writePump if necessary
					// or if client.conn.Close() leads to writePump exit.
					// close(client.Send) // This was in the original, but writePump/readPump defer closing conn/unregistering
					if len(boardClients) == 0 {
						delete(h.clients, client.BoardID)
						log.Printf("Board %d has no more clients, removing from hub.", client.BoardID)
					}
					log.Printf("Client unregistered from board %d. Total clients on board: %d", client.BoardID, len(boardClients))
				}
			}
		case wsMessage := <-h.broadcast:
			boardClients, ok := h.clients[wsMessage.BoardID]
			if !ok {
				log.Printf("No clients registered for board %d, message not sent.", wsMessage.BoardID)
				continue
			}

			// Prepare the message to be sent to clients (Type and Payload)
			clientMessage := struct {
				Type    string      `json:"type"`
				Payload interface{} `json:"payload,omitempty"`
			}{
				Type:    wsMessage.Type,
				Payload: wsMessage.Payload,
			}

			messageBytes, err := json.Marshal(clientMessage)
			if err != nil {
				log.Printf("Error marshalling WebSocket message for board %d: %v", wsMessage.BoardID, err)
				continue // Skip broadcasting this message if marshalling fails
			}

			log.Printf("Broadcasting message type '%s' to %d clients on board %d (Originating UserID: %d)", wsMessage.Type, len(boardClients), wsMessage.BoardID, wsMessage.UserID)
			for client := range boardClients {
				// Don't send message back to the originating user if UserID is set on the message
				// and the client has a UserID.
				if wsMessage.UserID != 0 && client.UserID == wsMessage.UserID {
					log.Printf("Skipping broadcast to originating client (User: %d) for message type '%s' on board %d", client.UserID, wsMessage.Type, wsMessage.BoardID)
					continue
				}

				select {
				case client.Send <- messageBytes:
				default:
					// If the client's send buffer is full, it implies the client is slow
					// or disconnected. readPump/writePump will handle actual unregistration.
					// We should not block the hub or unregister from here directly.
					log.Printf("Client send buffer full for board %d. Client: %p. Message type: %s. Client will be cleaned up by its own pumps.", client.BoardID, client, wsMessage.Type)
					// Removing the client here can lead to race conditions if client unregisters itself simultaneously.
					// delete(boardClients, client)
					// close(client.Send)
					// if len(boardClients) == 0 {
					// 	delete(h.clients, wsMessage.BoardID)
					// }
				}
			}
		}
	}
}
