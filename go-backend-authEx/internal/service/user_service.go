package service

import (
	"errors"
	"fmt"
	"log"

	"go-auth-example/internal/model"
	"go-auth-example/internal/repository"
)

// UserService interface
type UserService interface {
	GetUserProfile(userID int) (*model.User, error)
}

// userService struct
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService constructor
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// GetUserProfile implementation
func (s *userService) GetUserProfile(userID int) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		log.Printf("Error fetching user data for profile (ID: %d): %v", userID, err)
		return nil, fmt.Errorf("failed to fetch user profile data")
	}
	if user == nil {
		// User di token tidak ada di DB? Kasus aneh.
		return nil, errors.New("user associated with token not found")
	}

	// Jangan kirim password hash
	user.PasswordHash = ""
	return user, nil
}
