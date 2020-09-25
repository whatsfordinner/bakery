package orders

import (
	"github.com/google/uuid"
)

// An Order represents all the details of a customer's order
type Order struct {
	Pastry   string `json:"pastry"`
	Customer string `json:"customer"`
	OrderID  string `json:"orderId"`
	Status   string `json:"status"`
}

// ToSlice returns a slice of strings for use with radix commands
func (o *Order) ToSlice() []string {
	return []string{"Pastry", o.Pastry, "Customer", o.Customer, "OrderID", o.OrderID, "Status", o.Status}
}

// NewOrder takes a customer and a pastry and returns an *Order
func NewOrder(customer string, pastry string) (*Order, error) {
	orderUUID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	order := new(Order)
	order.OrderID = string(orderUUID[:])
	order.Customer = customer
	order.Pastry = pastry
	order.Status = "pending"

	return order, nil
}
