package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/whatsfordinner/bakery/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	tracer "github.com/whatsfordinner/bakery/pkg/trace"
)

// OrderMessage contains the body of the queue message
type OrderMessage struct {
	TraceContext tracer.ContextCarrier `json:"traceContext"`
	OrderKey     string                `json:"orderKey"`
	Pastry       string                `json:"pastry"`
}

// OrderQueue manages the connection to RabbitMQ
type OrderQueue struct {
	Connection *amqp.Connection
	QueueName  string
	tracer     trace.Tracer
}

func NewOrderQueue(c *config.Config) (*OrderQueue, error) {
	orderQueue := new(OrderQueue)

	orderQueue.tracer = otel.Tracer("rabbitmq")
	if err := orderQueue.Connect(c.RabbitHost, c.RabbitUsername, c.RabbitPassword); err != nil {
		return nil, err
	}

	return orderQueue, nil
}

// Connect will connect to the RabbitMQ instance
func (q *OrderQueue) Connect(host string, username string, password string) error {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s", username, password, host)
	connection, err := amqp.Dial(connectionString)

	if err != nil {
		return err
	}

	q.Connection = connection

	err = q.DeclareQueue()

	if err != nil {
		return err
	}

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
func (q *OrderQueue) PublishOrderMessage(ctx context.Context, orderKey string, pastry string) error {
	_, span := q.tracer.Start(ctx, "publish-order")
	defer span.End()

	span.SetAttributes(
		attribute.String("bakery.order_key", orderKey),
		attribute.String("bakery.pastry", pastry),
		attribute.Bool("queue.success", false),
	)

	if q.Connection == nil {
		return errors.New("connection to RabbitMQ hasn't been established")
	}

	channel, err := q.Connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	orderMessage := new(OrderMessage)
	orderMessage.TraceContext = tracer.ContextCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, orderMessage.TraceContext)
	orderMessage.OrderKey = orderKey
	orderMessage.Pastry = pastry
	messageBody, err := json.Marshal(*orderMessage)

	if err != nil {
		return err
	}

	if err = channel.Publish(
		"",
		q.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	); err != nil {
		return err
	}

	span.SetAttributes(attribute.Bool("queue.success", true))
	return nil
}

// ConsumeOrderQueue creates a blocking connection to the order queue
func (q *OrderQueue) ConsumeOrderQueue(ctx context.Context, processFunction func(context.Context, *OrderMessage) error, errorFunc func(context.Context, error)) error {
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
			ctx, span := q.tracer.Start(ctx, "receive-order")
			span.SetAttributes(attribute.Bool("queue.success", false))

			orderMessage := new(OrderMessage)
			err := json.Unmarshal(order.Body, orderMessage)

			if err != nil {
				errorFunc(ctx, err)
			} else {
				ctx = otel.GetTextMapPropagator().Extract(ctx, orderMessage.TraceContext)
				span.SetAttributes(
					attribute.String("bakery.order_key", orderMessage.OrderKey),
					attribute.String("bakery.pastry", orderMessage.Pastry),
				)
				err = processFunction(ctx, orderMessage)

				if err != nil {
					errorFunc(ctx, err)
				}
			}

			err = order.Ack(false)

			if err != nil {
				errorFunc(ctx, err)
			}
			span.SetAttributes(attribute.Bool("queue.success", true))
			span.End()
		}
	}()

	<-ctx.Done()

	return nil
}
