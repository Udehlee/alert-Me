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
	CurrentPrice string    `json:"current_price"`
	CreatedAt    time.Time `json:"created_at"`
}

type Product struct {
	ItemID string `json:"itemId"`
	Title  string `json:"title"`
	Price  struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"price"`
}

type SearchResponse struct {
	Products []Product `json:"products"`
}
