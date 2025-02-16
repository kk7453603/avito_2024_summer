package authentication

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kk7453603/avito_2024_summer/internal/models"
	"github.com/kk7453603/avito_2024_summer/internal/modules/authentication/mocks"
)

func TestAuthService_GetOrRegUser(t *testing.T) {
	tests := []struct {
		name         string
		existingUser *models.User
		username     string
		password     string
		hashPassword string
		shouldExist  bool
	}{
		{
			name:         "User already exists",
			existingUser: &models.User{Username: "testUser", Password: "hashedPasswd"},
			username:     "testUser",
			password:     "testPasswd",
			shouldExist:  true,
		},
		{
			name:         "User does not exist and is registered",
			existingUser: nil,
			username:     "newUser",
			password:     "newPasswd",
			hashPassword: "hashedNewPasswd",
			shouldExist:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.DataBase)
			mockHasher := new(mocks.Hasher)

			if tt.existingUser == nil {
				mockDB.On("GetUserByUsername", mock.Anything, tt.username).Return((*models.User)(nil), nil)
				mockHasher.On("Hash", tt.password).Return(tt.hashPassword, nil)
				mockDB.On("SaveUser", mock.Anything, mock.Anything).Return(nil)
			} else {
				mockDB.On("GetUserByUsername", mock.Anything, tt.username).Return(tt.existingUser, nil)
			}

			service := New(mockDB, mockHasher)
			ctx, ctxCancel := context.WithCancel(context.Background())

			user, exists, err := service.GetOrRegUser(ctx, tt.username, tt.password)

			require.NoError(t, err)
			require.Equal(t, tt.shouldExist, exists)
			require.Equal(t, tt.username, user.Username)
			if !tt.shouldExist {
				require.Equal(t, tt.hashPassword, user.Password)
			}

			mockDB.AssertExpectations(t)
			mockHasher.AssertExpectations(t)

			ctxCancel()
		})
	}
}

func TestAuthService_GetOrRegUserError(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		mockDBSetup   func(db *mocks.DataBase)
		mockHashSetup func(hasher *mocks.Hasher)
		expectError   bool
	}{
		{
			name:     "DataBase error on GetUserByUsername",
			username: "errorUser",
			password: "testPasswd",
			mockDBSetup: func(db *mocks.DataBase) {
				db.On("GetUserByUsername", mock.Anything, "errorUser").Return(nil, errors.New("db error"))
			},
			mockHashSetup: func(hasher *mocks.Hasher) {},
			expectError:   true,
		},
		{
			name:     "Error on password hashing",
			username: "newUser",
			password: "testPasswd",
			mockDBSetup: func(db *mocks.DataBase) {
				db.On("GetUserByUsername", mock.Anything, "newUser").Return(nil, nil)
			},
			mockHashSetup: func(hasher *mocks.Hasher) {
				hasher.On("Hash", "testPasswd").Return("", errors.New("hash error"))
			},
			expectError: true,
		},
		{
			name:     "Error on saving user",
			username: "newUser",
			password: "testPasswd",
			mockDBSetup: func(db *mocks.DataBase) {
				db.On("GetUserByUsername", mock.Anything, "newUser").Return(nil, nil)
				db.On("SaveUser", mock.Anything, mock.Anything).Return(errors.New("save error"))
			},
			mockHashSetup: func(hasher *mocks.Hasher) {
				hasher.On("Hash", "testPasswd").Return("hashedPasswd", nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.DataBase)
			mockHasher := new(mocks.Hasher)

			tt.mockDBSetup(mockDB)
			tt.mockHashSetup(mockHasher)

			service := New(mockDB, mockHasher)
			ctx, ctxCancel := context.WithCancel(context.Background())

			_, _, err := service.GetOrRegUser(ctx, tt.username, tt.password)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
			mockHasher.AssertExpectations(t)

			ctxCancel()
		})
	}
}

func TestAuthService_ComparePassword(t *testing.T) {
	mockHasher := new(mocks.Hasher)

	service := &AuthService{
		passwd: mockHasher,
	}

	hashedPassword := "hashed_password"
	password := "password"

	mockHasher.On("Compare", hashedPassword, password).Return(true)

	result := service.ComparePassword(hashedPassword, password)
	require.True(t, result)
	mockHasher.AssertCalled(t, "Compare", hashedPassword, password)

	mockHasher.On("Compare", hashedPassword, "wrong_password").Return(false)

	result = service.ComparePassword(hashedPassword, "wrong_password")
	require.False(t, result)
	mockHasher.AssertCalled(t, "Compare", hashedPassword, "wrong_password")
}
