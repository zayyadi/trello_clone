package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/zayyadi/trello/dto" // Changed import
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"
	"github.com/zayyadi/trello/realtime"

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
	hub             *realtime.Hub
}

func NewCardService(
	cardRepo repositories.CardRepositoryInterface,
	listRepo repositories.ListRepositoryInterface,
	boardRepo repositories.BoardRepositoryInterface,
	boardMemberRepo repositories.BoardMemberRepositoryInterface,
	userRepo repositories.UserRepositoryInterface, // Added
	hub *realtime.Hub,
) CardServiceInterface { // Return interface type
	return &CardService{
		cardRepo:        cardRepo,
		listRepo:        listRepo,
		boardRepo:       boardRepo,
		boardMemberRepo: boardMemberRepo,
		userRepo:        userRepo, // Added
		hub:             hub,
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
	boardID, err := s.checkAccessViaList(currentUserID, listID)
	if err != nil {
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

	createdCard, err := s.cardRepo.FindByID(card.ID) // Will now preload supervisor too
	if err != nil {
		return nil, err
	}

	// Broadcast card creation
	broadcastMessage(
		s.hub,
		boardID, // Obtained from checkAccessViaList
		realtime.MessageTypeCardCreated,
		dto.MapCardToResponse(createdCard, true), // Use dto mapper
		currentUserID,
	)

	return createdCard, nil
}

func (s *CardService) GetCardByID(cardID uint, currentUserID uint) (*models.Card, error) {
	boardID, _, err := s.checkAccessViaCard(currentUserID, cardID)
	if err != nil {
		return nil, err
	}

	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBoardNotFound
		}
		return nil, err
	}

	isOwner := (board.OwnerID == currentUserID)
	isCollaboratorOrAssignee, collabErr := s.cardRepo.IsUserCollaboratorOrAssignee(cardID, currentUserID)
	if collabErr != nil && !errors.Is(collabErr, gorm.ErrRecordNotFound) { // Allow if card not found for IsCollaboratorOrAssignee if it implies no relations
		return nil, collabErr
	}

	if !isOwner && !isCollaboratorOrAssignee {
		return nil, ErrForbidden
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
	boardID, listID, err := s.checkAccessViaCard(currentUserID, cardID)
	if err != nil {
		return nil, err
	}

	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBoardNotFound
		}
		return nil, err
	}

	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}

	isOwner := (board.OwnerID == currentUserID)
	isCollaboratorOrAssignee, collabErr := s.cardRepo.IsUserCollaboratorOrAssignee(cardID, currentUserID)
	if collabErr != nil && !errors.Is(collabErr, gorm.ErrRecordNotFound) {
		return nil, collabErr // Propagate actual DB errors
	}

	if title != nil {
		if !isOwner {
			return nil, ErrPermissionDenied
		}
		card.Title = *title
	}
	if description != nil {
		if !isOwner && !isCollaboratorOrAssignee {
			return nil, ErrPermissionDenied
		}
		card.Description = *description
	}
	if dueDate != nil { // No double pointer, direct update or keep old
		if !isOwner && !isCollaboratorOrAssignee {
			return nil, ErrPermissionDenied
		}
		card.DueDate = dueDate
	}
	if assignedUserID != nil {
		if !isOwner && !isCollaboratorOrAssignee {
			return nil, ErrPermissionDenied
		}
		card.AssignedUserID = *assignedUserID
	}
	if supervisorID != nil { // New field
		if !isOwner && !isCollaboratorOrAssignee {
			return nil, ErrPermissionDenied
		}
		card.SupervisorID = *supervisorID
	}
	if status != nil { // New field
		if !isOwner && !isCollaboratorOrAssignee {
			return nil, ErrPermissionDenied
		}
		// Basic validation for status
		switch *status {
		case models.StatusToDo, models.StatusPending, models.StatusDone, models.StatusUndone:
			card.Status = *status
		default:
			return nil, fmt.Errorf("%w: invalid card status '%s'", ErrInvalidInput, *status)
		}
	}
	if color != nil { // Add color update
		if !isOwner && !isCollaboratorOrAssignee {
			return nil, ErrPermissionDenied
		}
		if *color == "" { // Allow clearing the color
			card.Color = nil
		} else {
			card.Color = color
		}
	}

	// Handle position update within the same list
	if newPosition != nil && card.Position != *newPosition {
		if !isOwner && !isCollaboratorOrAssignee {
			return nil, ErrPermissionDenied
		}
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

	updatedCard, err := s.cardRepo.FindByID(card.ID)
	if err != nil {
		return nil, err
	}

	// Broadcast card update
	broadcastMessage(
		s.hub,
		boardID, // Obtained from checkAccessViaCard
		realtime.MessageTypeCardUpdated,
		dto.MapCardToResponse(updatedCard, true), // Use dto mapper
		currentUserID,
	)

	// Handle assignment/unassignment messages
	if assignedUserID != nil {
		if *assignedUserID == nil { // Unassigned
			broadcastMessage(s.hub, boardID, realtime.MessageTypeCardUnassigned, realtime.CardBasicInfo{ID: cardID, ListID: listID, BoardID: boardID}, currentUserID)
		} else { // Assigned or changed assignee
			// We might want to include more assignee details in the payload here
			assigneePayload := struct {
				CardID    uint  `json:"cardId"`
				UserID    uint  `json:"userId"`
				BoardID   uint  `json:"boardId"`
				ListID    uint  `json:"listId"`
			}{cardID, **assignedUserID, boardID, listID}
			broadcastMessage(s.hub, boardID, realtime.MessageTypeCardAssigned, assigneePayload, currentUserID)
		}
	}


	return updatedCard, nil
}

func (s *CardService) DeleteCard(cardID uint, currentUserID uint) error {
	boardID, listID, err := s.checkAccessViaCard(currentUserID, cardID)
	if err != nil {
		return err
	}

	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBoardNotFound
		}
		return err
	}

	if board.OwnerID != currentUserID {
		return ErrForbidden
	}

	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return ErrCardNotFound
	}

	err = s.cardRepo.PerformTransaction(func(tx *gorm.DB) error {
		// Shift positions of subsequent cards in the same list
		if err := tx.Model(&models.Card{}).
			Where("list_id = ? AND position > ?", listID, card.Position).
			Update("position", gorm.Expr("position - 1")).Error; err != nil {
			return err
		}
		// Delete the card
		return tx.Delete(&models.Card{}, cardID).Error
	})

	if err == nil {
		// Broadcast card deletion
		broadcastMessage(
			s.hub,
			boardID, // Obtained from checkAccessViaCard
			realtime.MessageTypeCardDeleted,
			realtime.CardBasicInfo{ID: cardID, ListID: listID, BoardID: boardID},
			currentUserID,
		)
	}
	return err
}

func (s *CardService) MoveCard(cardID uint, targetListID uint, newPosition uint, currentUserID uint) (*models.Card, error) {
	boardID, originalListID, err := s.checkAccessViaCard(currentUserID, cardID)
	if err != nil {
		return nil, err
	}
	targetBoardID, err := s.checkAccessViaList(currentUserID, targetListID) // Check access to target list
	if err != nil {
		return nil, err
	}

	if boardID != targetBoardID {
		// Moving card across boards is not supported by this message structure easily,
		// as broadcasting is per-board. This implies an error or more complex handling needed.
		// For now, assume targetListID is on the same board.
		// If they can differ, the broadcasting logic needs to be revisited (e.g., two messages).
		return nil, errors.New("moving card to a different board is not directly supported for simple broadcast")
	}


	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, ErrCardNotFound
	}
	originalPosition := card.Position // Capture original position before move

	// Validate newPosition (basic: >=1)
	if newPosition < 1 {
		return nil, ErrPositionOutOfBound
	}
	// More complex validation: newPosition <= max cards in target list + 1

	err = s.cardRepo.MoveCard(cardID, originalListID, targetListID, newPosition)
	if err != nil {
		return nil, err
	}

	movedCard, err := s.cardRepo.FindByID(card.ID)
	if err != nil {
		return nil, err
	}

	// Broadcast card move
	payload := realtime.CardMovedPayload{
		CardID:      cardID,
		OldListID:   originalListID,
		NewListID:   targetListID,
		OldPosition: originalPosition,
		NewPosition: movedCard.Position, // Use the final position from the moved card
		BoardID:     boardID,
		// UpdatedCards could be populated here if the MoveCard repo method returned them
	}
	broadcastMessage(
		s.hub,
		boardID,
		realtime.MessageTypeCardMoved,
		payload,
		currentUserID,
	)

	return movedCard, nil
}

// AddCollaboratorToCard adds a user as a collaborator to a card.
func (s *CardService) AddCollaboratorToCard(cardID uint, currentUserID uint, targetUserEmail string, targetUserIDInput *uint) (*models.User, error) {
	boardID, _, err := s.checkAccessViaCard(currentUserID, cardID)
	if err != nil {
		return nil, err
	}

	board, err := s.boardRepo.FindByID(boardID) // Note: This was the line with the error, ensure 'err' is properly assigned.
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBoardNotFound
		}
		return nil, err
	}

	if board.OwnerID != currentUserID {
		return nil, ErrForbidden
	}

	var targetUser *models.User
	var terr error // Use a different error variable name for this scope

	if targetUserEmail != "" {
		targetUser, terr = s.userRepo.FindByEmail(targetUserEmail)
	} else if targetUserIDInput != nil {
		targetUser, terr = s.userRepo.FindByID(*targetUserIDInput)
	} else {
		return nil, fmt.Errorf("%w: either target user email or ID must be provided", ErrInvalidInput)
	}

	if terr != nil {
		if errors.Is(terr, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, terr
	}

	// Check if user is already a collaborator
	isCollab, terr := s.cardRepo.IsCollaborator(cardID, targetUser.ID)
	if terr != nil {
		return nil, terr
	}

	if !isCollab {
		if terr = s.cardRepo.AddCollaborator(cardID, targetUser.ID); terr != nil {
			return nil, terr
		}
	}

	// Broadcast collaborator addition (even if already a collab, for idempotency or client sync)
	// Or only broadcast if !isCollab
	if !isCollab { // Only broadcast if it's a new addition
		collabPayload := realtime.CardCollaboratorPayload{
			CardID:   cardID,
			UserID:   targetUser.ID,
			BoardID:  boardID,
			UserName: targetUser.Username, // Corrected line
		}
		broadcastMessage(
			s.hub,
			boardID,
			realtime.MessageTypeCardCollaboratorAdded,
			collabPayload,
			currentUserID,
		)
	}
	return targetUser, nil
}

// RemoveCollaboratorFromCard removes a collaborator from a card.
func (s *CardService) RemoveCollaboratorFromCard(cardID uint, currentUserID uint, targetUserID uint) error {
	boardID, _, err := s.checkAccessViaCard(currentUserID, cardID)
	if err != nil {
		return err // Ensure current user has access to the card's board
	}

	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBoardNotFound // Return specific error
		}
		return err
	}

	if board.OwnerID != currentUserID {
		return ErrForbidden
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

	err = s.cardRepo.RemoveCollaborator(cardID, targetUserID)
	if err == nil {
		// Broadcast collaborator removal
		collabPayload := realtime.CardCollaboratorPayload{
			CardID:  cardID,
			UserID:  targetUserID,
			BoardID: boardID,
			// UserName might not be easily available here post-removal without another query
		}
		broadcastMessage(
			s.hub,
			boardID,
			realtime.MessageTypeCardCollaboratorRemoved,
			collabPayload,
			currentUserID,
		)
	}
	return err
}

// GetCardCollaborators retrieves all collaborators for a card.
func (s *CardService) GetCardCollaborators(cardID uint, currentUserID uint) ([]models.User, error) {
	if _, _, err := s.checkAccessViaCard(currentUserID, cardID); err != nil {
		return nil, err // Ensure current user has access
	}
	return s.cardRepo.GetCollaboratorsByCardID(cardID)
}
