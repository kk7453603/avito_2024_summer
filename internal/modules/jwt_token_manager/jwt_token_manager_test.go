package jwt_token_manager

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestTokenManager_NewToken(t *testing.T) {
	cfg := &Config{
		TTL:    "1h",
		secret: "mySecret",
	}

	manager, err := New(cfg)
	require.NoError(t, err)

	userID := "12345"
	username := "testUser"
	token, err := manager.NewToken(userID, username)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signature method: %v", token.Header["alg"])
		}
		return manager.secret, nil
	})
	require.NoError(t, err)
	require.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Errorf("Expected claims to be of type *CustomClaims, but got: %T", parsedToken.Claims)
	}

	require.True(t, ok)
	require.Equal(t, userID, claims["sub"])
	require.Equal(t, username, claims["username"])

	expectedExpiration := time.Now().Add(manager.TTL).Truncate(time.Second)
	actualExpiration, ok := claims["exp"].(float64)
	require.True(t, ok, "exp should be of type float64")
	require.Equal(t, expectedExpiration.Unix(), int64(actualExpiration))
}

func TestTokenManager_ParseToken(t *testing.T) {
	cfg := &Config{
		TTL:    "1h",
		secret: "mySecret",
	}

	manager, err := New(cfg)
	require.NoError(t, err)

	userID := "12345"
	username := "testUser"
	tokenString, err := manager.NewToken(userID, username)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	parsedToken, err := manager.ParseToken(tokenString)
	require.NoError(t, err)
	require.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Errorf("Expected claims to be of type jwt.MapClaims, but got: %T", parsedToken.Claims)
	}

	require.Equal(t, userID, claims["sub"])
	require.Equal(t, username, claims["username"])

	expectedExpiration := time.Now().Add(manager.TTL).Truncate(time.Second)
	actualExpiration, ok := claims["exp"].(float64)
	require.True(t, ok, "exp should be of type float64")
	require.Equal(t, expectedExpiration.Unix(), int64(actualExpiration))
}

func TestTokenManager_ParseTokenError(t *testing.T) {
	cfg := &Config{
		TTL:    "1h",
		secret: "mySecret",
	}

	manager, err := New(cfg)
	require.NoError(t, err)

	userID := "12345"
	username := "testUser"
	tokenString, err := manager.NewToken(userID, username)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "Valid token",
			token:       tokenString,
			expectError: false,
		},
		{
			name:        "Invalid token (modified signature)",
			token:       tokenString[:len(tokenString)-1] + "qwerty",
			expectError: true,
		},
		{
			name:        "Expired token",
			token:       generateExpiredToken(manager, userID, username),
			expectError: true,
		},
		{
			name:        "Empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "Invalid token format (random string)",
			token:       "randomString",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = manager.ParseToken(tt.token)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTokenManager_ParseClaims(t *testing.T) {
	cfg := &Config{
		TTL:    "1h",
		secret: "mySecret",
	}

	manager, err := New(cfg)
	require.NoError(t, err)

	userID := "12345"
	username := "testUser"
	token, err := manager.NewToken(userID, username)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := manager.ParseClaims(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	require.Equal(t, userID, (*claims)["sub"])
	require.Equal(t, username, (*claims)["username"])

	expectedExpiration := time.Now().Add(manager.TTL).Truncate(time.Second)
	actualExpiration, ok := (*claims)["exp"].(float64)
	require.True(t, ok, "exp should be of type float64")
	require.Equal(t, expectedExpiration.Unix(), int64(actualExpiration))
}

func generateExpiredToken(manager *TokenManager, userID, username string) string {
	expiredClaims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"exp":      time.Now().Add(-time.Hour).Unix(), // Устанавливаем время в прошлом
	}

	// Создаем новый токен с истекшими данными.
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString(manager.secret)
	if err != nil {
		panic("Failed to generate expired token: " + err.Error())
	}
	return tokenString
}
