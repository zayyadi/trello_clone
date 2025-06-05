package models

import (
	"time"

	"gorm.io/gorm"
)

// List model (Kanban list within a board)
type List struct {
	gorm.Model
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	BoardID   uint   `gorm:"not null" json:"boardID"`
	Board     Board  `gorm:"foreignKey:BoardID" json:"-"`        // Belongs to Board
	Position  uint   `gorm:"not null;default:0" json:"position"` // Order of the list within the board
	Cards     []Card `gorm:"foreignKey:ListID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cards,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
