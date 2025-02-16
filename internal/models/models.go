package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id" binding:"required"`
	Username  string    `json:"username" db:"username" binding:"required"`
	Password  string    `json:"password" db:"password" binding:"required"`
	Coins     int       `json:"coins" db:"coins" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" binding:"required"`
}

type Login struct {
	Username string `json:"username" binding:"required,min=8,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

type Merch struct {
	Type     string `json:"type" db:"item_slug"`
	Quantity int    `json:"quantity" db:"quantity"`
}

type Receiving struct {
	User   string `json:"fromUser" db:"username"`
	Amount int    `json:"amount" db:"coins"`
}

type Sending struct {
	User   string `json:"toUser" db:"username" binding:"required,min=8,alphanum"`
	Amount int    `json:"amount" db:"coins" binding:"required,gte=1"`
}

type CoinHistory struct {
	Receiving *[]Receiving `json:"received"`
	Sending   *[]Sending   `json:"sent"`
}

type Item struct {
	Slug  string `json:"slug" db:"slug"`
	Title string `json:"title" db:"title"`
	Price int    `json:"price" db:"price"`
}
