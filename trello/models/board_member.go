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

// package models

// import (
// 	"gorm.io/gorm"
// )

// // BoardMember represents the many-to-many relationship between Users and Boards
// type BoardMember struct {
// 	gorm.Model
// 	BoardID uint `gorm:"primaryKey;autoIncrement:false" json:"boardID"` // Composite primary key with UserID
// 	UserID  uint `gorm:"primaryKey;autoIncrement:false" json:"userID"`  // Composite primary key with BoardID
// 	User    User `gorm:"foreignKey:UserID" json:"user,omitempty"`       // Eager load user details if needed
// 	// Role    string `gorm:"default:'member'" json:"role"` // e.g., "admin", "member", "viewer" - for future enhancement
// }
