package orders

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/streadway/amqp"
	"github.com/whatsfordinner/bakery/pkg/config"
	"go.opentelemetry.io/otel"
)

func TestRabbitMQConnect(t *testing.T) {
	tests := map[string]struct {
		host      string
		username  string
		password  string
		shouldErr bool
	}{
		"valid host and credentials": {"localhost:5672", "guest", "guest", false},
		"invalid host":               {"fakehost:5672", "guest", "guest", true},
		"invalid credentials":        {"localhost:5672", "fake", "fake", true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewOrderQueue((&config.Config{
				RabbitHost:     test.host,
				RabbitUsername: test.username,
				RabbitPassword: test.password,
			}))

			if err != nil && !test.shouldErr {
				t.Fatalf("Expected no error but got %s", err.Error())
			}

			if err == nil && test.shouldErr {
				t.Fatalf("Expected error but got no error")
			}
		})
	}
}

func TestRabbitMQDisconnect(t *testing.T) {

}

func TestDeclareQueue(t *testing.T) {
	tests := map[string]struct {
		connected bool
		shouldErr bool
	}{
		"connected to RabbitMQ":     {true, false},
		"not connected to RabbitMQ": {false, true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tearDownRabbitMQ := setUpRabbitMQ()
			defer tearDownRabbitMQ()

			q := new(OrderQueue)
			q.tracer = otel.Tracer("test")

			if test.connected {
				q.Connect("localhost:5672", "guest", "guest")
			}

			err := q.DeclareQueue()

			if err != nil && !test.shouldErr {
				t.Fatalf("Expected no error but got %s", err.Error())
			}

			if err == nil && test.shouldErr {
				t.Fatalf("Exected error but got no error")
			}

			if err == nil && !test.shouldErr {
				channel, err := q.Connection.Channel()

				if err != nil {
					t.Fatalf("Error verifying queue exists: %s", err.Error())
				}

				queue, err := channel.QueueInspect("orders")

				if err != nil {
					t.Fatalf("Error verifying queue exists: %s", err.Error())
				}

				if queue.Name != "orders" {
					t.Fatalf("Queue does not have correct name.\nGot: %s\nExpected: %s", queue.Name, "orders")
				}
			}

			if q.Connection != nil {
				q.Disconnect()
			}
		})
	}

}

func TestPublishOrderMessage(t *testing.T) {
	tests := map[string]struct {
		connected bool
		shouldErr bool
	}{
		"connected to RabbitMQ":     {true, false},
		"not connected to RabbitMQ": {false, true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tearDownRabbitMQ := setUpRabbitMQ()
			defer tearDownRabbitMQ()

			q := new(OrderQueue)
			q.tracer = otel.Tracer("test")

			if test.connected {
				q.Connect("localhost:5672", "guest", "guest")
				err := q.DeclareQueue()

				if err != nil {
					t.Fatalf("Error declaring queue: %s", err.Error())
				}
			}

			err := q.PublishOrderMessage(context.Background(), "testkey", "la bombe")

			if err != nil && !test.shouldErr {
				t.Fatalf("Expected no error but got %s", err.Error())
			}

			if err == nil && test.shouldErr {
				t.Fatalf("Expected error but got no error")
			}

			if err == nil && test.shouldErr {
				channel, err := q.Connection.Channel()

				if err != nil {
					t.Fatalf("Error verifying queue exists: %s", err.Error())
				}

				queue, err := channel.QueueInspect("orders")

				if err != nil {
					t.Fatalf("Error verifying queue exists: %s", err.Error())
				}

				if queue.Messages == 0 {
					t.Fatal("No message was published to the queue")
				}
			}
		})
	}
}

func TestConsumeOrderQueue(t *testing.T) {
	goodOrderMessage := &OrderMessage{makeKey(NewOrder("homer", "la bombe")), "la bombe"}
	goodMessage, _ := json.Marshal(goodOrderMessage)
	badMessage := []byte("foobarbaz")
	tests := map[string]struct {
		connected    bool
		shouldErr    bool
		goodMessages int
		badMessages  int
	}{
		"no messages":                      {true, false, 0, 0},
		"one good message":                 {true, false, 1, 0},
		"many good messages":               {true, false, 5, 0},
		"one bad message":                  {true, false, 0, 1},
		"mixture of good and bad messages": {true, false, 2, 2},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			goodFunc := func(ctx context.Context, order *OrderMessage) error {
				if order.OrderKey != goodOrderMessage.OrderKey || order.Pastry != goodOrderMessage.Pastry {
					t.Fatalf("Received message does not match good message\nExpected: %+v\nGot: %+v", *goodOrderMessage, *order)
				}
				return nil
			}

			badFunc := func(ctx context.Context, err error) {}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tearDownRabbitMQ := setUpRabbitMQ()
			defer tearDownRabbitMQ()

			q := new(OrderQueue)
			q.tracer = otel.Tracer("test")
			q.Connect("localhost:5672", "guest", "guest")

			err := q.DeclareQueue()

			if err != nil {
				t.Fatalf("Error declaring queue: %s", err.Error())
			}

			go q.ConsumeOrderQueue(ctx, goodFunc, badFunc)

			channel, err := q.Connection.Channel()

			if err != nil {
				t.Fatalf(err.Error())
			}

			defer channel.Close()

			for i := 0; i < test.goodMessages; i++ {
				err = channel.Publish(
					"",
					q.QueueName,
					false,
					false,
					amqp.Publishing{
						ContentType: "application/json",
						Body:        goodMessage,
					},
				)

				if err != nil {
					t.Fatalf(err.Error())
				}
			}

			for i := 0; i < test.badMessages; i++ {
				err = channel.Publish(
					"",
					q.QueueName,
					false,
					false,
					amqp.Publishing{
						ContentType: "application/json",
						Body:        badMessage,
					},
				)

				if err != nil {
					t.Fatalf(err.Error())
				}
			}
		})
	}

}

func setUpRabbitMQ() func() {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")

	if err != nil {
		panic(err)
	}

	return func() {
		defer connection.Close()
		channel, err := connection.Channel()

		if err != nil {
			panic(err)
		}

		_, err = channel.QueueDelete("orders", false, false, false)

		if err != nil {
			panic(err)
		}
	}
}
