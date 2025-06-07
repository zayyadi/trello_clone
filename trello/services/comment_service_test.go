package services

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"
	"gorm.io/gorm"
)

// --- MockCommentRepository ---
type MockCommentRepository struct {
	CreateFunc       func(comment *models.Comment) error
	FindByCardIDFunc func(cardID uint) ([]models.Comment, error)
	FindByIDFunc     func(id uint) (*models.Comment, error)

	CreateCalledWith *models.Comment
}

func (m *MockCommentRepository) Create(comment *models.Comment) error {
	m.CreateCalledWith = comment
	if m.CreateFunc != nil {
		return m.CreateFunc(comment)
	}
	return nil
}
func (m *MockCommentRepository) FindByCardID(cardID uint) ([]models.Comment, error) {
	if m.FindByCardIDFunc != nil {
		return m.FindByCardIDFunc(cardID)
	}
	return nil, errors.New("FindByCardIDFunc not implemented")
}
func (m *MockCommentRepository) FindByID(id uint) (*models.Comment, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc not implemented")
}

var _ repositories.CommentRepositoryInterface = (*MockCommentRepository)(nil)

// --- MockCardRepositoryForCommentService ---
type MockCardRepositoryForCommentService struct {
	repositories.CardRepositoryInterface
	GetListIDByCardIDFunc                func(cardID uint) (uint, error)
	FindByIDFunc                         func(id uint) (*models.Card, error)
	IsUserCollaboratorOrAssigneeFunc func(cardID uint, userID uint) (bool, error) // Added
}

func (m *MockCardRepositoryForCommentService) GetListIDByCardID(cardID uint) (uint, error) {
	if m.GetListIDByCardIDFunc != nil {
		return m.GetListIDByCardIDFunc(cardID)
	}
	return 0, errors.New("GetListIDByCardIDFunc not implemented")
}
func (m *MockCardRepositoryForCommentService) FindByID(id uint) (*models.Card, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc not implemented for MockCardRepositoryForCommentService")
}
func (m *MockCardRepositoryForCommentService) IsUserCollaboratorOrAssignee(cardID uint, userID uint) (bool, error) {
	if m.IsUserCollaboratorOrAssigneeFunc != nil {
		return m.IsUserCollaboratorOrAssigneeFunc(cardID, userID)
	}
	return false, errors.New("IsUserCollaboratorOrAssigneeFunc not implemented")
}


func (m *MockCardRepositoryForCommentService) Create(card *models.Card) error {
	return errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) FindByListID(listID uint) ([]models.Card, error) {
	return nil, errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) Update(card *models.Card) error {
	return errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) Delete(id uint) error {
	return errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) PerformTransaction(fn func(tx *gorm.DB) error) error {
	return errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) MoveCard(cardID, oldListID, newListID uint, newPosition uint) error {
	return errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) GetMaxPosition(listID uint) (uint, error) {
	return 0, errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) ShiftPositions(listID uint, startPosition uint, shiftAmount int, excludedCardID *uint) error {
	return errors.New("not implemented")
}

func (m *MockCardRepositoryForCommentService) AddCollaborator(cardID uint, userID uint) error {
	return errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) RemoveCollaborator(cardID uint, userID uint) error {
	return errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) GetCollaboratorsByCardID(cardID uint) ([]models.User, error) {
	return nil, errors.New("not implemented")
}
func (m *MockCardRepositoryForCommentService) IsCollaborator(cardID uint, userID uint) (bool, error) {
	return false, errors.New("not implemented")
}


// --- MockListRepositoryForCommentService ---
type MockListRepositoryForCommentService struct {
	repositories.ListRepositoryInterface
	GetBoardIDByListIDFunc               func(listID uint) (uint, error)
}

func (m *MockListRepositoryForCommentService) GetBoardIDByListID(listID uint) (uint, error) {
	if m.GetBoardIDByListIDFunc != nil {
		return m.GetBoardIDByListIDFunc(listID)
	}
	return 0, errors.New("GetBoardIDByListIDFunc not implemented")
}
func (m *MockListRepositoryForCommentService) Create(list *models.List) error {
	return errors.New("not implemented")
}
func (m *MockListRepositoryForCommentService) FindByID(id uint) (*models.List, error) {
	return nil, errors.New("not implemented")
}
func (m *MockListRepositoryForCommentService) FindByBoardID(boardID uint) ([]models.List, error) {
	return nil, errors.New("not implemented")
}
func (m *MockListRepositoryForCommentService) Update(list *models.List) error {
	return errors.New("not implemented")
}
func (m *MockListRepositoryForCommentService) Delete(id uint) error {
	return errors.New("not implemented")
}
func (m *MockListRepositoryForCommentService) GetMaxPosition(boardID uint) (uint, error) {
	return 0, errors.New("not implemented")
}
func (m *MockListRepositoryForCommentService) GetDB() *gorm.DB { return nil }
func (m *MockListRepositoryForCommentService) PerformTransaction(fn func(tx *gorm.DB) error) error { return errors.New("not implemented") }


// --- MockBoardRepositoryForCommentService ---
type MockBoardRepositoryForCommentService struct {
	repositories.BoardRepositoryInterface
	FindByIDFunc func(id uint) (*models.Board, error)
}

func (m *MockBoardRepositoryForCommentService) FindByID(id uint) (*models.Board, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc on MockBoardRepositoryForCommentService not implemented")
}
func (m *MockBoardRepositoryForCommentService) Create(board *models.Board) error {
	return errors.New("not implemented")
}
func (m *MockBoardRepositoryForCommentService) FindByOwnerOrMember(userID uint) ([]models.Board, error) {
	return nil, errors.New("not implemented")
}
func (m *MockBoardRepositoryForCommentService) Update(board *models.Board) error {
	return errors.New("not implemented")
}
func (m *MockBoardRepositoryForCommentService) Delete(id uint) error {
	return errors.New("not implemented")
}
func (m *MockBoardRepositoryForCommentService) IsOwner(boardID uint, userID uint) (bool, error) {
	return false, errors.New("not implemented")
}

// --- MockBoardMemberRepositoryForCommentService ---
type MockBoardMemberRepositoryForCommentService struct {
	repositories.BoardMemberRepositoryInterface
	IsMemberFunc func(boardID uint, userID uint) (bool, error)
}

func (m *MockBoardMemberRepositoryForCommentService) IsMember(boardID uint, userID uint) (bool, error) {
	if m.IsMemberFunc != nil {
		return m.IsMemberFunc(boardID, userID)
	}
	return false, errors.New("IsMemberFunc on MockBoardMemberRepositoryForCommentService not implemented")
}
func (m *MockBoardMemberRepositoryForCommentService) AddMember(member *models.BoardMember) error {
	return errors.New("not implemented")
}
func (m *MockBoardMemberRepositoryForCommentService) RemoveMember(boardID uint, userID uint) error {
	return errors.New("not implemented")
}
func (m *MockBoardMemberRepositoryForCommentService) FindMembersByBoardID(boardID uint) ([]models.BoardMember, error) {
	return nil, errors.New("not implemented")
}
func (m *MockBoardMemberRepositoryForCommentService) FindByBoardIDAndUserID(boardID uint, userID uint) (*models.BoardMember, error) {
	return nil, errors.New("not implemented")
}

func TestCommentService_CreateComment_Success(t *testing.T) {
	mockCommentRepo := &MockCommentRepository{}
	mockCardRepo := &MockCardRepositoryForCommentService{}
	mockListRepo := &MockListRepositoryForCommentService{}
	mockBoardRepo := &MockBoardRepositoryForCommentService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCommentService{}

	commentService := NewCommentService(
		mockCommentRepo, mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo,
	)

	cardID := uint(1)
	userID := uint(2)
	boardID := uint(3)
	listID := uint(4)
	content := "This is a test comment"

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) {
		assert.Equal(t, cardID, cID)
		return listID, nil
	}
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) {
		assert.Equal(t, listID, lID)
		return boardID, nil
	}
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		assert.Equal(t, boardID, bID)
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil
	}

	var createdCommentID uint = 100
	mockCommentRepo.CreateFunc = func(comment *models.Comment) error {
		assert.Equal(t, content, comment.Content)
		assert.Equal(t, cardID, comment.CardID)
		assert.Equal(t, userID, comment.UserID)
		comment.ID = createdCommentID
		comment.CreatedAt = time.Now()
		comment.UpdatedAt = time.Now()
		return nil
	}

	mockCommentRepo.FindByIDFunc = func(id uint) (*models.Comment, error) {
		assert.Equal(t, createdCommentID, id)
		return &models.Comment{
			Model:   gorm.Model{ID: id, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Content: content,
			CardID:  cardID,
			UserID:  userID,
			User:    models.User{Model: gorm.Model{ID: userID}, Username: "testuser"},
		}, nil
	}

	comment, err := commentService.CreateComment(cardID, userID, content)

	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, createdCommentID, comment.ID)
	assert.Equal(t, content, comment.Content)
	assert.Equal(t, cardID, comment.CardID)
	assert.Equal(t, userID, comment.UserID)
	assert.NotNil(t, comment.User)
	assert.Equal(t, userID, comment.User.ID)
	assert.Equal(t, "testuser", comment.User.Username)

	assert.NotNil(t, mockCommentRepo.CreateCalledWith)
	if mockCommentRepo.CreateCalledWith != nil {
		assert.Equal(t, content, mockCommentRepo.CreateCalledWith.Content)
	}
}

func TestCommentService_CreateComment_ContentEmpty(t *testing.T) {
	mockCommentRepo := &MockCommentRepository{}
	mockCardRepo := &MockCardRepositoryForCommentService{}
	mockListRepo := &MockListRepositoryForCommentService{}
	mockBoardRepo := &MockBoardRepositoryForCommentService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCommentService{}

	commentService := NewCommentService(
		mockCommentRepo, mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo,
	)

	cardID := uint(1)
	userID := uint(2)
	boardID := uint(3)
	listID := uint(4)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil
	}

	comment, err := commentService.CreateComment(cardID, userID, "")

	assert.Error(t, err)
	assert.Equal(t, "comment content cannot be empty", err.Error())
	assert.Nil(t, comment)
}

func TestCommentService_CreateComment_CardNotFound(t *testing.T) {
	mockCommentRepo := &MockCommentRepository{}
	mockCardRepo := &MockCardRepositoryForCommentService{}
	mockListRepo := &MockListRepositoryForCommentService{}
	mockBoardRepo := &MockBoardRepositoryForCommentService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCommentService{}

	commentService := NewCommentService(
		mockCommentRepo, mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo,
	)

	cardID := uint(1)
	userID := uint(2)
	content := "Test comment"

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) {
		assert.Equal(t, cardID, cID)
		return 0, gorm.ErrRecordNotFound
	}

	comment, err := commentService.CreateComment(cardID, userID, content)

	assert.Error(t, err)
	assert.Equal(t, ErrCardNotFound, err)
	assert.Nil(t, comment)
}

func TestCommentService_CreateComment_PermissionDenied_Forbidden(t *testing.T) {
	mockCommentRepo := &MockCommentRepository{}
	mockCardRepo := &MockCardRepositoryForCommentService{}
	mockListRepo := &MockListRepositoryForCommentService{}
	mockBoardRepo := &MockBoardRepositoryForCommentService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCommentService{}

	commentService := NewCommentService(
		mockCommentRepo, mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo,
	)

	cardID := uint(1)
	userID := uint(2)
	boardID := uint(3)
	listID := uint(4)
	actualOwnerID := uint(5)
	content := "Test comment"

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, userID, uID)
		return false, nil
	}

	comment, err := commentService.CreateComment(cardID, userID, content)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, comment)
}

func TestCommentService_GetCommentsByCardID_Success(t *testing.T) {
	mockCommentRepo := &MockCommentRepository{}
	mockCardRepo := &MockCardRepositoryForCommentService{}
	mockListRepo := &MockListRepositoryForCommentService{}
	mockBoardRepo := &MockBoardRepositoryForCommentService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCommentService{}

	commentService := NewCommentService(
		mockCommentRepo, mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo,
	)

	cardID := uint(1)
	userID := uint(2)
	boardID := uint(3)
	listID := uint(4)

	expectedComments := []models.Comment{
		{Model: gorm.Model{ID: 100}, Content: "First comment", CardID: cardID, UserID: userID, User: models.User{Model: gorm.Model{ID: userID}, Username: "testuser"}},
		{Model: gorm.Model{ID: 101}, Content: "Second comment", CardID: cardID, UserID: uint(6), User: models.User{Model: gorm.Model{ID: uint(6)}, Username: "anotheruser"}},
	}

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil
	}

	mockCommentRepo.FindByCardIDFunc = func(cID uint) ([]models.Comment, error) {
		assert.Equal(t, cardID, cID)
		return expectedComments, nil
	}

	comments, err := commentService.GetCommentsByCardID(cardID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.Len(t, comments, 2)
	assert.Equal(t, expectedComments, comments)
}

func TestCommentService_GetCommentsByCardID_NoComments(t *testing.T) {
	mockCommentRepo := &MockCommentRepository{}
	mockCardRepo := &MockCardRepositoryForCommentService{}
	mockListRepo := &MockListRepositoryForCommentService{}
	mockBoardRepo := &MockBoardRepositoryForCommentService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCommentService{}

	commentService := NewCommentService(
		mockCommentRepo, mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo,
	)

	cardID := uint(1)
	userID := uint(2)
	boardID := uint(3)
	listID := uint(4)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil
	}

	mockCommentRepo.FindByCardIDFunc = func(cID uint) ([]models.Comment, error) {
		assert.Equal(t, cardID, cID)
		return []models.Comment{}, nil
	}

	comments, err := commentService.GetCommentsByCardID(cardID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.Len(t, comments, 0)
}

func TestCommentService_GetCommentsByCardID_PermissionDenied_Forbidden(t *testing.T) {
	mockCommentRepo := &MockCommentRepository{}
	mockCardRepo := &MockCardRepositoryForCommentService{}
	mockListRepo := &MockListRepositoryForCommentService{}
	mockBoardRepo := &MockBoardRepositoryForCommentService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCommentService{}

	commentService := NewCommentService(
		mockCommentRepo, mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo,
	)

	cardID := uint(1)
	userID := uint(2)
	boardID := uint(3)
	listID := uint(4)
	actualOwnerID := uint(5)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, nil
	}

	var findByCardIDCalled bool
	mockCommentRepo.FindByCardIDFunc = func(cID uint) ([]models.Comment, error) {
		findByCardIDCalled = true
		return nil, nil
	}

	comments, err := commentService.GetCommentsByCardID(cardID, userID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, comments)
	assert.False(t, findByCardIDCalled)
}

func TestCommentService_GetCommentsByCardID_RepoError(t *testing.T) {
	mockCommentRepo := &MockCommentRepository{}
	mockCardRepo := &MockCardRepositoryForCommentService{}
	mockListRepo := &MockListRepositoryForCommentService{}
	mockBoardRepo := &MockBoardRepositoryForCommentService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCommentService{}

	commentService := NewCommentService(
		mockCommentRepo, mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo,
	)

	cardID := uint(1)
	userID := uint(2)
	boardID := uint(3)
	listID := uint(4)
	expectedError := errors.New("DB error on FindByCardID")

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil
	}

	mockCommentRepo.FindByCardIDFunc = func(cID uint) ([]models.Comment, error) {
		return nil, expectedError
	}

	comments, err := commentService.GetCommentsByCardID(cardID, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, comments)
}
