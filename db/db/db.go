package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Udehlee/alert-Me/models"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DbName   string
}

type Conn struct {
	DB *sql.DB
}

func LoadDBEnv() (DBConfig, error) {
	cfg := DBConfig{}
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return cfg, err
	}

	cfg = DBConfig{
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     port,
		DbName:   os.Getenv("POSTGRES_DB"),
	}

	return cfg, nil
}

func InitConnectDB() (Conn, error) {
	db := Conn{}
	config, err := LoadDBEnv()
	if err != nil {
		return db, err
	}

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		config.Username, config.Password, config.Host, config.Port, config.DbName)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, fmt.Errorf("error creating database connection: %w", err)
	}

	db.DB = conn

	log.Println("Database connected successfully")
	return db, nil
}

// Save saves user details to db
func (c Conn) Save(user models.User) error {
	query := `INSERT INTO users (email, pass_word)
	          VALUES ($1, $2, $3, $4)
	          RETURNING user_id, email,created_at
             `
	row := c.DB.QueryRow(query, user.Email, user.Password, &user.CreatedAt)
	if err := row.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
		return fmt.Errorf("error scanning row: %w", err)
	}

	return nil
}

// SaveProduct saves products selected by user to watch
func (c Conn) SaveProduct(product models.SelectedProduct) error {
	query := `INSERT INTO SelectedProduct(name_,price,product_url)
	          VALUES ($1, $2, $3)
	          RETURNING id, title, price,product_url, created_at, last_checked
             `

	row := c.DB.QueryRow(query, product.ID, product.Name, product.Price, product.URL)
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.URL, &product.CreatedAt, &product.LastChecked)
	if err != nil {
		return fmt.Errorf("error scanning row: %w", err)
	}

	return nil
}

// PendingProduct retrives all products that has the status = pending
func (c Conn) PendingProduct() ([]models.SelectedProduct, error) {
	query := "SELECT id, name, price,product_url FROM selectedProduct Where status = 'pending'"

	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.SelectedProduct
	for rows.Next() {
		var p models.SelectedProduct
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.URL); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
