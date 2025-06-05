package models

import (
	"time"

	"gorm.io/gorm"
)

// CardStatus type for predefined statuses
type CardStatus string

const (
	StatusToDo    CardStatus = "to_do"
	StatusPending CardStatus = "pending" // Could be "in_progress"
	StatusDone    CardStatus = "done"
	StatusUndone  CardStatus = "undone" // Perhaps a synonym for "to_do" or for reopened tasks
)

// Card model (Task card within a list)
type Card struct {
	gorm.Model
	Title          string     `gorm:"not null" json:"title"`
	Description    string     `json:"description"`
	ListID         uint       `gorm:"not null" json:"listID"`
	List           List       `gorm:"foreignKey:ListID" json:"-"`
	Position       uint       `gorm:"not null;default:0" json:"position"`
	DueDate        *time.Time `json:"dueDate,omitempty"` // Already exists, ensure it's used
	Status         CardStatus `gorm:"type:varchar(20);default:'to_do'" json:"status"`
	AssignedUserID *uint      `json:"assignedUserID,omitempty"` // User doing the task
	AssignedUser   *User      `gorm:"foreignKey:AssignedUserID" json:"assignedUser,omitempty"`
	SupervisorID   *uint      `json:"supervisorID,omitempty"` // User supervising the task
	Supervisor     *User      `gorm:"foreignKey:SupervisorID" json:"supervisor,omitempty"`
	Comments       []Comment  `gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE;" json:"comments,omitempty"`
	Collaborators  []*User    `gorm:"many2many:card_collaborators;constraint:OnDelete:CASCADE;" json:"collaborators,omitempty"`
	Color          *string    `gorm:"type:varchar(7)" json:"color,omitempty"` // Hex color like #RRGGBB
}
