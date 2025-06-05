package models

import "time"

// CardCollaborator defines the join table for the many-to-many relationship
// between Cards and Users (Collaborators).
type CardCollaborator struct {
	CardID    uint      `gorm:"primaryKey;autoIncrement:false" json:"cardID"`
	UserID    uint      `gorm:"primaryKey;autoIncrement:false" json:"userID"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	// Optional: For preloading or direct queries on this join model.
	// User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// Card Card `gorm:"foreignKey:CardID" json:"card,omitempty"`
}
