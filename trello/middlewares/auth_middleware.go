package middleware

import (
	"net/http"
	"strings"

	"github.com/zayyadi/trello/handlers" // For RespondWithError
	"github.com/zayyadi/trello/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handlers.RespondWithError(c, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			handlers.RespondWithError(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateJWT(tokenString, jwtSecret)
		if err != nil {
			handlers.RespondWithError(c, http.StatusUnauthorized, "Invalid or expired token: "+err.Error())
			return
		}

		c.Set("userID", claims.UserID) // Store userID in context for handlers
		c.Next()
	}
}
