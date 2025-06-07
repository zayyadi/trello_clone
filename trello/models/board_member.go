package models

import (
	"time"
)

// BoardMember represents the many-to-many relationship between Users and Boards
type BoardMember struct {
	BoardID   uint      `gorm:"primaryKey" json:"boardID"`
	UserID    uint      `gorm:"primaryKey" json:"userID"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Board     Board     `gorm:"foreignKey:BoardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // To avoid circular refs in JSON
	CreatedAt time.Time `json:"createdAt"`
	// Role    string `gorm:"default:'member'" json:"role"` // e.g., "admin", "member", "viewer"
}
