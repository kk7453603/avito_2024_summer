package service

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kk7453603/avito_2024_summer/internal/models"
	"github.com/kk7453603/avito_2024_summer/pkg/middleware"
	"github.com/kk7453603/avito_2024_summer/pkg/utils"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error

	GetUser(ctx context.Context, username string, password string) (*models.User, error)

	GetUserByUsername(ctx context.Context, username string) (*models.User, error)

	GetUserById(ctx context.Context, userId string) (*models.User, error)

	GetUserInfo(ctx context.Context, userID string) (*models.InfoResponse, error)

	TransferCoins(ctx context.Context, fromUserID, toUserID, amount int) error

	BuyMerch(ctx context.Context, userID int, merchName string, quantity int) error
}

type Service struct {
	repo Repository
	elog echo.Logger
}

func New(repo Repository, logger echo.Logger) *Service {
	return &Service{repo: repo, elog: logger}
}

func (s *Service) Authentificate(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error) {

	user, err := s.repo.GetUser(ctx, req.Username, req.Password)
	if err == pgx.ErrNoRows {

		passhash, err := utils.HashPassword(req.Password)
		if err != nil {
			return &models.AuthResponse{}, err
		}

		user = &models.User{
			Username:    req.Username,
			Password:    passhash,
			CoinBalance: 1000,
			CreatedAt:   time.Now(),
		}
		if err := s.repo.CreateUser(ctx, user); err != nil {
			s.elog.Errorf("Ошибка создания пользователя %s: %v", req.Username, err)
			return nil, err
		}
	} else if err != nil {
		return &models.AuthResponse{}, err
	} else {

		if user.Username != req.Username {
			return &models.AuthResponse{}, errors.New("wrong username")
		}

		if ok := utils.CheckPasswordHash(req.Password, user.Password); !ok {
			return &models.AuthResponse{}, errors.New("wrong password")
		}

	}
	token, err := middleware.GenerateToken(
		middleware.JWTMiddlewareConfig{SecretKey: []byte(os.Getenv("JWT_SECRET"))},
		strconv.Itoa(user.ID),
	)
	if err != nil {
		s.elog.Errorf("Ошибка генерации JWT для пользователя %s: %v", req.Username, err)
		return nil, err
	}

	return &models.AuthResponse{Token: token}, nil
}

// GetUser возвращает данные пользователя по его имени.
func (s *Service) GetUser(ctx context.Context, username string) (*models.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		s.elog.Errorf("Ошибка получения пользователя %s: %v", username, err)
		return nil, err
	}
	return user, nil
}

func (s *Service) GetInfo(ctx context.Context, userID string) (*models.InfoResponse, error) {
	info, err := s.repo.GetUserInfo(ctx, userID)
	if err != nil {
		s.elog.Errorf("Ошибка получения информации для пользователя по id %s: %v", userID, err)
		return nil, err
	}
	return info, nil
}

// GetUserInfo возвращает полную информацию о пользователе (баланс, инвентарь, историю переводов).
/*
func (s *Service) GetUserInfo(ctx context.Context, username string) (*models.InfoResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		s.elog.Errorf("Ошибка получения пользователя %s: %v", username, err)
		return nil, err
	}
	info, err := s.repo.GetUserInfo(ctx, (user.ID)
	if err != nil {
		s.elog.Errorf("Ошибка получения информации для пользователя %s: %v", username, err)
		return nil, err
	}
	return info, nil
}
*/
// TransferCoins выполняет перевод монет от одного пользователя к другому.
func (s *Service) TransferCoins(ctx context.Context, fromUsername, toUsername string, amount int) error {
	if amount <= 0 {
		return errors.New("количество монет для перевода должно быть положительным")
	}

	sender, err := s.repo.GetUserByUsername(ctx, fromUsername)
	if err != nil {
		s.elog.Errorf("Ошибка получения отправителя %s: %v", fromUsername, err)
		return err
	}
	receiver, err := s.repo.GetUserByUsername(ctx, toUsername)
	if err != nil {
		s.elog.Errorf("Ошибка получения получателя %s: %v", toUsername, err)
		return err
	}

	if err = s.repo.TransferCoins(ctx, sender.ID, receiver.ID, amount); err != nil {
		s.elog.Errorf("Ошибка перевода %d монет от %s к %s: %v", amount, fromUsername, toUsername, err)
		return err
	}

	s.elog.Infof("Перевод %d монет от %s к %s выполнен успешно", amount, fromUsername, toUsername)
	return nil
}

// BuyMerch позволяет пользователю приобрести мерч, списывая монеты, регистрируя покупку и обновляя инвентарь.
func (s *Service) BuyMerch(ctx context.Context, userId, merchName string, quantity int) error {
	if quantity <= 0 {
		return errors.New("количество мерча должно быть положительным")
	}

	user, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		s.elog.Errorf("Ошибка получения пользователя c id %s: %v", userId, err)
		return err
	}

	if err = s.repo.BuyMerch(ctx, user.ID, merchName, quantity); err != nil {
		s.elog.Errorf("Ошибка покупки мерча %s (кол-во %d) для пользователя c id %s: %v", merchName, quantity, user.ID, err)
		return err
	}

	s.elog.Infof("Пользователь c id %s успешно приобрел %d единиц мерча %s", user.ID, quantity, merchName)
	return nil
}
