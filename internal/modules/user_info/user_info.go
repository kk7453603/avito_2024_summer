//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --all --output=./mocks

// Package user_info provides functionality for retrieving user-related information
// such as coin balance, inventory, and coin transaction history.
package user_info

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kk7453603/avito_2024_summer/internal/models"
)

// DataBase interface defines methods for retrieving user-related data.
type DataBase interface {
	GetCoinsByUserID(ctx context.Context, userID int) (int, error)
	GetInventoryByUserID(ctx context.Context, userID int) (*[]models.Merch, error)
	GetCoinHistoryByUserID(ctx context.Context, userID int) (*models.CoinHistory, error)
}

// UserInfoService provides functionality for retrieving user-related information.
type UserInfoService struct {
	storage DataBase
}

// New creates a new instance of UserInfoService with the given DataBase storage.
func New(storage DataBase) *UserInfoService {
	return &UserInfoService{storage}
}

// GetCoins retrieves the number of coins for a specific user.
func (s *UserInfoService) GetCoins(ctx context.Context, userID int) (int, error) {
	coins, err := s.storage.GetCoinsByUserID(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return coins, nil
}

// GetInventory retrieves the inventory of a specific user, returning an empty list if none exists.
func (s *UserInfoService) GetInventory(ctx context.Context, userID int) (*[]models.Merch, error) {
	inventory, err := s.storage.GetInventoryByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if inventory == nil {
		return &[]models.Merch{}, nil
	}
	return inventory, nil
}

// GetCoinHistory retrieves the coin transaction history of a specific user, ensuring non-nil fields.
func (s *UserInfoService) GetCoinHistory(ctx context.Context, userID int) (*models.CoinHistory, error) {
	coinHistory, err := s.storage.GetCoinHistoryByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if coinHistory == nil {
		coinHistory = &models.CoinHistory{}
	}

	if coinHistory.Receiving == nil {
		coinHistory.Receiving = &[]models.Receiving{}
	}

	if coinHistory.Sending == nil {
		coinHistory.Sending = &[]models.Sending{}
	}

	return coinHistory, nil
}
