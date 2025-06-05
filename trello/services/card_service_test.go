package services

import (
	"errors"
	"testing"
	"time" // Added for DueDate tests

	"github.com/stretchr/testify/assert"
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"

	// "gorm.io/driver/sqlite" // Removed as setupTestDB is now centralized
	"gorm.io/gorm"
	// "gorm.io/gorm/logger" // Removed as setupTestDB is now centralized
)

// setupTestDB is now in test_helpers_test.go

// --- MockCardRepository ---
type MockCardRepository struct {
	repositories.CardRepositoryInterface
	CreateFunc                   func(card *models.Card) error
	FindByIDFunc                 func(id uint) (*models.Card, error)
	FindByListIDFunc             func(listID uint) ([]models.Card, error)
	UpdateFunc                   func(card *models.Card) error
	DeleteFunc                   func(id uint) error
	GetListIDByCardIDFunc        func(cardID uint) (uint, error)
	PerformTransactionFunc       func(fn func(tx *gorm.DB) error) error
	MoveCardFunc                 func(cardID, oldListID, newListID uint, newPosition uint) error
	GetMaxPositionFunc           func(listID uint) (uint, error)
	ShiftPositionsFunc           func(listID uint, startPosition uint, shiftAmount int, excludedCardID *uint) error
	AddCollaboratorFunc          func(cardID uint, userID uint) error
	RemoveCollaboratorFunc       func(cardID uint, userID uint) error
	GetCollaboratorsByCardIDFunc func(cardID uint) ([]models.User, error)
	IsCollaboratorFunc           func(cardID uint, userID uint) (bool, error)
}

func (m *MockCardRepository) Create(card *models.Card) error {
	if card.Position == 0 {
		if m.GetMaxPositionFunc != nil {
			maxPos, _ := m.GetMaxPosition(card.ListID)
			card.Position = maxPos + 1
		} else {
			card.Position = 1
		}
	}
	if m.CreateFunc != nil {
		return m.CreateFunc(card)
	}
	return errors.New("CreateFunc not implemented")
}
func (m *MockCardRepository) FindByID(id uint) (*models.Card, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc not implemented")
}
func (m *MockCardRepository) FindByListID(listID uint) ([]models.Card, error) {
	if m.FindByListIDFunc != nil {
		return m.FindByListIDFunc(listID)
	}
	return nil, errors.New("FindByListIDFunc not implemented")
}
func (m *MockCardRepository) Update(card *models.Card) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(card)
	}
	return errors.New("UpdateFunc not implemented")
}
func (m *MockCardRepository) Delete(id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return errors.New("DeleteFunc not implemented")
}
func (m *MockCardRepository) GetListIDByCardID(cardID uint) (uint, error) {
	if m.GetListIDByCardIDFunc != nil {
		return m.GetListIDByCardIDFunc(cardID)
	}
	return 0, errors.New("GetListIDByCardIDFunc not implemented")
}
func (m *MockCardRepository) PerformTransaction(fn func(tx *gorm.DB) error) error {
	if m.PerformTransactionFunc != nil {
		return m.PerformTransactionFunc(fn)
	}
	return fn(&gorm.DB{})
}
func (m *MockCardRepository) MoveCard(cardID, oldListID, newListID uint, newPosition uint) error {
	if m.MoveCardFunc != nil {
		return m.MoveCardFunc(cardID, oldListID, newListID, newPosition)
	}
	return errors.New("MoveCardFunc not implemented")
}
func (m *MockCardRepository) GetMaxPosition(listID uint) (uint, error) {
	if m.GetMaxPositionFunc != nil {
		return m.GetMaxPositionFunc(listID)
	}
	return 0, errors.New("GetMaxPositionFunc not implemented")
}
func (m *MockCardRepository) ShiftPositions(listID uint, startPosition uint, shiftAmount int, excludedCardID *uint) error {
	if m.ShiftPositionsFunc != nil {
		return m.ShiftPositionsFunc(listID, startPosition, shiftAmount, excludedCardID)
	}
	return errors.New("ShiftPositionsFunc not implemented")
}
func (m *MockCardRepository) AddCollaborator(cardID uint, userID uint) error {
	if m.AddCollaboratorFunc != nil {
		return m.AddCollaboratorFunc(cardID, userID)
	}
	return errors.New("AddCollaboratorFunc not implemented")
}
func (m *MockCardRepository) RemoveCollaborator(cardID uint, userID uint) error {
	if m.RemoveCollaboratorFunc != nil {
		return m.RemoveCollaboratorFunc(cardID, userID)
	}
	return errors.New("RemoveCollaboratorFunc not implemented")
}
func (m *MockCardRepository) GetCollaboratorsByCardID(cardID uint) ([]models.User, error) {
	if m.GetCollaboratorsByCardIDFunc != nil {
		return m.GetCollaboratorsByCardIDFunc(cardID)
	}
	return nil, errors.New("GetCollaboratorsByCardIDFunc not implemented")
}
func (m *MockCardRepository) IsCollaborator(cardID uint, userID uint) (bool, error) {
	if m.IsCollaboratorFunc != nil {
		return m.IsCollaboratorFunc(cardID, userID)
	}
	return false, errors.New("IsCollaboratorFunc not implemented")
}

var _ repositories.CardRepositoryInterface = (*MockCardRepository)(nil)

// --- MockListRepositoryForCardService ---
type MockListRepositoryForCardService struct {
	repositories.ListRepositoryInterface
	GetBoardIDByListIDFunc func(listID uint) (uint, error)
	FindByIDFunc           func(id uint) (*models.List, error)
}

func (m *MockListRepositoryForCardService) GetBoardIDByListID(listID uint) (uint, error) {
	if m.GetBoardIDByListIDFunc != nil {
		return m.GetBoardIDByListIDFunc(listID)
	}
	return 0, errors.New("GetBoardIDByListIDFunc not implemented")
}
func (m *MockListRepositoryForCardService) FindByID(id uint) (*models.List, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc on MockListRepositoryForCardService not implemented")
}
func (m *MockListRepositoryForCardService) Create(list *models.List) error {
	return errors.New("not implemented")
}
func (m *MockListRepositoryForCardService) FindByBoardID(boardID uint) ([]models.List, error) {
	return nil, errors.New("not implemented")
}
func (m *MockListRepositoryForCardService) Update(list *models.List) error {
	return errors.New("not implemented")
}
func (m *MockListRepositoryForCardService) Delete(id uint) error {
	return errors.New("not implemented")
}
func (m *MockListRepositoryForCardService) GetMaxPosition(boardID uint) (uint, error) {
	return 0, errors.New("not implemented")
}
func (m *MockListRepositoryForCardService) GetDB() *gorm.DB { return nil }
func (m *MockListRepositoryForCardService) PerformTransaction(fn func(tx *gorm.DB) error) error {
	return errors.New("not implemented")
}

// --- MockBoardRepositoryForCardService ---
type MockBoardRepositoryForCardService struct {
	repositories.BoardRepositoryInterface
	FindByIDFunc func(id uint) (*models.Board, error)
}

func (m *MockBoardRepositoryForCardService) FindByID(id uint) (*models.Board, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc on MockBoardRepositoryForCardService not implemented")
}
func (m *MockBoardRepositoryForCardService) Create(board *models.Board) error {
	return errors.New("not implemented")
}
func (m *MockBoardRepositoryForCardService) FindByOwnerOrMember(userID uint) ([]models.Board, error) {
	return nil, errors.New("not implemented")
}
func (m *MockBoardRepositoryForCardService) Update(board *models.Board) error {
	return errors.New("not implemented")
}
func (m *MockBoardRepositoryForCardService) Delete(id uint) error {
	return errors.New("not implemented")
}
func (m *MockBoardRepositoryForCardService) IsOwner(boardID uint, userID uint) (bool, error) {
	return false, errors.New("not implemented")
}

// --- MockBoardMemberRepositoryForCardService ---
type MockBoardMemberRepositoryForCardService struct {
	repositories.BoardMemberRepositoryInterface
	IsMemberFunc func(boardID uint, userID uint) (bool, error)
}

func (m *MockBoardMemberRepositoryForCardService) IsMember(boardID uint, userID uint) (bool, error) {
	if m.IsMemberFunc != nil {
		return m.IsMemberFunc(boardID, userID)
	}
	return false, errors.New("IsMemberFunc on MockBoardMemberRepositoryForCardService not implemented")
}
func (m *MockBoardMemberRepositoryForCardService) AddMember(member *models.BoardMember) error {
	return errors.New("not implemented")
}
func (m *MockBoardMemberRepositoryForCardService) RemoveMember(boardID uint, userID uint) error {
	return errors.New("not implemented")
}
func (m *MockBoardMemberRepositoryForCardService) FindMembersByBoardID(boardID uint) ([]models.BoardMember, error) {
	return nil, errors.New("not implemented")
}
func (m *MockBoardMemberRepositoryForCardService) FindByBoardIDAndUserID(boardID uint, userID uint) (*models.BoardMember, error) {
	return nil, errors.New("not implemented")
}

// --- MockUserRepositoryForCardService ---
type MockUserRepositoryForCardService struct {
	repositories.UserRepositoryInterface
	CreateFunc      func(user *models.User) error
	FindByEmailFunc func(email string) (*models.User, error)
	FindByIDFunc    func(id uint) (*models.User, error)
}

func (m *MockUserRepositoryForCardService) Create(user *models.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return errors.New("CreateFunc not implemented in MockUserRepositoryForCardService")
}
func (m *MockUserRepositoryForCardService) FindByEmail(email string) (*models.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, errors.New("FindByEmailFunc not implemented in MockUserRepositoryForCardService")
}
func (m *MockUserRepositoryForCardService) FindByID(id uint) (*models.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc not implemented in MockUserRepositoryForCardService")
}

var _ repositories.UserRepositoryInterface = (*MockUserRepositoryForCardService)(nil)

func TestPlaceholder_CardService(t *testing.T) {
	assert.True(t, true, "This is a placeholder test.")
}

func TestCardService_CreateCard_WithColor(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{} // CardService depends on this now

	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(1)
	listID := uint(10)
	boardID := uint(100)
	cardTitle := "Card With Color"
	cardColor := "#123456"

	// Mock for checkAccessViaList
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var createdCardModel models.Card
	mockCardRepo.CreateFunc = func(card *models.Card) error {
		card.ID = 1                     // Simulate ID set by DB
		card.Position = 1               // Simulate position set by repo
		card.Status = models.StatusToDo // Default status
		createdCardModel = *card        // Capture card passed to Create
		return nil
	}
	mockCardRepo.GetMaxPositionFunc = func(lID uint) (uint, error) { return 0, nil } // For Create logic in mock

	// Mock for final FindByID
	mockCardRepo.FindByIDFunc = func(id uint) (*models.Card, error) {
		// Return the captured card, which should include the color
		assert.Equal(t, createdCardModel.ID, id)
		return &createdCardModel, nil
	}

	card, err := cardService.CreateCard(listID, cardTitle, "", nil, nil, nil, nil, &cardColor, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, cardTitle, card.Title)
	assert.NotNil(t, card.Color)
	assert.Equal(t, cardColor, *card.Color)
	assert.Equal(t, uint(1), card.ID) // Corrected: Use the ID set by mock CreateFunc
}

func TestCardService_CreateCard_WithoutColor(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}

	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(1)
	listID := uint(10)
	boardID := uint(100)
	cardTitle := "Card No Color"

	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var createdCardModel models.Card
	mockCardRepo.CreateFunc = func(card *models.Card) error {
		card.ID = 2
		card.Position = 1
		card.Status = models.StatusToDo
		createdCardModel = *card
		return nil
	}
	mockCardRepo.GetMaxPositionFunc = func(lID uint) (uint, error) { return 0, nil }
	mockCardRepo.FindByIDFunc = func(id uint) (*models.Card, error) { return &createdCardModel, nil }

	card, err := cardService.CreateCard(listID, cardTitle, "", nil, nil, nil, nil, nil, currentUserID) // Color is nil

	assert.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, cardTitle, card.Title)
	assert.Nil(t, card.Color)
	assert.Equal(t, uint(2), card.ID) // Corrected: Use the ID set by mock CreateFunc
}

func TestCardService_AddCollaboratorToCard_SuccessByEmail(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}

	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(1)
	cardID := uint(100)
	listID := uint(10)
	boardID := uint(1)
	targetUserEmail := "collaborator@example.com"
	targetUserID := uint(5)

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
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	expectedTargetUser := &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID, Email: targetUserEmail, Username: "collabUser"}
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		assert.Equal(t, targetUserEmail, email)
		return expectedTargetUser, nil
	}
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) {
		return false, nil
	}
	mockCardRepo.AddCollaboratorFunc = func(cID uint, uID uint) error {
		return nil
	}

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, targetUserEmail, nil)

	assert.NoError(t, err)
	assert.NotNil(t, addedUser)
	assert.Equal(t, targetUserID, addedUser.ID)
}

func TestCardService_AddCollaboratorToCard_SuccessByUserID(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}

	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(1)
	cardID := uint(100)
	listID := uint(10)
	boardID := uint(1)
	targetUserID := uint(5)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	expectedTargetUser := &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID, Username: "collabUserByID"}
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) { return expectedTargetUser, nil }
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return false, nil }
	mockCardRepo.AddCollaboratorFunc = func(cID uint, uID uint) error { return nil }

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, "", &targetUserID)
	assert.NoError(t, err)
	assert.NotNil(t, addedUser)
	assert.Equal(t, targetUserID, addedUser.ID)
}

func TestCardService_AddCollaboratorToCard_AlreadyCollaborator(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, targetUserID := uint(1), uint(100), uint(10), uint(1), uint(5)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) {
		return &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID}, nil
	}
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return true, nil }
	addCollabCalled := false
	mockCardRepo.AddCollaboratorFunc = func(cID uint, uID uint) error { addCollabCalled = true; return nil }

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, "", &targetUserID)
	assert.NoError(t, err)
	assert.NotNil(t, addedUser)
	assert.False(t, addCollabCalled)
}

func TestCardService_AddCollaboratorToCard_PermissionDenied(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, ownerOfBoardID, targetUserEmail := uint(1), uint(100), uint(10), uint(1), uint(2), "c@e.com"

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerOfBoardID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil }

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, targetUserEmail, nil)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, addedUser)
}

func TestCardService_AddCollaboratorToCard_TargetUserNotFoundByEmail(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, targetUserEmail := uint(1), uint(100), uint(10), uint(1), "non@ex.com"

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return nil, gorm.ErrRecordNotFound }

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, targetUserEmail, nil)
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, addedUser)
}

func TestCardService_AddCollaboratorToCard_TargetUserNotFoundByID(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, targetUserID := uint(1), uint(100), uint(10), uint(1), uint(999)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) { return nil, gorm.ErrRecordNotFound }

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, "", &targetUserID)
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, addedUser)
}

func TestCardService_AddCollaboratorToCard_RepoErrorOnAdd(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, targetUserID, expectedError := uint(1), uint(100), uint(10), uint(1), uint(5), errors.New("DB error")

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) {
		return &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID}, nil
	}
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return false, nil }
	mockCardRepo.AddCollaboratorFunc = func(cID uint, uID uint) error { return expectedError }

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, "", &targetUserID)
	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, addedUser)
}

func TestCardService_RemoveCollaboratorFromCard_Success(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, targetUserID := uint(1), uint(100), uint(10), uint(1), uint(5)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return true, nil }
	removeCalled := false
	mockCardRepo.RemoveCollaboratorFunc = func(cID uint, uID uint) error { removeCalled = true; return nil }

	err := cardService.RemoveCollaboratorFromCard(cardID, currentUserID, targetUserID)
	assert.NoError(t, err)
	assert.True(t, removeCalled)
}

func TestCardService_RemoveCollaboratorFromCard_PermissionDenied(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, ownerOfBoardID, targetUserID := uint(1), uint(100), uint(10), uint(1), uint(2), uint(5)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerOfBoardID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil }

	err := cardService.RemoveCollaboratorFromCard(cardID, currentUserID, targetUserID)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCardService_RemoveCollaboratorFromCard_NotCollaborator(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, targetUserID := uint(1), uint(100), uint(10), uint(1), uint(5)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return false, nil }

	err := cardService.RemoveCollaboratorFromCard(cardID, currentUserID, targetUserID)
	assert.ErrorIs(t, err, ErrUserNotCollaborator)
}

func TestCardService_RemoveCollaboratorFromCard_RepoErrorOnRemove(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, targetUserID, expectedError := uint(1), uint(100), uint(10), uint(1), uint(5), errors.New("DB error")

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return true, nil }
	mockCardRepo.RemoveCollaboratorFunc = func(cID uint, uID uint) error { return expectedError }

	err := cardService.RemoveCollaboratorFromCard(cardID, currentUserID, targetUserID)
	assert.ErrorIs(t, err, expectedError)
}

func TestCardService_GetCardCollaborators_Success(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID := uint(1), uint(100), uint(10), uint(1)
	expectedUsers := []models.User{{ID: 5}, {ID: 6}}

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockCardRepo.GetCollaboratorsByCardIDFunc = func(cID uint) ([]models.User, error) { return expectedUsers, nil }

	users, err := cardService.GetCardCollaborators(cardID, currentUserID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestCardService_GetCardCollaborators_PermissionDenied(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, ownerOfBoardID := uint(1), uint(100), uint(10), uint(1), uint(2)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerOfBoardID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil }

	users, err := cardService.GetCardCollaborators(cardID, currentUserID)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, users)
}

func TestCardService_CreateCard_WithColorValue(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(1)
	listID := uint(10)
	boardID := uint(100)
	cardTitle := "Card With Color"
	cardColor := "#FF0000"
	expectedCardID := uint(1)

	// Mocks for checkAccessViaList
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var capturedCard models.Card
	mockCardRepo.CreateFunc = func(card *models.Card) error {
		assert.NotNil(t, card.Color)
		assert.Equal(t, cardColor, *card.Color)
		card.ID = expectedCardID // Simulate DB setting ID
		card.Position = 1        // Simulate repo setting position
		card.Status = models.StatusToDo
		capturedCard = *card // Capture for final FindByID mock
		return nil
	}
	mockCardRepo.GetMaxPositionFunc = func(lID uint) (uint, error) { return 0, nil } // Called by Create mock
	mockCardRepo.FindByIDFunc = func(id uint) (*models.Card, error) {
		assert.Equal(t, expectedCardID, id)
		return &capturedCard, nil // Return the card state after Create
	}

	card, err := cardService.CreateCard(listID, cardTitle, "Desc", nil, nil, nil, nil, &cardColor, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, expectedCardID, card.ID)
	assert.NotNil(t, card.Color)
	assert.Equal(t, cardColor, *card.Color)
}

func TestCardService_CreateCard_WithoutColorValue(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(1)
	listID := uint(10)
	boardID := uint(100)
	cardTitle := "Card Without Color"
	expectedCardID := uint(2)

	// Mocks for checkAccessViaList
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var capturedCard models.Card
	mockCardRepo.CreateFunc = func(card *models.Card) error {
		assert.Nil(t, card.Color) // Color should be nil as input is nil
		card.ID = expectedCardID
		card.Position = 1
		card.Status = models.StatusToDo
		capturedCard = *card
		return nil
	}
	mockCardRepo.GetMaxPositionFunc = func(lID uint) (uint, error) { return 0, nil }
	mockCardRepo.FindByIDFunc = func(id uint) (*models.Card, error) {
		assert.Equal(t, expectedCardID, id)
		return &capturedCard, nil
	}

	card, err := cardService.CreateCard(listID, cardTitle, "Desc", nil, nil, nil, nil, nil, currentUserID) // Color is nil

	assert.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, expectedCardID, card.ID)
	assert.Nil(t, card.Color)
}

func TestCardService_GetCardCollaborators_RepoError(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID, expectedError := uint(1), uint(100), uint(10), uint(1), errors.New("DB error")

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockCardRepo.GetCollaboratorsByCardIDFunc = func(cID uint) ([]models.User, error) { return nil, expectedError }

	users, err := cardService.GetCardCollaborators(cardID, currentUserID)
	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, users)
}

func TestCardService_UpdateCard_SetColorFirstTime(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID, cardID, listID, boardID := uint(1), uint(50), uint(10), uint(100)
	initialCard := &models.Card{Model: gorm.Model{ID: cardID}, Title: "Original", ListID: listID, Color: nil}
	newColor := "#FFAABB"

	// Mocks for checkAccessViaCard
	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var cardStateAfterInitialFind models.Card
	// var cardStateForUpdate models.Card // Removed: declared and not used
	var cardStateForFinalFind models.Card   // This will store the card state for the *final* FindByID call
	var cardStateCapturedByUpdate models.Card // Stores the state of the card as captured by the UpdateFunc

	// Initial FindByID
	mockCardRepo.FindByIDFunc = func(cID uint) (*models.Card, error) {
		if cID == cardID && cardStateAfterInitialFind.ID == 0 { // First call in service
			cardStateAfterInitialFind = *initialCard
			return &cardStateAfterInitialFind, nil
		}
		// Final FindByID call in service, should use the state captured by UpdateFunc
		assert.Equal(t, cardID, cID)
		cardStateForFinalFind = cardStateCapturedByUpdate
		return &cardStateForFinalFind, nil
	}

	mockCardRepo.UpdateFunc = func(card *models.Card) error {
		assert.NotNil(t, card.Color)
		assert.Equal(t, newColor, *card.Color)
		cardStateCapturedByUpdate = *card // Capture state passed to Update
		return nil
	}

	updatedCard, err := cardService.UpdateCard(cardID, nil, nil, nil, nil, nil, nil, nil, &newColor, currentUserID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedCard)
	assert.NotNil(t, updatedCard.Color)
	assert.Equal(t, newColor, *updatedCard.Color)
}

func TestCardService_UpdateCard_ChangeExistingColor(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID, cardID, listID, boardID := uint(1), uint(50), uint(10), uint(100)
	initialColor := "#OLDCLR"
	initialCard := &models.Card{Model: gorm.Model{ID: cardID}, Title: "Original", ListID: listID, Color: &initialColor}
	newColor := "#NEWCLR"

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var cardStateAfterUpdate models.Card
	var findCallCount int
	mockCardRepo.FindByIDFunc = func(cID uint) (*models.Card, error) {
		findCallCount++
		if findCallCount == 1 {
			tCard := *initialCard
			return &tCard, nil
		}
		return &cardStateAfterUpdate, nil
	}
	mockCardRepo.UpdateFunc = func(card *models.Card) error {
		assert.NotNil(t, card.Color)
		assert.Equal(t, newColor, *card.Color)
		cardStateAfterUpdate = *card
		return nil
	}

	updatedCard, err := cardService.UpdateCard(cardID, nil, nil, nil, nil, nil, nil, nil, &newColor, currentUserID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedCard.Color)
	assert.Equal(t, newColor, *updatedCard.Color)
}

func TestCardService_UpdateCard_ClearColorWithEmptyString(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID, cardID, listID, boardID := uint(1), uint(50), uint(10), uint(100)
	initialColor := "#EXISTING"
	initialCard := &models.Card{Model: gorm.Model{ID: cardID}, Title: "Original", ListID: listID, Color: &initialColor}
	emptyColor := ""

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var cardStateAfterUpdate models.Card
	var findCallCount int
	mockCardRepo.FindByIDFunc = func(cID uint) (*models.Card, error) {
		findCallCount++
		if findCallCount == 1 {
			tCard := *initialCard
			return &tCard, nil
		}
		return &cardStateAfterUpdate, nil
	}
	mockCardRepo.UpdateFunc = func(card *models.Card) error {
		assert.Nil(t, card.Color) // Service should set card.Color to nil if input *color is ""
		cardStateAfterUpdate = *card
		return nil
	}

	updatedCard, err := cardService.UpdateCard(cardID, nil, nil, nil, nil, nil, nil, nil, &emptyColor, currentUserID)
	assert.NoError(t, err)
	assert.Nil(t, updatedCard.Color)
}

func TestCardService_UpdateCard_NoColorInRequest(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID, cardID, listID, boardID := uint(1), uint(50), uint(10), uint(100)
	initialColor := "#SHOULDSTAY"
	initialCard := &models.Card{Model: gorm.Model{ID: cardID}, Title: "Original", ListID: listID, Color: &initialColor}
	newTitle := "New Title"

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var cardStateAfterUpdate models.Card
	var findCallCount int
	mockCardRepo.FindByIDFunc = func(cID uint) (*models.Card, error) {
		findCallCount++
		if findCallCount == 1 {
			tCard := *initialCard
			return &tCard, nil
		}
		return &cardStateAfterUpdate, nil
	}
	mockCardRepo.UpdateFunc = func(card *models.Card) error {
		assert.NotNil(t, card.Color) // Color should be unchanged
		assert.Equal(t, initialColor, *card.Color)
		cardStateAfterUpdate = *card
		return nil
	}

	updatedCard, err := cardService.UpdateCard(cardID, &newTitle, nil, nil, nil, nil, nil, nil, nil, currentUserID) // Color param is nil
	assert.NoError(t, err)
	assert.NotNil(t, updatedCard.Color)
	assert.Equal(t, initialColor, *updatedCard.Color)
	assert.Equal(t, newTitle, updatedCard.Title)
}

// --- DueDate Tests ---

func TestCardService_CreateCard_WithDueDate(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(1)
	listID := uint(10)
	boardID := uint(100)
	cardTitle := "Card With DueDate"
	testDueDate := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second) // Example DueDate

	// Mocks for checkAccessViaList
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var capturedCard models.Card
	mockCardRepo.CreateFunc = func(card *models.Card) error {
		assert.NotNil(t, card.DueDate)
		assert.True(t, testDueDate.Equal(*card.DueDate), "Expected DueDate %v, got %v", testDueDate, *card.DueDate)
		card.ID = 1
		card.Position = 1
		card.Status = models.StatusToDo
		capturedCard = *card
		return nil
	}
	mockCardRepo.GetMaxPositionFunc = func(lID uint) (uint, error) { return 0, nil }
	mockCardRepo.FindByIDFunc = func(id uint) (*models.Card, error) {
		return &capturedCard, nil
	}

	createdCard, err := cardService.CreateCard(listID, cardTitle, "Desc", nil, &testDueDate, nil, nil, nil, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, createdCard)
	assert.NotNil(t, createdCard.DueDate)
	assert.True(t, testDueDate.Equal(*createdCard.DueDate))
}

func TestCardService_CreateCard_WithoutDueDate(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(1)
	listID := uint(10)
	boardID := uint(100)
	cardTitle := "Card Without DueDate"

	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var capturedCard models.Card
	mockCardRepo.CreateFunc = func(card *models.Card) error {
		assert.Nil(t, card.DueDate)
		card.ID = 2
		card.Position = 1
		card.Status = models.StatusToDo
		capturedCard = *card
		return nil
	}
	mockCardRepo.GetMaxPositionFunc = func(lID uint) (uint, error) { return 0, nil }
	mockCardRepo.FindByIDFunc = func(id uint) (*models.Card, error) {
		return &capturedCard, nil
	}

	createdCard, err := cardService.CreateCard(listID, cardTitle, "Desc", nil, nil, nil, nil, nil, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, createdCard)
	assert.Nil(t, createdCard.DueDate)
}

func TestCardService_UpdateCard_SetDueDateFirstTime(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID, cardID, listID, boardID := uint(1), uint(50), uint(10), uint(100)
	initialCard := &models.Card{Model: gorm.Model{ID: cardID}, Title: "Original", ListID: listID, DueDate: nil}
	newDueDate := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var cardStateCapturedByUpdate models.Card
	var findByIDCallCount int
	mockCardRepo.FindByIDFunc = func(cID uint) (*models.Card, error) {
		findByIDCallCount++
		if findByIDCallCount == 1 { // First call in UpdateCard, before update
			tempInitialCard := *initialCard // Use a copy
			return &tempInitialCard, nil
		}
		// Second call in UpdateCard, after update
		return &cardStateCapturedByUpdate, nil
	}

	mockCardRepo.UpdateFunc = func(card *models.Card) error {
		assert.NotNil(t, card.DueDate)
		assert.True(t, newDueDate.Equal(*card.DueDate))
		cardStateCapturedByUpdate = *card
		return nil
	}

	updatedCard, err := cardService.UpdateCard(cardID, nil, nil, nil, &newDueDate, nil, nil, nil, nil, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, updatedCard)
	assert.NotNil(t, updatedCard.DueDate)
	assert.True(t, newDueDate.Equal(*updatedCard.DueDate))
}

func TestCardService_UpdateCard_ChangeExistingDueDate(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID, cardID, listID, boardID := uint(1), uint(51), uint(10), uint(100)
	initialDueDate := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
	initialCard := &models.Card{Model: gorm.Model{ID: cardID}, Title: "Original", ListID: listID, DueDate: &initialDueDate}
	newDueDate := time.Now().Add(72 * time.Hour).UTC().Truncate(time.Second)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var cardStateCapturedByUpdate models.Card
	var findByIDCallCount int
	mockCardRepo.FindByIDFunc = func(cID uint) (*models.Card, error) {
		findByIDCallCount++
		if findByIDCallCount == 1 {
			tempInitialCard := *initialCard
			return &tempInitialCard, nil
		}
		return &cardStateCapturedByUpdate, nil
	}

	mockCardRepo.UpdateFunc = func(card *models.Card) error {
		assert.NotNil(t, card.DueDate)
		assert.True(t, newDueDate.Equal(*card.DueDate))
		cardStateCapturedByUpdate = *card
		return nil
	}

	updatedCard, err := cardService.UpdateCard(cardID, nil, nil, nil, &newDueDate, nil, nil, nil, nil, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, updatedCard)
	assert.NotNil(t, updatedCard.DueDate)
	assert.True(t, newDueDate.Equal(*updatedCard.DueDate))
}

func TestCardService_UpdateCard_WithNilDueDateParameter(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID, cardID, listID, boardID := uint(1), uint(52), uint(10), uint(100)
	initialDueDate := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
	initialCard := &models.Card{Model: gorm.Model{ID: cardID}, Title: "Original", ListID: listID, DueDate: &initialDueDate}
	newTitle := "Updated Title"

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var cardStateCapturedByUpdate models.Card
	var findByIDCallCount int
	mockCardRepo.FindByIDFunc = func(cID uint) (*models.Card, error) {
		findByIDCallCount++
		if findByIDCallCount == 1 {
			tempInitialCard := *initialCard
			return &tempInitialCard, nil
		}
		return &cardStateCapturedByUpdate, nil
	}

	updateCalled := false
	mockCardRepo.UpdateFunc = func(card *models.Card) error {
		updateCalled = true
		// Check that DueDate was NOT changed, as input param was nil
		assert.NotNil(t, card.DueDate)
		assert.True(t, initialDueDate.Equal(*card.DueDate))
		assert.Equal(t, newTitle, card.Title) // Ensure other fields can still be updated
		cardStateCapturedByUpdate = *card
		return nil
	}

	// Pass nil for dueDate parameter
	updatedCard, err := cardService.UpdateCard(cardID, &newTitle, nil, nil, nil, nil, nil, nil, nil, currentUserID)

	assert.NoError(t, err)
	assert.True(t, updateCalled, "UpdateFunc should be called")
	assert.NotNil(t, updatedCard)
	assert.NotNil(t, updatedCard.DueDate, "DueDate should not have been cleared")
	assert.True(t, initialDueDate.Equal(*updatedCard.DueDate), "DueDate should be unchanged")
	assert.Equal(t, newTitle, updatedCard.Title)
}
