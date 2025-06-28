package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Udehlee/alert-Me/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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

	if err := conn.Ping(); err != nil {
		return db, fmt.Errorf("connection is not active: %w", err)
	}

	db.DB = conn

	if err := runMigrations(conn); err != nil {
		return db, fmt.Errorf("migration unsuccessful: %w", err)
	}

	log.Println("database connected successfully")
	return db, nil
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./internals/db/migrations",
		"postgres", driver)

	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}

// SaveProduct saves products selected by user to watch
func (c Conn) SaveProduct(product models.Product) error {
	query := `INSERT INTO Products(name_,price,product_url) VALUES ($1, $2, $3)
	          RETURNING id, name_, price,product_url, created_at, last_checked
             `

	row := c.DB.QueryRow(query, product.Name, product.Price, product.URL)
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.URL, &product.CreatedAt, &product.LastChecked)
	if err != nil {
		return fmt.Errorf("error scanning row: %w", err)
	}

	return nil
}

// PendingProduct retrives all products that has the status = pending
func (c Conn) PendingProduct() ([]models.Product, error) {
	query := `
		SELECT id, name_, price, product_url, last_checked FROM Products
		WHERE status_ = 'pending'
	`
	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.URL, &p.LastChecked); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating row: %w", err)
	}

	return products, nil
}
