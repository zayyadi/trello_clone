package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zayyadi/trello/models"
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
	CreateFunc      func(user *models.User) error // Likely not used by BoardService directly
	FindByEmailFunc func(email string) (*models.User, error)
	FindByIDFunc    func(id uint) (*models.User, error)

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

var _ repositories.UserRepositoryInterface = (*MockUserRepositoryForBoardService)(nil)

// Placeholder for actual tests
func TestBoardService_CreateBoard_Success(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{} // Not directly used by CreateBoard logic itself but needed for service init
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	ownerID := uint(1)
	boardName := "Test Board"
	boardDescription := "Test Description"

	// Configure mock responses
	mockBoardRepo.CreateFunc = func(board *models.Board) error {
		assert.Equal(t, boardName, board.Name)
		assert.Equal(t, boardDescription, board.Description)
		assert.Equal(t, ownerID, board.OwnerID)
		board.ID = 100 // Simulate DB assigning an ID
		return nil
	}

	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		assert.Equal(t, uint(100), member.BoardID) // Check if board ID from Create is used
		assert.Equal(t, ownerID, member.UserID)
		// BoardMember does not have its own ID field.
		return nil
	}

	// FindByID is called at the end of CreateBoard to return the created board
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, uint(100), id)
		return &models.Board{
			Model:       gorm.Model{ID: id}, // Correct: Set ID in gorm.Model
			Name:        boardName,
			Description: boardDescription,
			OwnerID:     ownerID,
			Owner:       models.User{Model: gorm.Model{ID: ownerID}, ID: ownerID, Username: "owner"}, // Correct: models.User
			Members: []models.BoardMember{ // Simulate preload
				{BoardID: id, UserID: ownerID, User: models.User{Model: gorm.Model{ID: ownerID}, ID: ownerID, Username: "owner"}}, // Correct: models.User, no BoardMember.ID
			},
		}, nil
	}

	createdBoard, err := boardService.CreateBoard(boardName, boardDescription, ownerID)

	assert.NoError(t, err)
	assert.NotNil(t, createdBoard)
	assert.Equal(t, uint(100), createdBoard.ID)
	assert.Equal(t, boardName, createdBoard.Name)
	assert.Equal(t, ownerID, createdBoard.OwnerID)
	assert.Equal(t, ownerID, createdBoard.Owner.ID) // Assert Owner.ID
	assert.Len(t, createdBoard.Members, 1)
	if len(createdBoard.Members) > 0 {
		assert.Equal(t, ownerID, createdBoard.Members[0].UserID)
		assert.Equal(t, ownerID, createdBoard.Members[0].User.ID) // Assert User.ID within BoardMember
	}

	// Check that the mock functions were called
	assert.NotNil(t, mockBoardRepo.CreateCalledWith)
	assert.NotNil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
	assert.Equal(t, uint(100), mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_CreateBoard_ErrOnBoardCreate(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	ownerID := uint(1)
	boardName := "Test Board"
	boardDescription := "Test Description"
	expectedError := errors.New("DB error on board create")

	// Configure mock responses
	mockBoardRepo.CreateFunc = func(board *models.Board) error {
		return expectedError // Simulate error
	}

	// AddMemberFunc and FindByIDFunc should not be called
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

	// Check that Create was called, but others were not
	assert.NotNil(t, mockBoardRepo.CreateCalledWith)             // Create was attempted
	assert.Nil(t, mockBoardMemberRepo.AddMemberCalledWithMember) // AddMember was not called
	assert.Equal(t, uint(0), mockBoardRepo.FindByIDCalledWith)   // FindByID was not called (or called with 0 if default value matters)
}

func TestBoardService_CreateBoard_ErrOnAddOwnerAsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	ownerID := uint(1)
	boardName := "Test Board"
	boardDescription := "Test Description"
	expectedError := errors.New("DB error on add member")

	// Configure mock responses
	mockBoardRepo.CreateFunc = func(board *models.Board) error {
		board.ID = 100 // Simulate DB assigning an ID
		return nil
	}

	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		return expectedError // Simulate error
	}

	// FindByIDFunc should not be called if adding member fails
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		t.Error("FindByID should not be called if AddMember fails")
		return nil, nil
	}

	createdBoard, err := boardService.CreateBoard(boardName, boardDescription, ownerID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, createdBoard)

	// Check that Create and AddMember were called, but FindByID was not
	assert.NotNil(t, mockBoardRepo.CreateCalledWith)
	assert.NotNil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
	assert.Equal(t, uint(0), mockBoardRepo.FindByIDCalledWith) // FindByID was not called
}

func TestBoardService_GetBoardByID_SuccessAsOwner(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{} // Not directly used if user is owner

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10) // This is the current user ID

	expectedBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:    "Owned Board",
		OwnerID: ownerID,                                                  // Current user is the owner
		Owner:   models.User{Model: gorm.Model{ID: ownerID}, ID: ownerID}, // Correct: models.User
	}

	// Configure mock responses
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return expectedBoard, nil
	}

	// IsMember should not be called if user is owner
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
	// Verify IsMember was not called
	assert.Equal(t, uint(0), mockBoardMemberRepo.IsMemberCalledWithBoardID) // Check default/zero value
}

func TestBoardService_GetBoardByID_SuccessAsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	boardOwnerID := uint(10)
	currentUserID := uint(20) // This user is a member, not the owner

	expectedBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:    "Member Board",
		OwnerID: boardOwnerID,                                                       // Current user is not the owner
		Owner:   models.User{Model: gorm.Model{ID: boardOwnerID}, ID: boardOwnerID}, // Correct: models.User
	}

	// Configure mock responses
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return expectedBoard, nil
	}

	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return true, nil // User is a member
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

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	boardOwnerID := uint(10)
	currentUserID := uint(20) // This user is neither owner nor member
	expectedError := errors.New("DB error on IsMember")

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:    "Other's Board",
		OwnerID: boardOwnerID, // Not currentUserID
	}

	// Configure mock responses
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return false, expectedError // Error checking membership
	}

	board, err := boardService.GetBoardByID(boardID, currentUserID)

	assert.Error(t, err)
	// The service code returns ErrForbidden if IsMember returns an error.
	// This might be something to refine in the service (e.g., return the actual DB error or a generic server error)
	// For now, testing current behavior.
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

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	userID := uint(1)
	expectedBoards := []models.Board{
		{Model: gorm.Model{ID: 1}, Name: "Board 1", OwnerID: userID},
		{Model: gorm.Model{ID: 2}, Name: "Board 2", OwnerID: uint(2), Members: []models.BoardMember{{UserID: userID}}},
	}

	mockBoardRepo.FindByOwnerOrMemberFunc = func(uID uint) ([]models.Board, error) {
		assert.Equal(t, userID, uID)
		// Ensure the mock returns boards with their gorm.Model ID set, as the real repo would.
		return expectedBoards, nil
	}

	boards, err := boardService.GetBoardsForUser(userID)

	assert.NoError(t, err)
	assert.NotNil(t, boards)
	assert.Len(t, boards, 2)
	if len(boards) == 2 { // Avoid panic on nil or short slice
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

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	userID := uint(1)
	expectedBoards := []models.Board{} // Empty slice

	mockBoardRepo.FindByOwnerOrMemberFunc = func(uID uint) ([]models.Board, error) {
		assert.Equal(t, userID, uID)
		return expectedBoards, nil // Return empty slice, no error
	}

	boards, err := boardService.GetBoardsForUser(userID)

	assert.NoError(t, err)
	assert.NotNil(t, boards) // Should be an empty slice, not nil
	assert.Len(t, boards, 0)
	assert.Equal(t, userID, mockBoardRepo.FindByOwnerOrMemberCalledWith)
}

func TestBoardService_GetBoardsForUser_Error(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	userID := uint(1)
	expectedError := errors.New("DB error FindByOwnerOrMember")

	mockBoardRepo.FindByOwnerOrMemberFunc = func(uID uint) ([]models.Board, error) {
		assert.Equal(t, userID, uID)
		return nil, expectedError // Return error
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

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	currentUserID := uint(10) // Owner
	newName := "Updated Board Name"
	newDescription := "Updated Description"

	originalBoard := &models.Board{
		Model:       gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:        "Original Name",
		Description: "Original Description",
		OwnerID:     currentUserID,
	}

	// FindByID (first call)
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		if id == boardID && mockBoardRepo.UpdateCalledWith == nil {
			// First call, before Update is called
			return originalBoard, nil
		} else if id == boardID && mockBoardRepo.UpdateCalledWith != nil {
			// Second call, after Update is called
			updatedBoardFromMock := *mockBoardRepo.UpdateCalledWith
			updatedBoardFromMock.Model.ID = id // Ensure returned mock has ID
			return &updatedBoardFromMock, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	// Update
	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		assert.Equal(t, boardID, board.ID)
		assert.Equal(t, newName, board.Name)
		assert.Equal(t, newDescription, board.Description)
		assert.Equal(t, currentUserID, board.OwnerID) // Ensure OwnerID is not changed
		return nil
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, &newDescription, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, updatedBoard)
	assert.Equal(t, boardID, updatedBoard.ID)
	assert.Equal(t, newName, updatedBoard.Name)
	assert.Equal(t, newDescription, updatedBoard.Description)

	assert.NotNil(t, mockBoardRepo.UpdateCalledWith)
	// FindByIDCalledWith will be boardID due to the sequence of calls, test logic handles this.
	// Need a better way to track calls if more FindByID are added. For now this is okay.
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_UpdateBoard_BoardNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	currentUserID := uint(10)
	newName := "Updated Name"

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return nil, gorm.ErrRecordNotFound // Simulate board not found
	}

	// UpdateFunc should not be called
	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		t.Error("Update should not be called if board not found")
		return nil
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardNotFound, err)
	assert.Nil(t, updatedBoard)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith) // First FindByID was called
	assert.Nil(t, mockBoardRepo.UpdateCalledWith)              // Update was not called
}

func TestBoardService_UpdateBoard_Forbidden(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(5)
	currentUserID := uint(10) // Not the owner
	newName := "Updated Name"

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:    "Original Name",
		OwnerID: ownerID, // Different from currentUserID
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	// UpdateFunc should not be called
	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		t.Error("Update should not be called if user is not owner")
		return nil
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, updatedBoard)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith) // FindByID was called
	assert.Nil(t, mockBoardRepo.UpdateCalledWith)              // Update was not called
}

func TestBoardService_UpdateBoard_ErrOnUpdate(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	currentUserID := uint(10) // Owner
	newName := "Updated Name"
	expectedError := errors.New("DB error on update")

	originalBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:    "Original Name",
		OwnerID: currentUserID,
	}

	var findByIDCallCount = 0
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		findByIDCallCount++
		if id == boardID && findByIDCallCount == 1 { // First call
			return originalBoard, nil
		}
		// Second call to FindByID (after Update) should not happen if Update fails
		if findByIDCallCount == 2 {
			t.Error("Second FindByID should not be called if Update fails")
		}
		return nil, gorm.ErrRecordNotFound
	}

	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		return expectedError // Simulate error
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, updatedBoard)

	assert.Equal(t, 1, findByIDCallCount)            // Initial FindByID was called
	assert.NotNil(t, mockBoardRepo.UpdateCalledWith) // Update was attempted
}

func TestBoardService_UpdateBoard_ErrOnFinalFindByID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	currentUserID := uint(10) // Owner
	newName := "Updated Board Name"
	expectedError := errors.New("DB error on final FindByID")

	originalBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:    "Original Name",
		OwnerID: currentUserID,
	}

	var findByIDCallCount = 0
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		findByIDCallCount++
		if id == boardID && findByIDCallCount == 1 { // First call
			return originalBoard, nil
		}
		if id == boardID && findByIDCallCount == 2 { // Second call (after successful update)
			return nil, expectedError // Error on this call
		}
		return nil, gorm.ErrRecordNotFound // Should not happen in this test flow
	}

	mockBoardRepo.UpdateFunc = func(board *models.Board) error {
		// Successful update
		return nil
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, &newName, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, updatedBoard)

	assert.Equal(t, 2, findByIDCallCount) // Both FindByID calls were made
	assert.NotNil(t, mockBoardRepo.UpdateCalledWith)
}

func TestBoardService_DeleteBoard_Success(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	currentUserID := uint(10) // Owner

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:    "Board to Delete",
		OwnerID: currentUserID,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardRepo.DeleteFunc = func(id uint) error {
		assert.Equal(t, boardID, id)
		return nil // Successful deletion
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

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	currentUserID := uint(10)

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return nil, gorm.ErrRecordNotFound // Simulate board not found
	}

	// DeleteFunc should not be called
	mockBoardRepo.DeleteFunc = func(id uint) error {
		t.Error("Delete should not be called if board not found")
		return nil
	}

	err := boardService.DeleteBoard(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardNotFound, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, uint(0), mockBoardRepo.DeleteCalledWith) // Delete was not called
}

func TestBoardService_DeleteBoard_Forbidden(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(5)
	currentUserID := uint(10) // Not the owner

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:    "Board to Delete",
		OwnerID: ownerID, // Different from currentUserID
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	// DeleteFunc should not be called
	mockBoardRepo.DeleteFunc = func(id uint) error {
		t.Error("Delete should not be called if user is not owner")
		return nil
	}

	err := boardService.DeleteBoard(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, uint(0), mockBoardRepo.DeleteCalledWith) // Delete was not called
}

func TestBoardService_DeleteBoard_ErrOnDelete(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	currentUserID := uint(10) // Owner
	expectedError := errors.New("DB error on delete")

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Correct: Set ID in gorm.Model
		Name:    "Board to Delete",
		OwnerID: currentUserID,
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardRepo.DeleteFunc = func(id uint) error {
		assert.Equal(t, boardID, id)
		return expectedError // Simulate error
	}

	err := boardService.DeleteBoard(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, boardID, mockBoardRepo.DeleteCalledWith) // Delete was attempted
}

func TestBoardService_DeleteBoard_ErrOnFindByID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	currentUserID := uint(10)
	expectedError := errors.New("DB error on FindByID for delete")

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return nil, expectedError // Error on find
	}

	mockBoardRepo.DeleteFunc = func(id uint) error {
		t.Error("Delete should not be called if FindByID fails")
		return nil
	}

	err := boardService.DeleteBoard(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err) // Should be the error from FindByID
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, uint(0), mockBoardRepo.DeleteCalledWith)
}

func TestBoardService_AddMemberToBoard_SuccessByEmail(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10) // Current user, owner of the board
	targetUserEmail := "newmember@example.com"
	targetUserID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID, Email: targetUserEmail}
	// BoardMember has no ID field. User field is not a pointer.
	addedMemberWithUser := &models.BoardMember{
		BoardID: boardID, UserID: targetUserID, User: *userToFind,
	}

	// Mock setup
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
		return false, nil // Not a member yet
	}
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		assert.Equal(t, boardID, member.BoardID)
		assert.Equal(t, targetUserID, member.UserID)
		// No ID on BoardMember
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
	assert.Equal(t, targetUserEmail, newMember.User.Email) // User is not a pointer

	// Verify calls
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, targetUserEmail, mockUserRepo.FindByEmailCalledWith)
	assert.Equal(t, uint(0), mockUserRepo.FindByIDCalledWith) // FindByID not called
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

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10) // Current user, owner of the board
	targetUserID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID, Email: "foundbyid@example.com"}
	addedMemberWithUser := &models.BoardMember{
		BoardID: boardID, UserID: targetUserID, User: *userToFind,
	}

	// Mock setup
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) {
		assert.Equal(t, targetUserID, id)
		return userToFind, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, nil // Not a member yet
	}
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		// No ID on BoardMember
		return nil
	}
	mockBoardMemberRepo.FindByBoardIDAndUserIDFunc = func(bID uint, uID uint) (*models.BoardMember, error) {
		return addedMemberWithUser, nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, nil, &targetUserID, ownerID)

	assert.NoError(t, err)
	assert.NotNil(t, newMember)
	assert.Equal(t, targetUserID, newMember.UserID)
	assert.Equal(t, targetUserID, newMember.User.ID) // User is not a pointer

	// Verify calls
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, "", mockUserRepo.FindByEmailCalledWith) // FindByEmail not called
	assert.Equal(t, targetUserID, mockUserRepo.FindByIDCalledWith)
	assert.NotNil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
}

func TestBoardService_AddMemberToBoard_BoardNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"

	// No need to fully populate boardToReturn as it's not found
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return nil, gorm.ErrRecordNotFound
	}

	// None of the other repo functions should be called
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		t.Error("UserRepo.FindByEmail should not be called if board not found")
		return nil, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("BoardMemberRepo.IsMember should not be called if board not found")
		return false, nil
	}
	// ... and so on for other functions that should not be called.

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

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	actualOwnerID := uint(5)
	currentUserID := uint(10) // Not the owner
	targetUserEmail := "test@example.com"

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { // Should not be called
		t.Error("UserRepo.FindByEmailFunc should not be called in Forbidden case")
		return nil, nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, newMember)
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, "", mockUserRepo.FindByEmailCalledWith) // User repo not called
}

func TestBoardService_AddMemberToBoard_UserToAddNotFoundByEmail(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "nonexistent@example.com"

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		assert.Equal(t, targetUserEmail, email)
		return nil, gorm.ErrRecordNotFound // User not found
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err) // Ensure this error is defined
	assert.Nil(t, newMember)
	assert.Equal(t, targetUserEmail, mockUserRepo.FindByEmailCalledWith)
}

func TestBoardService_AddMemberToBoard_UserToAddNotFoundByID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserID := uint(999) // Non-existent user ID

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) {
		assert.Equal(t, targetUserID, id)
		return nil, gorm.ErrRecordNotFound // User not found
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

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "member@example.com"
	targetUserID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID, Email: targetUserEmail}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, targetUserID, uID)
		return true, nil // User IS already a member
	}

	// AddMember and FindByBoardIDAndUserID should not be called
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		t.Error("AddMember should not be called if user already member")
		return nil
	}

	newMember, err := boardService.AddMemberToBoard(boardID, &targetUserEmail, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrUserAlreadyMember, err) // Ensure this error is defined
	assert.Nil(t, newMember)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, targetUserID, mockBoardMemberRepo.IsMemberCalledWithUserID)
	assert.Nil(t, mockBoardMemberRepo.AddMemberCalledWithMember)
}

func TestBoardService_AddMemberToBoard_OwnerCannotBeAddedAsMember(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10) // This is also the target user
	targetUserEmail := "owner@example.com"

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: ownerID}, ID: ownerID, Email: targetUserEmail} // Target user is the owner

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }

	// IsMember and AddMember should not be called if target is owner
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
	assert.Equal(t, ErrUserAlreadyMember, err) // Or a more specific error like ErrCannotAddOwnerAsMember
	assert.Nil(t, newMember)
}

func TestBoardService_AddMemberToBoard_InvalidInput(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}

	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)

	// Board FindByID will be called
	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}

	// Call with nil email and nil userID
	newMember, err := boardService.AddMemberToBoard(boardID, nil, nil, ownerID)

	assert.Error(t, err)
	assert.Equal(t, "Either email or userID must be provided", err.Error())
	assert.Nil(t, newMember)

	// Call with empty email and zero userID
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
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"
	expectedError := errors.New("DB error on FindByEmail")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, expectedError // Error from FindByEmail
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
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserID := uint(20)
	expectedError := errors.New("DB error on FindByID")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByIDFunc = func(id uint) (*models.User, error) {
		return nil, expectedError // Error from FindByID
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
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"
	targetUserID := uint(20)
	expectedError := errors.New("DB error on IsMember")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID, Email: targetUserEmail}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, expectedError // Error from IsMember
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
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"
	targetUserID := uint(20)
	expectedError := errors.New("DB error on AddMember")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID, Email: targetUserEmail}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil } // Not a member
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error {
		return expectedError // Error from AddMember
	}
	// FindByBoardIDAndUserID should not be called if AddMember fails
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
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	targetUserEmail := "test@example.com"
	targetUserID := uint(20)
	expectedError := errors.New("DB error on FindByBoardIDAndUserID")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	userToFind := &models.User{Model: gorm.Model{ID: targetUserID}, ID: targetUserID, Email: targetUserEmail}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockUserRepo.FindByEmailFunc = func(email string) (*models.User, error) { return userToFind, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return false, nil } // Not a member
	mockBoardMemberRepo.AddMemberFunc = func(member *models.BoardMember) error { return nil }       // AddMember succeeds
	mockBoardMemberRepo.FindByBoardIDAndUserIDFunc = func(bID uint, uID uint) (*models.BoardMember, error) {
		return nil, expectedError // Error from final FindByBoardIDAndUserID
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
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)          // Current user, owner of the board
	memberToRemoveID := uint(20) // User to be removed

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	// Mock setup
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, memberToRemoveID, uID)
		return true, nil // User is a member
	}
	mockBoardMemberRepo.RemoveMemberFunc = func(bID uint, uID uint) error {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, memberToRemoveID, uID)
		return nil // Successful removal
	}

	err := boardService.RemoveMemberFromBoard(boardID, memberToRemoveID, ownerID)

	assert.NoError(t, err)

	// Verify calls
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
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	memberToRemoveID := uint(20)

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return nil, gorm.ErrRecordNotFound // Board not found
	}
	// IsMember and RemoveMember should not be called
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
	currentUserID := uint(10) // Not the owner
	memberToRemoveID := uint(20)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	// IsMember and RemoveMember should not be called
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
	ownerID := uint(10) // Current user and also the member to remove (owner)

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		return boardToReturn, nil
	}
	// IsMember and RemoveMember should not be called if trying to remove owner
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called when trying to remove owner")
		return false, nil
	}

	err := boardService.RemoveMemberFromBoard(boardID, ownerID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrCannotRemoveOwner, err) // Ensure this error is defined
	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
}

func TestBoardService_RemoveMemberFromBoard_MemberNotFound(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10)
	memberToRemoveID := uint(20) // Not actually a member

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, nil // User is NOT a member
	}
	// RemoveMember should not be called
	mockBoardMemberRepo.RemoveMemberFunc = func(bID uint, uID uint) error {
		t.Error("RemoveMember should not be called if user is not a member")
		return nil
	}

	err := boardService.RemoveMemberFromBoard(boardID, memberToRemoveID, ownerID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardMemberNotFound, err) // Ensure this error is defined
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
		return false, expectedError // Error from IsMember
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
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) { return true, nil } // User is a member
	mockBoardMemberRepo.RemoveMemberFunc = func(bID uint, uID uint) error {
		return expectedError // Error from RemoveMember
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
	ownerID := uint(10) // Current user is owner

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}
	expectedMembers := []models.BoardMember{
		{BoardID: boardID, UserID: ownerID, User: models.User{Model: gorm.Model{ID: ownerID}, ID: ownerID}},
		{BoardID: boardID, UserID: 20, User: models.User{Model: gorm.Model{ID: 20}, ID: 20}},
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	// IsMember should not be called for owner
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
	currentUserID := uint(10) // Current user is a member, not owner

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}
	expectedMembers := []models.BoardMember{
		{BoardID: boardID, UserID: actualOwnerID, User: models.User{Model: gorm.Model{ID: actualOwnerID}, ID: actualOwnerID}},
		{BoardID: boardID, UserID: currentUserID, User: models.User{Model: gorm.Model{ID: currentUserID}, ID: currentUserID}},
	}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return true, nil // Current user is a member
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
		return nil, gorm.ErrRecordNotFound // Board not found
	}
	// IsMember and FindMembersByBoardID should not be called
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
	currentUserID := uint(10) // Not owner, not member

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return false, nil // Current user is NOT a member
	}
	// FindMembersByBoardID should not be called
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
	currentUserID := uint(10) // Not owner
	expectedError := errors.New("DB error on IsMember")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: actualOwnerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		return false, expectedError // Error from IsMember
	}
	mockBoardMemberRepo.FindMembersByBoardIDFunc = func(bID uint) ([]models.BoardMember, error) {
		t.Error("FindMembersByBoardID should not be called if IsMember errors")
		return nil, nil
	}

	members, err := boardService.GetBoardMembers(boardID, currentUserID)

	assert.Error(t, err)
	// Service currently returns ErrForbidden if IsMember fails.
	// Test this behavior, but it could be refined to return the actual error.
	assert.Equal(t, ErrForbidden, err)
	assert.Nil(t, members)
}

func TestBoardService_GetBoardMembers_ErrOnFindMembersByBoardID(t *testing.T) {
	mockBoardRepo := &MockBoardRepository{}
	mockUserRepo := &MockUserRepositoryForBoardService{}
	mockBoardMemberRepo := &MockBoardMemberRepository{}
	boardService := NewBoardService(mockBoardRepo, mockUserRepo, mockBoardMemberRepo)

	boardID := uint(1)
	ownerID := uint(10) // Current user is owner, so IsMember won't be called
	expectedError := errors.New("DB error on FindMembersByBoardID")

	boardToReturn := &models.Board{Model: gorm.Model{ID: boardID}, OwnerID: ownerID}

	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) { return boardToReturn, nil }
	mockBoardMemberRepo.FindMembersByBoardIDFunc = func(bID uint) ([]models.BoardMember, error) {
		return nil, expectedError // Error from FindMembersByBoardID
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

	// Configure mock responses
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return nil, gorm.ErrRecordNotFound // Simulate board not found
	}

	// IsMember should not be called if board is not found
	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		t.Error("IsMember should not be called if board not found")
		return false, nil
	}

	board, err := boardService.GetBoardByID(boardID, userID)

	assert.Error(t, err)
	assert.Equal(t, ErrBoardNotFound, err) // Ensure this error is defined and used
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
	currentUserID := uint(20) // This user is neither owner nor member

	foundBoard := &models.Board{
		Model:   gorm.Model{ID: boardID}, // Corrected
		Name:    "Other's Board",
		OwnerID: boardOwnerID, // Not currentUserID
	}

	// Configure mock responses
	mockBoardRepo.FindByIDFunc = func(id uint) (*models.Board, error) {
		assert.Equal(t, boardID, id)
		return foundBoard, nil
	}

	mockBoardMemberRepo.IsMemberFunc = func(bID uint, uID uint) (bool, error) {
		assert.Equal(t, boardID, bID)
		assert.Equal(t, currentUserID, uID)
		return false, nil // User is NOT a member
	}

	board, err := boardService.GetBoardByID(boardID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, ErrForbidden, err) // Ensure this error is defined and used
	assert.Nil(t, board)

	assert.Equal(t, boardID, mockBoardRepo.FindByIDCalledWith)
	assert.Equal(t, boardID, mockBoardMemberRepo.IsMemberCalledWithBoardID)
	assert.Equal(t, currentUserID, mockBoardMemberRepo.IsMemberCalledWithUserID)
}
