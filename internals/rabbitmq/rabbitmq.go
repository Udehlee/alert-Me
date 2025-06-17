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

// Close safely closes the RabbitMQ channel and connection
func (r *RabbitMQ) Close() {
	if r.Ch != nil {
		r.Ch.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
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
		log.Printf("Failed to publish to queue %q: %v", queueName, err)
		return fmt.Errorf("failed to publish message: %v", err)
	}

	log.Printf("Message published to queue %q: %s", queueName, body)
	return nil
}

// Consumer listens and  processes incoming messages from queue
func (r *RabbitMQ) Consumer(queueName string, msgHandler func([]byte) error) error {
	_, err := r.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := r.Ch.Qos(5, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
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
			log.Printf("Received message: %s", msg.Body)

			err := msgHandler(msg.Body)
			if err != nil {
				log.Printf("Failed to process message: %v", err)
				msg.Nack(false, false)
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Printf(" Failed to ack message: %v", err)
			}
		}
	}()

	select {}
}
