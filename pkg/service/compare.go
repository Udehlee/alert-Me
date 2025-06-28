package service

import (
	"fmt"
	"log"

	"github.com/Udehlee/alert-Me/models"
	"github.com/Udehlee/alert-Me/pkg/utils"
)

// comparePrice compares product prices and sends an alert if the price drops
func (s *Service) ComparePrice(queueName string) error {
	compare := func(body []byte) error {
		var DBProduct models.Product
		if err := utils.UnmarshalJSON(body, &DBProduct); err != nil {
			return err
		}

		if DBProduct.URL == "" {
			return fmt.Errorf("received product with empty URL")
		}

		recheckedProduct, err := utils.ExtractProduct(DBProduct.URL)
		if err != nil {
			return fmt.Errorf("failed to extract product for rechecking: %s %w", DBProduct.URL, err)
		}

		if recheckedProduct.Price < DBProduct.Price {
			log.Printf("Price drop alert \n  %q\n: Old_Price: %.2f\n, Current_Price: %.2f\n",
				recheckedProduct.Name, DBProduct.Price, recheckedProduct.Price)
			err := s.Rabbit.PublishToQueue("price_drop_alert", []byte("price drops"))
			if err != nil {
				log.Printf("Failed to publish alert: %v", err)
			}
		} else {
			log.Printf("currently, No price change for \n product name %q\n Current price: %.2f\n",
				recheckedProduct.Name, recheckedProduct.Price)
		}

		return nil
	}

	if err := s.Rabbit.Consumer(queueName, compare); err != nil {
		return fmt.Errorf("failed to start price checker: %w", err)
	}

	return nil
}
