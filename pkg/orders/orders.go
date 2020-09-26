package orders

import (
	"time"
)

// An Order represents all the details of a customer's order
type Order struct {
	Pastry    string `json:"pastry"`
	Customer  string `json:"customer"`
	OrderTime string `json:"orderTime"`
	Status    string `json:"status"`
}

// NewOrder takes a customer and a pastry and returns an *Order
func NewOrder(customer string, pastry string) (*Order, error) {
	order := new(Order)
	order.Customer = customer
	order.OrderTime = time.Now().Format(time.RFC3339)
	order.Pastry = pastry
	order.Status = "pending"

	return order, nil
}
