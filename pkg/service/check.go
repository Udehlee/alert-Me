package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Udehlee/alert-Me/internals/db/db"
	"github.com/Udehlee/alert-Me/internals/rabbitmq"
	"github.com/Udehlee/alert-Me/models"
	"github.com/Udehlee/alert-Me/pkg/utils"
)

type Service struct {
	Db     db.Conn
	Rabbit *rabbitmq.RabbitMQ
}

func NewService(db db.Conn, rabbit *rabbitmq.RabbitMQ) *Service {
	return &Service{
		Db:     db,
		Rabbit: rabbit,
	}
}

// StartConsumer starts a consumer that listens to a queue
// and saves scraped product name and price from submitted URLs database.
func (s *Service) StartConsumer() error {
	log.Println("Starting consumer for queue")
	HandleMsg := func(body []byte) error {
		var payload models.UrlRequest
		if err := utils.UnmarshalJSON(body, &payload); err != nil {
			return fmt.Errorf("failed to parse incoming  product_url request: %w", err)
		}

		product, err := utils.ExtractProduct(payload.URL)
		if err != nil {
			return fmt.Errorf("failed to extract product_url %s: %w", payload.URL, err)
		}

		if err := s.Db.SaveProduct(product); err != nil {
			return fmt.Errorf("error saving product to database: %w", err)
		}

		return nil
	}

	if err := s.Rabbit.Consumer("product_url_queue", HandleMsg); err != nil {
		return fmt.Errorf("failed to start product URL consumer: %w", err)
	}

	log.Println("Consumer started successfully")
	return nil
}

// SendForRecheck periodically fetches pending products from database
// and republishes them to a queue for processing.
func (s *Service) SendForRecheck(ctx context.Context, queueName string) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			products, err := s.Db.PendingProduct()
			if err != nil {
				log.Printf(" Error fetching pending product: %v", err)
				continue
			}

			for _, p := range products {
				body, _ := json.Marshal(p.URL)
				err := s.Rabbit.PublishToQueue(queueName, body)
				if err != nil {
					log.Printf(" Failed to send product for recheck: %v", err)
				} else {
					log.Printf("Sent product %s for recheck", p.URL)
				}
			}
		case <-ctx.Done():
			log.Println("Stopping recheck process")
			return
		}
	}
}
