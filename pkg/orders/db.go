package orders

import (
	"fmt"

	"github.com/mediocregopher/radix/v3"
)

// OrderDB manages the connection pool to Redis
type OrderDB struct {
	Pool *radix.Pool
}

// Connect creates a connection pool to the specified host and stores it in the OrderDB
func (db *OrderDB) Connect(host string) error {
	newPool, err := radix.NewPool("tcp", host, 5)

	if err != nil {
		return err
	}

	db.Pool = newPool
	return nil
}

// Disconnect closes the connection pool
func (db *OrderDB) Disconnect() error {
	err := db.Pool.Close()

	if err != nil {
		return err
	}

	return nil
}

// CreateOrder writes a full hash to Redis
func (db *OrderDB) CreateOrder(order *Order) error {
	err := db.Pool.Do(radix.FlatCmd(nil, "HSET", order.OrderID, order.ToSlice()))

	if err != nil {
		return err
	}

	return nil
}

// ReadOrder returns the *Order associated with the provided orderID or an error
func (db *OrderDB) ReadOrder(orderID string) (*Order, error) {
	order := new(Order)
	err := db.Pool.Do(radix.Cmd(order, "HGETALL", orderID))

	if err != nil {
		return nil, err
	}

	if order.Customer == "" {
		return nil, fmt.Errorf("order with ID %s not found", orderID)
	}

	return order, nil
}

// UpdateOrder updates the status of an order in Redis
func (db *OrderDB) UpdateOrder(orderID string, newStatus string) error {
	// Check that the order exists first
	_, err := db.ReadOrder(orderID)

	if err != nil {
		return fmt.Errorf("order with ID %s not found", orderID)
	}

	return db.Pool.Do(radix.Cmd(nil, "HSET", orderID, "Status", newStatus))
}
