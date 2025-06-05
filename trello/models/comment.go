package models

import "gorm.io/gorm"

// Comment represents a comment made by a user on a card.
type Comment struct {
	gorm.Model        // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Content    string `gorm:"not null" json:"content"`
	CardID     uint   `gorm:"not null;index" json:"cardID"` // Index for faster lookups
	UserID     uint   `gorm:"not null" json:"userID"`
	User       User   `gorm:"foreignKey:UserID" json:"user"` // For preloading author details
	// Card       Card   `gorm:"foreignKey:CardID" json:"-"` // Optional: if you need Comment.Card back-reference
}
