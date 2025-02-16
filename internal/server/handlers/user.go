// Package handlers provides HTTP handlers for user-related operations, including authentication,
// retrieving user information, transferring coins, and purchasing items.
// Each handler interacts with the appropriate service layer to perform its tasks and returns HTTP responses
// in JSON format. The package uses the gin framework for routing and request handling,
// ensuring a clean and modular structure for handling user requests.
package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kk7453603/avito_2024_summer/internal/models"
)

// ErrInDB is a common error message for database-related issues.
var ErrInDB = errors.New("something happened to the database")

// UserHandlers provides HTTP handlers for user-related operations.
type UserHandlers struct {
	ctx       context.Context    // Context for managing request-scoped values and cancellation.
	authSrv   AuthService        // Service for authentication-related operations.
	tknMng    TokenManager       // Manager for JWT token operations.
	usrInfSrv UserInfoService    // Service for retrieving user information.
	txSrv     TransactionService // Service for handling coin transactions.
	buyItmSrv BuyItemService     // Service for handling item purchases.
}

// NewUserHandlers creates a new instance of UserHandlers with the provided dependencies.
func NewUserHandlers(ctx context.Context,
	authSrv AuthService, tknMng TokenManager, usrInfSrv UserInfoService,
	txSrv TransactionService, buyItmSrv BuyItemService) *UserHandlers {
	return &UserHandlers{
		ctx:       ctx,
		authSrv:   authSrv,
		tknMng:    tknMng,
		usrInfSrv: usrInfSrv,
		txSrv:     txSrv,
		buyItmSrv: buyItmSrv,
	}
}

// AuthHandler handles user authentication and token generation.
func (uh *UserHandlers) AuthHandler(c *gin.Context) {
	// switch c.GetHeader("Accept") {
	// case "application/json":
	// 	// continue
	// default:
	// 	c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "the ‘accept’ header is not application/json"})
	// 	return
	// }

	var login models.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok, err := uh.authSrv.GetOrRegUser(uh.ctx, login.Username, login.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	} else if ok {
		if !uh.authSrv.ComparePassword(user.Password, login.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
			return
		}
	}

	tokenString, err := uh.tknMng.NewToken(strconv.Itoa(user.ID), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failure"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// InfoHandler retrieves and returns user information, including coins, inventory, and coin history.
func (uh *UserHandlers) InfoHandler(c *gin.Context) {
	// switch c.GetHeader("Accept") {
	// case "application/json":
	// 	// continue
	// default:
	// 	c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "the ‘accept’ header is not application/json"})
	// 	return
	// }
	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "context parsing failure"})
		return
	}

	coins, err := uh.usrInfSrv.GetCoins(uh.ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	}

	inventory, err := uh.usrInfSrv.GetInventory(uh.ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	}

	coinHistory, err := uh.usrInfSrv.GetCoinHistory(uh.ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	}

	type Response struct {
		Coins       int                 `json:"coins"`
		Inventory   *[]models.Merch     `json:"inventory"`
		CoinHistory *models.CoinHistory `json:"coinHistory"`
	}

	c.JSON(http.StatusOK, Response{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	})
}

// SendCoinsHandler handles the transfer of coins from one user to another.
func (uh *UserHandlers) SendCoinsHandler(c *gin.Context) {
	var send models.Sending
	if err := c.ShouldBindJSON(&send); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// if send.User == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "`toUser` must not be empty"})
	// 	return
	// } else if send.Amount <= 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "`amount` must be positive"})
	// 	return
	// }

	recipientID, err := uh.txSrv.GetIDRecipient(uh.ctx, send.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	} else if recipientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "`toUser` is not found"})
		return
	}

	senderIDStr, _ := c.Get("user_id")
	senderID, err := strconv.Atoi(senderIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "context parsing failure"})
		return
	}
	if senderCoins, err := uh.txSrv.GetSenderCoins(uh.ctx, senderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	} else if senderCoins < send.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you don't have enough coins"})
		return
	}

	if err = uh.txSrv.SendCoinsToUser(uh.ctx, senderID, recipientID, send.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// BuyItemHandler handles the purchase of an item by a user.
func (uh *UserHandlers) BuyItemHandler(c *gin.Context) {
	itemSlug := c.Param("item")
	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "context parsing failure"})
		return
	}

	item, err := uh.buyItmSrv.GetItem(uh.ctx, itemSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	} else if item == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item not found"})
		return
	}

	buyerCoins, err := uh.buyItmSrv.GetBuyerCoins(uh.ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	} else if buyerCoins < item.Price {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you don't have enough coins"})
		return
	}

	if err = uh.buyItmSrv.BuyItem(uh.ctx, userID, item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInDB.Error()})
		return
	}

	c.Status(http.StatusOK)
}
