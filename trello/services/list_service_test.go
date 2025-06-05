package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"
	"gorm.io/gorm"
	// No longer importing sqlite or logger here
)

// --- MockGormDB for transaction testing ---
// This is a simplified mock for gorm.DB to handle the Transaction method.
type MockGormDB struct {
	gorm.DB
	TransactionFunc func(fc func(tx *gorm.DB) error) error
}

func (m *MockGormDB) Transaction(fc func(tx *gorm.DB) error) error {
	if m.TransactionFunc != nil {
		return m.TransactionFunc(fc)
	}
	return fc(&gorm.DB{})
}


// --- MockListRepository ---
type MockListRepository struct {
	CreateFunc             func(list *models.List) error
	FindByIDFunc           func(id uint) (*models.List, error)
	FindByBoardIDFunc      func(boardID uint) ([]models.List, error)
	UpdateFunc             func(list *models.List) error
	DeleteFunc             func(id uint) error
	GetMaxPositionFunc     func(boardID uint) (uint, error)
	GetBoardIDByListIDFunc func(listID uint) (uint, error)
	PerformTransactionFunc func(fn func(tx *gorm.DB) error) error

	CreateCalledWithList         *models.List
	FindByIDCalledWithID         uint
	FindByBoardIDCalledWithID    uint
	UpdateCalledWithList         *models.List
	DeleteCalledWithID           uint
	GetMaxPositionCalledWithID   uint
	GetBoardIDByListIDCalledWith uint
}

func (m *MockListRepository) GetBoardIDByListID(listID uint) (uint, error) {
	m.GetBoardIDByListIDCalledWith = listID
	if m.GetBoardIDByListIDFunc != nil {
		return m.GetBoardIDByListIDFunc(listID)
	}
	return 0, errors.New("GetBoardIDByListIDFunc not implemented in mock")
}

func (m *MockListRepository) Create(list *models.List) error {
	m.CreateCalledWithList = list
	if list.Position == 0 {
		if m.GetMaxPositionFunc != nil {
			maxPos, _ := m.GetMaxPosition(list.BoardID)
			list.Position = maxPos + 1
		} else {
			list.Position = 1
		}
	}
	if m.CreateFunc != nil {
		return m.CreateFunc(list)
	}
	return nil
}
func (m *MockListRepository) FindByID(id uint) (*models.List, error) {
	m.FindByIDCalledWithID = id
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc not implemented")
}
func (m *MockListRepository) FindByBoardID(boardID uint) ([]models.List, error) {
	m.FindByBoardIDCalledWithID = boardID
	if m.FindByBoardIDFunc != nil {
		return m.FindByBoardIDFunc(boardID)
	}
	return nil, errors.New("FindByBoardIDFunc not implemented")
}
func (m *MockListRepository) Update(list *models.List) error {
	m.UpdateCalledWithList = list
	if m.UpdateFunc != nil {
		return m.UpdateFunc(list)
	}
	return nil
}
func (m *MockListRepository) Delete(id uint) error {
	m.DeleteCalledWithID = id
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockListRepository) GetMaxPosition(boardID uint) (uint, error) {
	m.GetMaxPositionCalledWithID = boardID
	if m.GetMaxPositionFunc != nil {
		return m.GetMaxPositionFunc(boardID)
	}
	return 0, errors.New("GetMaxPositionFunc not implemented")
}
func (m *MockListRepository) PerformTransaction(fn func(tx *gorm.DB) error) error {
	if m.PerformTransactionFunc != nil {
		return m.PerformTransactionFunc(fn)
	}
	return fn(&gorm.DB{})
}
var _ repositories.ListRepositoryInterface = (*MockListRepository)(nil)


// --- MockBoardRepositoryForListService ---
type MockBoardRepositoryForListService struct {
	repositories.BoardRepositoryInterface
	FindByIDFunc                          func(id uint) (*models.Board, error)
	IsOwnerFunc                           func(boardID uint, userID uint) (bool, error)
}
func (m *MockBoardRepositoryForListService) FindByID(id uint) (*models.Board, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, errors.New("FindByIDFunc on MockBoardRepositoryForListService not implemented")
}
func (m *MockBoardRepositoryForListService) IsOwner(boardID uint, userID uint) (bool, error) {
	if m.IsOwnerFunc != nil { return m.IsOwnerFunc(boardID, userID) }
	return false, errors.New("IsOwnerFunc on MockBoardRepositoryForListService not implemented")
}
func (m *MockBoardRepositoryForListService) Create(board *models.Board) error { return errors.New("not implemented") }
func (m *MockBoardRepositoryForListService) FindByOwnerOrMember(userID uint) ([]models.Board, error) { return nil, errors.New("not implemented") }
func (m *MockBoardRepositoryForListService) Update(board *models.Board) error { return errors.New("not implemented") }
func (m *MockBoardRepositoryForListService) Delete(id uint) error { return errors.New("not implemented") }


// --- MockBoardMemberRepositoryForListService ---
type MockBoardMemberRepositoryForListService struct {
	repositories.BoardMemberRepositoryInterface
	IsMemberFunc                                func(boardID uint, userID uint) (bool, error)
}
func (m *MockBoardMemberRepositoryForListService) IsMember(boardID uint, userID uint) (bool, error) {
	if m.IsMemberFunc != nil { return m.IsMemberFunc(boardID, userID) }
	return false, errors.New("IsMemberFunc on MockBoardMemberRepositoryForListService not implemented")
}
func (m *MockBoardMemberRepositoryForListService) AddMember(member *models.BoardMember) error { return errors.New("not implemented") }
func (m *MockBoardMemberRepositoryForListService) RemoveMember(boardID uint, userID uint) error { return errors.New("not implemented") }
func (m *MockBoardMemberRepositoryForListService) FindMembersByBoardID(boardID uint) ([]models.BoardMember, error) { return nil, errors.New("not implemented") }
func (m *MockBoardMemberRepositoryForListService) FindByBoardIDAndUserID(boardID uint, userID uint) (*models.BoardMember, error) { return nil, errors.New("not implemented") }


// TestListService_CreateList_Success and other tests remain the same until UpdateList/DeleteList tests

func TestListService_CreateList_Success(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}

	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)

	userID := uint(1)
	boardID := uint(10)
	listName := "New List"

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil
	}
	createdList := &models.List{Model: gorm.Model{ID: 100}, Name: listName, BoardID: boardID, Position: 1}
	mockListRepo.CreateFunc = func(list *models.List) error {
		list.ID = createdList.ID
		list.Position = createdList.Position
		return nil
	}
	mockListRepo.GetMaxPositionFunc = func(bID uint) (uint, error) {
		return 0, nil
	}
	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) {
		return createdList, nil
	}
	list, err := listService.CreateList(listName, boardID, userID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, createdList.ID, list.ID)
}

func TestListService_CreateList_PermissionDenied_NotOwnerNotMember(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, ownerID, listName := uint(1), uint(10), uint(2), "New List"

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, nil
	}
	mockListRepo.CreateFunc = func(list *models.List) error {
		t.Error("listRepo.Create should not be called")
		return nil
	}
	list, err := listService.CreateList(listName, boardID, userID, nil)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, list)
}

func TestListService_CreateList_BoardNotFound(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, listName := uint(1), uint(10), "New List"

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return nil, gorm.ErrRecordNotFound
	}
	mockListRepo.CreateFunc = func(list *models.List) error {
		t.Error("listRepo.Create should not be called")
		return nil
	}
	list, err := listService.CreateList(listName, boardID, userID, nil)
	assert.ErrorIs(t, err, ErrBoardNotFound)
	assert.Nil(t, list)
}

func TestListService_CreateList_ErrOnCreate(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, listName, expectedError := uint(1), uint(10), "New List", errors.New("DB error")

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil
	}
	mockListRepo.CreateFunc = func(list *models.List) error { return expectedError }
	mockListRepo.GetMaxPositionFunc = func(bID uint) (uint, error) { return 0, nil }

	list, err := listService.CreateList(listName, boardID, userID, nil)
	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, list)
}

func TestListService_CreateList_ErrOnFinalFindByID(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, listName, expectedError := uint(1), uint(10), "New List", errors.New("DB error")
	createdListID := uint(100)

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil
	}
	mockListRepo.CreateFunc = func(list *models.List) error { list.ID = createdListID; return nil }
	mockListRepo.GetMaxPositionFunc = func(bID uint) (uint, error) { return 0, nil }
	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) { return nil, expectedError }

	list, err := listService.CreateList(listName, boardID, userID, nil)
	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, list)
}

func TestListService_GetListsByBoardID_Success(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID := uint(1), uint(10)
	expectedLists := []models.List{{Model: gorm.Model{ID:1}}}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil }
	mockListRepo.FindByBoardIDFunc = func(bID uint) ([]models.List, error) { return expectedLists, nil }

	lists, err := listService.GetListsByBoardID(boardID, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedLists, lists)
}

func TestListService_GetListsByBoardID_PermissionDenied(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, ownerID := uint(1), uint(10), uint(2)

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil }

	lists, err := listService.GetListsByBoardID(boardID, userID)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, lists)
}

func TestListService_GetListsByBoardID_RepoError(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, expectedError := uint(1), uint(10), errors.New("DB error")

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil }
	mockListRepo.FindByBoardIDFunc = func(bID uint) ([]models.List, error) { return nil, expectedError }

	lists, err := listService.GetListsByBoardID(boardID, userID)
	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, lists)
}

func TestListService_GetListByID_Success(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, listID := uint(1), uint(10), uint(100)
	expectedList := &models.List{Model: gorm.Model{ID: listID}, BoardID: boardID}

	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) { return expectedList, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil }

	list, err := listService.GetListByID(listID, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, list)
}

func TestListService_GetListByID_ListNotFound(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, listID := uint(1), uint(100)

	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) { return nil, gorm.ErrRecordNotFound }

	list, err := listService.GetListByID(listID, userID)
	assert.ErrorIs(t, err, ErrListNotFound)
	assert.Nil(t, list)
}

func TestListService_GetListByID_PermissionDenied(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, listID, actualOwnerID := uint(1), uint(10), uint(100), uint(2)
	foundList := &models.List{Model: gorm.Model{ID: listID}, BoardID: boardID}

	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) { return foundList, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil }

	list, err := listService.GetListByID(listID, userID)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, list)
}

func TestListService_UpdateList_TransactionError(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, listID, boardID, originalPosition, newPosition := uint(1), uint(100), uint(10), uint(1), uint(2)
	expectedError := errors.New("transaction failed")
	originalList := &models.List{Model: gorm.Model{ID: listID}, Name: "Test List", BoardID: boardID, Position: originalPosition}

	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) { listCopy := *originalList; return &listCopy, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil }
	mockListRepo.GetMaxPositionFunc = func(bID uint) (uint, error) { return 3, nil }

	mockListRepo.PerformTransactionFunc = func(fn func(tx *gorm.DB) error) error {
		return expectedError
	}

	list, err := listService.UpdateList(listID, nil, &newPosition, userID)
	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, list)
}

func TestListService_DeleteList_Success(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, listID, boardID := uint(1), uint(100), uint(10)
	listToDelete := &models.List{Model: gorm.Model{ID: listID}, BoardID: boardID, Position: 1}

	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) { return listToDelete, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil }

	dbForTx := setupTestDB(t)
	mockListRepo.PerformTransactionFunc = func(fn func(tx *gorm.DB) error) error {
		return fn(dbForTx)
	}
	deleteCalledOnRepo := false
	mockListRepo.DeleteFunc = func(id uint) error { deleteCalledOnRepo = true; return nil }

	err := listService.DeleteList(listID, userID)
	assert.NoError(t, err)
	assert.True(t, deleteCalledOnRepo)
}

func TestListService_UpdateList_SuccessNameChange(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, listID, originalName, newName := uint(1), uint(10), uint(100), "Original", "Updated"
	originalList := &models.List{Model: gorm.Model{ID: listID}, Name: originalName, BoardID: boardID, Position: 1}
	updatedList := &models.List{Model: gorm.Model{ID: listID}, Name: newName, BoardID: boardID, Position: 1}

	var findByIdCallCount int
	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) {
		findByIdCallCount++
		if findByIdCallCount == 1 { listCopy := *originalList; return &listCopy, nil }
		return updatedList, nil
	}
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil }
	mockListRepo.UpdateFunc = func(list *models.List) error { return nil }

	list, err := listService.UpdateList(listID, &newName, nil, userID)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, newName, list.Name)
}

func TestListService_UpdateList_SuccessPositionChange(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, boardID, listID, originalPosition, newPos := uint(1), uint(10), uint(100), uint(1), uint(2)
	originalList := &models.List{Model: gorm.Model{ID: listID}, Name: "Test List", BoardID: boardID, Position: originalPosition}
	listAfterTxSave := &models.List{Model: gorm.Model{ID: listID}, Name: "Test List", BoardID: boardID, Position: newPos}

	var initialFindByIDCalled bool
	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) {
		if !initialFindByIDCalled { initialFindByIDCalled = true; listCopy := *originalList; return &listCopy, nil }
		return listAfterTxSave, nil
	}
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: userID}, nil }
	mockListRepo.GetMaxPositionFunc = func(bID uint) (uint, error) { return 3, nil }

	dbForTx := setupTestDB(t)
	mockListRepo.PerformTransactionFunc = func(fn func(tx *gorm.DB) error) error {
		return fn(dbForTx)
	}
	mockListRepo.UpdateFunc = func(l *models.List) error { t.Error("listRepo.Update should not be called"); return nil }

	list, err := listService.UpdateList(listID, nil, &newPos, userID)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, newPos, list.Position)
}

func TestListService_UpdateList_ListNotFound(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, listID, newName := uint(1), uint(100), "New Name"

	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) { return nil, gorm.ErrRecordNotFound }

	list, err := listService.UpdateList(listID, &newName, nil, userID)
	assert.ErrorIs(t, err, ErrListNotFound)
	assert.Nil(t, list)
}

func TestListService_UpdateList_PermissionDenied(t *testing.T) {
	mockListRepo := &MockListRepository{}
	mockBoardRepo := &MockBoardRepositoryForListService{}
	mockBoardMemberRepo := &MockBoardMemberRepositoryForListService{}
	listService := NewListService(mockListRepo, mockBoardRepo, mockBoardMemberRepo)
	userID, listID, boardID, actualOwnerID, newName := uint(1), uint(100), uint(10), uint(2), "New Name"
	foundList := &models.List{Model: gorm.Model{ID: listID}, BoardID: boardID}

	mockListRepo.FindByIDFunc = func(id uint) (*models.List, error) { return foundList, nil }
	mockListRepo.GetBoardIDByListIDFunc = func(lID uint) (uint, error) { return boardID, nil }
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil }

	list, err := listService.UpdateList(listID, &newName, nil, userID)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, list)
}
