//go:build integration

package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kk7453603/avito_2024_summer/internal/models"
	"github.com/kk7453603/avito_2024_summer/internal/server/handlers/mocks"
	"github.com/kk7453603/avito_2024_summer/internal/server/middlewares"
)

var validToken = "validToken"

// dummyTokenManager – простая реализация TokenManager для тестирования.
// При получении токена "validToken" возвращает корректные claims.
type dummyTokenManager struct{}

func (d *dummyTokenManager) NewToken(userID string, username string) (string, error) {
	return validToken, nil
}

func (d *dummyTokenManager) ParseClaims(token string) (*jwt.MapClaims, error) {
	if token == validToken {
		claims := jwt.MapClaims{
			"sub":      "1",        // идентификатор пользователя (строкой)
			"username": "testUser", // имя пользователя
		}
		return &claims, nil
	}
	return nil, errors.New("invalid token")
}

// TestUserHandlers_SendCoinsHandler проверяет сценарий успешной передачи монеток другому сотруднику.
func TestUserHandlers_SendCoinsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	recipientUser := models.User{
		ID:       2,
		Username: "otherUser",
		Coins:    100,
	}

	senderUser := models.User{
		ID:       1,
		Username: "testUser",
		Coins:    250,
	}

	amountCoins := 50

	mTxSvc := mocks.NewTransactionService(t)
	// Ожидаем следующее поведение:
	// При вызове GetIDRecipient с "otherUser" возвращаем идентификатор 2.
	mTxSvc.
		On("GetIDRecipient", mock.Anything, recipientUser.Username).
		Return(recipientUser.ID, nil)
	// При вызове GetSenderCoins для пользователя с id=1 возвращаем 250 монет.
	mTxSvc.
		On("GetSenderCoins", mock.Anything, senderUser.ID).
		Return(senderUser.Coins, nil)
	// При вызове SendCoinsToUser с параметрами (1, 2, 50) возвращаем nil.
	mTxSvc.
		On("SendCoinsToUser", mock.Anything, senderUser.ID, recipientUser.ID, amountCoins).
		Return(nil)

	dTokenMng := &dummyTokenManager{}

	// Создаём обработчики, передавая TransactionService в соответствующий параметр.
	uh := NewUserHandlers(context.Background(), nil, dTokenMng, nil, mTxSvc, nil)

	meddlers := middlewares.NewMiddlewares(dTokenMng)
	authorized := router.Group("/", meddlers.JWTMiddleware())
	{
		authorized.POST("/sendCoin", uh.SendCoinsHandler)
	}

	// Готовим тело запроса в формате JSON.
	body := fmt.Sprintf(`{"toUser": "%s", "Amount": %d}`, recipientUser.Username, amountCoins)
	req, err := http.NewRequest(http.MethodPost, "/sendCoin", strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+validToken)

	// Выполняем запрос.
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем, что статус ответа 200 OK.
	require.Equal(t, http.StatusOK, w.Code)
}

// TestUserHandlers_BuyItemHandler проверяет сценарий успешной покупки мерча.
func TestUserHandlers_BuyItemHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	user := models.User{
		ID:       1,
		Username: "testUser",
		Coins:    150,
	}

	item := &models.Item{
		Slug:  "merch123",
		Price: 100,
	}

	mBuyItemSvc := mocks.NewBuyItemService(t)
	mBuyItemSvc.
		On("GetItem", mock.Anything, item.Slug).
		Return(item, nil)
	// При запросе количества монет для покупателя (user_id = 1) возвращаем 150.
	mBuyItemSvc.
		On("GetBuyerCoins", mock.Anything, user.ID).
		Return(user.Coins, nil)
	// При вызове BuyItem возвращаем nil (успех).
	mBuyItemSvc.
		On("BuyItem", mock.Anything, user.ID, item).
		Return(nil)

	dTokenMng := &dummyTokenManager{}

	// Создаём обработчики с необходимыми зависимостями.
	// Для неиспользуемых сервисов можно передавать nil.
	uh := NewUserHandlers(context.Background(), nil, dTokenMng, nil, nil, mBuyItemSvc)

	// Настраиваем группу маршрутов с JWT-мидлваром.
	meddlers := middlewares.NewMiddlewares(dTokenMng)
	authorized := router.Group("/", meddlers.JWTMiddleware())
	{
		authorized.GET("/buy/:item", uh.BuyItemHandler)
	}

	// Формируем HTTP‑запрос на покупку мерча с корректным токеном.
	req, err := http.NewRequest(http.MethodGet, "/buy/merch123", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+validToken)
	req.Header.Set("Accept", "application/json")

	// Выполняем запрос.
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем, что статус ответа 200 OK.
	require.Equal(t, http.StatusOK, w.Code)
}
