package orders

import (
	"testing"

	"github.com/streadway/amqp"
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
			q := new(OrderQueue)
			err := q.Connect(test.host, test.username, test.password)

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

	testOrder := NewOrder("homer", "la bombe")
	testOrderKey := makeKey(testOrder)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tearDownRabbitMQ := setUpRabbitMQ()
			defer tearDownRabbitMQ()

			q := new(OrderQueue)

			if test.connected {
				q.Connect("localhost:5672", "guest", "guest")
				err := q.DeclareQueue()

				if err != nil {
					t.Fatalf("Error declaring queue: %s", err.Error())
				}
			}

			err := q.PublishOrderMessage(&OrderMessage{testOrderKey, "la bombe"})

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
