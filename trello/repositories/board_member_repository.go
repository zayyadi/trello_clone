package repositories

import (
	"github.com/zayyadi/trello/models"
	"gorm.io/gorm"
)

type BoardMemberRepository struct {
	db *gorm.DB
}

func NewBoardMemberRepository(db *gorm.DB) BoardMemberRepositoryInterface {
	return &BoardMemberRepository{db: db}
}

func (r *BoardMemberRepository) AddMember(member *models.BoardMember) error {
	return r.db.Create(member).Error
}

func (r *BoardMemberRepository) RemoveMember(boardID, userID uint) error {
	return r.db.Where("board_id = ? AND user_id = ?", boardID, userID).Delete(&models.BoardMember{}).Error
}

func (r *BoardMemberRepository) IsMember(boardID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.BoardMember{}).Where("board_id = ? AND user_id = ?", boardID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *BoardMemberRepository) FindMembersByBoardID(boardID uint) ([]models.BoardMember, error) {
	var members []models.BoardMember
	err := r.db.Preload("User").Where("board_id = ?", boardID).Find(&members).Error
	return members, err
}

func (r *BoardMemberRepository) FindByBoardIDAndUserID(boardID, userID uint) (*models.BoardMember, error) {
	var boardMember models.BoardMember
	err := r.db.Preload("User").Where("board_id = ? AND user_id = ?", boardID, userID).First(&boardMember).Error
	return &boardMember, err
}
