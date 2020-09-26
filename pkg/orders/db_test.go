package orders

import (
	"reflect"
	"testing"

	"github.com/mediocregopher/radix/v3"
)

func TestConnect(t *testing.T) {
	tests := map[string]struct {
		host      string
		shouldErr bool
	}{
		"valid host":   {"127.0.0.1:6379", false},
		"invalid host": {"fakehost", true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := new(OrderDB)
			err := db.Connect(test.host)

			if err != nil && !test.shouldErr {
				t.Fatalf("Expected no error but got %s", err.Error())
			}

			if err == nil && test.shouldErr {
				t.Fatalf("Expected error but got no error")
			}
		})
	}
}

func TestDisconnect(t *testing.T) {
	tests := map[string]struct {
		connect   bool
		shouldErr bool
	}{
		"pool is connected":     {true, false},
		"pool is not connected": {false, true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := new(OrderDB)

			err := db.Connect("127.0.0.1:6379")

			if err != nil {
				t.Fatalf("Error connecting to redis: %s", err.Error())
			}

			if !test.connect {
				db.Pool.Close()
			}

			err = db.Disconnect()

			if err != nil && !test.shouldErr {
				t.Fatalf("Expected no error but got %s", err.Error())
			}

			if err == nil && test.shouldErr {
				t.Fatalf("Expected error but got no error")
			}
		})
	}
}

func TestCreateOrder(t *testing.T) {
	tearDown := setUp()
	defer tearDown()
}

func TestReadOrder(t *testing.T) {
	tests := map[string]struct {
		expected  *Order
		orderID   string
		shouldErr bool
	}{
		"existing order": {
			&Order{"cookie", "dina", "order1", "pending"},
			"order1",
			false,
		},
		"non-existent order": {
			nil,
			"fakeorder",
			true,
		},
	}

	db := new(OrderDB)
	err := db.Connect("127.0.0.1:6379")

	if err != nil {
		panic(err)
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			tearDown := setUp()
			defer tearDown()

			result, err := db.ReadOrder(test.orderID)

			if err != nil && !test.shouldErr {
				t.Fatalf("Expected no error but got %s", err.Error())
			}

			if err == nil && test.shouldErr {
				t.Fatalf("Expected errer but got no error")
			}

			if err == nil && !test.shouldErr && !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("Results did not match.\nGot: %+v\nExpected: %+v", result, test.expected)
			}
		})
	}
}

func TestUpdateOrder(t *testing.T) {}

func setUp() func() {
	orders := []*Order{
		{"cookie", "dina", "order1", "pending"},
		{"brownie", "claude", "order2", "complete"},
		{"panini", "omar", "order3", "pending"},
	}
	db, err := radix.NewPool("tcp", "127.0.0.1:6379", 1)

	if err != nil {
		panic(err)
	}

	for _, order := range orders {
		err = db.Do(radix.FlatCmd(nil, "HSET", order.OrderID, order.ToSlice()))

		if err != nil {
			panic(err)
		}
	}

	return func() {
		for _, order := range orders {
			err = db.Do(radix.Cmd(nil, "DEL", order.OrderID))

			if err != nil {
				panic(err)
			}
		}

		defer db.Close()
	}
}
