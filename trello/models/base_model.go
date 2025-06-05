package models

import (
	"gorm.io/gorm"
)

// BaseModel defines common fields for GORM models
// GORM's gorm.Model already provides ID, CreatedAt, UpdatedAt, DeletedAt.
// We can use it directly or define our own if we need to customize.
// For simplicity, we will use gorm.Model.
// If you needed custom types for ID (e.g. UUID) or different naming,
// you would define your own Base here.

type BaseModel = gorm.Model // Alias for convenience

// Example of custom base if needed:
// type BaseModel struct {
// 	ID        uint           `gorm:"primarykey" json:"id"`
// 	CreatedAt time.Time      `json:"createdAt"`
// 	UpdatedAt time.Time      `json:"updatedAt"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Use json:"-" to hide from API by default
// }
