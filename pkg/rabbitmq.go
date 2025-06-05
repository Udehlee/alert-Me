package event

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

// ConnectRabbitMQ initializes the connection and channel
func ConnectRabbitMQ() (RabbitMQ, error) {
	rb := RabbitMQ{}

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
