package orders

// An Order represents all the details of a customer's order
type Order struct {
	Pastry   string `json:"pastry"`
	Customer string `json:"customer"`
	OrderID  string `json:"orderId"`
	Status   string `json:"status"`
}

// GetOrderID returns a hash to be used as the key in the DB
func GetOrderID(customer string, pastry string) string {
	return ""
}

// NewOrder generates a new Order, writes it to the DB and returns a *Order
func NewOrder(customer string, pastry string) (*Order, error) {
	order := new(Order)
	order.OrderID = GetOrderID(customer, pastry)
	order.Customer = customer
	order.Pastry = pastry
	order.Status = "pending"

	return order, nil
}

// UpdateOrder changes that status of an order with the matching ID
func UpdateOrder(orderID string, status string) error {
	return nil
}

// GetOrder returns a *Order with the data matching the provided order ID
func GetOrder(orderID string) (*Order, error) {
	return new(Order), nil
}
