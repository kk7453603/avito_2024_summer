package models

// InfoResponse представляет ответ для запроса /api/info.
type InfoResponse struct {
	Coins       int         `json:"coins"`       // Количество доступных монет.
	Inventory   []Item      `json:"inventory"`   // Список купленных товаров.
	CoinHistory CoinHistory `json:"coinHistory"` // История транзакций по монетам.
}

// Item описывает отдельный товар в инвентаре.
type Item struct {
	Type     string `json:"type"`     // Тип предмета.
	Quantity int    `json:"quantity"` // Количество предметов.
}

// CoinHistory содержит две группы транзакций: полученные и отправленные.
type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"` // Транзакции по полученным монетам.
	Sent     []SentTransaction     `json:"sent"`     // Транзакции по отправленным монетам.
}

// ReceivedTransaction описывает транзакцию, когда монеты получены.
type ReceivedTransaction struct {
	FromUser string `json:"fromUser"` // Имя пользователя, отправившего монеты.
	Amount   int    `json:"amount"`   // Количество полученных монет.
}

// SentTransaction описывает транзакцию, когда монеты отправлены.
type SentTransaction struct {
	ToUser string `json:"toUser"` // Имя пользователя, которому отправлены монеты.
	Amount int    `json:"amount"` // Количество отправленных монет.
}

// ErrorResponse представляет стандартный ответ при ошибке.
type ErrorResponse struct {
	Errors string `json:"errors"` // Сообщение об ошибке.
}

// AuthRequest представляет запрос аутентификации на /api/auth.
type AuthRequest struct {
	Username string `json:"username"` // Имя пользователя для аутентификации.
	Password string `json:"password"` // Пароль для аутентификации.
}

// AuthResponse возвращается при успешной аутентификации и содержит JWT-токен.
type AuthResponse struct {
	Token string `json:"token"` // JWT-токен для доступа к защищённым ресурсам.
}

// SendCoinRequest представляет запрос на перевод монет (/api/sendCoin).
type SendCoinRequest struct {
	ToUser string `json:"toUser"` // Имя пользователя, которому нужно отправить монеты.
	Amount int    `json:"amount"` // Количество монет для перевода.
}
