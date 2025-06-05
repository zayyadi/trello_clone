package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zayyadi/trello/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	// Using "file::memory:?cache=shared" allows multiple connections to the same in-memory DB if needed,
	// though for these tests, each test usually gets a fresh DB.
	// Simpler in-memory: "file:test.db?mode=memory&cache=shared" or just sqlite.Open(":memory:")
	// Using ":memory:" for simplicity as each test will get its own repo instance.
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		// Disable verbose logging in tests unless debugging
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate schema
	// Need to migrate all related models if there are foreign key constraints
	// For UserRepository, only User model is directly interacted with,
	// but User has relations to Board, BoardMember, Card.
	// For simplicity here, only migrating User. If FK issues arise, migrate others.
	err = db.AutoMigrate(&models.User{}, &models.Board{}, &models.List{}, &models.Card{}, &models.BoardMember{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestUserRepository_Create_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db) // This should return UserRepositoryInterface

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID, "User ID should be set after creation")

	// Verify by finding
	foundUser, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user1 := &models.User{
		Username: "user1",
		Email:    "duplicate@example.com",
		Password: "password123",
	}
	err := repo.Create(user1)
	assert.NoError(t, err)

	user2 := &models.User{
		Username: "user2",
		Email:    "duplicate@example.com", // Same email
		Password: "password456",
	}
	err = repo.Create(user2)
	assert.Error(t, err, "Should return an error for duplicate email")
	// GORM often returns a generic error that wraps the driver-specific error.
	// For SQLite, this is typically something like "UNIQUE constraint failed: users.email".
	// Checking for a specific GORM error type might be too brittle if not available.
	// For now, just assert.Error is a good start.
	// A more robust check might involve string containment if GORM doesn't provide a typed error.
	// assert.Contains(t, err.Error(), "UNIQUE constraint failed: users.email") // This is driver specific
}

func TestUserRepository_Create_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user1 := &models.User{
		Username: "duplicateuser",
		Email:    "user1@example.com",
		Password: "password123",
	}
	err := repo.Create(user1)
	assert.NoError(t, err)

	user2 := &models.User{
		Username: "duplicateuser", // Same username
		Email:    "user2@example.com",
		Password: "password456",
	}
	err = repo.Create(user2)
	assert.Error(t, err, "Should return an error for duplicate username")
	// Similar to email, check for error, specific error message check is driver-dependent
	// assert.Contains(t, err.Error(), "UNIQUE constraint failed: users.username")
}

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	emailToFind := "findme@example.com"
	user := &models.User{
		Username: "findmeuser",
		Email:    emailToFind,
		Password: "password123",
	}
	err := repo.Create(user)
	assert.NoError(t, err)

	foundUser, err := repo.FindByEmail(emailToFind)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// No users created, or try to find one that doesn't exist
	foundUser, err := repo.FindByEmail("nonexistent@example.com")
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "Error should be gorm.ErrRecordNotFound")
	assert.Nil(t, foundUser)
}

func TestUserRepository_FindByID_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Username: "findbyiduser",
		Email:    "findbyid@example.com",
		Password: "password123",
	}
	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID) // Ensure ID was set

	foundUser, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// No users created, or find a non-existent ID
	nonExistentID := uint(999)
	foundUser, err := repo.FindByID(nonExistentID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "Error should be gorm.ErrRecordNotFound")
	assert.Nil(t, foundUser)
}
