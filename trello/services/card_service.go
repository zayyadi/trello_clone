package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"

	"gorm.io/gorm"
)

// CardServiceInterface defines the contract for card service operations
type CardServiceInterface interface {
	CreateCard(listID uint, title, description string, position *uint, dueDate *time.Time, assignedUserID *uint, supervisorID *uint, color *string, currentUserID uint) (*models.Card, error)
	GetCardByID(cardID uint, currentUserID uint) (*models.Card, error)
	GetCardsByListID(listID uint, currentUserID uint) ([]models.Card, error)
	UpdateCard(cardID uint, title, description *string, newPosition *uint, dueDate *time.Time, assignedUserID **uint, supervisorID **uint, status *models.CardStatus, color *string, currentUserID uint) (*models.Card, error)
	DeleteCard(cardID uint, currentUserID uint) error
	MoveCard(cardID uint, targetListID uint, newPosition uint, currentUserID uint) (*models.Card, error)
	AddCollaboratorToCard(cardID uint, currentUserID uint, targetUserEmail string, targetUserIDInput *uint) (*models.User, error)
	RemoveCollaboratorFromCard(cardID uint, currentUserID uint, targetUserID uint) error
	GetCardCollaborators(cardID uint, currentUserID uint) ([]models.User, error)
}

type CardService struct {
	cardRepo        repositories.CardRepositoryInterface
	listRepo        repositories.ListRepositoryInterface
	boardRepo       repositories.BoardRepositoryInterface // For permission checks
	boardMemberRepo repositories.BoardMemberRepositoryInterface
	userRepo        repositories.UserRepositoryInterface // Added for collaborator methods
}

func NewCardService(
	cardRepo repositories.CardRepositoryInterface,
	listRepo repositories.ListRepositoryInterface,
	boardRepo repositories.BoardRepositoryInterface,
	boardMemberRepo repositories.BoardMemberRepositoryInterface,
	userRepo repositories.UserRepositoryInterface, // Added
) CardServiceInterface { // Return interface type
	return &CardService{
		cardRepo:        cardRepo,
		listRepo:        listRepo,
		boardRepo:       boardRepo,
		boardMemberRepo: boardMemberRepo,
		userRepo:        userRepo, // Added
	}
}

// Helper to check board access via list
func (s *CardService) checkAccessViaList(userID, listID uint) (uint, error) {
	boardID, err := s.listRepo.GetBoardIDByListID(listID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrListNotFound
		}
		return 0, err
	}

	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrBoardNotFound
		}
		return 0, err
	}
	if board.OwnerID == userID {
		return boardID, nil
	}
	isMember, err := s.boardMemberRepo.IsMember(boardID, userID)
	if err != nil || !isMember {
		return 0, ErrForbidden
	}
	return boardID, nil
}

// Helper to check board access via card
func (s *CardService) checkAccessViaCard(userID, cardID uint) (uint, uint, error) { // Returns boardID, listID
	listID, err := s.cardRepo.GetListIDByCardID(cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, ErrCardNotFound
		}
		return 0, 0, err
	}
	boardID, err := s.checkAccessViaList(userID, listID)
	return boardID, listID, err
}

func (s *CardService) CreateCard(listID uint, title, description string, position *uint, dueDate *time.Time, assignedUserID *uint, supervisorID *uint, color *string, currentUserID uint) (*models.Card, error) {
	if _, err := s.checkAccessViaList(currentUserID, listID); err != nil {
		return nil, err
	}

	card := &models.Card{
		ListID:         listID,
		Title:          title,
		Description:    description,
		DueDate:        dueDate,
		AssignedUserID: assignedUserID,
		SupervisorID:   supervisorID,
		Status:         models.StatusToDo, // Default status
		Color:          color,             // Add color
	}
	// Position handling will be done by the repository's Create method typically
	if position != nil { // If a specific position is requested by client (less common for create)
		card.Position = *position // Note: This might need more complex logic if repo doesn't handle reordering on create
	}

	if err := s.cardRepo.Create(card); err != nil {
		return nil, err
	}
	// ...
	return s.cardRepo.FindByID(card.ID) // Will now preload supervisor too
}

func (s *CardService) GetCardByID(cardID uint, currentUserID uint) (*models.Card, error) {
	if _, _, err := s.checkAccessViaCard(currentUserID, cardID); err != nil {
		return nil, err
	}
	return s.cardRepo.FindByID(cardID)
}

func (s *CardService) GetCardsByListID(listID uint, currentUserID uint) ([]models.Card, error) {
	if _, err := s.checkAccessViaList(currentUserID, listID); err != nil {
		return nil, err
	}
	return s.cardRepo.FindByListID(listID)
}

func (s *CardService) UpdateCard(
	cardID uint,
	title, description *string,
	newPosition *uint,
	dueDate *time.Time,
	assignedUserID **uint, // Pointer to pointer for explicit null
	supervisorID **uint, // New field
	status *models.CardStatus, // New field
	color *string, // Add color
	currentUserID uint,
) (*models.Card, error) {
	_, listID, err := s.checkAccessViaCard(currentUserID, cardID) // Modified to get listID
	if err != nil {
		return nil, err
	}

	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}

	if title != nil {
		card.Title = *title
	}
	if description != nil {
		card.Description = *description
	}
	if dueDate != nil { // No double pointer, direct update or keep old
		card.DueDate = dueDate
	}
	if assignedUserID != nil {
		card.AssignedUserID = *assignedUserID
	}
	if supervisorID != nil { // New field
		card.SupervisorID = *supervisorID
	}
	if status != nil { // New field
		// Basic validation for status
		switch *status {
		case models.StatusToDo, models.StatusPending, models.StatusDone, models.StatusUndone:
			card.Status = *status
		default:
			return nil, fmt.Errorf("%w: invalid card status '%s'", ErrInvalidInput, *status)
		}
	}
	if color != nil { // Add color update
		if *color == "" { // Allow clearing the color
			card.Color = nil
		} else {
			card.Color = color
		}
	}

	// Handle position update within the same list
	if newPosition != nil && card.Position != *newPosition {
		currentPosition := card.Position
		targetPosition := *newPosition

		err = s.cardRepo.PerformTransaction(func(tx *gorm.DB) error {
			// Adjust positions of other cards
			if targetPosition < currentPosition { // Moving up
				if err := tx.Model(&models.Card{}).
					Where("list_id = ? AND position >= ? AND position < ?", listID, targetPosition, currentPosition).
					Where("id != ?", cardID). // Exclude the card being moved
					Update("position", gorm.Expr("position + 1")).Error; err != nil {
					return err
				}
			} else { // Moving down
				if err := tx.Model(&models.Card{}).
					Where("list_id = ? AND position > ? AND position <= ?", listID, currentPosition, targetPosition).
					Where("id != ?", cardID). // Exclude the card being moved
					Update("position", gorm.Expr("position - 1")).Error; err != nil {
					return err
				}
			}
			// Update the current card's position and other fields
			card.Position = targetPosition
			return tx.Save(card).Error
		})
		if err != nil {
			return nil, err
		}
	} else { // No position change, or only other fields (including color) updated
		if err := s.cardRepo.Update(card); err != nil {
			return nil, err
		}
	}
	return s.cardRepo.FindByID(card.ID)
}

func (s *CardService) DeleteCard(cardID uint, currentUserID uint) error {
	_, listID, err := s.checkAccessViaCard(currentUserID, cardID)
	if err != nil {
		return err
	}
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return ErrCardNotFound
	}

	return s.cardRepo.PerformTransaction(func(tx *gorm.DB) error {
		// Shift positions of subsequent cards in the same list
		if err := tx.Model(&models.Card{}).
			Where("list_id = ? AND position > ?", listID, card.Position).
			Update("position", gorm.Expr("position - 1")).Error; err != nil {
			return err
		}
		// Delete the card
		return s.cardRepo.Delete(cardID)
	})
}

func (s *CardService) MoveCard(cardID uint, targetListID uint, newPosition uint, currentUserID uint) (*models.Card, error) {
	_, originalListID, err := s.checkAccessViaCard(currentUserID, cardID)
	if err != nil {
		return nil, err
	}
	if _, err := s.checkAccessViaList(currentUserID, targetListID); err != nil { // Check access to target list
		return nil, err
	}

	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}

	// Validate newPosition (basic: >=1)
	if newPosition < 1 {
		return nil, ErrPositionOutOfBound
	}
	// More complex validation: newPosition <= max cards in target list + 1

	err = s.cardRepo.MoveCard(cardID, originalListID, targetListID, newPosition)
	if err != nil {
		return nil, err
	}

	return s.cardRepo.FindByID(card.ID) // Return updated card
}

// AddCollaboratorToCard adds a user as a collaborator to a card.
func (s *CardService) AddCollaboratorToCard(cardID uint, currentUserID uint, targetUserEmail string, targetUserIDInput *uint) (*models.User, error) {
	// Authorize current user can access the card (implicitly means they are on the board)
	if _, _, err := s.checkAccessViaCard(currentUserID, cardID); err != nil {
		return nil, err
	}

	var targetUser *models.User
	var err error

	if targetUserEmail != "" {
		targetUser, err = s.userRepo.FindByEmail(targetUserEmail)
	} else if targetUserIDInput != nil {
		targetUser, err = s.userRepo.FindByID(*targetUserIDInput)
	} else {
		return nil, fmt.Errorf("%w: either target user email or ID must be provided", ErrInvalidInput)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Check if user is already a collaborator
	isCollab, err := s.cardRepo.IsCollaborator(cardID, targetUser.ID)
	if err != nil {
		return nil, err
	}
	if isCollab {
		// Optionally, return specific error like ErrAlreadyCollaborator
		// For now, just returning the user as if successfully added (idempotency)
		return targetUser, nil
	}

	if err := s.cardRepo.AddCollaborator(cardID, targetUser.ID); err != nil {
		return nil, err
	}
	return targetUser, nil
}

// RemoveCollaboratorFromCard removes a collaborator from a card.
func (s *CardService) RemoveCollaboratorFromCard(cardID uint, currentUserID uint, targetUserID uint) error {
	if _, _, err := s.checkAccessViaCard(currentUserID, cardID); err != nil {
		return err // Ensure current user has access to the card's board
	}

	// Check if target user is actually a collaborator
	isCollab, err := s.cardRepo.IsCollaborator(cardID, targetUserID)
	if err != nil {
		return err
	}
	if !isCollab {
		return ErrUserNotCollaborator // New error to be defined
	}

	// Board owner or the collaborator themselves can remove.
	// (Additional permission logic can be added here if needed, e.g. only board owner or card creator)
	// For now, if currentUserID has access to the card (board member), they can remove anyone.
	// A stricter rule might be: board owner can remove anyone, card assignees/collaborators can remove themselves.

	return s.cardRepo.RemoveCollaborator(cardID, targetUserID)
}

// GetCardCollaborators retrieves all collaborators for a card.
func (s *CardService) GetCardCollaborators(cardID uint, currentUserID uint) ([]models.User, error) {
	if _, _, err := s.checkAccessViaCard(currentUserID, cardID); err != nil {
		return nil, err // Ensure current user has access
	}
	return s.cardRepo.GetCollaboratorsByCardID(cardID)
}
