package repository

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kk7453603/avito_2024_summer/internal/models"
	"github.com/labstack/echo/v4"
)

var ErrUserNotFound = errors.New("user not found")
var ErrFailedToStopTask = errors.New("failed to stop task")

type SqlHandler struct {
	db   *pgxpool.Pool
	elog echo.Logger
	dsn  string
}

func New(ctx context.Context, e echo.Logger) *SqlHandler {
	var pool *pgxpool.Pool
	var err error

	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	for attempt := 1; attempt <= 5; attempt++ {
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			e.Printf("Попытка подключения к БД %d/%d не удалась: %v", attempt, 5, err)
		} else {
			ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
			err = pool.Ping(ctxPing)
			cancel()
			if err != nil {
				e.Printf("Попытка пинга БД %d/%d не удалась: %v", attempt, 5, err)
			} else {
				e.Printf("Успешное подключение к БД на попытке %d", attempt)
				return &SqlHandler{db: pool, dsn: dsn, elog: e}
			}
		}
		time.Sleep(time.Second)
	}
	panic("DB connection failed!")
}

func (h *SqlHandler) Migrate() {
	m, err := migrate.New(os.Getenv("DB_MIGRATIONS_PATH"), h.dsn+"?sslmode=disable")
	h.elog.Debugf("sourceURL: %s , DSN: %s", os.Getenv("DB_MIGRATIONS_PATH"), h.dsn)
	if err != nil {
		h.elog.Errorf("migration error: %v", err)
	}
	if err := m.Up(); err != nil && errors.Is(err, errors.New("migration error: no change")) {
		h.elog.Errorf("migration error: %v", err)
	}
	h.elog.Info("Миграции выполнены")
}

// CreateUser создаёт нового пользователя с начальным балансом.
func (r *SqlHandler) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, coin_balance, password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	r.elog.Printf("Создаём пользователя: %s", user.Username)
	err := r.db.QueryRow(ctx, query, user.Username, user.CoinBalance, user.Password, time.Now()).Scan(&user.ID)
	if err != nil {
		r.elog.Printf("Ошибка при создании пользователя %s: %v", user.Username, err)
	}
	return err
}

// GetUserByUsername возвращает пользователя по имени.
func (r *SqlHandler) GetUser(ctx context.Context, username string, password string) (*models.User, error) {
	query := `SELECT id, username,password, coin_balance, created_at FROM users WHERE username=$1 AND password=$2`
	r.elog.Printf("Запрос пользователя по имени: %s", username)
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, username, password).Scan(&user.ID, &user.Username, &user.Password, &user.CoinBalance, &user.CreatedAt)
	if err != nil {
		r.elog.Printf("Пользователь %s не найден: %v", username, err)
		return nil, err
	}
	return user, nil
}

func (r *SqlHandler) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, coin_balance, created_at FROM users WHERE username=$1`
	r.elog.Printf("Запрос пользователя по имени: %s", username)
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.CoinBalance, &user.CreatedAt)
	if err != nil {
		r.elog.Printf("Пользователь %s не найден: %v", username, err)
		return nil, err
	}
	return user, nil
}

func (r *SqlHandler) GetUserById(ctx context.Context, userId string) (*models.User, error) {
	query := `SELECT id, username, coin_balance, created_at FROM users WHERE id=$1`
	r.elog.Printf("Запрос пользователя по id: %s", userId)
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, userId).Scan(&user.ID, &user.Username, &user.CoinBalance, &user.CreatedAt)
	if err != nil {
		r.elog.Printf("Пользователь c id %s не найден: %v", userId, err)
		return nil, err
	}
	return user, nil
}

// GetUserInfo собирает информацию о пользователе: баланс, инвентарь и историю переводов.
func (r *SqlHandler) GetUserInfo(ctx context.Context, userID string) (*models.InfoResponse, error) {
	info := &models.InfoResponse{}

	r.elog.Printf("Получение информации для пользователя с id=%d", userID)
	// Получаем баланс.
	err := r.db.QueryRow(ctx, "SELECT coin_balance FROM users WHERE id=$1", userID).Scan(&info.Coins)
	if err != nil {
		r.elog.Printf("Ошибка получения баланса для пользователя id=%d: %v", userID, err)
		return nil, err
	}

	// Получаем инвентарь.
	invQuery := `SELECT merch_name, quantity FROM inventory WHERE user_id = $1`
	rows, err := r.db.Query(ctx, invQuery, userID)
	if err != nil {
		r.elog.Printf("Ошибка получения инвентаря для пользователя id=%d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			r.elog.Printf("Ошибка сканирования инвентаря для пользователя id=%d: %v", userID, err)
			return nil, err
		}
		items = append(items, item)
	}
	info.Inventory = items

	// Получаем историю полученных монет.
	recvQuery := `
		SELECT u.username, ct.amount 
		FROM coin_transfers ct
		JOIN users u ON ct.from_user_id = u.id
		WHERE ct.to_user_id = $1
		ORDER BY ct.created_at DESC
	`
	rows, err = r.db.Query(ctx, recvQuery, userID)
	if err != nil {
		r.elog.Printf("Ошибка получения истории полученных монет для пользователя id=%d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var received []models.ReceivedTransaction
	for rows.Next() {
		var rt models.ReceivedTransaction
		if err := rows.Scan(&rt.FromUser, &rt.Amount); err != nil {
			r.elog.Printf("Ошибка сканирования истории полученных монет для пользователя id=%d: %v", userID, err)
			return nil, err
		}
		received = append(received, rt)
	}

	// Получаем историю отправленных монет.
	sentQuery := `
		SELECT u.username, ct.amount 
		FROM coin_transfers ct
		JOIN users u ON ct.to_user_id = u.id
		WHERE ct.from_user_id = $1
		ORDER BY ct.created_at DESC
	`
	rows, err = r.db.Query(ctx, sentQuery, userID)
	if err != nil {
		r.elog.Printf("Ошибка получения истории отправленных монет для пользователя id=%d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var sent []models.SentTransaction
	for rows.Next() {
		var st models.SentTransaction
		if err := rows.Scan(&st.ToUser, &st.Amount); err != nil {
			r.elog.Printf("Ошибка сканирования истории отправленных монет для пользователя id=%d: %v", userID, err)
			return nil, err
		}
		sent = append(sent, st)
	}

	info.CoinHistory = models.CoinHistory{
		Received: received,
		Sent:     sent,
	}

	return info, nil
}

// TransferCoins выполняет перевод монет между пользователями с использованием транзакции.
func (r *SqlHandler) TransferCoins(ctx context.Context, fromUserID, toUserID, amount int) error {
	r.elog.Printf("Перевод %d монет от пользователя id=%d к пользователю id=%d", amount, fromUserID, toUserID)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.elog.Printf("Ошибка начала транзакции для перевода: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			r.elog.Printf("Откат транзакции при ошибке перевода: %v", err)
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Проверяем баланс отправителя.
	var balance int
	err = tx.QueryRow(ctx, "SELECT coin_balance FROM users WHERE id=$1 FOR UPDATE", fromUserID).Scan(&balance)
	if err != nil {
		r.elog.Printf("Ошибка получения баланса отправителя id=%d: %v", fromUserID, err)
		return err
	}
	if balance < amount {
		err = errors.New("недостаточно монет для перевода")
		r.elog.Printf("Ошибка: %v", err)
		return err
	}

	// Списываем монеты у отправителя.
	_, err = tx.Exec(ctx, "UPDATE users SET coin_balance = coin_balance - $1 WHERE id=$2", amount, fromUserID)
	if err != nil {
		r.elog.Printf("Ошибка списания монет у отправителя id=%d: %v", fromUserID, err)
		return err
	}
	// Зачисляем монеты получателю.
	_, err = tx.Exec(ctx, "UPDATE users SET coin_balance = coin_balance + $1 WHERE id=$2", amount, toUserID)
	if err != nil {
		r.elog.Printf("Ошибка зачисления монет получателю id=%d: %v", toUserID, err)
		return err
	}
	// Регистрируем перевод.
	_, err = tx.Exec(ctx, "INSERT INTO coin_transfers (from_user_id, to_user_id, amount, created_at) VALUES ($1, $2, $3, $4)",
		fromUserID, toUserID, amount, time.Now())
	if err != nil {
		r.elog.Printf("Ошибка регистрации перевода монет: %v", err)
		return err
	}
	r.elog.Printf("Перевод успешно выполнен")
	return nil
}

// BuyMerch реализует покупку мерча: списываем монеты, регистрируем покупку и обновляем инвентарь.
func (r *SqlHandler) BuyMerch(ctx context.Context, userID int, merchName string, quantity int) error {
	r.elog.Printf("Покупка мерча: пользователь id=%d покупает %d единиц товара %s", userID, quantity, merchName)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.elog.Printf("Ошибка начала транзакции для покупки мерча: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			r.elog.Printf("Откат транзакции покупки мерча: %v", err)
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Получаем стоимость товара.
	var price int
	err = tx.QueryRow(ctx, "SELECT price FROM merch_items WHERE name=$1", merchName).Scan(&price)
	if err != nil {
		r.elog.Printf("Ошибка получения цены товара %s: %v", merchName, err)
		return err
	}
	totalPrice := price * quantity

	// Проверяем баланс пользователя.
	var balance int
	err = tx.QueryRow(ctx, "SELECT coin_balance FROM users WHERE id=$1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		r.elog.Printf("Ошибка получения баланса пользователя id=%d: %v", userID, err)
		return err
	}
	if balance < totalPrice {
		err = errors.New("недостаточно монет для покупки мерча")
		r.elog.Printf("Ошибка: %v", err)
		return err
	}

	// Списываем монеты.
	_, err = tx.Exec(ctx, "UPDATE users SET coin_balance = coin_balance - $1 WHERE id=$2", totalPrice, userID)
	if err != nil {
		r.elog.Printf("Ошибка списания монет при покупке мерча для пользователя id=%d: %v", userID, err)
		return err
	}

	// Регистрируем покупку.
	_, err = tx.Exec(ctx, `
		INSERT INTO merch_purchases (user_id, merch_name, quantity, total_price, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, merchName, quantity, totalPrice, time.Now())
	if err != nil {
		r.elog.Printf("Ошибка регистрации покупки мерча: %v", err)
		return err
	}

	// Обновляем инвентарь.
	_, err = tx.Exec(ctx, `
		INSERT INTO inventory (user_id, merch_name, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, merch_name) DO UPDATE
		SET quantity = inventory.quantity + $3
	`, userID, merchName, quantity)
	if err != nil {
		r.elog.Printf("Ошибка обновления инвентаря для пользователя id=%d: %v", userID, err)
		return err
	}

	r.elog.Printf("Покупка мерча успешно выполнена")
	return nil
}
