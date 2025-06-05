package services

import (
	"testing"

	// "github.com/stretchr/testify/assert" // Not needed for setupTestDB itself
	"github.com/zayyadi/trello/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB initializes an in-memory SQLite database for testing for the services package.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to database for test: %v", err)
	}

	// Migrate all known models
	err = db.AutoMigrate(
		&models.User{},
		&models.Board{},
		&models.List{},
		&models.Card{},
		&models.BoardMember{},
		&models.Comment{},
		&models.CardCollaborator{}, // Ensure this is also migrated
	)
	if err != nil {
		t.Fatalf("Failed to migrate database for test: %v", err)
	}
	return db
}

// To make this file a valid test file, add at least one test function.
func TestMainServices(t *testing.T) {
    // This function can be used for package-level setup/teardown if needed,
    // or just exist to make the go compiler treat this as a test file.
}
