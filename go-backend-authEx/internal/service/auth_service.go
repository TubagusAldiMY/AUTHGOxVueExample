package service

import (
	"errors"
	"fmt"
	"log"

	"go-auth-example/internal/auth" // <- Import auth utils
	"go-auth-example/internal/model"
	"go-auth-example/internal/repository" // <- Import repository interface
)

// AuthService interface mendefinisikan operasi otentikasi
type AuthService interface {
	Register(input model.RegisterInput) (*model.User, error)
	Login(input model.LoginInput) (string, error) // Return JWT string
}

// authService struct mengimplementasikan AuthService
type authService struct {
	userRepo repository.UserRepository // <- Dependensi ke interface repo
}

// NewAuthService adalah constructor untuk authService
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Implementasi Register
func (s *authService) Register(input model.RegisterInput) (*model.User, error) {
	// Cek apakah email sudah ada (contoh business logic sederhana)
	existingUser, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		// Log error database internal, tapi jangan ekspos detailnya ke user
		log.Printf("Error checking email existence: %v", err)
		return nil, fmt.Errorf("failed to check user existence")
	}
	if existingUser != nil {
		return nil, errors.New("email already registered") // Error spesifik bisnis
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		log.Printf("Error hashing password during registration: %v", err)
		return nil, fmt.Errorf("failed to process registration")
	}

	// Buat user baru
	newUser := &model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	// Simpan user ke database via repository
	err = s.userRepo.Create(newUser)
	if err != nil {
		// Error dari repository (misal, username conflict) sudah cukup deskriptif
		return nil, err // Teruskan error dari repo
	}

	// Penting: Hapus hash password sebelum dikembalikan
	newUser.PasswordHash = ""
	return newUser, nil
}

// Implementasi Login
func (s *authService) Login(input model.LoginInput) (string, error) {
	// Cari user berdasarkan email via repository
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		log.Printf("Database error during login for email %s: %v", input.Email, err)
		return "", errors.New("an error occurred during login")
	}
	if user == nil {
		return "", errors.New("invalid email or password") // Pesan error generik
	}

	// Cek password
	if !auth.CheckPasswordHash(input.Password, user.PasswordHash) {
		return "", errors.New("invalid email or password") // Pesan error generik
	}

	// Buat token JWT
	token, err := auth.GenerateJWT(*user)
	if err != nil {
		log.Printf("Error generating JWT for user %d: %v", user.ID, err)
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
