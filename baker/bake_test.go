package main

import (
	"context"
	"testing"
	"time"

	"github.com/whatsfordinner/bakery/pkg/config"
	"github.com/whatsfordinner/bakery/pkg/orders"
)

func TestBakePastry(t *testing.T) {
	tests := map[string]struct {
		pastryName  string
		minDuration time.Duration
	}{
		"no pastry":   {"", 0},
		"some pastry": {"la bombe", 3000},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a, teardown := setUpApp()
			defer teardown()

			startTime := time.Now()
			a.bakePastry(context.Background(), test.pastryName)
			finishTime := time.Since(startTime)

			if test.minDuration.Milliseconds() > finishTime.Milliseconds() {
				t.Fatalf("Expected to sleep for %dms but only slept for %dms", test.minDuration.Milliseconds(), finishTime.Milliseconds())
			}
		})
	}
}

func TestBakeOrder(t *testing.T) {
	tests := map[string]struct {
		shouldErr  bool
		pastryName string
	}{
		"baking \"la bombe\"": {false, "la bombe"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a, teardown := setUpApp()
			defer teardown()

			key, err := a.DB.CreateOrder(context.Background(), orders.NewOrder("Homer", test.pastryName))
			if err != nil {
				t.Fatalf(err.Error())
			}

			err = a.bakeOrder(context.Background(), &orders.OrderMessage{OrderKey: key, Pastry: test.pastryName})

			if err != nil && !test.shouldErr {
				t.Fatalf("Expected no error but got %s", err.Error())
			}

			if err == nil && test.shouldErr {
				t.Fatalf("Expected error but go no error")
			}

			if err == nil && !test.shouldErr {
				order, err := a.DB.ReadOrder(context.Background(), key)
				if err != nil {
					t.Fatal(err.Error())
				}

				if order.Status != "finished" {
					t.Fatalf("Expected order status to be finished but got %s", order.Status)
				}
			}

		})
	}

}

func setUpApp() (*app, func()) {
	a := new(app)
	c := config.GetConfig(context.Background())
	a.init(c)

	return a, func() {
		err := a.DB.Disconnect()
		if err != nil {
			panic(err.Error())
		}

		err = a.Queue.Disconnect()
		if err != nil {
			panic(err.Error())
		}
	}
}
