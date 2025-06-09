package models

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the custom claims for the JWT
type Claims struct {
	UserID uint `json:"userID"`
	jwt.RegisteredClaims
}
