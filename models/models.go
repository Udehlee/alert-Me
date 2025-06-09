package models

import "time"

type User struct {
	ID        string    `json:"user_id"`
	Email     string    `json:"email"`
	Password  string    `json:"pass_word"`
	CreatedAt time.Time `json:"created_at"`
}

type SelectedProduct struct {
	ID          int       `json:"id"`
	Name        string    `json:"name_"`
	Price       string    `json:"price"`
	URL         string    `json:"product_url"`
	CreatedAt   time.Time `json:"created_at"`
	LastChecked time.Time `json:"last_checked"`
}

type UrlRequest struct {
	URL string `json:"product_url"`
}
