// internal/service/auth_service.go
package service

import (
	"errors"
	"fmt"
	// "log" // Dihapus, diganti dengan logger kustom

	"go-auth-example/internal/auth"
	"go-auth-example/internal/logger" // <-- IMPORT LOGGER KUSTOM
	"go-auth-example/internal/model"
	"go-auth-example/internal/repository"

	"github.com/sirupsen/logrus" // <-- Impor logrus untuk Fields
)

// AuthService interface mendefinisikan operasi otentikasi
type AuthService interface {
	Register(input model.RegisterInput) (*model.User, error)
	Login(input model.LoginInput) (string, error) // Return JWT string
}

// authService struct mengimplementasikan AuthService
type authService struct {
	userRepo repository.UserRepository // Dependensi ke interface repo
}

// NewAuthService adalah constructor untuk authService
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Implementasi Register
func (s *authService) Register(input model.RegisterInput) (*model.User, error) {
	// Definisikan field log yang umum untuk method ini
	logFields := logrus.Fields{
		"service":  "AuthService",
		"method":   "Register",
		"email":    input.Email,
		"username": input.Username,
	}

	// Cek apakah email sudah ada
	existingUserByEmail, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		// Log error database internal
		logger.Log.WithFields(logFields).Errorf("Error checking email existence: %v", err)
		// Kembalikan error generik ke handler, handler akan memetakannya ke APIError
		return nil, fmt.Errorf("failed to check user existence")
	}
	if existingUserByEmail != nil {
		logger.Log.WithFields(logFields).Info("Registration attempt with existing email.")
		return nil, errors.New("email already registered") // Error spesifik bisnis
	}

	// (Opsional) Cek apakah username sudah ada jika repository Anda memiliki GetByUsername
	// existingUserByUsername, err := s.userRepo.GetByUsername(input.Username)
	// if err != nil {
	//    logger.Log.WithFields(logFields).Errorf("Error checking username existence: %v", err)
	//    return nil, fmt.Errorf("failed to check user existence")
	// }
	// if existingUserByUsername != nil {
	//    logger.Log.WithFields(logFields).Info("Registration attempt with existing username.")
	//    return nil, errors.New("username already registered")
	// }
	// Catatan: Penanganan error duplikasi dari DB (seperti di user_repo.go) juga penting.

	// Hash password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		logger.Log.WithFields(logFields).Errorf("Error hashing password during registration: %v", err)
		return nil, fmt.Errorf("failed to process registration")
	}

	// Buat user baru
	newUser := &model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		// CreatedAt akan di-generate oleh DB atau di user_repo.go
	}

	// Simpan user ke database via repository
	err = s.userRepo.Create(newUser) // Asumsi Create akan mengisi newUser.ID dan newUser.CreatedAt
	if err != nil {
		// Error dari repository (misal, username/email conflict yang lolos cek sebelumnya atau error DB lain)
		// Pesan error dari repo sudah cukup deskriptif jika ada (seperti "username already exists")
		logger.Log.WithFields(logFields).Errorf("Error creating user in repository: %v", err)
		return nil, err // Teruskan error dari repo
	}

	// Tambahkan user_id ke log setelah berhasil dibuat
	logFields["user_id"] = newUser.ID
	logger.Log.WithFields(logFields).Info("User successfully registered by service.")

	// Penting: Hapus hash password sebelum dikembalikan
	newUser.PasswordHash = ""
	return newUser, nil
}

// Implementasi Login
func (s *authService) Login(input model.LoginInput) (string, error) {
	logFields := logrus.Fields{
		"service": "AuthService",
		"method":  "Login",
		"email":   input.Email,
	}

	// Cari user berdasarkan email via repository
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		logger.Log.WithFields(logFields).Errorf("Database error during login for email %s: %v", input.Email, err)
		// Kembalikan error generik, handler akan memetakannya
		return "", errors.New("an error occurred during login")
	}
	if user == nil {
		logger.Log.WithFields(logFields).Warn("Login attempt for non-existent email.")
		return "", errors.New("invalid email or password") // Pesan error generik
	}

	// Cek password
	if !auth.CheckPasswordHash(input.Password, user.PasswordHash) {
		// Tambahkan user_id ke log jika user ditemukan tapi password salah
		logFields["user_id_attempted"] = user.ID
		logger.Log.WithFields(logFields).Warn("Invalid password attempt for existing user.")
		return "", errors.New("invalid email or password") // Pesan error generik
	}

	// Buat token JWT
	// Kita akan mengirimkan seluruh user model ke GenerateJWT, jadi pastikan tidak ada info sensitif selain yang dibutuhkan claims
	token, err := auth.GenerateJWT(*user) // GenerateJWT ada di internal/auth/jwt.go
	if err != nil {
		logFields["user_id"] = user.ID
		logger.Log.WithFields(logFields).Errorf("Error generating JWT for user %d: %v", user.ID, err)
		return "", errors.New("failed to generate token")
	}

	logFields["user_id"] = user.ID
	logger.Log.WithFields(logFields).Info("User successfully logged in by service.")
	return token, nil
}
