package services

import (
	"errors"

	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"
	"gorm.io/gorm"
)

// CommentServiceInterface defines the contract for comment service operations.
type CommentServiceInterface interface {
	CreateComment(cardID uint, userID uint, content string) (*models.Comment, error)
	GetCommentsByCardID(cardID uint, userID uint) ([]models.Comment, error)
}

// CommentService handles business logic related to comments.
type CommentService struct {
	commentRepo     repositories.CommentRepositoryInterface
	cardRepo        repositories.CardRepositoryInterface
	listRepo        repositories.ListRepositoryInterface
	boardRepo       repositories.BoardRepositoryInterface
	boardMemberRepo repositories.BoardMemberRepositoryInterface
}

// NewCommentService creates a new CommentService.
func NewCommentService(
	commentRepo repositories.CommentRepositoryInterface,
	cardRepo repositories.CardRepositoryInterface,
	listRepo repositories.ListRepositoryInterface,
	boardRepo repositories.BoardRepositoryInterface,
	boardMemberRepo repositories.BoardMemberRepositoryInterface,
) CommentServiceInterface {
	return &CommentService{
		commentRepo:     commentRepo,
		cardRepo:        cardRepo,
		listRepo:        listRepo,
		boardRepo:       boardRepo,
		boardMemberRepo: boardMemberRepo,
	}
}

// Helper function to check if a user has access to the board a card belongs to.
func (s *CommentService) checkCardBoardAccess(userID uint, cardID uint) error {
	listID, err := s.cardRepo.GetListIDByCardID(cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCardNotFound // Card implies list, implies board
		}
		return err
	}

	boardID, err := s.listRepo.GetBoardIDByListID(listID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrListNotFound // List implies board
		}
		return err
	}

	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBoardNotFound
		}
		return err
	}

	if board.OwnerID == userID {
		return nil // User is the owner of the board
	}

	isMember, err := s.boardMemberRepo.IsMember(boardID, userID)
	if err != nil {
		return err // DB error during IsMember check
	}
	if !isMember {
		return ErrForbidden // User is not owner and not a member
	}

	return nil // User has access
}

// CreateComment creates a new comment on a card.
func (s *CommentService) CreateComment(cardID uint, userID uint, content string) (*models.Comment, error) {
	if err := s.checkCardBoardAccess(userID, cardID); err != nil {
		return nil, err
	}

	if content == "" {
		return nil, errors.New("comment content cannot be empty")
	}

	comment := &models.Comment{
		Content: content,
		CardID:  cardID,
		UserID:  userID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// To return the comment with User preloaded, we might need a FindByID in CommentRepository
	// For now, commentRepo.Create likely doesn't preload. If it does, great.
	// If not, and preloading is desired here, an extra FindByID call would be needed.
	// Let's assume Create does not preload User by default.
	// To ensure User is preloaded for the response, fetch it:
	// This assumes CommentRepository has a FindByID that preloads User.
	// If not, the FindByCardID preloads User, so the created comment might not have it
	// unless the Create method in the repo is updated or we fetch it here.
	// For simplicity, returning the comment as is from Create.
	// The FindByCardID will ensure User is preloaded for GET requests.
	// If Comment model's User field is needed immediately after Create, this would need enhancement.
	// Let's fetch the user separately for now to populate.
	// NOTE: This is not ideal; ideally, Create would return the preloaded user or FindByID would be called.
	// But to avoid adding FindByID to CommentRepoInterface right now:
	// --- Correction: FindByID is now added to the interface and repo ---
	// Fetch the created comment with preloaded User details.
	return s.commentRepo.FindByID(comment.ID)
}

// GetCommentsByCardID retrieves all comments for a given card.
func (s *CommentService) GetCommentsByCardID(cardID uint, userID uint) ([]models.Comment, error) {
	if err := s.checkCardBoardAccess(userID, cardID); err != nil {
		return nil, err
	}

	return s.commentRepo.FindByCardID(cardID)
}
