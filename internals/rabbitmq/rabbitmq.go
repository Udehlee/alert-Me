package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewRabbitMQ(conn *amqp.Connection, ch *amqp.Channel) *RabbitMQ {
	return &RabbitMQ{
		Conn: conn,
		Ch:   ch,
	}
}

// ConnectRabbitMQ initializes the rabbitmq connection and channel
func ConnectRabbitMQ() (*RabbitMQ, error) {
	rb := RabbitMQ{}

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	rb.Conn = conn
	rb.Ch = ch

	fmt.Println("RabbitMQ connected sucessfully")
	return &rb, nil
}

// PublishToQueue sends queueName and product_url to Queue
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
		return fmt.Errorf("failed to declare message: %v", err)

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

	log.Printf("Message published to queue %q: %s", queueName, body)
	return nil
}

// Consumer listens and  processes incoming messages from queue
func (r *RabbitMQ) Consumer(queueName string, handleMsg func([]byte) error) error {
	_, err := r.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	msgs, err := r.Ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			err := handleMsg(msg.Body)
			if err != nil {
				log.Printf("Failed to process message: %v", err)
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Printf(" Failed to ack message: %v", err)
			}
		}
	}()

	return nil
}
