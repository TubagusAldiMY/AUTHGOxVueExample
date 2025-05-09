// internal/repository/user_repo.go
package repository // <- Ubah package

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"go-auth-example/internal/model" // <- Import model baru
)

// Definisikan interface untuk UserRepository
type UserRepository interface {
	Create(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id int) (*model.User, error)
}

// Implementasi UserRepository untuk PostgreSQL
type postgresUserRepository struct {
	db *sql.DB // Simpan koneksi DB di struct
}

// NewPostgresUserRepository adalah constructor untuk membuat instance repository
func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

// Ubah fungsi menjadi method dari struct postgresUserRepository
// Gunakan p.db, bukan variabel global DB

func (p *postgresUserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (username, email, password_hash, created_at)
	          VALUES ($1, $2, $3, $4) RETURNING id`

	// Gunakan p.db
	err := p.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, time.Now()).Scan(&user.ID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "users_username_key") {
				return fmt.Errorf("username already exists")
			}
			if strings.Contains(err.Error(), "users_email_key") {
				return fmt.Errorf("email already exists")
			}
		}
		return fmt.Errorf("could not create user: %w", err)
	}
	return nil
}

func (p *postgresUserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1`
	// Gunakan p.db
	err := p.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting user by email %s: %v", email, err)
		return nil, fmt.Errorf("could not get user by email: %w", err)
	}
	return user, nil
}

func (p *postgresUserRepository) GetByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1`
	// Gunakan p.db
	err := p.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting user by ID %d: %v", id, err)
		return nil, fmt.Errorf("could not get user by ID: %w", err)
	}
	return user, nil
}
