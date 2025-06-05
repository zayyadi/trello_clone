package repositories

import (
	"github.com/zayyadi/trello/models"
	"gorm.io/gorm"
)

type BoardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepositoryInterface {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(board *models.Board) error {
	return r.db.Create(board).Error
}

func (r *BoardRepository) FindByID(id uint) (*models.Board, error) {
	var board models.Board
	// Preload Owner and Members with their User details
	err := r.db.Preload("Owner").Preload("Members.User").First(&board, id).Error
	return &board, err
}

func (r *BoardRepository) FindByOwnerOrMember(userID uint) ([]models.Board, error) {
	var boards []models.Board
	// Find boards where user is owner OR is a member
	// This query is a bit more complex.
	// Option 1: Two queries and merge
	// Option 2: Subquery or JOIN
	// Let's try a JOIN approach
	err := r.db.Joins("LEFT JOIN board_members on board_members.board_id = boards.id").
		Where("boards.owner_id = ? OR board_members.user_id = ?", userID, userID).
		Preload("Owner").Preload("Members.User").Distinct().Find(&boards).Error
	return boards, err
}

func (r *BoardRepository) Update(board *models.Board) error {
	return r.db.Save(board).Error
}

func (r *BoardRepository) Delete(id uint) error {
	return r.db.Delete(&models.Board{}, id).Error
}

func (r *BoardRepository) IsOwner(boardID, userID uint) (bool, error) {
	var board models.Board
	if err := r.db.Select("owner_id").First(&board, boardID).Error; err != nil {
		return false, err
	}
	return board.OwnerID == userID, nil
}
