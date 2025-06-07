package models

import (
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model                       // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Username           string        `gorm:"uniqueIndex;not null" json:"username"`
	Email              string        `gorm:"uniqueIndex;not null" json:"email"`
	Password           string        `gorm:"not null" json:"-"`                  // json:"-" to hide password hash
	Boards             []Board       `gorm:"foreignKey:OwnerID" json:"-"`        // Boards owned by this user
	MemberOfBoards     []BoardMember `gorm:"foreignKey:UserID" json:"-"`         // Boards this user is a member of
	AssignedCards      []Card        `gorm:"foreignKey:AssignedUserID" json:"-"` // Cards assigned to this user
	CollaboratingCards []*Card       `gorm:"many2many:card_collaborators;" json:"-"`
}
