package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Udehlee/alert-Me/models"
	"github.com/Udehlee/alert-Me/pkg/rabbitmq"
	"github.com/Udehlee/alert-Me/pkg/utils"
)

type Service struct {
	Rabbit *rabbitmq.RabbitMQ
}

func NewService(rabbit *rabbitmq.RabbitMQ) *Service {
	return &Service{
		Rabbit: rabbit,
	}
}

// StartConsumer starts a consumer that listens to a queue
// and saves scraped product name and price from submitted URLs database.
func (s *Service) StartConsumer() error {
	return s.Rabbit.Consumer("product_url_queue", func(body []byte) error {
		var payload models.UrlRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			return err
		}

		product, err := utils.ExtractProduct(payload.URL)
		if err != nil {
			return err
		}

		return s.Rabbit.DB.SaveProduct(product)
	})

}

// PeriodicCheck periodically fetches pending products from database
// and republishes them to a queue for processing.
func (s *Service) PeriodicCheck(ctx context.Context, queueName string) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			products, err := s.Rabbit.DB.PendingProduct()
			if err != nil {
				log.Println("error fetching pending product:", err)
				continue
			}

			for _, p := range products {
				msg := models.SelectedProduct{
					ID:    p.ID,
					URL:   p.URL,
					Price: p.Price,
				}
				body, _ := json.Marshal(msg)
				_ = s.Rabbit.PublishToQueue(queueName, body)
			}

		case <-ctx.Done():
			log.Println("Stopped periodic check for pending product")
			return
		}
	}
}
