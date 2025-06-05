package db

import (
	"fmt"
	"log"

	"github.com/zayyadi/trello/config"
	"github.com/zayyadi/trello/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection and performs migrations
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Or logger.Silent for less verbose logs
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established.")

	// Auto-migrate schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Board{},
		&models.List{},
		&models.Card{},
		&models.BoardMember{},
		&models.Comment{},
		&models.CardCollaborator{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}
	log.Println("Database migration completed.")

	DB = db // Store the instance globally if needed, or pass it around
	return db, nil
}
