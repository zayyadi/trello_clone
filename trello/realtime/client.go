package realtime

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub // Changed to Hub to match ws_handler.go
	// The websocket connection.
	Conn *websocket.Conn // Changed to Conn to match ws_handler.go
	// Buffered channel of outbound messages.
	Send chan []byte // Changed to Send to match ws_handler.go
	// The ID of the board this client is interested in.
	BoardID uint // Changed to BoardID to match ws_handler.go
	// UserID of the connected user.
	UserID uint
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() { // Changed to ReadPump to match ws_handler.go
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
		log.Printf("Client disconnected from board %d (readPump cleanup)", c.BoardID)
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("Client read error on board %d, disconnecting: %v", c.BoardID, err)
			break
		}
		// For now, messages from clients are discarded.
		// If client-to-server messages were needed (e.g., user actions),
		// they would be processed here and potentially passed to the hub.
		log.Printf("Received message from client on board %d (discarded): %s", c.BoardID, message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() { // Changed to WritePump to match ws_handler.go
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		log.Printf("Client disconnected from board %d (writePump cleanup)", c.BoardID)
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Printf("Hub closed channel for client on board %d. Sending close message.", c.BoardID)
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing message to client on board %d: %v", c.BoardID, err)
				return
			}
			// log.Printf("Sent message to client on board %d: %s", c.BoardID, message) // Can be too verbose

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping to client on board %d: %v", c.BoardID, err)
				return
			}
			// log.Printf("Sent ping to client on board %d", c.BoardID) // Can be too verbose
		}
	}
}
