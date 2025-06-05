package services

import (
	"errors"

	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"

	"gorm.io/gorm"
)

type BoardService struct {
	boardRepo       repositories.BoardRepositoryInterface
	userRepo        repositories.UserRepositoryInterface
	boardMemberRepo repositories.BoardMemberRepositoryInterface
}

func NewBoardService(
	boardRepo repositories.BoardRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	boardMemberRepo repositories.BoardMemberRepositoryInterface,
) *BoardService {
	return &BoardService{
		boardRepo:       boardRepo,
		userRepo:        userRepo,
		boardMemberRepo: boardMemberRepo,
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
	return s.boardRepo.FindByID(board.ID) // Fetch with owner preloaded
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
	return s.boardRepo.FindByID(board.ID) // Fetch updated board
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
	return s.boardRepo.Delete(boardID)
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
	return s.boardMemberRepo.FindByBoardIDAndUserID(boardID, targetUserID) // Fetch with user preloaded
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

	return s.boardMemberRepo.RemoveMember(boardID, memberUserID)
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
