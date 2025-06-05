package services

import (
	"errors"

	"github.com/zayyadi/trello/models"
	"github.com/zayyadi/trello/repositories"
	"github.com/zayyadi/trello/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo     repositories.UserRepositoryInterface
	jwtSecretKey string
}

func NewAuthService(userRepo repositories.UserRepositoryInterface, jwtSecretKey string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecretKey: jwtSecretKey}
}

func (s *AuthService) Register(username, email, password string) (*models.User, string, error) {
	// Check if email or username already exists
	_, err := s.userRepo.FindByEmail(email)
	if err == nil { // nil error means user found
		return nil, "", ErrEmailExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) { // Other DB error
		return nil, "", err
	}
	// Similar check for username if your model has unique username (it does)
	// For brevity, assuming FindByEmail is the primary check for duplicates for now.

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateJWT(user.ID, s.jwtSecretKey)
	if err != nil {
		// If token generation fails during registration, return nil for user.
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(user.ID, s.jwtSecretKey)
	if err != nil {
		// Return the user object even if token generation fails,
		// as authentication (email/password check) itself was successful.
		return user, "", err
	}

	return user, token, nil
}
