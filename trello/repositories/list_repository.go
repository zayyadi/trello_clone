package repositories

import (
	"github.com/zayyadi/trello/models"

	"gorm.io/gorm"
)

type ListRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) ListRepositoryInterface {
	return &ListRepository{db: db}
}

func (r *ListRepository) Create(list *models.List) error {
	// Set position to max_position + 1 for the board
	var maxPosition uint
	r.db.Model(&models.List{}).Where("board_id = ?", list.BoardID).Select("COALESCE(MAX(position), 0)").Row().Scan(&maxPosition)
	list.Position = maxPosition + 1
	return r.db.Create(list).Error
}

func (r *ListRepository) FindByID(id uint) (*models.List, error) {
	var list models.List
	// Preload cards sorted by position
	err := r.db.Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Order("cards.position ASC")
	}).First(&list, id).Error
	return &list, err
}

func (r *ListRepository) FindByBoardID(boardID uint) ([]models.List, error) {
	var lists []models.List
	// Preload cards for each list, sorted by position
	err := r.db.Where("board_id = ?", boardID).Order("position ASC").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Order("cards.position ASC")
		}).
		Find(&lists).Error
	return lists, err
}

func (r *ListRepository) Update(list *models.List) error {
	return r.db.Save(list).Error
}

// UpdatePositions updates positions of multiple lists (e.g., after drag-drop)
// listsWithNewPositions is a map of listID to newPosition
func (r *ListRepository) UpdatePositions(boardID uint, listsWithNewPositions map[uint]uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for listID, newPosition := range listsWithNewPositions {
			if err := tx.Model(&models.List{}).Where("id = ? AND board_id = ?", listID, boardID).Update("position", newPosition).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ListRepository) Delete(id uint) error {
	// Related cards are deleted by CASCADE constraint if DB supports it and GORM is configured.
	// Otherwise, delete cards manually or ensure cascade is set up.
	return r.db.Delete(&models.List{}, id).Error
}

func (r *ListRepository) GetBoardIDByListID(listID uint) (uint, error) {
	var list models.List
	if err := r.db.Select("board_id").First(&list, listID).Error; err != nil {
		return 0, err
	}
	return list.BoardID, nil
}

func (r *ListRepository) GetMaxPosition(boardID uint) (uint, error) {
	var maxPosition uint
	result := r.db.Model(&models.List{}).Where("board_id = ?", boardID).Select("COALESCE(MAX(position), 0)").Row()
	if err := result.Scan(&maxPosition); err != nil {
		return 0, err
	}
	return maxPosition, nil
}

// PerformTransaction executes the given function within a database transaction.
func (r *ListRepository) PerformTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
