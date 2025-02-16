//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --all --output=./mocks

// Package buy_item provides functionality for handling the purchase of items by users.
// It includes methods for retrieving item details, checking a buyer's coin balance,
// and processing purchases.
package buy_item

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kk7453603/avito_2024_summer/internal/models"
)

// DataBase interface defines methods for handling item purchases and user data.
type DataBase interface {
	GetItemBySlug(ctx context.Context, slug string) (*models.Item, error)
	GetCoinsByUserID(ctx context.Context, userID int) (int, error)
	MakePurchaseByUserID(ctx context.Context, userID int, item *models.Item) error
}

// BuyItemService provides functionality for handling item purchases.
type BuyItemService struct {
	storage DataBase
}

// New creates a new instance of BuyItemService with the given storage.
func New(storage DataBase) *BuyItemService {
	return &BuyItemService{storage}
}

// GetItem retrieves an item by its slug, handling DataBase errors.
func (s *BuyItemService) GetItem(ctx context.Context, slug string) (*models.Item, error) {
	item, err := s.storage.GetItemBySlug(ctx, slug)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return item, nil
}

// GetBuyerCoins retrieves the number of coins a buyer has by their ID.
func (s *BuyItemService) GetBuyerCoins(ctx context.Context, userID int) (int, error) {
	coins, err := s.storage.GetCoinsByUserID(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return coins, nil
}

// BuyItem processes the purchase of an item by a user.
func (s *BuyItemService) BuyItem(ctx context.Context, userID int, item *models.Item) error {
	return s.storage.MakePurchaseByUserID(ctx, userID, item)
}
