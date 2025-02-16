//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --all --output=./mocks

// Package authentication provides functionality for user authentication,
// including user retrieval,registration, and password management.
package authentication

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kk7453603/avito_2024_summer/internal/models"
)

// DataBase interface defines methods for interacting with the user storage.
type DataBase interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	SaveUser(ctx context.Context, user *models.User) error
}

// Hasher interface defines methods for password hashing and comparison.
type Hasher interface {
	Hash(passwd string) (string, error)
	Compare(hashedPasswd, passwd string) bool
}

// AuthService provides authentication-related functionality.
type AuthService struct {
	storage DataBase
	passwd  Hasher
}

// New creates a new instance of AuthService with the given storage and Hasher.
func New(storage DataBase, passwd Hasher) *AuthService {
	return &AuthService{storage, passwd}
}

// GetOrRegUser retrieves an existing user or registers a new one if they don't exist.
func (s *AuthService) GetOrRegUser(ctx context.Context, username, password string) (*models.User, bool, error) {
	user, err := s.storage.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, false, err
	}

	if user != nil {
		return user, true, nil
	}

	hashedPasswd, err := s.passwd.Hash(password)
	if err != nil {
		return nil, false, err
	}

	user = &models.User{
		Username: username,
		Password: hashedPasswd,
	}

	err = s.storage.SaveUser(ctx, user)
	if err != nil {
		return nil, false, err
	}

	return user, false, nil
}

// ComparePassword checks if the provided password matches the hashed password.
func (s *AuthService) ComparePassword(hashedPasswd, passwd string) bool {
	return s.passwd.Compare(hashedPasswd, passwd)
}
