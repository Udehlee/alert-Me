package models

import "time"

type User struct {
	ID        string    `json:"user_id"`
	Email     string    `json:"email"`
	Password  string    `json:"pass_word"`
	CreatedAt time.Time `json:"created_at"`
}

type SelectedProduct struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name"`
	CurrentPrice float64   `json:"current_price"`
	CreatedAt    time.Time `json:"created_at"`
}
