package models

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name_"`
	Price       float64   `json:"price"`
	URL         string    `json:"product_url"`
	Status      string    `json:"status_"`
	CreatedAt   time.Time `json:"created_at"`
	LastChecked time.Time `json:"last_checked"`
}

type UrlRequest struct {
	URL string `json:"product_url"`
}
