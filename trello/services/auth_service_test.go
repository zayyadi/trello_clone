package services

import (
	"errors"
	"testing"

	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"
	"github.com/zayyadi/trello/utils"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepositoryInterface
type MockUserRepository struct {
	CreateFunc      func(user *models.User) error
	FindByEmailFunc func(email string) (*models.User, error)
	FindByIDFunc    func(id uint) (*models.User, error)

	// Store calls if needed for assertions
	CreateCalledWith      *models.User
	FindByEmailCalledWith string
	FindByIDCalledWith    uint
}

func (m *MockUserRepository) Create(user *models.User) error {
	m.CreateCalledWith = user // Store the passed user
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	// Default behavior
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	m.FindByEmailCalledWith = email // Store the passed email
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	// Default behavior: return nil, and a "not found" error to mimic gorm.ErrRecordNotFound
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	m.FindByIDCalledWith = id // Store the passed ID
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	// Default behavior: return nil, and a "not found" error to mimic gorm.ErrRecordNotFound
	return nil, gorm.ErrRecordNotFound
}

// Ensure MockUserRepository implements UserRepositoryInterface
var _ repositories.UserRepositoryInterface = (*MockUserRepository)(nil)

func TestAuthService_Register_EmailExists(t *testing.T) {
	mockRepo := &MockUserRepository{}
	authService := NewAuthService(mockRepo, "test-secret")

	// Configure mock responses
	existingUser := &models.User{Model: gorm.Model{ID: 1}, Username: "existinguser", Email: "test@example.com"}
	mockRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		if email == "test@example.com" {
			return existingUser, nil // Simulate user found
		}
		return nil, gorm.ErrRecordNotFound
	}

	username := "newuser"
	email := "test@example.com" // Same email as existing user
	password := "password123"

	user, token, err := authService.Register(username, email, password)

	assert.Error(t, err)
	assert.Equal(t, ErrEmailExists, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	// Assert that FindByEmail was called
	assert.Equal(t, email, mockRepo.FindByEmailCalledWith)
	// Assert that Create was not called
	assert.Nil(t, mockRepo.CreateCalledWith)
}

func TestAuthService_Register_CreateUserError(t *testing.T) {
	mockRepo := &MockUserRepository{}
	authService := NewAuthService(mockRepo, "test-secret")

	// Configure mock responses
	mockRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, gorm.ErrRecordNotFound // Simulate user not found
	}
	expectedError := errors.New("DB create error")
	mockRepo.CreateFunc = func(user *models.User) error {
		return expectedError // Simulate error during user creation
	}

	username := "testuser"
	email := "test@example.com"
	password := "password123"

	user, token, err := authService.Register(username, email, password)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	// Assert that FindByEmail and Create were called
	assert.Equal(t, email, mockRepo.FindByEmailCalledWith)
	assert.NotNil(t, mockRepo.CreateCalledWith)
	if mockRepo.CreateCalledWith != nil {
		assert.Equal(t, username, mockRepo.CreateCalledWith.Username)
		assert.Equal(t, email, mockRepo.CreateCalledWith.Email)
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	authService := NewAuthService(mockRepo, "test-secret")

	// Configure mock responses
	mockRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, gorm.ErrRecordNotFound // Simulate user not found
	}
	mockRepo.CreateFunc = func(user *models.User) error {
		// Simulate successful creation by assigning an ID to the user
		user.ID = 1
		return nil
	}

	username := "testuser"
	email := "test@example.com"
	password := "password123"

	user, token, err := authService.Register(username, email, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.NotEqual(t, password, user.Password) // Password should be hashed
	assert.Equal(t, uint(1), user.ID)           // Check if ID was assigned

	// Assert that mock methods were called
	assert.Equal(t, email, mockRepo.FindByEmailCalledWith)
	assert.NotNil(t, mockRepo.CreateCalledWith)
	if mockRepo.CreateCalledWith != nil {
		assert.Equal(t, username, mockRepo.CreateCalledWith.Username)
		assert.Equal(t, email, mockRepo.CreateCalledWith.Email)
	}
}

// TestAuthService_Register_PasswordHashingError (TODO)
// This test is tricky because utils.HashPassword itself would need to be fallible
// in a way that's easy to mock or trigger. If HashPassword can't fail for
// AuthService's inputs (e.g. empty password string is caught by validation before),
// this test might be more appropriate at the utils level.
// For now, we assume HashPassword is reliable or tested separately.
// If it could return a specific error we want to check:
/*
func TestAuthService_Register_PasswordHashingError(t *testing.T) {
	mockRepo := &MockUserRepository{}
	authService := NewAuthService(mockRepo, "test-secret")

	mockRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, gorm.ErrRecordNotFound
	}

	// To make this test work, we'd need to be able to make utils.HashPassword fail.
	// This might involve changing utils.HashPassword or using a mock for it,
	// which is beyond the scope of only mocking UserRepository.
	// For demonstration, let's assume a hypothetical error:
	// utils.SetNextHashPasswordError(errors.New("bcrypt error"))

	username := "testuser"
	email := "testerror@example.com"
	password := "password123"

	user, token, err := authService.Register(username, email, password)

	// assert.Error(t, err)
	// assert.Contains(t, err.Error(), "bcrypt error") // or specific error type
	assert.Nil(t, user)
	assert.Empty(t, token)
}
*/

// TestAuthService_Register_TokenGenerationError (TODO)
// Similar to password hashing, this depends on utils.GenerateJWT failing.
// utils.GenerateJWT could fail if the secret key is invalid (e.g. too short for HS256)
// or if user.ID is 0 (if we add such a check, though typically it's positive).
// We can test this by providing a problematic JWT secret key.
func TestAuthService_Register_TokenGenerationError(t *testing.T) {
	mockRepo := &MockUserRepository{}
	// Use the special secret key to force GenerateJWT to fail
	authService := NewAuthService(mockRepo, "FORCE_JWT_ERROR_FOR_TEST")

	mockRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, gorm.ErrRecordNotFound
	}
	mockRepo.CreateFunc = func(user *models.User) error {
		user.ID = 1 // Simulate successful creation
		return nil
	}

	username := "testuser"
	email := "testtokenfail@example.com"
	password := "password123"

	user, token, err := authService.Register(username, email, password)

	assert.Error(t, err) // Expect an error from GenerateJWT
	assert.Nil(t, user)  // User might be populated by CreateUser, but overall op fails
	assert.Empty(t, token)
	// We might want to check for a specific error type if GenerateJWT returns one
	// For example: assert.Equal(t, utils.ErrTokenGeneration, err)
}


func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	authService := NewAuthService(mockRepo, "test-secret-login")

	testPassword := "password123"
	hashedTestPassword, _ := utils.HashPassword(testPassword)

	expectedUser := &models.User{
		Model:    gorm.Model{ID: 1},
		Username: "testloginuser",
		Email:    "login@example.com",
		Password: hashedTestPassword,
	}

	mockRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		if email == expectedUser.Email {
			return expectedUser, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	user, token, err := authService.Login(expectedUser.Email, testPassword)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)

	assert.Equal(t, expectedUser.Email, mockRepo.FindByEmailCalledWith)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{}
	authService := NewAuthService(mockRepo, "test-secret-login")

	mockRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		return nil, gorm.ErrRecordNotFound // Simulate user not found
	}

	email := "nonexistent@example.com"
	password := "password123"

	user, token, err := authService.Login(email, password)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	assert.Equal(t, email, mockRepo.FindByEmailCalledWith)
}

func TestAuthService_Login_IncorrectPassword(t *testing.T) {
	mockRepo := &MockUserRepository{}
	authService := NewAuthService(mockRepo, "test-secret-login")

	correctPassword := "correctPassword"
	incorrectPassword := "incorrectPassword"
	hashedCorrectPassword, _ := utils.HashPassword(correctPassword)

	existingUser := &models.User{
		Model:    gorm.Model{ID: 1},
		Username: "testloginuser",
		Email:    "login@example.com",
		Password: hashedCorrectPassword,
	}

	mockRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		if email == existingUser.Email {
			return existingUser, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	user, token, err := authService.Login(existingUser.Email, incorrectPassword)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	assert.Equal(t, existingUser.Email, mockRepo.FindByEmailCalledWith)
}

// TestAuthService_Login_TokenGenerationError (TODO)
func TestAuthService_Login_TokenGenerationError(t *testing.T) {
	mockRepo := &MockUserRepository{}
	// Use the special secret key to force GenerateJWT to fail
	authService := NewAuthService(mockRepo, "FORCE_JWT_ERROR_FOR_TEST")

	testPassword := "password123"
	hashedTestPassword, _ := utils.HashPassword(testPassword)

	existingUser := &models.User{
		Model:    gorm.Model{ID: 1},
		Username: "testloginuser",
		Email:    "login@example.com",
		Password: hashedTestPassword,
	}

	mockRepo.FindByEmailFunc = func(email string) (*models.User, error) {
		if email == existingUser.Email {
			return existingUser, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	user, token, err := authService.Login(existingUser.Email, testPassword)

	assert.Error(t, err) // Expect an error from GenerateJWT
	// User object might be available as FindByEmail and CheckPasswordHash succeeded
	assert.NotNil(t, user)
	assert.Empty(t, token)
	// We might want to check for a specific error type if GenerateJWT returns one
	// For example: assert.Equal(t, utils.ErrTokenGeneration, err)
}
