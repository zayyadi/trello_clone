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
		log.Printf("INFO [WebSocket ReadPump]: Finished cleanup for client (User: %d, Board: %d, Remote: %s)", c.UserID, c.BoardID, c.Conn.RemoteAddr())
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR [WebSocket ReadPump]: Unexpected close error for User: %d, Board: %d, Client: %s: %v", c.UserID, c.BoardID, c.Conn.RemoteAddr(), err)
			} else {
				log.Printf("INFO [WebSocket ReadPump]: Read error / connection closed for User: %d, Board: %d, Client: %s: %v", c.UserID, c.BoardID, c.Conn.RemoteAddr(), err)
			}
			break // Exit loop
		}
		// For now, messages from clients are discarded.
		log.Printf("INFO [WebSocket ReadPump]: Received message from client (User: %d, Board: %d, Remote: %s) (currently discarded): %s", c.UserID, c.BoardID, c.Conn.RemoteAddr(), message)
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
		log.Printf("INFO [WebSocket WritePump]: Finished cleanup for client (User: %d, Board: %d, Remote: %s)", c.UserID, c.BoardID, c.Conn.RemoteAddr())
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the c.Send channel.
				log.Printf("INFO [WebSocket WritePump]: Hub closed send channel for client (User: %d, Board: %d, Remote: %s). Sending close message.", c.UserID, c.BoardID, c.Conn.RemoteAddr())
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return // Exit goroutine
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("ERROR [WebSocket WritePump]: Error getting next writer for client (User: %d, Board: %d, Remote: %s): %v", c.UserID, c.BoardID, c.Conn.RemoteAddr(), err)
				return // Exit goroutine
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'}) // Add a newline if sending multiple messages in one frame
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				log.Printf("ERROR [WebSocket WritePump]: Error closing writer for client (User: %d, Board: %d, Remote: %s): %v", c.UserID, c.BoardID, c.Conn.RemoteAddr(), err)
				return // Exit goroutine
			}
			// log.Printf("INFO [WebSocket WritePump]: Sent message to client (User: %d, Board: %d, Remote: %s)", c.UserID, c.BoardID, c.Conn.RemoteAddr())

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ERROR [WebSocket WritePump]: Ping failed for client (User: %d, Board: %d, Remote: %s): %v. Closing connection.", c.UserID, c.BoardID, c.Conn.RemoteAddr(), err)
				return // Exit goroutine
			}
			// log.Printf("INFO [WebSocket WritePump]: Sent ping to client (User: %d, Board: %d, Remote: %s)", c.UserID, c.BoardID, c.Conn.RemoteAddr())
		}
	}
}
