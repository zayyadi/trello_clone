package repositories

import (
	"github.com/zayyadi/trello/models"
	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepositoryInterface {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(card *models.Card) error {
	// Set position to max_position + 1 for the list
	var maxPosition uint
	r.db.Model(&models.Card{}).Where("list_id = ?", card.ListID).Select("COALESCE(MAX(position), 0)").Row().Scan(&maxPosition)
	card.Position = maxPosition + 1
	return r.db.Create(card).Error
}

func (r *CardRepository) FindByID(id uint) (*models.Card, error) {
	var card models.Card
	// Preload AssignedUser and Supervisor
	err := r.db.Preload("AssignedUser").Preload("Supervisor").First(&card, id).Error
	return &card, err
}

func (r *CardRepository) FindByListID(listID uint) ([]models.Card, error) {
	var cards []models.Card
	// Preload AssignedUser and Supervisor for each card
	err := r.db.Where("list_id = ?", listID).Order("position ASC").
		Preload("AssignedUser").Preload("Supervisor").
		Find(&cards).Error
	return cards, err
}

func (r *CardRepository) Update(card *models.Card) error {
	return r.db.Save(card).Error
}

func (r *CardRepository) Delete(id uint) error {
	return r.db.Delete(&models.Card{}, id).Error
}

func (r *CardRepository) GetListIDByCardID(cardID uint) (uint, error) {
	var card models.Card
	if err := r.db.Select("list_id").First(&card, cardID).Error; err != nil {
		return 0, err
	}
	return card.ListID, nil
}

func (r *CardRepository) UpdateCardPosition(cardID, listID, newPosition uint) error {
	return r.db.Model(&models.Card{}).Where("id = ?", cardID).Updates(map[string]interface{}{"list_id": listID, "position": newPosition}).Error
}

// ShiftPositions adjusts positions of cards in a list, typically when a card is added, moved, or removed.
// listID: The ID of the list whose cards' positions need adjustment.
// startPosition: The position from which to start shifting.
// shiftAmount: The amount to shift by (positive to increase positions, negative to decrease).
// excludedCardID: (Optional) A card ID to exclude from shifting (e.g., the card being moved).
func (r *CardRepository) ShiftPositions(listID uint, startPosition uint, shiftAmount int, excludedCardID *uint) error {
	query := r.db.Model(&models.Card{}).Where("list_id = ? AND position >= ?", listID, startPosition)
	if excludedCardID != nil {
		query = query.Where("id != ?", *excludedCardID)
	}
	return query.UpdateColumn("position", gorm.Expr("position + ?", shiftAmount)).Error
}

func (r *CardRepository) GetMaxPosition(listID uint) (uint, error) {
	var maxPosition uint
	result := r.db.Model(&models.Card{}).Where("list_id = ?", listID).Select("COALESCE(MAX(position), 0)").Row()
	if err := result.Scan(&maxPosition); err != nil {
		return 0, err
	}
	return maxPosition, nil
}

// MoveCard updates the ListID and Position of a card.
// It also handles reordering in the source and destination lists.
func (r *CardRepository) PerformTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *CardRepository) MoveCard(cardID, oldListID, newListID uint, newPosition uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var cardToMove models.Card
		if err := tx.First(&cardToMove, cardID).Error; err != nil {
			return err // Card not found
		}
		currentPosition := cardToMove.Position

		// 1. Shift cards in the old list (if moving from a list)
		if cardToMove.ListID == oldListID { // Ensure it's actually in oldListID
			// Decrement positions of cards that were after the moved card in the old list
			if err := tx.Model(&models.Card{}).
				Where("list_id = ? AND position > ?", oldListID, currentPosition).
				Update("position", gorm.Expr("position - 1")).Error; err != nil {
				return err
			}
		}

		// 2. Shift cards in the new list
		// Increment positions of cards at or after the new position in the new list
		if err := tx.Model(&models.Card{}).
			Where("list_id = ? AND position >= ?", newListID, newPosition).
			Update("position", gorm.Expr("position + 1")).Error; err != nil {
			return err
		}

		// 3. Update the card itself
		cardToMove.ListID = newListID
		cardToMove.Position = newPosition
		if err := tx.Save(&cardToMove).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *CardRepository) AddCollaborator(cardID uint, userID uint) error {
	collaborator := models.CardCollaborator{CardID: cardID, UserID: userID}
	// Using FirstOrCreate to prevent duplicate entries if the relationship already exists.
	// If it already exists, it does nothing and returns no error.
	// If you need to know if it was newly created or already existed, FirstOrCreate is not ideal.
	// In that case, a check followed by Create would be better.
	// For adding, ensuring the link exists is usually sufficient.
	return r.db.FirstOrCreate(&collaborator).Error
}

func (r *CardRepository) RemoveCollaborator(cardID uint, userID uint) error {
	collaborator := models.CardCollaborator{CardID: cardID, UserID: userID}
	// Delete the specific association.
	// GORM's Delete with a struct containing primary key values will delete that record.
	result := r.db.Delete(&collaborator)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Or a custom "collaborator not found" error
	}
	return nil
}

func (r *CardRepository) GetCollaboratorsByCardID(cardID uint) ([]models.User, error) {
	var users []models.User
	// To get collaborators (Users) for a card, we need to join through card_collaborators.
	// Then join with users table.
	err := r.db.Joins("JOIN card_collaborators cc ON cc.user_id = users.id").
		Where("cc.card_id = ?", cardID).
		Find(&users).Error
	return users, err
}

func (r *CardRepository) IsCollaborator(cardID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.CardCollaborator{}).
		Where("card_id = ? AND user_id = ?", cardID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
