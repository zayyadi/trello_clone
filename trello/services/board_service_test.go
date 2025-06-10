package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/realtime"
	"github.com/zayyadi/trello/repositories"
	"gorm.io/gorm"
)

// MockBoardRepository is a mock implementation of BoardRepositoryInterface
type MockBoardRepository struct {
	CreateFunc              func(board *models.Board) error
	FindByIDFunc            func(id uint) (*models.Board, error)
	FindByOwnerOrMemberFunc func(userID uint) ([]models.Board, error)
	UpdateFunc              func(board *models.Board) error
	DeleteFunc              func(id uint) error
	IsOwnerFunc             func(boardID uint, userID uint) (bool, error)

	// Store calls
	CreateCalledWith              *models.Board
	FindByIDCalledWith            uint
	FindByOwnerOrMemberCalledWith uint
	UpdateCalledWith              *models.Board
	DeleteCalledWith              uint
	IsOwnerCalledWithBoardID      uint
	IsOwnerCalledWithUserID       uint
}

func (m *MockBoardRepository) Create(board *models.Board) error {
	m.CreateCalledWith = board
	if m.CreateFunc != nil {
		return m.CreateFunc(board)
	}
	return nil
}
func (m *MockBoardRepository) FindByID(id uint) (*models.Board, error) {
	m.FindByIDCalledWith = id
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc not implemented")
}
func (m *MockBoardRepository) FindByOwnerOrMember(userID uint) ([]models.Board, error) {
	m.FindByOwnerOrMemberCalledWith = userID
	if m.FindByOwnerOrMemberFunc != nil {
		return m.FindByOwnerOrMemberFunc(userID)
	}
	return nil, errors.New("FindByOwnerOrMemberFunc not implemented")
}
func (m *MockBoardRepository) Update(board *models.Board) error {
	m.UpdateCalledWith = board
	if m.UpdateFunc != nil {
		return m.UpdateFunc(board)
	}
	return nil
}
func (m *MockBoardRepository) Delete(id uint) error {
	m.DeleteCalledWith = id
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockBoardRepository) IsOwner(boardID uint, userID uint) (bool, error) {
	m.IsOwnerCalledWithBoardID = boardID
	m.IsOwnerCalledWithUserID = userID
	if m.IsOwnerFunc != nil {
		return m.IsOwnerFunc(boardID, userID)
	}
	return false, errors.New("IsOwnerFunc not implemented")
}

var _ repositories.BoardRepositoryInterface = (*MockBoardRepository)(nil)

// MockBoardMemberRepository is a mock implementation of BoardMemberRepositoryInterface
type MockBoardMemberRepository struct {
	AddMemberFunc              func(member *models.BoardMember) error
	RemoveMemberFunc           func(boardID uint, userID uint) error
	IsMemberFunc               func(boardID uint, userID uint) (bool, error)
	FindMembersByBoardIDFunc   func(boardID uint) ([]models.BoardMember, error)
	FindByBoardIDAndUserIDFunc func(boardID uint, userID uint) (*models.BoardMember, error)

	// Store calls
	AddMemberCalledWithMember               *models.BoardMember
	RemoveMemberCalledWithBoardID           uint
	RemoveMemberCalledWithUserID            uint
	IsMemberCalledWithBoardID               uint
	IsMemberCalledWithUserID                uint
	FindMembersByBoardIDCalledWith          uint
	FindByBoardIDAndUserIDCalledWithBoardID uint
	FindByBoardIDAndUserIDCalledWithUserID  uint
}

func (m *MockBoardMemberRepository) AddMember(member *models.BoardMember) error {
	m.AddMemberCalledWithMember = member
	if m.AddMemberFunc != nil {
		return m.AddMemberFunc(member)
	}
	return nil
}
func (m *MockBoardMemberRepository) RemoveMember(boardID uint, userID uint) error {
	m.RemoveMemberCalledWithBoardID = boardID
	m.RemoveMemberCalledWithUserID = userID
	if m.RemoveMemberFunc != nil {
		return m.RemoveMemberFunc(boardID, userID)
	}
	return nil
}
func (m *MockBoardMemberRepository) IsMember(boardID uint, userID uint) (bool, error) {
	m.IsMemberCalledWithBoardID = boardID
	m.IsMemberCalledWithUserID = userID
	if m.IsMemberFunc != nil {
		return m.IsMemberFunc(boardID, userID)
	}
	return false, errors.New("IsMemberFunc not implemented")
}
func (m *MockBoardMemberRepository) FindMembersByBoardID(boardID uint) ([]models.BoardMember, error) {
	m.FindMembersByBoardIDCalledWith = boardID
	if m.FindMembersByBoardIDFunc != nil {
		return m.FindMembersByBoardIDFunc(boardID)
	}
	return nil, errors.New("FindMembersByBoardIDFunc not implemented")
}
func (m *MockBoardMemberRepository) FindByBoardIDAndUserID(boardID uint, userID uint) (*models.BoardMember, error) {
	m.FindByBoardIDAndUserIDCalledWithBoardID = boardID
	m.FindByBoardIDAndUserIDCalledWithUserID = userID
	if m.FindByBoardIDAndUserIDFunc != nil {
		return m.FindByBoardIDAndUserIDFunc(boardID, userID)
	}
	return nil, errors.New("FindByBoardIDAndUserIDFunc not implemented")
}

var _ repositories.BoardMemberRepositoryInterface = (*MockBoardMemberRepository)(nil)

// MockUserRepositoryForBoardService is a mock implementation of UserRepositoryInterface for BoardService tests
type MockUserRepositoryForBoardService struct {
	CreateFunc      func(user *models.User) error
	FindByEmailFunc func(email string) (*models.User, error)
	FindByIDFunc    func(id uint) (*models.User, error)
	FindAllFunc     func() ([]models.User, error) // Add FindAllFunc

	// Store calls
	CreateCalledWithUser  *models.User
	FindByEmailCalledWith string
	FindByIDCalledWith    uint
}

func (m *MockUserRepositoryForBoardService) Create(user *models.User) error {
	m.CreateCalledWithUser = user
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return errors.New("CreateFunc not implemented by MockUserRepositoryForBoardService")
}
func (m *MockUserRepositoryForBoardService) FindByEmail(email string) (*models.User, error) {
	m.FindByEmailCalledWith = email
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, errors.New("FindByEmailFunc not implemented")
}
func (m *MockUserRepositoryForBoardService) FindByID(id uint) (*models.User, error) {
	m.FindByIDCalledWith = id
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("FindByIDFunc not implemented")
}

func (m *MockUserRepositoryForBoardService) FindAll() ([]models.User, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, nil // Default behavior: return empty slice, no error
}

// MockHub is a mock implementation of realtime.Hub for testing purposes
type MockHub struct {
	SubmitFunc func(msg *realtime.WebSocketMessage)
	RunFunc    func()
	Register   chan *realtime.Client
	Unregister chan *realtime.Client
}

func (m *MockHub) Submit(msg *realtime.WebSocketMessage) {
	if m.SubmitFunc != nil {
		m.SubmitFunc(msg)
	}
}

func (m *MockHub) Run() {
	if m.RunFunc != nil {
		m.RunFunc()
	}
}

var _ repositories.UserRepositoryInterface = (*MockUserRepositoryForBoardService)(nil)

func TestBoardService_CreateBoard_Success(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	ownerID := uint(1)
	boardName := "Test Board"
	boardDescription := "Test Description"

	mockBoardRepo.CreateFunc = func(board *models.Board) error {
		assert.Equal(t, boardName, board.Name)
		assert.Equal(t, boardDescription, board.Description)
		assert.Equal(t, ownerID, board.OwnerID)
		board.ID = 100
		return nil
	}

	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		assert.Equal(t, uint(100), member.BoardID)
		assert.Equal(t, ownerID, member.UserID)
		return nil
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, uint(100), id)
		return &models.Board{
			Model:       gorm.Model{ID: id},
			Name:        boardName,
			Description: boardDescription,
			OwnerID:     ownerID,
			Owner:       models.User{Model: gorm.Model{ID: ownerID}, Username: "owner"},
			Members: []models.BoardMember{
				{BoardID: id, UserID: ownerID, User: models.User{Model: gorm.Model{ID: ownerID}, Username: "owner"}},
			},
		}, nil
	}

	createdBoard, err := boardService.CreateBoard(boardName, boardDescription, ownerID)

	assert.NoError(t, err)
	assert.NotNil(t, createdBoard)
	assert.Equal(t, uint(100), createdBoard.ID)
	assert.Equal(t, boardName, createdBoard.Name)
	assert.Equal(t, ownerID, createdBoard.OwnerID)
	assert.Equal(t, ownerID, createdBoard.Owner.ID)
	assert.Len(t, createdBoard.Members, 1)
	if len(createdBoard.Members) > 0 {
		assert.Equal(t, ownerID, createdBoard.Members[0].UserID)
		assert.Equal(t, ownerID, createdBoard.Members[0].User.ID)
	}

	assert.NotNil(t, mockBoardRepo.CreateCalledWith)
	assert.NotNil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
	assert.Equal(t, uint(100), mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_CreateBoard_ErrOnBoardCreate(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	ownerID := uint(1)
	boardName := "Test Board"
	boardDescription := "Test Description"
	expectedError := errors.New("DB error on board create")

	mockBoardRepo.CreateFunc = func(board *models.Board) error {
		return expectedError
	}

	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		t.Error("AddMember should not be called if board creation fails")
		return nil
	}
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		t.Error("FindByID should not be called if board creation fails")
		return nil, nil
	}

	createdBoard, err := boardService.CreateBoard(boardName, boardDescription, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, createdBoard)

	assert.NotNil(t, mockBoardRepo.CreateCalledWith)
	assert.Nil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
	assert.Equal(t, uint(0), mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_CreateBoard_ErrOnAddOwnerAsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	ownerID := uint(1)
	boardName := "Test Board"
	boardDescription := "Test Description"
	expectedError := errors.New("DB error on add member")

	mockBoardRepo.CreateFunc = func(board *models.Board) error {
		board.ID = 100
		return nil
	}

	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		return expectedError
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		t.Error("FindByID should not be called if AddMember fails")
		return nil, nil
	}

	createdBoard, err := boardService.CreateBoard(boardName, boardDescription, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, createdBoard)

	assert.NotNil(t, mockBoardRepo.CreateCalledWith)
	assert.NotNil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
	assert.Equal(t, uint(0), mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_GetBoardByID_SuccessAsOwner(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)

	expectedBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Owned Board",
		OwnerID: ownerID,
		Owner:   models.User{Model: gorm.Model{ID: ownerID}},
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return expectedBoard, nil
	}

	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called if user is owner")
		return false, nil
	}

	board, err := boardService.GetBoardByID(boardID, ownerID)

	assert.NoError(t, err)
	assert.NotNil(t, board)
	assert.Equal(t, expectedBoard.ID, board.ID)
	assert.Equal(t, expectedBoard.Name, board.Name)
	assert.Equal(t, ownerID, board.OwnerID)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, uint(0), mockBoardMemberRepo.IsMemberCalledWithBoardID)
}

func TestBoardService_GetBoardByID_SuccessAsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	boardOwnerID := uint(10)
	currentUserID := uint(20)

	expectedBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Member Board",
		OwnerID: boardOwnerID,
		Owner:   models.User{Model: gorm.Model{ID: boardOwnerID}},
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return expectedBoard, nil
	}

	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return true, nil
	}

	board, err := boardService.GetBoardByID(boardID, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, board)
	assert.Equal(t, expectedBoard.ID, board.ID)
	assert.Equal(t, expectedBoard.Name, board.Name)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, currentUserID, mockBoardMemberRepo.IsMemberCalledWithUserID)
}

func TestBoardService_GetBoardByID_IsMemberError(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	boardOwnerID := uint(10)
	currentUserID := uint(20)
	expectedError := errors.New("DB error on IsMember")

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Other's Board",
		OwnerID: boardOwnerID,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return false, expectedError
	}

	board, err := boardService.GetBoardByID(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, board)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, currentUserID, mockBoardMemberRepo.IsMemberCalledWithUserID)
}

func TestBoardService_GetBoardsForUser_Success(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	userID := uint(1)
	expectedBoards := []models.Board{
		{Model: gorm.Model{ID: 1}, Name: "Board 1", OwnerID: userID},
		{Model: gorm.Model{ID: 2}, Name: "Board 2", OwnerID: uint(2), Members: []models.BoardMember{{UserID: userID}}},
	}

	mockBoardRepo.FindByOwnerOrMemberFunc = func(uID uint) ([]models.Board, error) {
		assert.Equal(t, userID, uID)
		return expectedBoards, nil
	}

	boards, err := boardService.GetBoardsForUser(userID)

	assert.NoError(t, err)
	assert.NotNil(t, boards)
	assert.Len(t, boards, 2)
	if len(boards) == 2 {
		assert.Equal(t, expectedBoards[0].ID, boards[0].ID)
		assert.Equal(t, expectedBoards[0].Name, boards[0].Name)
		assert.Equal(t, expectedBoards[1].ID, boards[1].ID)
		assert.Equal(t, expectedBoards[1].Name, boards[1].Name)
	}
	assert.Equal(t, userID, mockBoardRepo.FindByOwnerOrMemberCalledWith)
}

func TestBoardService_GetBoardsForUser_NoBoards(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	userID := uint(1)
	expectedBoards := []models.Board{}

	mockBoardRepo.FindByOwnerOrMemberFunc = func(uID uint) ([]models.Board, error) {
		assert.Equal(t, userID, uID)
		return expectedBoards, nil
	}

	boards, err := boardService.GetBoardsForUser(userID)

	assert.NoError(t, err)
	assert.NotNil(t, boards)
	assert.Len(t, boards, 0)
	assert.Equal(t, userID, mockBoardRepo.FindByOwnerOrMemberCalledWith)
}

func TestBoardService_GetBoardsForUser_Error(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	userID := uint(1)
	expectedError := errors.New("DB error FindByOwnerOrMember")

	mockBoardRepo.FindByOwnerOrMemberFunc = func(uID uint) ([]models.Board, error) {
		assert.Equal(t, userID, uID)
		return nil, expectedError
	}

	boards, err := boardService.GetBoardsForUser(userID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, boards)
	assert.Equal(t, userID, mockBoardRepo.FindByOwnerOrMemberCalledWith)
}

func TestBoardService_UpdateBoard_Success(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	currentUserID := uint(10)
	newName := "Updated Board Name"
	newDescription := "Updated Description"

	originalBoard := &models.Board{
		Model:       gorm.Model{ID: boardID},
		Name:        "Original Name",
		Description: "Original Description",
		OwnerID:     currentUserID,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		if id == boardID && mockBoardRepo.UpdateCalledWith == nil {
			return originalBoard, nil
		} else if id == boardID && mockBoardRepo.UpdateCalledWith != nil {
			updatedBoardFromMock := *mockBoardRepo.UpdateCalledWith
			updatedBoardFromMock.Model.ID = id
			return &updatedBoardFromMock, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		assert.Equal(t, boardID, board.ID)
		assert.Equal(t, newName, board.Name)
		assert.Equal(t, newDescription, board.Description)
		assert.Equal(t, currentUserID, board.OwnerID)
		return nil
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, &newDescription, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, updatedBoard)
	assert.Equal(t, boardID, updatedBoard.ID)
	assert.Equal(t, newName, updatedBoard.Name)
	assert.Equal(t, newDescription, updatedBoard.Description)

	assert.NotNil(t, mockBoardRepo.UpdateCalledWith)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_UpdateBoard_BoardNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	currentUserID := uint(10)
	newName := "Updated Name"

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return nil, gorm.ErrRecordNotFound
	}

	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		t.Error("Update should not be called if board not found")
		return nil
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardNotFound, err)
	assert.Nil(t, updatedBoard)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Nil(t, mockBoardRepo.UpdateCalledWith)
}

func TestBoardService_UpdateBoard_Forbidden(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(5)
	currentUserID := uint(10)
	newName := "Updated Name"

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Original Name",
		OwnerID: ownerID,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		t.Error("Update should not be called if user is not owner")
		return nil
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, updatedBoard)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Nil(t, mockBoardRepo.UpdateCalledWith)
}

func TestBoardService_UpdateBoard_ErrOnUpdate(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	currentUserID := uint(10)
	newName := "Updated Name"
	expectedError := errors.New("DB error on update")

	originalBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Original Name",
		OwnerID: currentUserID,
	}

	var findByIDCallCount = 0
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		findByIDCallCount++
		if id == boardID && findByIDCallCount == 1 {
			return originalBoard, nil
		}
		if findByIDCallCount == 2 {
			t.Error("Second FindByID should not be called if Update fails")
		}
		return nil, gorm.ErrRecordNotFound
	}

	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		return expectedError
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, updatedBoard)

	assert.Equal(t, 1, findByIDCallCount)
	assert.NotNil(t, mockBoardRepo.UpdateCalledWith)
}

func TestBoardService_UpdateBoard_ErrOnFinalFindByID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	currentUserID := uint(10)
	newName := "Updated Board Name"
	expectedError := errors.New("DB error on final FindByID")

	originalBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Original Name",
		OwnerID: currentUserID,
	}

	var findByIDCallCount = 0
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		findByIDCallCount++
		if id == boardID && findByIDCallCount == 1 {
			return originalBoard, nil
		}
		if id == boardID && findByIDCallCount == 2 {
			return nil, expectedError
		}
		return nil, gorm.ErrRecordNotFound
	}

	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		return nil
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, updatedBoard)

	assert.Equal(t, 2, findByIDCallCount)
	assert.NotNil(t, mockBoardRepo.UpdateCalledWith)
}

func TestBoardService_DeleteBoard_Success(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	currentUserID := uint(10)

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Board to Delete",
		OwnerID: currentUserID,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardRepo.DeleteFunc = func(id uint) error {
		assert.Equal(t, boardID, id)
		return nil
	}

	err := boardService.DeleteBoard(boardID, currentUserID)

	assert.NoError(t, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, boardID, mockBoardRepo.DeleteCalledWith)
}

func TestBoardService_DeleteBoard_BoardNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	currentUserID := uint(10)

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return nil, gorm.ErrRecordNotFound
	}

	mockBoardRepo.DeleteFunc = func(id uint) error {
		t.Error("Delete should not be called if board not found")
		return nil
	}

	err := boardService.DeleteBoard(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardNotFound, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, uint(0), mockBoardRepo.DeleteCalledWith)
}

func TestBoardService_DeleteBoard_Forbidden(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(5)
	currentUserID := uint(10)

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Board to Delete",
		OwnerID: ownerID,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardRepo.DeleteFunc = func(id uint) error {
		t.Error("Delete should not be called if user is not owner")
		return nil
	}

	err := boardService.DeleteBoard(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, uint(0), mockBoardRepo.DeleteCalledWith)
}

func TestBoardService_DeleteBoard_ErrOnDelete(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	currentUserID := uint(10)
	expectedError := errors.New("DB error on delete")

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Board to Delete",
		OwnerID: currentUserID,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardRepo.DeleteFunc = func(id uint) error {
		assert.Equal(t, boardID, id)
		return expectedError
	}

	err := boardService.DeleteBoard(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, boardID, mockBoardRepo.DeleteCalledWith)
}

func TestBoardService_DeleteBoard_ErrOnFindByID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	currentUserID := uint(10)
	expectedError := errors.New("DB error on FindByID for delete")

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return nil, expectedError
	}

	mockBoardRepo.DeleteFunc = func(id uint) error {
		t.Error("Delete should not be called if FindByID fails")
		return nil
	}

	err := boardService.DeleteBoard(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, uint(0), mockBoardRepo.DeleteCalledWith)
}

func TestBoardService_AddMemberToBoard_SuccessByEmail(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "newmember@example.com"
	targetUserID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, Email: targetUserEmail}
	addedMemberWithUser := &models.BoardMember{
		BoardID: boardID, UserID: targetUserID, User: *userToFind,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return boardToReturn, nil
	}
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		assert.Equal(t, targetUserEmail, email)
		return userToFind, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, targetUserID, uID)
		return false, nil
	}
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		assert.Equal(t, boardID, member.BoardID)
		assert.Equal(t, targetUserID, member.UserID)
		return nil
	}
	mockBoardMemberRepo.FindByBoardIDAndUserIDFunc = func(bID uint, uID uint) (*models.BoardMember, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, targetUserID, uID)
		return addedMemberWithUser, nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.NoError(t, err)
	assert.NotNil(t, newMember)
	assert.Equal(t, boardID, newMember.BoardID)
	assert.Equal(t, targetUserID, newMember.UserID)
	assert.Equal(t, targetUserEmail, newMember.User.Email)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, targetUserEmail, mockUserRepo.FindByEmailCalledWith)
	assert.Equal(t, uint(0), mockUserRepo.FindByIDCalledWith)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, targetUserID, mockBoardMemberRepo.IsMemberCalledWithUserID)
	assert.NotNil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
	assert.Equal(t, boardID, mockBoardMemberRepo.FindByBoardIDAndUserIDCalledWithBoardID)
	assert.Equal(t, targetUserID, mockBoardMemberRepo.FindByBoardIDAndUserIDCalledWithUserID)
}

func TestBoardService_AddMemberToBoard_SuccessByUserID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, Email: "foundbyid@example.com"}
	addedMemberWithUser := &models.BoardMember{
		BoardID: boardID, UserID: targetUserID, User: *userToFind,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) {
		assert.Equal(t, targetUserID, id)
		return userToFind, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, nil
	}
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		return nil
	}
	mockBoardMemberRepo.FindByBoardIDAndUserIDFunc = func(bID uint, uID uint) (*models.BoardMember, error) {
		return addedMemberWithUser, nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, nil, &targetUserID, ownerID)

	assert.NoError(t, err)
	assert.NotNil(t, newMember)
	assert.Equal(t, targetUserID, newMember.UserID)
	assert.Equal(t, targetUserID, newMember.User.ID)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, "", mockUserRepo.FindByEmailCalledWith)
	assert.Equal(t, targetUserID, mockUserRepo.FindByIDCalledWith)
	assert.NotNil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
}

func TestBoardService_AddMemberToBoard_BoardNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return nil, gorm.ErrRecordNotFound
	}

	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		t.Error("UserRepo.FindByEmail should not be called if board not found")
		return nil, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("BoardMemberRepo.IsMember should not be called if board not found")
		return false, nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardNotFound, err)
	assert.Nil(t, newMember)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_AddMemberToBoard_Forbidden(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	actualOwnerID := uint(5)
	currentUserID := uint(10)
	targetUserEmail := "test@example.com"

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		t.Error("UserRepo.FindByEmailFunc should not be called in Forbidden case")
		return nil, nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, newMember)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, "", mockUserRepo.FindByEmailCalledWith)
}

func TestBoardService_AddMemberToBoard_UserToAddNotFoundByEmail(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "nonexistent@example.com"

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		assert.Equal(t, targetUserEmail, email)
		return nil, gorm.ErrRecordNotFound
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	assert.Nil(t, newMember)
	assert.Equal(t, targetUserEmail, mockUserRepo.FindByEmailCalledWith)
}

func TestBoardService_AddMemberToBoard_UserToAddNotFoundByID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserID := uint(999)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) {
		assert.Equal(t, targetUserID, id)
		return nil, gorm.ErrRecordNotFound
	}

	newMember, err := boardService.AddMemberToBoard(boardID, nil, &targetUserID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	assert.Nil(t, newMember)
	assert.Equal(t, targetUserID, mockUserRepo.FindByIDCalledWith)
}

func TestBoardService_AddMemberToBoard_UserAlreadyMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "member@example.com"
	targetUserID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, Email: targetUserEmail}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, targetUserID, uID)
		return true, nil
	}

	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		t.Error("AddMember should not be called if user already member")
		return nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrUserAlreadyMember, err)
	assert.Nil(t, newMember)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, targetUserID, mockBoardMemberRepo.IsMemberCalledWithUserID)
	assert.Nil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
}

func TestBoardService_AddMemberToBoard_OwnerCannotBeAddedAsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "owner@example.com"

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: ownerID}, Email: targetUserEmail}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }

	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called if target is owner")
		return false, nil
	}
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		t.Error("AddMember should not be called if target is owner")
		return nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrUserAlreadyMember, err)
	assert.Nil(t, newMember)
}

func TestBoardService_AddMemberToBoard_InvalidInput(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, nil, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, "Either email or userID must be provided", err.Error())
	assert.Nil(t, newMember)

	emptyEmail := ""
	zeroID := uint(0)
	newMember, err = boardService.AddMemberToBoard(boardID, &emptyEmail, &zeroID, ownerID)
	assert.Error(t, err)
	assert.Equal(t, "Either email or userID must be provided", err.Error())
	assert.Nil(t, newMember)
}

func TestBoardService_AddMemberToBoard_ErrOnUserFindByEmail(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"
	expectedError := errors.New("DB error on FindByEmail")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, expectedError
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, newMember)
	assert.Equal(t, targetUserEmail, mockUserRepo.FindByEmailCalledWith)
}

func TestBoardService_AddMemberToBoard_ErrOnUserFindByID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserID := uint(20)
	expectedError := errors.New("DB error on FindByID")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) {
		return nil, expectedError
	}

	newMember, err := boardService.AddMemberToBoard(boardID, nil, &targetUserID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, newMember)
	assert.Equal(t, targetUserID, mockUserRepo.FindByIDCalledWith)
}

func TestBoardService_AddMemberToBoard_ErrOnIsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"
	targetUserID := uint(20)
	expectedError := errors.New("DB error on IsMember")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, Email: targetUserEmail}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, expectedError
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, newMember)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, targetUserID, mockBoardMemberRepo.IsMemberCalledWithUserID)
}

func TestBoardService_AddMemberToBoard_ErrOnAddMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"
	targetUserID := uint(20)
	expectedError := errors.New("DB error on AddMember")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, Email: targetUserEmail}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil }
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		return expectedError
	}
	mockBoardMemberRepo.FindByBoardIDAndUserIDFunc = func(bID uint, uID uint) (*models.BoardMember, error) {
		t.Error("FindByBoardIDAndUserID should not be called if AddMember fails")
		return nil, nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, newMember)
	assert.NotNil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
}

func TestBoardService_AddMemberToBoard_ErrOnFinalFindByBoardIDAndUserID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"
	targetUserID := uint(20)
	expectedError := errors.New("DB error on FindByBoardIDAndUserID")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, Email: targetUserEmail}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil }
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error { return nil }
	mockBoardMemberRepo.FindByBoardIDAndUserIDFunc = func(bID uint, uID uint) (*models.BoardMember, error) {
		return nil, expectedError
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, newMember)
	assert.Equal(t, boardID, mockBoardMemberRepo.FindByBoardIDAndUserIDCalledWithBoardID)
	assert.Equal(t, targetUserID, mockBoardMemberRepo.FindByBoardIDAndUserIDCalledWithUserID)
}

func TestBoardService_RemoveMemberFromBoard_Success(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	memberToRemoveID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, memberToRemoveID, uID)
		return true, nil
	}
	mockBoardMemberRepo.RemoveMemberFunc = func(bID uint, uID uint) error {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, memberToRemoveID, uID)
		return nil
	}

	err := boardService.RemoveMemberFromBoard(boardID, memberToRemoveID, ownerID)

	assert.NoError(t, err)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, memberToRemoveID, mockBoardMemberRepo.IsMemberCalledWithUserID)
	assert.Equal(t, boardID, mockBoardMemberRepo.RemoveMemberCalledWithBoardID)
	assert.Equal(t, memberToRemoveID, mockBoardMemberRepo.RemoveMemberCalledWithUserID)
}

func TestBoardService_RemoveMemberFromBoard_BoardNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	mockHub := &MockHub{} // Initialize mock hub
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo, mockHub)

	boardID := uint(1)
	ownerID := uint(10)
	memberToRemoveID := uint(20)

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return nil, gorm.ErrRecordNotFound
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called if board not found")
		return false, nil
	}
	mockBoardMemberRepo.RemoveMemberFunc = func(bID uint, uID uint) error {
		t.Error("RemoveMember should not be called if board not found")
		return nil
	}

	err := boardService.RemoveMemberFromBoard(boardID, memberToRemoveID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardNotFound, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_RemoveMemberFromBoard_Forbidden(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	actualOwnerID := uint(5)
	currentUserID := uint(10)
	memberToRemoveID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called if user is not owner")
		return false, nil
	}

	err := boardService.RemoveMemberFromBoard(boardID, memberToRemoveID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_RemoveMemberFromBoard_CannotRemoveOwner(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called when trying to remove owner")
		return false, nil
	}

	err := boardService.RemoveMemberFromBoard(boardID, ownerID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrCannotRemoveOwner, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_RemoveMemberFromBoard_MemberNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	memberToRemoveID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, nil
	}
	mockBoardMemberRepo.RemoveMemberFunc = func(bID uint, uID uint) error {
		t.Error("RemoveMember should not be called if user is not a member")
		return nil
	}

	err := boardService.RemoveMemberFromBoard(boardID, memberToRemoveID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardMemberNotFound, err)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, memberToRemoveID, mockBoardMemberRepo.IsMemberCalledWithUserID)
}

func TestBoardService_RemoveMemberFromBoard_ErrOnIsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	memberToRemoveID := uint(20)
	expectedError := errors.New("DB error on IsMember")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, expectedError
	}
	mockBoardMemberRepo.RemoveMemberFunc = func(bID uint, uID uint) error {
		t.Error("RemoveMember should not be called if IsMember errors")
		return nil
	}

	err := boardService.RemoveMemberFromBoard(boardID, memberToRemoveID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestBoardService_RemoveMemberFromBoard_ErrOnRemoveMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	memberToRemoveID := uint(20)
	expectedError := errors.New("DB error on RemoveMember")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil }
	mockBoardMemberRepo.RemoveMemberFunc = func(bID uint, uID uint) error {
		return expectedError
	}

	err := boardService.RemoveMemberFromBoard(boardID, memberToRemoveID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, boardID, mockBoardMemberRepo.RemoveMemberCalledWithBoardID)
	assert.Equal(t, memberToRemoveID, mockBoardMemberRepo.RemoveMemberCalledWithUserID)
}

func TestBoardService_GetBoardMembers_SuccessAsOwner(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	expectedMembers := []models.BoardMember{
		{BoardID: boardID, UserID: ownerID, User: models.User{Model: gorm.Model{ID: ownerID}}},
		{BoardID: boardID, UserID: 20, User: models.User{Model: gorm.Model{ID: 20}}},
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called if current user is owner")
		return false, nil
	}
	mockBoardMemberRepo.FindMembersByBoardIDFunc = func(bID uint) ([]models.BoardMember, error) {
		assert.Equal(t, boardID, bID)
		return expectedMembers, nil
	}

	members, err := boardService.GetBoardMembers(boardID, ownerID)

	assert.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Equal(t, expectedMembers, members)
	assert.Equal(t, boardID, mockBoardMemberRepo.FindMembersByBoardIDCalledWith)
}

func TestBoardService_GetBoardMembers_SuccessAsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	actualOwnerID := uint(5)
	currentUserID := uint(10)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}
	expectedMembers := []models.BoardMember{
		{BoardID: boardID, UserID: actualOwnerID, User: models.User{Model: gorm.Model{ID: actualOwnerID}}},
		{BoardID: boardID, UserID: currentUserID, User: models.User{Model: gorm.Model{ID: currentUserID}}},
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return true, nil
	}
	mockBoardMemberRepo.FindMembersByBoardIDFunc = func(bID uint) ([]models.BoardMember, error) {
		assert.Equal(t, boardID, bID)
		return expectedMembers, nil
	}

	members, err := boardService.GetBoardMembers(boardID, currentUserID)

	assert.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Equal(t, expectedMembers, members)
	assert.Equal(t, boardID, mockBoardMemberRepo.FindMembersByBoardIDCalledWith)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, currentUserID, mockBoardMemberRepo.IsMemberCalledWithUserID)
}

func TestBoardService_GetBoardMembers_BoardNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	currentUserID := uint(10)

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return nil, gorm.ErrRecordNotFound
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called if board not found")
		return false, nil
	}
	mockBoardMemberRepo.FindMembersByBoardIDFunc = func(bID uint) ([]models.BoardMember, error) {
		t.Error("FindMembersByBoardID should not be called if board not found")
		return nil, nil
	}

	members, err := boardService.GetBoardMembers(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardNotFound, err)
	assert.Nil(t, members)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_GetBoardMembers_Forbidden(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	actualOwnerID := uint(5)
	currentUserID := uint(10)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return false, nil
	}
	mockBoardMemberRepo.FindMembersByBoardIDFunc = func(bID uint) ([]models.BoardMember, error) {
		t.Error("FindMembersByBoardID should not be called if user is forbidden")
		return nil, nil
	}

	members, err := boardService.GetBoardMembers(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, members)
}

func TestBoardService_GetBoardMembers_ErrOnIsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	actualOwnerID := uint(5)
	currentUserID := uint(10)
	expectedError := errors.New("DB error on IsMember")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, expectedError
	}
	mockBoardMemberRepo.FindMembersByBoardIDFunc = func(bID uint) ([]models.BoardMember, error) {
		t.Error("FindMembersByBoardID should not be called if IsMember errors")
		return nil, nil
	}

	members, err := boardService.GetBoardMembers(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, members)
}

func TestBoardService_GetBoardMembers_ErrOnFindMembersByBoardID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	expectedError := errors.New("DB error on FindMembersByBoardID")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.FindMembersByBoardIDFunc = func(bID uint) ([]models.BoardMember, error) {
		return nil, expectedError
	}

	members, err := boardService.GetBoardMembers(boardID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, members)
	assert.Equal(t, boardID, mockBoardMemberRepo.FindMembersByBoardIDCalledWith)
}

func TestBoardService_GetBoardByID_BoardNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	userID := uint(10)

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return nil, gorm.ErrRecordNotFound
	}

	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called if board not found")
		return false, nil
	}

	board, err := boardService.GetBoardByID(boardID, userID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardNotFound, err)
	assert.Nil(t, board)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, uint(0), mockBoardMemberRepo.IsMemberCalledWithBoardID)
}

func TestBoardService_GetBoardByID_Forbidden(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	boardOwnerID := uint(10)
	currentUserID := uint(20)

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID},
		Name:    "Other's Board",
		OwnerID: boardOwnerID,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return false, nil
	}

	board, err := boardService.GetBoardByID(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, board)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, currentUserID, mockBoardMemberRepo.IsMemberCalledWithUserID)
}
