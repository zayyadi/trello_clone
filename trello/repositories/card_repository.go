package repositories

import (
	"log" // Import log package

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
	err := r.db.Create(card).Error
	if err != nil {
		log.Printf("ERROR [CardRepository.Create]: Failed to create card in DB. Input: %+v, Error: %v\n", card, err)
	}
	return err
}

func (r *CardRepository) FindByID(id uint) (*models.Card, error) {
	var card models.Card
	// Preload AssignedUser, Supervisor and Collaborators
	err := r.db.Preload("AssignedUser").Preload("Supervisor").Preload("Collaborators").First(&card, id).Error
	return &card, err
}

func (r *CardRepository) FindByListID(listID uint) ([]models.Card, error) {
	var cards []models.Card
	// Preload AssignedUser, Supervisor and Collaborators for each card
	err := r.db.Where("list_id = ?", listID).Order("position ASC").
		Preload("AssignedUser").Preload("Supervisor").Preload("Collaborators").
		Find(&cards).Error
	return cards, err
}

func (r *CardRepository) Update(card *models.Card) error {
	err := r.db.Save(card).Error
	if err != nil {
		log.Printf("ERROR [CardRepository.Update]: Failed to update card in DB. Input: %+v, Error: %v\n", card, err)
	}
	return err
}

func (r *CardRepository) Delete(id uint) error {
	err := r.db.Delete(&models.Card{}, id).Error
	if err != nil {
		log.Printf("ERROR [CardRepository.Delete]: Failed to delete card with ID %d. Error: %v\n", id, err)
	}
	return err
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
			log.Printf("ERROR [CardRepository.MoveCard.Save]: Failed to save moved card %d. Error: %v\n", cardID, err)
			return err
		}

		return nil
	})
	// err is implicitly returned by the Transaction func if any step fails
	// Adding a generic log here if transaction itself fails might be useful too, but GORM handles that.
	return err // return the error from the transaction
}

func (r *CardRepository) AddCollaborator(cardID uint, userID uint) error {
	collaborator := models.CardCollaborator{CardID: cardID, UserID: userID}
	err := r.db.FirstOrCreate(&collaborator).Error
	if err != nil {
		log.Printf("ERROR [CardRepository.AddCollaborator]: Failed to add collaborator UserID %d to CardID %d. Error: %v\n", userID, cardID, err)
	}
	return err
}

func (r *CardRepository) RemoveCollaborator(cardID uint, userID uint) error {
	collaborator := models.CardCollaborator{CardID: cardID, UserID: userID}
	result := r.db.Delete(&collaborator)
	if result.Error != nil {
		log.Printf("ERROR [CardRepository.RemoveCollaborator]: Failed to remove collaborator UserID %d from CardID %d. Error: %v\n", userID, cardID, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("WARN [CardRepository.RemoveCollaborator]: No collaborator found for UserID %d on CardID %d to remove.\n", userID, cardID)
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

// IsUserCollaboratorOrAssignee checks if a user is either directly assigned to the card
// or is listed as a collaborator.
func (r *CardRepository) IsUserCollaboratorOrAssignee(cardID uint, userID uint) (bool, error) {
	var card models.Card
	// Check if the user is the AssignedUserID
	// We select only AssignedUserID to make the query lightweight.
	if err := r.db.Select("assigned_user_id").First(&card, cardID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, gorm.ErrRecordNotFound // Card itself not found
		}
		return false, err // Other database error
	}

	if card.AssignedUserID != nil && *card.AssignedUserID == userID {
		return true, nil // User is the assignee
	}

	// If not the assignee, check if the user is a collaborator
	return r.IsCollaborator(cardID, userID)
}
