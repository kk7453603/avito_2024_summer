package delivery

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/kk7453603/avito_2024_summer/internal/models"
	custommiddleware "github.com/kk7453603/avito_2024_summer/pkg/middleware"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Authentificate(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error)

	GetInfo(ctx context.Context, userID string) (*models.InfoResponse, error)

	TransferCoins(ctx context.Context, fromUsername, toUsername string, amount int) error

	BuyMerch(ctx context.Context, userId, merchName string, quantity int) error
}

type Delivery struct {
	serv   Service
	logger echo.Logger
}

func New(serv Service, logger echo.Logger) *Delivery {
	return &Delivery{serv: serv, logger: logger}
}

func (d *Delivery) InitRoutes(g *echo.Group) {
	secret := os.Getenv("JWT_SECRET")
	jwtConfig := custommiddleware.JWTMiddlewareConfig{SecretKey: []byte(secret)}

	g.POST("/api/auth", d.Auth)
	g.GET("/api/info", d.Info, custommiddleware.JWTMiddleware(jwtConfig))
	g.POST("/api/sendCoin", d.SendCoin, custommiddleware.JWTMiddleware(jwtConfig))
	g.GET("/api/buy/:item", d.BuyMerch, custommiddleware.JWTMiddleware(jwtConfig))
}

// Auth godoc
//
//	@Summary		Аутентификация и получение JWT-токена
//	@Description	Аутентифицирует пользователя и возвращает JWT-токен. При первой аутентификации пользователь создается автоматически.
//	@Accept			application/json
//	@Produce		application/json
//	@Param			authRequest	body		models.AuthRequest	true	"Данные для аутентификации"
//	@Success		200			{object}	models.AuthResponse
//	@Failure		400			{object}	models.ErrorResponse
//	@Failure		500			{object}	models.ErrorResponse
//	@Router			/api/auth [post]
func (d *Delivery) Auth(c echo.Context) error {
	req := models.AuthRequest{}
	if err := c.Bind(req); err != nil {
		d.logger.Errorf("Ошибка привязки данных для аутентификации: %v", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Errors: "Неверный формат запроса"})
	}
	ctx := c.Request().Context()
	resp, err := d.serv.Authentificate(ctx, &req)
	if err != nil {
		d.logger.Errorf("Ошибка аутентификации: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

// Info godoc
//
//	@Summary		Получение информации о монетах, инвентаре и истории транзакций
//	@Description	Возвращает баланс монет, список купленных товаров и историю переводов (полученные и отправленные).
//	@Security		BearerAuth
//	@Produce		application/json
//	@Success		200	{object}	models.InfoResponse
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/api/info [get]
func (d *Delivery) Info(c echo.Context) error {

	userIDInterface := c.Get("userID")
	if userIDInterface == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Errors: "Пользователь не авторизован"})
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: "Ошибка идентификации пользователя"})
	}
	ctx := c.Request().Context()
	info, err := d.serv.GetInfo(ctx, userID)
	if err != nil {
		d.logger.Errorf("Ошибка получения информации для пользователя id=%d: %v", userID, err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: err.Error()})
	}
	return c.JSON(http.StatusOK, info)
}

// SendCoin godoc
//
//	@Summary		Отправить монеты другому пользователю
//	@Description	Переводит указанное количество монет от авторизованного пользователя к другому.
//	@Security		BearerAuth
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sendCoinRequest	body		models.SendCoinRequest	true	"Данные для перевода монет"
//	@Success		200				{string}	string					"Перевод выполнен успешно"
//	@Failure		400				{object}	models.ErrorResponse
//	@Failure		401				{object}	models.ErrorResponse
//	@Failure		500				{object}	models.ErrorResponse
//	@Router			/api/sendCoin [post]
func (d *Delivery) SendCoin(c echo.Context) error {
	userIDInterface := c.Get("userID")
	if userIDInterface == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Errors: "Пользователь не авторизован"})
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: "Ошибка идентификации пользователя"})
	}
	req := new(models.SendCoinRequest)
	if err := c.Bind(req); err != nil {
		d.logger.Errorf("Ошибка привязки данных для перевода монет: %v", err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Errors: "Неверный формат запроса"})
	}
	ctx := c.Request().Context()
	if err := d.serv.TransferCoins(ctx, userID, req.ToUser, req.Amount); err != nil {
		d.logger.Errorf("Ошибка перевода монет: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "Перевод выполнен успешно"})
}

// BuyMerch godoc
//
//	@Summary		Купить мерч за монетки
//	@Description	Покупает мерч для авторизованного пользователя. Имя товара передается в параметре пути.
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			item	path		string	true	"Наименование мерча"
//	@Param			count	query		int		1		"Количество мерча"
//	@Success		200		{string}	string	"Покупка выполнена успешно"
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/api/buy/{item} [get]
func (d *Delivery) BuyMerch(c echo.Context) error {
	userIDInterface := c.Get("userID")
	if userIDInterface == nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Errors: "Пользователь не авторизован"})
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: "Ошибка идентификации пользователя"})
	}
	merchName := c.Param("item")
	if merchName == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Errors: "Название мерча не указано"})
	}
	count := c.Param("count")
	if count == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Errors: "Количество покупаемого мерча не указано"})
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Errors: "Количество покупаемого мерча не указано или указано неверно"})
	}
	ctx := c.Request().Context()
	if err := d.serv.BuyMerch(ctx, userID, merchName, countInt); err != nil {
		d.logger.Errorf("Ошибка покупки мерча %s для пользователя id=%d: %v", merchName, userID, err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "Покупка выполнена успешно"})
}
