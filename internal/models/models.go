package models

import "time"

// User представляет запись из таблицы users.
type User struct {
	ID          int       `json:"id"`       // Идентификатор пользователя (SERIAL).
	Username    string    `json:"username"` // Уникальное имя пользователя.
	Password    string    `json:"password"`
	CoinBalance int       `json:"coin_balance"` // Баланс монет (неотрицательный).
	CreatedAt   time.Time `json:"created_at"`   // Дата и время создания записи.
}

// MerchItem представляет запись из таблицы merch_items (справочник мерча).
type MerchItem struct {
	Name  string `json:"name"`  // Название товара (PRIMARY KEY).
	Price int    `json:"price"` // Цена товара.
}

// InventoryItem представляет запись из таблицы inventory – купленные товары пользователя.
type InventoryItem struct {
	ID        int    `json:"id"`         // Идентификатор записи инвентаря.
	UserID    int    `json:"user_id"`    // Внешний ключ к пользователю.
	MerchName string `json:"merch_name"` // Название товара (из merch_items).
	Quantity  int    `json:"quantity"`   // Количество купленного товара.
}

// CoinTransfer представляет запись из таблицы coin_transfers – перевод монет между пользователями.
type CoinTransfer struct {
	ID         int       `json:"id"`           // Идентификатор перевода.
	FromUserID int       `json:"from_user_id"` // Идентификатор пользователя-отправителя.
	ToUserID   int       `json:"to_user_id"`   // Идентификатор пользователя-получателя.
	Amount     int       `json:"amount"`       // Количество переведённых монет.
	CreatedAt  time.Time `json:"created_at"`   // Дата и время перевода.
}

// MerchPurchase представляет запись из таблицы merch_purchases – покупка мерча пользователем.
type MerchPurchase struct {
	ID         int       `json:"id"`          // Идентификатор покупки.
	UserID     int       `json:"user_id"`     // Внешний ключ к пользователю.
	MerchName  string    `json:"merch_name"`  // Название товара (из merch_items).
	Quantity   int       `json:"quantity"`    // Количество купленного товара.
	TotalPrice int       `json:"total_price"` // Итоговая стоимость покупки (price * quantity).
	CreatedAt  time.Time `json:"created_at"`  // Дата и время покупки.
}
