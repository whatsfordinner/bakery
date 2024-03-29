package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/whatsfordinner/bakery/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	tracing "github.com/whatsfordinner/bakery/pkg/trace"
)

// OrderMessage contains the body of the queue message
type OrderMessage struct {
	TraceContext tracing.ContextCarrier `json:"traceContext"`
	TimeEnqueued string                 `json:"timeEnqueued"`
	OrderKey     string                 `json:"orderKey"`
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
func (q *OrderQueue) PublishOrderMessage(ctx context.Context, orderKey string) error {
	_, span := q.tracer.Start(ctx, "publish-order", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	span.SetAttributes(
		attribute.String("bakery.order_key", orderKey),
		attribute.Bool("queue.publish.success", false),
	)

	if q.Connection == nil {
		return errors.New("connection to RabbitMQ hasn't been established")
	}

	channel, err := q.Connection.Channel()

	if err != nil {
		span.AddEvent(
			fmt.Sprintf("Failed to open channel to RabbitMQ: %s", err.Error()),
		)
		return err
	}

	defer channel.Close()

	orderMessage := new(OrderMessage)
	orderMessage.TraceContext = tracing.ContextCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, orderMessage.TraceContext)
	orderMessage.TimeEnqueued = time.Now().Format(time.RFC3339)
	orderMessage.OrderKey = orderKey
	messageBody, err := json.Marshal(*orderMessage)

	if err != nil {
		span.AddEvent(
			fmt.Sprintf("Failed to marshal message: %s", err.Error()),
		)
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
		span.AddEvent(
			fmt.Sprintf("Failed to publish message to queue: %s", err.Error()),
		)
		return err
	}

	span.SetAttributes(attribute.Bool("queue.publish.success", true))
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
			orderMessage := new(OrderMessage)
			err := json.Unmarshal(order.Body, orderMessage)

			if err != nil {
				ctx, span := q.tracer.Start(
					ctx, "receive-order",
					trace.WithSpanKind(trace.SpanKindConsumer),
				)

				span.SetAttributes(attribute.Bool("queue.consume.success", false))
				span.AddEvent(
					fmt.Sprintf("Unable to unmarshal message content: %s", err.Error()),
				)

				errorFunc(ctx, err)
				span.End()

				return
			}

			remoteCtx := otel.GetTextMapPropagator().Extract(ctx, orderMessage.TraceContext)
			ctx, span := q.tracer.Start(
				ctx, "receive-order",
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithLinks(trace.LinkFromContext(remoteCtx)),
			)
			span.SetAttributes(
				attribute.Bool("queue.consume.success", false),
				attribute.String("bakery.order_key", orderMessage.OrderKey),
			)

			timeEnqueued, err := time.Parse(time.RFC3339, orderMessage.TimeEnqueued)

			if err != nil {
				span.AddEvent(
					fmt.Sprintf("Failed to parse time to determine time on queue: %s", err.Error()),
				)
			} else {
				timeOnQueue := time.Since(timeEnqueued)
				span.SetAttributes(
					attribute.Int64("queue.time_on_queue_ms", timeOnQueue.Milliseconds()),
				)
			}

			err = processFunction(ctx, orderMessage)

			if err != nil {
				span.AddEvent(
					fmt.Sprintf("Message processing error: %s", err.Error()),
				)
				errorFunc(ctx, err)
				return
			}

			span.AddEvent("Message processed")

			err = order.Ack(false)

			if err != nil {
				span.AddEvent(
					fmt.Sprintf("Unable to acknowledge message: %s", err.Error()),
				)
				errorFunc(ctx, err)
				return
			}

			span.AddEvent("Message acknoqledged")
			span.SetAttributes(attribute.Bool("queue.consume.success", true))
			span.End()
		}
	}()

	<-ctx.Done()

	return nil
}
