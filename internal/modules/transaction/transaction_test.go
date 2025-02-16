package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kk7453603/avito_2024_summer/internal/modules/transaction/mocks"
)

func TestTransactService_GetIDRecipient(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		wantUserID int
		wantErr    error
	}{
		{
			name:       "User found",
			username:   "engineer-e8",
			wantUserID: 123,
			wantErr:    nil,
		},
		{
			name:       "User not found",
			username:   "abracadabra",
			wantUserID: 0,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.DataBase)
			service := New(mockDB)
			ctx, ctxCancel := context.WithCancel(context.Background())

			mockDB.On("GetIDByUsername", mock.Anything, tt.username).Return(tt.wantUserID, tt.wantErr).Once()

			id, err := service.GetIDRecipient(ctx, tt.username)

			require.NoError(t, err)
			require.Equal(t, tt.wantUserID, id)

			mockDB.AssertExpectations(t)

			ctxCancel()
		})
	}
}

func TestTransactService_GetIDRecipientError(t *testing.T) {
	username := "johnDoe"
	wantUserID := 0
	wantErr := errors.New("database error")

	mockDB := new(mocks.DataBase)
	service := New(mockDB)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	mockDB.On("GetIDByUsername", mock.Anything, username).Return(wantUserID, wantErr).Once()

	id, err := service.GetIDRecipient(ctx, username)

	require.Error(t, err)
	require.Equal(t, wantUserID, id)

	mockDB.AssertExpectations(t)
}

func TestTransactService_GetSenderCoins(t *testing.T) {
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

			coins, err := service.GetSenderCoins(ctx, tt.userID)

			require.NoError(t, err)
			require.Equal(t, tt.wantCoins, coins)

			mockDB.AssertExpectations(t)

			ctxCancel()
		})
	}
}

func TestTransactService_GetSenderCoinsError(t *testing.T) {
	wantUserID := 0
	wantCoins := 0
	wantErr := errors.New("database error")

	mockDB := new(mocks.DataBase)
	service := New(mockDB)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	mockDB.On("GetCoinsByUserID", mock.Anything, wantUserID).Return(wantCoins, wantErr).Once()

	coins, err := service.GetSenderCoins(ctx, wantUserID)

	require.Error(t, err)
	require.Equal(t, wantUserID, coins)

	mockDB.AssertExpectations(t)
}

func TestTransactService_SendCoinsToUser(t *testing.T) {
	ErrInDB := errors.New("database error")
	tests := []struct {
		name        string
		senderID    int
		recipientID int
		coins       int
		wantErr     bool
		expErr      error
	}{
		{
			name:        "User found",
			senderID:    1,
			recipientID: 2,
			coins:       100,
			wantErr:     false,
			expErr:      nil,
		},
		{
			name:        "User not found",
			senderID:    1,
			recipientID: 0,
			coins:       100,
			wantErr:     true,
			expErr:      ErrInDB,
		},
		{
			name:        "User not found",
			senderID:    0,
			recipientID: 1,
			coins:       100,
			wantErr:     true,
			expErr:      ErrInDB,
		},
		{
			name:        "Sending to myself",
			senderID:    1,
			recipientID: 1,
			coins:       100,
			wantErr:     true,
			expErr:      ErrInDB,
		},
		{
			name:        "Sending 0 coins",
			senderID:    1,
			recipientID: 2,
			coins:       0,
			wantErr:     true,
			expErr:      ErrInDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.DataBase)
			service := New(mockDB)
			ctx, ctxCancel := context.WithCancel(context.Background())

			mockDB.On("TransferCoins", mock.Anything, tt.senderID, tt.recipientID, tt.coins).Return(tt.expErr).Once()

			err := service.SendCoinsToUser(ctx, tt.senderID, tt.recipientID, tt.coins)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockDB.AssertExpectations(t)

			ctxCancel()
		})
	}
}
