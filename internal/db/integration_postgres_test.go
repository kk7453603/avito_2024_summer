//go:build integration

package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/kk7453603/avito_2024_summer/internal/models"
)

var (
	ctx = context.Background()
	cfg = &Config{
		Host:     "localhost",
		Port:     "2345",
		Name:     "test_db",
		User:     "test_user",
		Password: "test_password",
	}
	storage *Storage
	pool    *pgxpool.Pool
)

func clearDataBase(t *testing.T) {
	_, err := pool.Exec(ctx, "TRUNCATE TABLE users, inventory, transactions CASCADE")
	require.NoError(t, err)
}

func TestMain(m *testing.M) {
	var err error
	storage, err = NewPostgresPool(ctx, cfg)
	if err != nil {
		panic(err)
	}

	dsn := getPsqlDsn(cfg)
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	// Запуск тестов
	code := m.Run()

	storage.Close()
	os.Exit(code)
}

func TestStorage_NewPostgresPool(t *testing.T) {
	t.Parallel()

	st, err := NewPostgresPool(ctx, cfg)
	require.NoError(t, err, "NewPostgresPool should successfully create the connection")
	require.NotNil(t, st, "Storage (connection pool) must not be nil")

	err = st.pool.Ping(ctx)
	require.NoError(t, err, "Ping to the database must be successful")

	st.Close()

	err = st.pool.Ping(ctx)
	require.Error(t, err, "Ping after closing the connection should return an error")
}

func TestStorage_GetIDByUsername(t *testing.T) {
	clearDataBase(t)

	t.Run("ExistingUser", func(t *testing.T) {
		user := &models.User{
			Username: "testUser1",
			Password: "hashed_password_1",
		}

		err := pool.QueryRow(ctx, saveUser, user.Username, user.Password).Scan(
			&user.ID,
			&user.Coins,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		require.NoError(t, err)

		id, err := storage.GetIDByUsername(ctx, user.Username)
		require.NoError(t, err)
		require.Equal(t, user.ID, id)
	})

	t.Run("NonExistingUser", func(t *testing.T) {
		id, err := storage.GetIDByUsername(ctx, "non_existing_user")
		require.Error(t, err)
		require.Equal(t, id, 0)
	})
}

func TestStorage_GetUserByUsername(t *testing.T) {
	clearDataBase(t)

	t.Run("ExistingUser", func(t *testing.T) {
		user := &models.User{
			Username: "testUser2",
			Password: "hashed_password_2",
		}

		err := pool.QueryRow(ctx, saveUser, user.Username, user.Password).Scan(
			&user.ID,
			&user.Coins,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		require.NoError(t, err)

		fetchedUser, err := storage.GetUserByUsername(ctx, user.Username)
		require.NoError(t, err)

		require.Equal(t, user.ID, fetchedUser.ID)
		require.Equal(t, user.Username, fetchedUser.Username)
		require.Equal(t, user.Password, fetchedUser.Password)
		require.Equal(t, user.Coins, fetchedUser.Coins)
		require.WithinDuration(t, user.CreatedAt, fetchedUser.CreatedAt, time.Second)
		require.WithinDuration(t, user.UpdatedAt, fetchedUser.UpdatedAt, time.Second)
	})

	t.Run("NonExistingUser", func(t *testing.T) {
		fetchedUser, err := storage.GetUserByUsername(ctx, "non_existing_user")
		require.Error(t, err)
		require.Nil(t, fetchedUser)
	})
}
