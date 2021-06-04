package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

// OrderQueue manages the connection to RabbitMQ
type OrderQueue struct {
	Connection *amqp.Connection
	QueueName  string
}

// OrderMessage contains the body of the queue message
type OrderMessage struct {
	OrderKey string `json:"orderKey"`
	Pastry   string `json:"pastry"`
}

// Connect will connect to the RabbitMQ instance
func (q *OrderQueue) Connect(host string, username string, password string) error {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s", username, password, host)
	connection, err := amqp.Dial(connectionString)

	if err != nil {
		return err
	}

	q.Connection = connection
	return nil
}

// Disconnect closes the connection to RabbitMQ
func (q *OrderQueue) Disconnect() error {
	err := q.Connection.Close()

	if err != nil {
		return err
	}

	return nil
}

// DeclareQueue declares the queue that will be used for processing orders
func (q *OrderQueue) DeclareQueue() error {
	if q.Connection == nil {
		return errors.New("connection to RabbitMQ hasn't been established")
	}

	channel, err := q.Connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"orders",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	q.QueueName = queue.Name
	return nil
}

// PublishOrderMessage publishes an OrderMessage to the queue
func (q *OrderQueue) PublishOrderMessage(orderMessage *OrderMessage) error {
	if q.Connection == nil {
		return errors.New("connection to RabbitMQ hasn't been established")
	}

	channel, err := q.Connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	messageBody, err := json.Marshal(*orderMessage)

	if err != nil {
		return err
	}

	err = channel.Publish(
		"",
		q.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	)

	return err
}

// ConsumeOrderQueue creates a blocking connection to the order queue
func (q *OrderQueue) ConsumeOrderQueue(ctx context.Context, processFunction func(*OrderMessage) error, errorFunc func(error)) error {
	if q.Connection == nil {
		return errors.New("connection to RabbitMQ hasn't been established")
	}

	channel, err := q.Connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	orders, err := channel.Consume(
		q.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		for order := range orders {
			orderMessage := new(OrderMessage)
			err := json.Unmarshal(order.Body, orderMessage)

			if err != nil {
				errorFunc(err)
			}

			err = processFunction(orderMessage)

			if err != nil {
				errorFunc(err)
			}

			err = order.Ack(false)

			if err != nil {
				errorFunc(err)
			}
		}
	}()

	<-ctx.Done()

	return nil
}
