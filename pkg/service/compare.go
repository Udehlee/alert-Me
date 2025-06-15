package service

import (
	"fmt"
	"log"

	"github.com/Udehlee/alert-Me/models"
	"github.com/Udehlee/alert-Me/pkg/utils"
)

// PriceCheck looks at product prices and sends an alert if the price drops
func (s *Service) PriceCheck(queueName string) error {
	HandlePriceCheck := func(body []byte) error {
		var DBProduct models.Product
		if err := utils.UnmarshalJSON(body, &DBProduct); err != nil {
			return err
		}

		recheckedProduct, err := utils.ExtractProduct(DBProduct.URL)
		if err != nil {
			return fmt.Errorf("failed to extract product for rechecking: %s %w", DBProduct.URL, err)
		}

		if recheckedProduct.Price < DBProduct.Price {
			//for now, we would publish it indicating there was a price drop
			err := s.Rabbit.PublishToQueue("price_drop_alert", []byte("alertMsg"))
			if err != nil {
				log.Printf("Failed to publish alert: %v", err)
			}
		} else {
			log.Printf("No price change for %q Current: %.2f", recheckedProduct.Name, recheckedProduct.Price)
		}

		return nil
	}

	if err := s.Rabbit.Consumer(queueName, HandlePriceCheck); err != nil {
		return fmt.Errorf("failed to start price checker: %w", err)
	}

	return nil
}
