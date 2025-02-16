package user_info

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kk7453603/avito_2024_summer/internal/models"
	"github.com/kk7453603/avito_2024_summer/internal/modules/user_info/mocks"
)

func TestUserInfoService_GetCoins(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		wantCoins int
		wantErr   error
	}{
		{
			name:      "User found",
			userID:    1,
			wantCoins: 1000,
			wantErr:   nil,
		},
		{
			name:      "User not found",
			userID:    0,
			wantCoins: 0,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.DataBase)
			service := New(mockDB)
			ctx, ctxCancel := context.WithCancel(context.Background())

			mockDB.On("GetCoinsByUserID", mock.Anything, tt.userID).Return(tt.wantCoins, tt.wantErr).Once()

			coins, err := service.GetCoins(ctx, tt.userID)

			require.NoError(t, err)
			require.Equal(t, tt.wantCoins, coins)

			mockDB.AssertExpectations(t)

			ctxCancel()
		})
	}
}

func TestUserInfoService_GetCoinsError(t *testing.T) {
	wantUserID := 0
	wantCoins := 0
	wantErr := errors.New("database error")

	mockDB := new(mocks.DataBase)
	service := New(mockDB)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	mockDB.On("GetCoinsByUserID", mock.Anything, wantUserID).Return(wantCoins, wantErr).Once()

	coins, err := service.GetCoins(ctx, wantUserID)

	require.Error(t, err)
	require.Equal(t, wantCoins, coins)

	mockDB.AssertExpectations(t)
}

func TestUserInfoService_GetInventory(t *testing.T) {
	mockDB := new(mocks.DataBase)
	service := New(mockDB)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	validInventory := []models.Merch{
		{Type: "book", Quantity: 1},
		{Type: "pen", Quantity: 2},
	}

	tests := []struct {
		name           string
		userID         int
		mockInventory  *[]models.Merch
		mockError      error
		expectedResult *[]models.Merch
	}{
		{
			name:           "Successful retrieval of inventory",
			userID:         1,
			mockInventory:  &validInventory,
			mockError:      nil,
			expectedResult: &validInventory,
		},
		{
			name:           "Empty inventory",
			userID:         2,
			mockInventory:  &[]models.Merch{},
			mockError:      nil,
			expectedResult: &[]models.Merch{},
		},
		{
			name:           "Nil inventory",
			userID:         3,
			mockInventory:  nil,
			mockError:      nil,
			expectedResult: &[]models.Merch{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB.On("GetInventoryByUserID", mock.Anything, tt.userID).Return(tt.mockInventory, tt.mockError)

			result, err := service.GetInventory(ctx, tt.userID)

			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)

			mockDB.AssertExpectations(t)
		})
	}
}

func TestUserInfoService_GetInventoryError(t *testing.T) {
	mockDB := new(mocks.DataBase)
	service := New(mockDB)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	userID := 0
	wantErr := errors.New("database error")

	mockDB.On("GetInventoryByUserID", mock.Anything, userID).Return(nil, wantErr)

	result, err := service.GetInventory(ctx, userID)

	require.Error(t, err)
	require.Nil(t, result)

	mockDB.AssertExpectations(t)
}

func TestUserInfoService_GetCoinHistory(t *testing.T) {
	mockDB := new(mocks.DataBase)
	service := New(mockDB)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	validCoinHistory := models.CoinHistory{
		Receiving: &[]models.Receiving{
			{User: "johnDoe", Amount: 100},
		},
		Sending: &[]models.Sending{
			{User: "nickles-cage", Amount: 50},
		},
	}

	tests := []struct {
		name            string
		userID          int
		mockCoinHistory *models.CoinHistory
		mockError       error
		expectedResult  *models.CoinHistory
	}{
		{
			name:            "Successful retrieval of coin history with non-nil fields",
			userID:          1,
			mockCoinHistory: &validCoinHistory,
			mockError:       nil,
			expectedResult:  &validCoinHistory,
		},
		{
			name:   "Successful retrieval of coin history with nil fields",
			userID: 2,
			mockCoinHistory: &models.CoinHistory{
				Receiving: nil,
				Sending:   nil,
			},
			mockError: nil,
			expectedResult: &models.CoinHistory{
				Receiving: &[]models.Receiving{},
				Sending:   &[]models.Sending{},
			},
		},
		{
			name:            "Successful retrieval of coin history with invalid pointer",
			userID:          3,
			mockCoinHistory: nil,
			mockError:       nil,
			expectedResult: &models.CoinHistory{
				Receiving: &[]models.Receiving{},
				Sending:   &[]models.Sending{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB.On("GetCoinHistoryByUserID", mock.Anything, tt.userID).Return(tt.mockCoinHistory, tt.mockError)

			result, err := service.GetCoinHistory(ctx, tt.userID)

			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)

			mockDB.AssertExpectations(t)
		})
	}
}

func TestUserInfoService_GetCoinHistoryError(t *testing.T) {
	mockDB := new(mocks.DataBase)
	service := New(mockDB)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	userID := 0
	wantErr := errors.New("database error")

	mockDB.On("GetCoinHistoryByUserID", mock.Anything, userID).Return(nil, wantErr)

	result, err := service.GetCoinHistory(ctx, userID)

	require.Error(t, err)
	require.Nil(t, result)

	mockDB.AssertExpectations(t)
}
