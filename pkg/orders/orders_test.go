package orders

import (
	"reflect"
	"testing"
	"time"
)

func TestNewOrder(t *testing.T) {
	tests := map[string]struct {
		expected  *Order
		customer  string
		pastry    string
		shouldErr bool
	}{
		"valid order": {
			&Order{"cookie", "foobar", "some-time", "pending"},
			"foobar",
			"cookie",
			false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := NewOrder(test.customer, test.pastry)

			if err != nil && !test.shouldErr {
				t.Fatalf("Expected no error but got %s", err.Error())
			}

			if err == nil && test.shouldErr {
				t.Fatalf("Expected error but got none")
			}

			// Check that a valid time has been provided
			_, err = time.Parse(time.RFC3339, result.OrderTime)

			if err != nil {
				t.Fatalf("Invalid time in object: %s. %s", result.OrderTime, err.Error())
			}

			// Update the OrderID to something so we can use reflect
			result.OrderTime = "some-time"

			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("Results did not match.\nGot:%+v\nExpected:%+v", result, test.expected)
			}
		})
	}
}
