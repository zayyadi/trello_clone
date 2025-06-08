package handlers

import (
	"log"
	"net/http"
	"strconv"
	// "strings" // Not used yet, but might be useful later

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // Import jwt
	"github.com/gorilla/websocket"
	// "github.com/zayyadi/trello/config" // JWTSecret will be passed in NewWebSocketHandler
	"github.com/zayyadi/trello/models"  // Required for Claims
	"github.com/zayyadi/trello/realtime"
	"github.com/zayyadi/trello/services" // Will be needed for BoardService
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default
		// TODO: For production, you might want to check the origin of the request
		// against a whitelist of allowed origins.
		return true
	},
}

// WebSocketHandler handles WebSocket connections.
type WebSocketHandler struct {
	hub          *realtime.Hub
	boardService services.BoardServiceInterface // Use interface type
	jwtSecret    string
}

// NewWebSocketHandler creates a new WebSocketHandler.
func NewWebSocketHandler(hub *realtime.Hub, boardService services.BoardServiceInterface, jwtSecret string) *WebSocketHandler {
	return &WebSocketHandler{
		hub:          hub,
		boardService: boardService,
		jwtSecret:    jwtSecret,
	}
}

// HandleConnections upgrades HTTP connections to WebSocket connections
// and registers clients with the hub.
func (h *WebSocketHandler) HandleConnections(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		log.Println("WebSocket: Token not provided")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
		return
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		log.Printf("WebSocket: Invalid token: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	userID := claims.UserID

	boardIDStr := c.Query("boardID")
	if boardIDStr == "" {
		log.Println("WebSocket: boardID not provided")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "boardID query parameter is required"})
		return
	}

	boardIDUint64, err := strconv.ParseUint(boardIDStr, 10, 32)
	if err != nil {
		log.Printf("WebSocket: Invalid boardID format: %s", boardIDStr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid boardID format"})
		return
	}
	boardID := uint(boardIDUint64)

	// Authorization Check
	isMember, err := h.boardService.IsUserMemberOfBoard(userID, boardID)
	if err != nil {
	 log.Printf("WebSocket: Error checking board membership for user %d, board %d: %v", userID, boardID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not verify board membership"})
		return
	}
	if !isMember {
	 log.Printf("WebSocket: User %d not authorized for board %d", userID, boardID)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User not authorized for this board"})
		return
	}
	log.Printf("WebSocket: User %d authorized for board %d", userID, boardID)


	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket: Failed to upgrade connection for user %d, board %d: %v", userID, boardID, err)
		// upgrader.Upgrade automatically sends an HTTP error response on failure,
		// so no explicit c.AbortWithStatusJSON is needed here if it already sent one.
		// However, if it doesn't, or to be sure, we can add one.
		// For now, relying on upgrader's behavior.
		return
	}

	client := &realtime.Client{
		Hub:     h.hub,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		BoardID: boardID,
		UserID:  userID, // Set the UserID on the client
	}

	client.Hub.Register <- client // Pass client to register channel

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()

	log.Printf("WebSocket: Client (User: %d) connected to board %d", userID, boardID)
}
