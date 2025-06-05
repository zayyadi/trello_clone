package repositories

import (
	"github.com/zayyadi/trello/models"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepositoryInterface {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) FindByCardID(cardID uint) ([]models.Comment, error) {
	var comments []models.Comment
	// Preload User details for each comment, order by creation date
	err := r.db.Preload("User").Where("card_id = ?", cardID).Order("created_at asc").Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	return &comment, err
}
