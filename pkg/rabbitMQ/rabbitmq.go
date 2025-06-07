package rabbitMQ

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Udehlee/alert-Me/db/db"
	"github.com/Udehlee/alert-Me/models"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
	db   db.Conn
}

func NewRabbitMQ(conn *amqp.Connection, ch *amqp.Channel, db db.Conn) *RabbitMQ {
	return &RabbitMQ{
		Conn: conn,
		Ch:   ch,
		db:   db,
	}
}

// ConnectRabbitMQ initializes the connection and channel
func ConnectRabbitMQ(db db.Conn) (RabbitMQ, error) {
	rb := RabbitMQ{db: db}

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return rb, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return rb, fmt.Errorf("failed to create channel: %w", err)
	}

	rb.Conn = conn
	rb.Ch = ch

	fmt.Println("RabbitMQ connected sucessfully")
	return rb, nil
}

// PublishToQueue sends queueName and product_url to rabbitQueue
func (r *RabbitMQ) PublishToQueue(queueName string, body []byte) error {
	_, err := r.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)

	}

	err = r.Ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

// Consumer listens and  processes incoming messages from queue
// It scrapes product data from URLs and saves them to the database
func (r *RabbitMQ) Consumer(queueName string, scraper func(string) models.SelectedProduct) {
	msgs, err := r.Ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var payload models.UrlRequest

			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			product := scraper(payload.URL)

			if err := r.db.SaveProduct(product); err != nil {
				log.Printf("Failed to save productto db: %v", err)
				continue
			}

			log.Printf(" Scraped Product: %+v\n", product)
		}
	}()

	log.Println("Consumer running...")
	select {}
}
