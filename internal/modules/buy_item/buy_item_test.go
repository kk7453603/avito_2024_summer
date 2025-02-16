package buy_item

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kk7453603/avito_2024_summer/internal/models"
	"github.com/kk7453603/avito_2024_summer/internal/modules/buy_item/mocks"
)

func TestBuyItemService_GetItem(t *testing.T) {
	mockDB := new(mocks.DataBase)
	service := New(mockDB)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	slug := "test-item"
	expectedItem := &models.Item{
		Slug:  slug,
		Title: "Test Item",
		Price: 100,
	}

	mockDB.On("GetItemBySlug", mock.Anything, slug).Return(expectedItem, nil)

	item, err := service.GetItem(ctx, slug)

	require.NoError(t, err)
	require.NotNil(t, item)
	require.Equal(t, expectedItem, item)

	mockDB.AssertExpectations(t)
}

func TestBuyItemService_GetItemError(t *testing.T) {
	tests := []struct {
		name      string
		slug      string
		mockError error
		expectNil bool
		expectErr bool
	}{
		{
			name:      "Item not found",
			slug:      "non-existent-item",
			mockError: sql.ErrNoRows,
			expectNil: true,
			expectErr: false,
		},
		{
			name:      "Database error",
			slug:      "db-error-item",
			mockError: errors.New("database connection error"),
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.DataBase)
			service := New(mockDB)
			ctx, ctxCancel := context.WithCancel(context.Background())

			mockDB.On("GetItemBySlug", mock.Anything, tt.slug).Return(nil, tt.mockError)

			item, err := service.GetItem(ctx, tt.slug)

			if tt.expectNil {
				require.Nil(t, item)
			} else {
				require.NotNil(t, item)
			}

			if tt.expectErr {
				require.Error(t, err)
				require.Equal(t, tt.mockError, err)
			} else {
				require.NoError(t, err)
			}

			mockDB.AssertExpectations(t)

			ctxCancel()
		})
	}
}

func TestBuyItemService_GetBuyerCoins(t *testing.T) {
	mockDB := new(mocks.DataBase)
	service := New(mockDB)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	userID := 1
	expectedCoins := 1000

	mockDB.On("GetCoinsByUserID", mock.Anything, userID).Return(expectedCoins, nil)

	coins, err := service.GetBuyerCoins(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, expectedCoins, coins)

	mockDB.AssertExpectations(t)
}

func TestBuyItemService_GetBuyerCoinsError(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		mockError     error
		expectedErr   error
		expectedCoins int
	}{
		{
			name:          "User not found",
			userID:        999,
			mockError:     sql.ErrNoRows,
			expectedErr:   nil, // Ошибка sql. ErrNoRows обрабатывается внутри метода и не возвращается
			expectedCoins: 0,
		},
		{
			name:          "Database error",
			userID:        1,
			mockError:     errors.New("database error"),
			expectedErr:   errors.New("database error"),
			expectedCoins: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.DataBase)
			service := New(mockDB)
			ctx, ctxCancel := context.WithCancel(context.Background())

			mockDB.On("GetCoinsByUserID", mock.Anything, tt.userID).Return(0, tt.mockError)

			coins, err := service.GetBuyerCoins(ctx, tt.userID)
			require.Equal(t, tt.expectedCoins, coins)

			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}

			mockDB.AssertExpectations(t)

			ctxCancel()
		})
	}
}

func TestBuyItemService_BuyItem(t *testing.T) {
	item := &models.Item{
		Slug:  "valid-item",
		Title: "Valid Item",
		Price: 100,
	}

	tests := []struct {
		name        string
		userID      int
		item        *models.Item
		mockError   error
		expectedErr error
	}{
		{
			name:        "No errors",
			userID:      1,
			item:        item,
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:        "Database error",
			userID:      1,
			item:        item,
			mockError:   errors.New("database error"),
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.DataBase)
			service := New(mockDB)
			ctx, ctxCancel := context.WithCancel(context.Background())

			mockDB.On("MakePurchaseByUserID", mock.Anything, tt.userID, tt.item).Return(tt.mockError)

			err := service.BuyItem(ctx, tt.userID, tt.item)
			if tt.expectedErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockDB.AssertExpectations(t)

			ctxCancel()
		})
	}
}
