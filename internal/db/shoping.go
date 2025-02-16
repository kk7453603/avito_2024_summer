// Package db provides functionality for interacting with the PostgreSQL database.
package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/kk7453603/avito_2024_summer/internal/models"
)

const (
	getIDByUsername                = `SELECT id FROM users WHERE username=$1`
	getUserByUsername              = `SELECT * FROM users WHERE username=$1`
	getCoinsByUserID               = `SELECT coins FROM users WHERE id=$1`
	getInventoryByUserID           = `SELECT item_slug, quantity FROM inventory WHERE user_id = $1`
	getReceivedCoinHistoryByUserID = `SELECT u.username, t.coins FROM transactions t JOIN users u ON t.sender_id = u.id WHERE t.receiver_id = $1;`
	getSendingCoinHistoryByUserID  = `SELECT u.username, t.coins FROM transactions t JOIN users u ON t.receiver_id = u.id WHERE t.sender_id = $1;`
	saveUser                       = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, coins, created_at, updated_at;`
	subtractFromCoinsByUserID      = `UPDATE users SET coins = COALESCE(coins, 0) - $1 WHERE id = $2;`
	addToCoinsByUserID             = `UPDATE users SET coins = COALESCE(coins, 0) + $1 WHERE id = $2;`
	recordTransaction              = `INSERT INTO transactions (sender_id, receiver_id, coins) VALUES($1, $2, $3);`
	getItemBySlug                  = `SELECT * FROM store WHERE slug = $1;`
	addItemToInventoryByUserID     = `
		INSERT INTO inventory (user_id, item_slug, quantity, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, item_slug) 
		DO UPDATE SET quantity = inventory.quantity + excluded.quantity, updated_at = NOW();`
)

// GetIDByUsername retrieves the user ID associated with the given username.
func (s *Storage) GetIDByUsername(ctx context.Context, username string) (int, error) {
	id := 0
	err := s.pool.QueryRow(ctx, getIDByUsername, username).Scan(&id)
	return id, err
}

// GetUserByUsername retrieves a user's details by their username.
func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := s.pool.QueryRow(ctx, getUserByUsername, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Coins,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// GetCoinsByUserID retrieves the number of coins a user has by their ID.
func (s *Storage) GetCoinsByUserID(ctx context.Context, userID int) (int, error) {
	coins := 0
	err := s.pool.QueryRow(ctx, getCoinsByUserID, userID).Scan(&coins)
	return coins, err
}

// GetInventoryByUserID retrieves the inventory of a user by their ID.
func (s *Storage) GetInventoryByUserID(ctx context.Context, userID int) (*[]models.Merch, error) {
	rows, err := s.pool.Query(ctx, getInventoryByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Merch, 0, 8)
	for rows.Next() {
		var item models.Merch
		err := rows.Scan(&item.Type, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &items, nil
}

// GetCoinHistoryByUserID retrieves the coin transaction history of a user by their ID.
func (s *Storage) GetCoinHistoryByUserID(ctx context.Context, userID int) (*models.CoinHistory, error) {
	g, gCtx := errgroup.WithContext(ctx)

	// RECEIVED
	var recs *[]models.Receiving
	g.Go(func() error {
		data, err := fetchCoinHistory[models.Receiving](gCtx, s.pool, getReceivedCoinHistoryByUserID, userID)
		if err != nil {
			return err
		}
		recs = data
		return nil
	})

	// SENT
	var sends *[]models.Sending
	g.Go(func() error {
		data, err := fetchCoinHistory[models.Sending](gCtx, s.pool, getSendingCoinHistoryByUserID, userID)
		if err != nil {
			return err
		}
		sends = data
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	ch := models.CoinHistory{
		Receiving: recs,
		Sending:   sends,
	}
	return &ch, nil
}

// coinHistory is a generic constraint for coin history types (Receiving or Sending).
type coinHistory interface {
	models.Receiving | models.Sending
}

// fetchCoinHistory fetches coin history data (either Receiving or Sending) for a user.
func fetchCoinHistory[T coinHistory](ctx context.Context, pool *pgxpool.Pool, query string, userID int) (*[]T, error) {
	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])

	return &items, err
}

// SaveUser saves a new user to the database and updates the user struct with generated fields.
func (s *Storage) SaveUser(ctx context.Context, user *models.User) error {
	err := s.pool.QueryRow(ctx, saveUser, user.Username, user.Password).Scan(
		&user.ID,
		&user.Coins,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// TransferCoins transfers coins from one user to another and records the transaction.
func (s *Storage) TransferCoins(ctx context.Context, fromUserID, toUserID, coins int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	// Subtract money from the sender
	_, err = tx.Exec(ctx, subtractFromCoinsByUserID, coins, fromUserID)
	if err != nil {
		return err
	}

	// Adding money to the recipient
	_, err = tx.Exec(ctx, addToCoinsByUserID, coins, toUserID)
	if err != nil {
		return err
	}

	// Transaction record
	_, err = tx.Exec(ctx, recordTransaction, fromUserID, toUserID, coins)
	if err != nil {
		return err
	}

	return nil
}

// MakePurchaseByUserID processes a purchase of an item by a user.
func (s *Storage) MakePurchaseByUserID(ctx context.Context, userID int, item *models.Item) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	// Subtract money from the user
	_, err = tx.Exec(ctx, subtractFromCoinsByUserID, item.Price, userID)
	if err != nil {
		return err
	}

	// Add the item to the inventory
	_, err = tx.Exec(ctx, addItemToInventoryByUserID, userID, item.Slug, 1)
	if err != nil {
		return err
	}

	return nil
}

// GetItemBySlug retrieves an item's details by its slug.
func (s *Storage) GetItemBySlug(ctx context.Context, slug string) (*models.Item, error) {
	var item models.Item
	err := s.pool.QueryRow(ctx, getItemBySlug, slug).Scan(&item.Slug, &item.Title, &item.Price)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
