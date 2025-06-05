package models

import (
	"gorm.io/gorm"
)

// Board model
type Board struct {
	gorm.Model
	Name        string        `gorm:"not null" json:"name"`
	Description string        `json:"description"`
	OwnerID     uint          `gorm:"not null" json:"ownerID"`
	Owner       User          `gorm:"foreignKey:OwnerID" json:"owner"` // Belongs to Owner
	Lists       []List        `gorm:"foreignKey:BoardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"lists,omitempty"`
	Members     []BoardMember `gorm:"foreignKey:BoardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members,omitempty"`
}

// TableName returns the table name for the Board model
