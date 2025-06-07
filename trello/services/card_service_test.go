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
	IsUserCollaboratorOrAssigneeFunc func(cardID uint, userID uint) (bool, error)
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
	if m.PerformTransactionFunc != nil {
		return m.PerformTransactionFunc(fn)
	}
	// Default behavior: pass a non-nil but potentially non-functional DB.
	// Tests that rely on the transaction actually working need to override PerformTransactionFunc.
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
func (m *MockCardRepository) IsUserCollaboratorOrAssignee(cardID uint, userID uint) (bool, error) {
	if m.IsUserCollaboratorOrAssigneeFunc != nil {
		return m.IsUserCollaboratorOrAssigneeFunc(cardID, userID)
	}
	return false, errors.New("IsUserCollaboratorOrAssigneeFunc not implemented")
}

var _ repositories.CardRepositoryInterface = (*MockCardRepository)(nil)

type MockListRepositoryForCardService struct {
	repositories.ListRepositoryInterface
	GetBoardIDByListIDFunc func(listID uint) (uint, error)
	FindByIDFunc           func(id uint) (*models.List, error)
}

func (m *MockListRepositoryForCardService) GetBoardIDByListID(listID uint) (uint, error) {
	if m.GetBoardIDByListIDFunc != nil { return m.GetBoardIDByListIDFunc(listID) }
	return 0, errors.New("GetBoardIDByListIDFunc not implemented")
}
func (m *MockListRepositoryForCardService) FindByID(id uint) (*models.List, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, errors.New("FindByIDFunc on MockListRepositoryForCardService not implemented")
}
func (m *MockListRepositoryForCardService) Create(list *models.List) error { return errors.New("not implemented") }
func (m *MockListRepositoryForCardService) FindByBoardID(boardID uint) ([]models.List, error) { return nil, errors.New("not implemented") }
func (m *MockListRepositoryForCardService) Update(list *models.List) error { return errors.New("not implemented") }
func (m *MockListRepositoryForCardService) Delete(id uint) error           { return errors.New("not implemented") }
func (m *MockListRepositoryForCardService) GetMaxPosition(boardID uint) (uint, error) { return 0, errors.New("not implemented") }
func (m *MockListRepositoryForCardService) GetDB() *gorm.DB                 { return nil }
func (m *MockListRepositoryForCardService) PerformTransaction(fn func(tx *gorm.DB) error) error { return errors.New("not implemented") }

type MockBoardRepositoryForCardService struct {
	repositories.BoardRepositoryInterface
	FindByIDFunc func(id uint) (*models.Board, error)
	IsOwnerFunc  func(boardID uint, userID uint) (bool, error)
}

func (m *MockBoardRepositoryForCardService) FindByID(id uint) (*models.Board, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, errors.New("FindByIDFunc on MockBoardRepositoryForCardService not implemented")
}
func (m *MockBoardRepositoryForCardService) Create(board *models.Board) error { return errors.New("not implemented") }
func (m *MockBoardRepositoryForCardService) FindByOwnerOrMember(userID uint) ([]models.Board, error) { return nil, errors.New("not implemented") }
func (m *MockBoardRepositoryForCardService) Update(board *models.Board) error { return errors.New("not implemented") }
func (m *MockBoardRepositoryForCardService) Delete(id uint) error           { return errors.New("not implemented") }
func (m *MockBoardRepositoryForCardService) IsOwner(boardID uint, userID uint) (bool, error) {
	if m.IsOwnerFunc != nil { return m.IsOwnerFunc(boardID, userID) }
	return false, errors.New("not implemented")
}

type MockBoardMemberRepositoryForCardService struct {
	repositories.BoardMemberRepositoryInterface
	IsMemberFunc func(boardID uint, userID uint) (bool, error)
}

func (m *MockBoardMemberRepositoryForCardService) IsMember(boardID uint, userID uint) (bool, error) {
	if m.IsMemberFunc != nil { return m.IsMemberFunc(boardID, userID) }
	return false, errors.New("IsMemberFunc on MockBoardMemberRepositoryForCardService not implemented")
}
func (m *MockBoardMemberRepositoryForCardService) AddMember(member *models.BoardMember) error { return errors.New("not implemented") }
func (m *MockBoardMemberRepositoryForCardService) RemoveMember(boardID uint, userID uint) error { return errors.New("not implemented") }
func (m *MockBoardMemberRepositoryForCardService) FindMembersByBoardID(boardID uint) ([]models.BoardMember, error) { return nil, errors.New("not implemented") }
func (m *MockBoardMemberRepositoryForCardService) FindByBoardIDAndUserID(boardID uint, userID uint) (*models.BoardMember, error) { return nil, errors.New("not implemented") }

type MockUserRepositoryForCardService struct {
	repositories.UserRepositoryInterface
	CreateFunc      func(user *models.User) error
	FindByEmailFunc func(email string) (*models.User, error)
	FindByIDFunc    func(id uint) (*models.User, error)
}

func (m *MockUserRepositoryForCardService) Create(user *models.User) error {
	if m.CreateFunc != nil { return m.CreateFunc(user) }
	return errors.New("CreateFunc not implemented in MockUserRepositoryForCardService")
}
func (m *MockUserRepositoryForCardService) FindByEmail(email string) (*models.User, error) {
	if m.FindByEmailFunc != nil { return m.FindByEmailFunc(email) }
	return nil, errors.New("FindByEmailFunc not implemented in MockUserRepositoryForCardService")
}
func (m *MockUserRepositoryForCardService) FindByID(id uint) (*models.User, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, errors.New("FindByIDFunc not implemented in MockUserRepositoryForCardService")
}
var _ repositories.UserRepositoryInterface = (*MockUserRepositoryForCardService)(nil)


func TestCardService_GetCardByID_Permissions(t *testing.T) {
	cardID := uint(1)
	listID := uint(10)
	boardID := uint(100)
	ownerUserID := uint(1)
	collaboratorUserID := uint(2)
	memberUserID := uint(3)
	otherUserID := uint(4)

	mockCardResult := &models.Card{Model: gorm.Model{ID: cardID}, Title: "Test Card"}

	tests := []struct {
		name                           string
		currentUserID                  uint
		mockGetListIDByCardIDFunc      func(cID uint) (uint, error)
		mockGetBoardIDByListIDFunc     func(lID uint) (uint, error)
		mockBoardFindByIDFunc          func(bID uint) (*models.Board, error)
		mockIsMemberFunc               func(bID uint, uID uint) (bool, error)
		mockIsUserCollabOrAssigneeFunc func(cID uint, uID uint) (bool, error)
		mockCardFindByIDFunc           func(cID uint) (*models.Card, error)
		expectedError                  error
		expectCard                     bool
	}{
		{
			name:          "User is board owner",
			currentUserID: ownerUserID,
			mockGetListIDByCardIDFunc: func(cID uint) (uint, error) { return listID, nil },
			mockGetBoardIDByListIDFunc: func(lID uint) (uint, error) { return boardID, nil },
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) {
				return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil
			},
			mockIsUserCollabOrAssigneeFunc: func(cID uint, uID uint) (bool, error) { return false, nil },
			mockCardFindByIDFunc:           func(cID uint) (*models.Card, error) { return mockCardResult, nil },
			expectedError:                  nil,
			expectCard:                     true,
		},
		{
			name:          "User is collaborator/assignee",
			currentUserID: collaboratorUserID,
			mockGetListIDByCardIDFunc: func(cID uint) (uint, error) { return listID, nil },
			mockGetBoardIDByListIDFunc: func(lID uint) (uint, error) { return boardID, nil },
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) {
				return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil
			},
			mockIsMemberFunc: func(bID uint, uID uint) (bool, error) { return true, nil },
			mockIsUserCollabOrAssigneeFunc: func(cID uint, uID uint) (bool, error) { return true, nil },
			mockCardFindByIDFunc:           func(cID uint) (*models.Card, error) { return mockCardResult, nil },
			expectedError:                  nil,
			expectCard:                     true,
		},
		{
			name:          "User is board member, not owner/collaborator/assignee",
			currentUserID: memberUserID,
			mockGetListIDByCardIDFunc: func(cID uint) (uint, error) { return listID, nil },
			mockGetBoardIDByListIDFunc: func(lID uint) (uint, error) { return boardID, nil },
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) {
				return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil
			},
			mockIsMemberFunc: func(bID uint, uID uint) (bool, error) { return true, nil },
			mockIsUserCollabOrAssigneeFunc: func(cID uint, uID uint) (bool, error) { return false, nil },
			mockCardFindByIDFunc:           func(cID uint) (*models.Card, error) { return mockCardResult, nil },
			expectedError:                  ErrForbidden,
			expectCard:                     false,
		},
		{
			name:          "User is not board member",
			currentUserID: otherUserID,
			mockGetListIDByCardIDFunc: func(cID uint) (uint, error) { return listID, nil },
			mockGetBoardIDByListIDFunc: func(lID uint) (uint, error) { return boardID, nil },
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) {
				return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil
			},
			mockIsMemberFunc: func(bID uint, uID uint) (bool, error) { return false, nil },
			expectedError:    ErrForbidden,
			expectCard:       false,
		},
		{
			name:          "Card not found (repo level)",
			currentUserID: ownerUserID,
			mockGetListIDByCardIDFunc: func(cID uint) (uint, error) { return 0, gorm.ErrRecordNotFound },
			expectedError:             ErrCardNotFound,
			expectCard:                false,
		},
		{
			name:          "Board not found (repo level)",
			currentUserID: ownerUserID,
			mockGetListIDByCardIDFunc: func(cID uint) (uint, error) { return listID, nil },
			mockGetBoardIDByListIDFunc: func(lID uint) (uint, error) { return boardID, nil },
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) { return nil, gorm.ErrRecordNotFound },
			expectedError:             ErrBoardNotFound,
			expectCard:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCardRepo := &MockCardRepository{
				GetListIDByCardIDFunc:        tt.mockGetListIDByCardIDFunc,
				IsUserCollaboratorOrAssigneeFunc: tt.mockIsUserCollabOrAssigneeFunc,
				FindByIDFunc:                 tt.mockCardFindByIDFunc,
			}
			mockListRepo := &MockListRepositoryForCardService{GetBoardIDByListIDFunc: tt.mockGetBoardIDByListIDFunc}
			mockBoardRepo := &MockBoardRepositoryForCardService{FindByIDFunc: tt.mockBoardFindByIDFunc}
			mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{IsMemberFunc: tt.mockIsMemberFunc}
			mockUserRepo := &MockUserRepositoryForCardService{}

			service := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

			card, err := service.GetCardByID(cardID, tt.currentUserID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectCard {
				assert.NotNil(t, card)
				assert.Equal(t, mockCardResult.ID, card.ID)
			} else {
				assert.Nil(t, card)
			}
		})
	}
}


func TestCardService_UpdateCard_Permissions(t *testing.T) {
	cardID := uint(1)
	listID := uint(10)
	boardID := uint(100)
	ownerUserID := uint(1)
	collaboratorUserID := uint(2)
	memberUserID := uint(3)

	originalCard := &models.Card{
		Model:          gorm.Model{ID: cardID},
		Title:          "Original Title",
		Description:    "Original Description",
		ListID:         listID,
		AssignedUserID: &collaboratorUserID,
	}
	newTitle := "New Title"
	newDescription := "New Description"
	newDueDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name                           string
		currentUserID                  uint
		updatePayloadTitle             *string
		updatePayloadDescription       *string
		updatePayloadDueDate          *time.Time
		mockBoardFindByIDFunc          func(bID uint) (*models.Board, error)
		mockIsMemberFunc               func(bID uint, uID uint) (bool, error)
		mockIsUserCollabOrAssigneeFunc func(cID uint, uID uint) (bool, error)
		mockCardFindByIDFunc           func(cID uint) (*models.Card, error)
		mockCardUpdateFunc             func(card *models.Card) error
		expectedError                  error
		expectUpdateCall               bool
	}{
		{
			name:          "Owner updates title",
			currentUserID: ownerUserID,
			updatePayloadTitle: &newTitle,
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil },
			mockIsUserCollabOrAssigneeFunc: func(cID uint, uID uint) (bool, error) { return false, nil },
			mockCardFindByIDFunc: func(cID uint) (*models.Card, error) { cardCopy := *originalCard; return &cardCopy, nil },
			mockCardUpdateFunc:   func(card *models.Card) error { assert.Equal(t, newTitle, card.Title); return nil },
			expectedError:      nil,
			expectUpdateCall:   true,
		},
		{
			name:          "Collaborator updates description",
			currentUserID: collaboratorUserID,
			updatePayloadDescription: &newDescription,
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil },
			mockIsMemberFunc:      func(bID uint, uID uint) (bool, error) { return true, nil },
			mockIsUserCollabOrAssigneeFunc: func(cID uint, uID uint) (bool, error) { return true, nil },
			mockCardFindByIDFunc: func(cID uint) (*models.Card, error) { cardCopy := *originalCard; return &cardCopy, nil },
			mockCardUpdateFunc:   func(card *models.Card) error { assert.Equal(t, newDescription, card.Description); return nil },
			expectedError:      nil,
			expectUpdateCall:   true,
		},
		{
			name:          "Collaborator updates due date",
			currentUserID: collaboratorUserID,
			updatePayloadDueDate: &newDueDate,
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil },
			mockIsMemberFunc:      func(bID uint, uID uint) (bool, error) { return true, nil },
			mockIsUserCollabOrAssigneeFunc: func(cID uint, uID uint) (bool, error) { return true, nil },
			mockCardFindByIDFunc: func(cID uint) (*models.Card, error) { cardCopy := *originalCard; return &cardCopy, nil },
			mockCardUpdateFunc:   func(card *models.Card) error { assert.True(t, newDueDate.Equal(*card.DueDate)); return nil },
			expectedError:      nil,
			expectUpdateCall:   true,
		},
		{
			name:          "Collaborator fails to update title",
			currentUserID: collaboratorUserID,
			updatePayloadTitle: &newTitle,
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil },
			mockIsMemberFunc:      func(bID uint, uID uint) (bool, error) { return true, nil },
			mockIsUserCollabOrAssigneeFunc: func(cID uint, uID uint) (bool, error) { return true, nil },
			mockCardFindByIDFunc: func(cID uint) (*models.Card, error) { cardCopy := *originalCard; return &cardCopy, nil },
			expectedError:      ErrPermissionDenied,
			expectUpdateCall:   false,
		},
		{
			name:          "Board member (not owner/collab) fails to update description",
			currentUserID: memberUserID,
			updatePayloadDescription: &newDescription,
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil },
			mockIsMemberFunc:      func(bID uint, uID uint) (bool, error) { return true, nil },
			mockIsUserCollabOrAssigneeFunc: func(cID uint, uID uint) (bool, error) { return false, nil },
			mockCardFindByIDFunc: func(cID uint) (*models.Card, error) { cardCopy := *originalCard; return &cardCopy, nil },
			expectedError:      ErrPermissionDenied,
			expectUpdateCall:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCalled := false
			mockCardRepo := &MockCardRepository{
				GetListIDByCardIDFunc:        func(cID uint) (uint, error) { return listID, nil },
				IsUserCollaboratorOrAssigneeFunc: tt.mockIsUserCollabOrAssigneeFunc,
				FindByIDFunc:                 tt.mockCardFindByIDFunc,
				UpdateFunc: func(card *models.Card) error {
					updateCalled = true
					if tt.mockCardUpdateFunc != nil {
						return tt.mockCardUpdateFunc(card)
					}
					return nil
				},
			}
			mockListRepo := &MockListRepositoryForCardService{GetBoardIDByListIDFunc: func(lID uint) (uint, error) { return boardID, nil }}
			mockBoardRepo := &MockBoardRepositoryForCardService{FindByIDFunc: tt.mockBoardFindByIDFunc}
			mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{IsMemberFunc: tt.mockIsMemberFunc}
			mockUserRepo := &MockUserRepositoryForCardService{}

			service := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

			var assignedUserPtr **uint
			var supervisorPtr **uint
			var statusPtr *models.CardStatus
			var colorPtr *string
			var positionPtr *uint


			_, err := service.UpdateCard(cardID, tt.updatePayloadTitle, tt.updatePayloadDescription, positionPtr, tt.updatePayloadDueDate, assignedUserPtr, supervisorPtr, statusPtr, colorPtr, tt.currentUserID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectUpdateCall, updateCalled)
		})
	}
}


func TestCardService_DeleteCard_Permissions(t *testing.T) {
	cardID := uint(1)
	listID := uint(10)
	boardID := uint(100)
	ownerUserID := uint(1)
	collaboratorUserID := uint(2)
	memberUserID := uint(3)

	originalCard := &models.Card{Model: gorm.Model{ID: cardID}, ListID: listID, Position: 1}

	tests := []struct {
		name                    string
		currentUserID           uint
		mockBoardFindByIDFunc   func(bID uint) (*models.Board, error)
		mockIsMemberFunc        func(bID uint, uID uint) (bool, error)
		mockCardFindByIDFunc    func(cID uint) (*models.Card, error)
		mockPerformTxFunc       func(fn func(tx *gorm.DB) error) error
		expectedError           error
		expectDeleteTransaction bool
	}{
		{
			name:          "Owner deletes card",
			currentUserID: ownerUserID,
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil },
			mockIsMemberFunc: func(bID uint, uID uint) (bool, error) { return true, nil},
			mockCardFindByIDFunc: func(cID uint) (*models.Card, error) { return originalCard, nil },
			// For this specific test, we expect the transaction to proceed with a functional DB.
			mockPerformTxFunc: func(fn func(tx *gorm.DB) error) error {
				// This 't' is not available here. This approach needs 't' to be available
				// or setupTestDB to be callable without 't' or accessible globally.
				// Assuming setupTestDB is a global helper or can be called, though this is problematic.
				// The proper way is to set this up per test case using t.
				// For now, this will still cause issues if setupTestDB needs t.
				// The test case itself (inside t.Run) should define this.
				// So, this default mockPerformTxFunc will be overridden in the test case.
				// For the default here, let's keep it simple, but it won't work for the Owner_deletes_card case by itself.
				return fn(&gorm.DB{})
			},
			expectedError:           nil,
			expectDeleteTransaction: true,
		},
		{
			name:          "Collaborator (non-owner) fails to delete card",
			currentUserID: collaboratorUserID,
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil },
			mockIsMemberFunc: func(bID uint, uID uint) (bool, error) { return true, nil},
			expectedError:           ErrForbidden,
			expectDeleteTransaction: false,
		},
		{
			name:          "Board member (non-owner/collab) fails to delete card",
			currentUserID: memberUserID,
			mockBoardFindByIDFunc: func(bID uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerUserID}, nil },
			mockIsMemberFunc: func(bID uint, uID uint) (bool, error) { return true, nil},
			expectedError:           ErrForbidden,
			expectDeleteTransaction: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteTransactionCalled := false
			mockCardRepo := &MockCardRepository{
				GetListIDByCardIDFunc: func(cID uint) (uint, error) { return listID, nil },
				FindByIDFunc:          tt.mockCardFindByIDFunc,
				// DeleteFunc will be called by the actual service logic if PerformTransactionFunc passes a working DB
			}

			if tt.name == "Owner deletes card" { // Override PerformTransactionFunc specifically for this case
				mockCardRepo.PerformTransactionFunc = func(fn func(tx *gorm.DB) error) error {
					deleteTransactionCalled = true
					db := setupTestDB(t) // 't' is available here from the t.Run scope
					return db.Transaction(fn) // Execute the transaction on a real test DB
				}
				// We also need to ensure that if DeleteFunc is called, it's the mock's DeleteFunc
                // However, the service now uses tx.Delete, so the actual DB operation will occur.
                // So, no explicit mock for DeleteFunc is needed here if the DB op is expected to succeed.
			} else if tt.mockPerformTxFunc != nil { // For other cases that might mock tx differently (though not in this suite)
				mockCardRepo.PerformTransactionFunc = func(fn func(tx *gorm.DB) error) error {
					deleteTransactionCalled = true
					return tt.mockPerformTxFunc(fn)
				}
			} else { // Default if not owner and no specific mockPerformTxFunc (e.g. for forbidden cases)
				mockCardRepo.PerformTransactionFunc = func(fn func(tx *gorm.DB) error) error {
                    // This transaction should ideally not be called if permissions fail before it.
                    // If it is called, setting deleteTransactionCalled helps verify.
					deleteTransactionCalled = true
					return fn(&gorm.DB{}) // Default, might panic if used
				}
			}
			// The mockCardRepo.DeleteFunc is not directly used by the service anymore if tx.Delete is used.
			// So, we don't need to set it up unless PerformTransaction itself calls the repo's Delete.

			mockListRepo := &MockListRepositoryForCardService{GetBoardIDByListIDFunc: func(lID uint) (uint, error) { return boardID, nil }}
			mockBoardRepo := &MockBoardRepositoryForCardService{FindByIDFunc: tt.mockBoardFindByIDFunc}
			mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{IsMemberFunc: tt.mockIsMemberFunc}
			mockUserRepo := &MockUserRepositoryForCardService{}

			service := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
			err := service.DeleteCard(cardID, tt.currentUserID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectDeleteTransaction, deleteTransactionCalled)
		})
	}
}

func TestPlaceholder_CardService(t *testing.T) {
	assert.True(t, true, "This is a placeholder test.")
}

func TestCardService_CreateCard_WithColor(t *testing.T) {
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
	cardColor := "#123456"

	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}

	var createdCardModel models.Card
	mockCardRepo.CreateFunc = func(card *models.Card) error {
		card.ID = 1
		card.Position = 1
		card.Status = models.StatusToDo
		createdCardModel = *card
		return nil
	}
	mockCardRepo.GetMaxPositionFunc = func(lID uint) (uint, error) { return 0, nil }

	mockCardRepo.FindByIDFunc = func(id uint) (*models.Card, error) {
		assert.Equal(t, createdCardModel.ID, id)
		return &createdCardModel, nil
	}

	card, err := cardService.CreateCard(listID, cardTitle, "", nil, nil, nil, nil, &cardColor, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, cardTitle, card.Title)
	assert.NotNil(t, card.Color)
	assert.Equal(t, cardColor, *card.Color)
	assert.Equal(t, uint(1), card.ID)
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

	card, err := cardService.CreateCard(listID, cardTitle, "", nil, nil, nil, nil, nil, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, cardTitle, card.Title)
	assert.Nil(t, card.Color)
	assert.Equal(t, uint(2), card.ID)
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
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }


	expectedTargetUser := &models.User{Model: gorm.Model{ID: targetUserID}, Email: targetUserEmail, Username: "collabUser"}
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
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }

	expectedTargetUser := &models.User{Model: gorm.Model{ID: targetUserID}, Username: "collabUserByID"}
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) { return expectedTargetUser, nil }
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return false, nil }
	mockCardRepo.AddCollaboratorFunc = func(cID uint, uID uint) error { return nil }

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, "", &targetUserID)
	assert.NoError(t, err)
	assert.NotNil(t, addedUser)
	assert.Equal(t, targetUserID, addedUser.ID)
}

func TestCardService_AddCollaboratorToCard_PermissionDenied_NotOwner(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(2)
	ownerID := uint(1)
	cardID := uint(100)
	listID := uint(10)
	boardID := uint(1)
	targetUserEmail := "c@e.com"

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		if bID == boardID && uID == currentUserID {
			return true, nil
		}
		return false, nil
	}


	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, targetUserEmail, nil)
	assert.ErrorIs(t, err, ErrForbidden, "Expected ErrForbidden for non-owner trying to add collaborator")
	assert.Nil(t, addedUser)
}


func TestCardService_RemoveCollaboratorFromCard_Success(t *testing.T) {
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
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return true, nil }

	removeCalled := false
	mockCardRepo.RemoveCollaboratorFunc = func(cID uint, uID uint) error { removeCalled = true; return nil }

	err := cardService.RemoveCollaboratorFromCard(cardID, currentUserID, targetUserID)
	assert.NoError(t, err)
	assert.True(t, removeCalled)
}

func TestCardService_RemoveCollaboratorFromCard_PermissionDenied_NotOwner(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)

	currentUserID := uint(2)
	ownerID := uint(1)
	cardID := uint(100)
	listID := uint(10)
	boardID := uint(1)
	targetUserID := uint(5)

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		if bID == boardID && uID == currentUserID { return true, nil }
		return false, nil
	}

	err := cardService.RemoveCollaboratorFromCard(cardID, currentUserID, targetUserID)
	assert.ErrorIs(t, err, ErrForbidden)
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
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) {
		return &models.User{Model: gorm.Model{ID: targetUserID}}, nil
	}
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return true, nil }
	addCollabCalled := false
	mockCardRepo.AddCollaboratorFunc = func(cID uint, uID uint) error { addCollabCalled = true; return nil }

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, "", &targetUserID)
	assert.NoError(t, err)
	assert.NotNil(t, addedUser)
	assert.False(t, addCollabCalled, "AddCollaborator should not be called if user is already a collaborator")
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
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }
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
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) { return nil, gorm.ErrRecordNotFound }

	addedUser, err := cardService.AddCollaboratorToCard(cardID, currentUserID, "", &targetUserID)
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, addedUser)
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
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }
	mockCardRepo.IsCollaboratorFunc = func(cID uint, uID uint) (bool, error) { return false, nil }

	err := cardService.RemoveCollaboratorFromCard(cardID, currentUserID, targetUserID)
	assert.ErrorIs(t, err, ErrUserNotCollaborator)
}


func TestCardService_GetCardCollaborators_Success(t *testing.T) {
	mockCardRepo := &MockCardRepository{}
	mockListRepo := &MockListRepositoryForCardService{}
	mockBoardRepo := &MockBoardRepositoryForCardService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForCardService{}
	mockUserRepo := &MockUserRepositoryForCardService{}
	cardService := NewCardService(mockCardRepo, mockListRepo, mockBoardRepo, mockBoardMemberRepo, mockUserRepo)
	currentUserID, cardID, listID, boardID := uint(1), uint(100), uint(10), uint(1)
	expectedUsers := []models.User{{Model: gorm.Model{ID: 5}}, {Model: gorm.Model{ID: 6}}}

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }
	mockCardRepo.IsUserCollaboratorOrAssigneeFunc = func(cID uint, uID uint) (bool, error) { return true, nil }

	mockCardRepo.GetCollaboratorsByCardIDFunc = func(cID uint) ([]models.User, error) { return expectedUsers, nil }

	users, err := cardService.GetCardCollaborators(cardID, currentUserID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestCardService_GetCardCollaborators_PermissionDenied_NotMember(t *testing.T) {
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

	mockCardRepo.GetListIDByCardIDFunc = func(cID uint) (uint, error) { return listID, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(bID uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: currentUserID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }
	mockCardRepo.IsUserCollaboratorOrAssigneeFunc = func(cID uint, uID uint) (bool, error) { return true, nil }


	var cardStateAfterInitialFind models.Card
	var cardStateForFinalFind models.Card
	var cardStateCapturedByUpdate models.Card

	mockCardRepo.FindByIDFunc = func(cID uint) (*models.Card, error) {
		if cID == cardID && cardStateAfterInitialFind.ID == 0 {
			cardStateAfterInitialFind = *initialCard
			return &cardStateAfterInitialFind, nil
		}
		assert.Equal(t, cardID, cID)
		cardStateForFinalFind = cardStateCapturedByUpdate
		return &cardStateForFinalFind, nil
	}

	mockCardRepo.UpdateFunc = func(card *models.Card) error {
		assert.NotNil(t, card.Color)
		assert.Equal(t, newColor, *card.Color)
		cardStateCapturedByUpdate = *card
		return nil
	}

	updatedCard, err := cardService.UpdateCard(cardID, nil, nil, nil, nil, nil, nil, nil, &newColor, currentUserID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedCard)
	assert.NotNil(t, updatedCard.Color)
	assert.Equal(t, newColor, *updatedCard.Color)
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
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }
	mockCardRepo.IsUserCollaboratorOrAssigneeFunc = func(cID uint, uID uint) (bool, error) { return true, nil }

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
		assert.NotNil(t, card.DueDate)
		assert.True(t, initialDueDate.Equal(*card.DueDate))
		assert.Equal(t, newTitle, card.Title)
		cardStateCapturedByUpdate = *card
		return nil
	}

	updatedCard, err := cardService.UpdateCard(cardID, &newTitle, nil, nil, nil, nil, nil, nil, nil, currentUserID)

	assert.NoError(t, err)
	assert.True(t, updateCalled, "UpdateFunc should be called")
	assert.NotNil(t, updatedCard)
	assert.NotNil(t, updatedCard.DueDate, "DueDate should not have been cleared")
	assert.True(t, initialDueDate.Equal(*updatedCard.DueDate), "DueDate should be unchanged")
	assert.Equal(t, newTitle, updatedCard.Title)
}
