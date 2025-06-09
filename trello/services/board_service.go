package services

import (
	"errors"

	"github.com/zayyadi/trello/handlers" // For DTOs
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"
	"github.com/zayyadi/trello/realtime"

	"gorm.io/gorm"
)

type BoardService struct {
	boardRepo       repositories.BoardRepositoryInterface
	userRepo        repositories.UserRepositoryInterface
	boardMemberRepo repositories.BoardMemberRepositoryInterface
	hub             *realtime.Hub
}

// BoardServiceInterface defines methods for board service (including IsUserMemberOfBoard)
// This could be a more complete interface if needed elsewhere, or just define the methods used.
type BoardServiceInterface interface {
	CreateBoard(name, description string, ownerID uint) (*models.Board, error)
	GetBoardByID(boardID, userID uint) (*models.Board, error)
	GetBoardsForUser(userID uint) ([]models.Board, error)
	UpdateBoard(boardID uint, name, description *string, userID uint) (*models.Board, error)
	DeleteBoard(boardID, userID uint) error
	AddMemberToBoard(boardID uint, email *string, memberUserID *uint, currentUserID uint) (*models.BoardMember, error)
	RemoveMemberFromBoard(boardID, memberUserID, currentUserID uint) error
	GetBoardMembers(boardID, currentUserID uint) ([]models.BoardMember, error)
	IsUserMemberOfBoard(userID uint, boardID uint) (bool, error) // New method
}

func NewBoardService(
	boardRepo repositories.BoardRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	boardMemberRepo repositories.BoardMemberRepositoryInterface,
	hub *realtime.Hub,
) BoardServiceInterface { // Return interface type
	return &BoardService{
		boardRepo:       boardRepo,
		userRepo:        userRepo,
		boardMemberRepo: boardMemberRepo,
		hub:             hub,
	}
}

func (s *BoardService) CreateBoard(name, description string, ownerID uint) (*models.Board, error) {
	board := &models.Board{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}
	var err error // Declare err once
	if err = s.boardRepo.Create(board); err != nil {
		return nil, err
	}
	// Add owner as a member automatically
	if err = s.boardMemberRepo.AddMember(&models.BoardMember{BoardID: board.ID, UserID: ownerID}); err != nil {
		// Log error but don't fail board creation if member addition fails
		// Or, implement rollback for board creation if member addition is critical
		return nil, err
	}

	createdBoard, err := s.boardRepo.FindByID(board.ID) // Fetch with owner preloaded
	if err != nil {
		return nil, err
	}

	// Broadcast board creation
	broadcastMessage(
		s.hub,
		createdBoard.ID,
		realtime.MessageTypeBoardCreated,
		handlers.MapBoardToBoardResponse(createdBoard), // Use existing DTO mapper
		ownerID,
	)

	return createdBoard
}

func (s *BoardService) GetBoardByID(boardID, userID uint) (*models.Board, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBoardNotFound
		}
		return nil, err
	}
	if board.OwnerID == userID {
		return board, nil
	}
	isMember, err := s.boardMemberRepo.IsMember(boardID, userID)
	if err != nil || !isMember {
		return nil, ErrForbidden
	}
	return board, nil
}

func (s *BoardService) GetBoardsForUser(userID uint) ([]models.Board, error) {
	boards, err := s.boardRepo.FindByOwnerOrMember(userID)
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func (s *BoardService) UpdateBoard(boardID uint, name, description *string, userID uint) (*models.Board, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBoardNotFound
		}
		return nil, err
	}
	if board.OwnerID != userID {
		return nil, ErrForbidden // Only owner can update board details
	}

	if name != nil {
		board.Name = *name
	}
	if description != nil {
		board.Description = *description
	}

	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}
	updatedBoard, err := s.boardRepo.FindByID(board.ID) // Fetch updated board
	if err != nil {
		// Log or handle error if fetching fails, but board was updated.
		// For now, return the error, but consider if partial success should be handled.
		return nil, err
	}

	// Broadcast board update
	broadcastMessage(
		s.hub,
		updatedBoard.ID,
		realtime.MessageTypeBoardUpdated,
		handlers.MapBoardToBoardResponse(updatedBoard),
		userID,
	)
	return updatedBoard
}

func (s *BoardService) DeleteBoard(boardID, userID uint) error {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBoardNotFound
		}
		return err
	}
	if board.OwnerID != userID {
		return ErrForbidden // Only owner can delete board
	}
	err = s.boardRepo.Delete(boardID)
	if err == nil {
		// Broadcast board deletion
		broadcastMessage(
			s.hub,
			boardID,
			realtime.MessageTypeBoardDeleted,
			realtime.BoardBasicInfo{ID: boardID}, // Simple payload for deletion
			userID,
		)
	}
	return err
}

func (s *BoardService) AddMemberToBoard(boardID uint, email *string, memberUserID *uint, currentUserID uint) (*models.BoardMember, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBoardNotFound
		}
		return nil, err
	}
	if board.OwnerID != currentUserID {
		return nil, ErrForbidden // Only owner can add members
	}

	var targetUserID uint
	if email != nil && *email != "" {
		user, err := s.userRepo.FindByEmail(*email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrUserNotFound
			}
			return nil, err
		}
		targetUserID = user.ID
	} else if memberUserID != nil && *memberUserID != 0 {
		user, err := s.userRepo.FindByID(*memberUserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrUserNotFound
			}
			return nil, err
		}
		targetUserID = user.ID
	} else {
		return nil, errors.New("Either email or userID must be provided")
	}

	if targetUserID == board.OwnerID {
		return nil, ErrUserAlreadyMember // Owner is implicitly a member
	}

	isMember, err := s.boardMemberRepo.IsMember(boardID, targetUserID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, ErrUserAlreadyMember
	}

	boardMember := &models.BoardMember{
		BoardID: boardID,
		UserID:  targetUserID,
	}
	if err = s.boardMemberRepo.AddMember(boardMember); err != nil {
		return nil, err
	}
	addedMember, err := s.boardMemberRepo.FindByBoardIDAndUserID(boardID, targetUserID) // Fetch with user preloaded
	if err != nil {
		return nil, err
	}

	// Broadcast member addition
	broadcastMessage(
		s.hub,
		boardID,
		realtime.MessageTypeBoardMemberAdded,
		handlers.MapBoardMemberToResponse(addedMember), // Use existing DTO mapper
		currentUserID,
	)
	return addedMember
}

func (s *BoardService) RemoveMemberFromBoard(boardID, memberUserID, currentUserID uint) error {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBoardNotFound
		}
		return err
	}
	if board.OwnerID != currentUserID {
		return ErrForbidden // Only owner can remove members
	}
	if memberUserID == board.OwnerID {
		return ErrCannotRemoveOwner // Cannot remove the owner
	}

	isMember, err := s.boardMemberRepo.IsMember(boardID, memberUserID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrBoardMemberNotFound
	}

	err = s.boardMemberRepo.RemoveMember(boardID, memberUserID)
	if err == nil {
		// Broadcast member removal
		broadcastMessage(
			s.hub,
			boardID,
			realtime.MessageTypeBoardMemberRemoved,
			realtime.BoardMemberPayload{BoardID: boardID, UserID: memberUserID}, // Simple payload
			currentUserID,
		)
	}
	return err
}

func (s *BoardService) GetBoardMembers(boardID, currentUserID uint) ([]models.BoardMember, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBoardNotFound
		}
		return nil, err
	}
	if board.OwnerID != currentUserID {
		isMember, err := s.boardMemberRepo.IsMember(boardID, currentUserID)
		if err != nil || !isMember {
			return nil, ErrForbidden // Only owner or member can view members
		}
	}
	members, err := s.boardMemberRepo.FindMembersByBoardID(boardID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// IsUserMemberOfBoard checks if a user is the owner or an explicit member of the board.
func (s *BoardService) IsUserMemberOfBoard(userID uint, boardID uint) (bool, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrBoardNotFound // Or just return false, nil if board not found means not a member
		}
		return false, err
	}

	// Check if the user is the owner of the board
	if board.OwnerID == userID {
		return true, nil
	}

	// Check if the user is listed in the board_members table
	isMember, err := s.boardMemberRepo.IsMember(boardID, userID)
	if err != nil {
		return false, err
	}
	return isMember, nil
}
