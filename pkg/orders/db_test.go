package orders

import (
	"fmt"
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
	tests := map[string]struct {
		input     *Order
		shouldErr bool
	}{
		"new order": {
			&Order{"brioche", "casey", "order5", "pending"},
			false,
		},
		"overwriting order": {
			&Order{"panini", "omar", "order3", "complete"},
			false,
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

			key, err := db.CreateOrder(test.input)

			if err != nil && !test.shouldErr {
				t.Fatalf("Expected no error but got %s", err.Error())
			}

			if err == nil && test.shouldErr {
				t.Fatalf("Expected errer but got no error")
			}

			if err == nil && !test.shouldErr {
				result, err := db.ReadOrder(key)

				if err != nil {
					t.Fatalf("Error while validting created order: %s", err.Error())
				}

				if !reflect.DeepEqual(result, test.input) {
					t.Fatalf("Read object does not match input object.\nGot: %+v\nExpected: %+v", result, test.input)
				}
			}
		})
	}
}

func TestReadOrder(t *testing.T) {
	tests := map[string]struct {
		expected  *Order
		orderID   string
		shouldErr bool
	}{
		"existing order": {
			&Order{"cookie", "dina", "time1", "pending"},
			"order0",
			false,
		},
		"non-existent order": {
			nil,
			"fakeorder",
			false,
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
		{"cookie", "dina", "time1", "pending"},
		{"brownie", "claude", "time2", "complete"},
		{"panini", "omar", "time3", "pending"},
	}
	db, err := radix.NewPool("tcp", "127.0.0.1:6379", 1)

	if err != nil {
		panic(err)
	}

	for i, order := range orders {
		err = db.Do(radix.FlatCmd(nil, "HSET", fmt.Sprintf("order%d", i), *order))

		if err != nil {
			panic(err)
		}
	}

	return func() {
		keys := []string{}
		err = db.Do(radix.Cmd(&keys, "KEYS", "*"))
		for _, key := range keys {
			err = db.Do(radix.Cmd(nil, "DEL", key))

			if err != nil {
				panic(err)
			}
		}

		defer db.Close()
	}
}
