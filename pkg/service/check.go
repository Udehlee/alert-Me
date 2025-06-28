package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	db "github.com/Udehlee/alert-Me/internals/db/conn"
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
// and saves scraped product name and price from submitted URLs
func (s *Service) StartConsumer() error {
	msg := func(body []byte) error {
		var payload models.UrlRequest
		if err := utils.UnmarshalJSON(body, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal incoming  product_url request: %w", err)
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

	if err := s.Rabbit.Consumer("product_url_queue", msg); err != nil {
		return fmt.Errorf("failed to start product URL consumer: %w", err)
	}

	log.Println("Consumer started successfully")
	return nil
}

// SendForRecheck periodically fetches pending products from database
// and republishes them to a queue
func (s *Service) SendForRecheck(ctx context.Context, queueName string) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	log.Println("starting fetching pending product:")
	for {
		select {
		case <-ticker.C:
			products, err := s.Db.PendingProduct()
			if err != nil {
				log.Printf("Error fetching pending product: %v", err)
				continue
			}

			if len(products) == 0 {
				log.Println("currently, there is no product to watch for now")
				continue
			}

			for _, p := range products {
				body, err := json.Marshal(p)
				if err != nil {
					log.Printf("Failed to marshal URL %s: %v", p.URL, err)
					continue
				}

				err = s.Rabbit.PublishToQueue(queueName, body)
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
