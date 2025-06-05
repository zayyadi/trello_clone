package repositories

import (
	"github.com/zayyadi/trello/models"
	"gorm.io/gorm"
)

// UserRepositoryInterface defines the contract for user repository operations.
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}

// BoardRepositoryInterface defines the contract for board repository operations.
type BoardRepositoryInterface interface {
	Create(board *models.Board) error
	FindByID(id uint) (*models.Board, error)
	FindByOwnerOrMember(userID uint) ([]models.Board, error)
	Update(board *models.Board) error
	Delete(id uint) error
	IsOwner(boardID uint, userID uint) (bool, error)
}

// BoardMemberRepositoryInterface defines the contract for board member repository operations.
type BoardMemberRepositoryInterface interface {
	AddMember(member *models.BoardMember) error
	RemoveMember(boardID uint, userID uint) error
	IsMember(boardID uint, userID uint) (bool, error)
	FindMembersByBoardID(boardID uint) ([]models.BoardMember, error)
	FindByBoardIDAndUserID(boardID uint, userID uint) (*models.BoardMember, error)
}

// ListRepositoryInterface defines the contract for list repository operations.
type ListRepositoryInterface interface {
	Create(list *models.List) error
	FindByID(id uint) (*models.List, error)
	FindByBoardID(boardID uint) ([]models.List, error)
	Update(list *models.List) error
	Delete(id uint) error
	GetMaxPosition(boardID uint) (uint, error)
	// GetDB() *gorm.DB // For transaction handling - REMOVED
	PerformTransaction(fn func(tx *gorm.DB) error) error // ADDED
	GetBoardIDByListID(listID uint) (uint, error)
}

// CardRepositoryInterface defines the contract for card repository operations.
type CardRepositoryInterface interface {
	Create(card *models.Card) error
	FindByID(id uint) (*models.Card, error)
	FindByListID(listID uint) ([]models.Card, error)
	Update(card *models.Card) error
	Delete(id uint) error
	GetListIDByCardID(cardID uint) (uint, error)
	PerformTransaction(fn func(tx *gorm.DB) error) error
	MoveCard(cardID, oldListID, newListID uint, newPosition uint) error
	GetMaxPosition(listID uint) (uint, error)
	ShiftPositions(listID uint, startPosition uint, shiftAmount int, excludedCardID *uint) error
	AddCollaborator(cardID uint, userID uint) error
	RemoveCollaborator(cardID uint, userID uint) error
	GetCollaboratorsByCardID(cardID uint) ([]models.User, error)
	IsCollaborator(cardID uint, userID uint) (bool, error)
}

// CommentRepositoryInterface defines the contract for comment repository operations.
type CommentRepositoryInterface interface {
	Create(comment *models.Comment) error
	FindByCardID(cardID uint) ([]models.Comment, error)
	FindByID(id uint) (*models.Comment, error)
}
