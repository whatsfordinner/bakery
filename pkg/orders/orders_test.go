package orders

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestToSlice(t *testing.T) {
	tests := map[string]struct {
		expected []string
		object   *Order
	}{
		"uninitialised order": {
			[]string{"Pastry", "", "Customer", "", "OrderID", "", "Status", ""},
			new(Order),
		},
		"initialised order": {
			[]string{"Pastry", "cookie", "Customer", "foobar", "OrderID", "some-uuid", "Status", "radical"},
			&Order{"cookie", "foobar", "some-uuid", "radical"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := test.object.ToSlice()

			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("Results did not match.\nGot:%+v\nExpected:%+v", result, test.expected)
			}
		})
	}
}

func TestNewOrder(t *testing.T) {
	tests := map[string]struct {
		expected  *Order
		customer  string
		pastry    string
		shouldErr bool
	}{
		"valid order": {
			&Order{"cookie", "foobar", "some-uuid", "pending"},
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

			// Check that a valid UUID has been provided
			_, err = uuid.Parse(result.OrderID)

			if err != nil {
				t.Fatalf("Invalid UUID in object: %s. %s", result.OrderID, err.Error())
			}

			// Update the OrderID to something so we can use reflect
			result.OrderID = "some-uuid"

			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("Results did not match.\nGot:%+v\nExpected:%+v", result, test.expected)
			}
		})
	}
}
