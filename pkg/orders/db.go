package orders

import (
	"crypto/md5"
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

// CreateOrder writes a full hash to Redis and returns the key
func (db *OrderDB) CreateOrder(order *Order) (string, error) {
	key := makeKey(order)
	err := db.Pool.Do(radix.FlatCmd(nil, "HSET", key, *order))

	if err != nil {
		return "", err
	}

	return key, nil
}

// ReadOrder returns the *Order associated with the provided orderID or an error
func (db *OrderDB) ReadOrder(orderKey string) (*Order, error) {
	order := new(Order)
	err := db.Pool.Do(radix.Cmd(order, "HGETALL", orderKey))

	if err != nil {
		return nil, err
	}

	if order.Customer == "" {
		return nil, nil
	}

	return order, nil
}

// UpdateOrder updates the status of an order in Redis
func (db *OrderDB) UpdateOrder(orderKey string, newStatus string) error {
	// Check that the order exists first
	order, err := db.ReadOrder(orderKey)

	if err != nil {
		return err
	}

	if order == nil {
		return fmt.Errorf("Order with key %s does not exist", orderKey)
	}

	return db.Pool.Do(radix.Cmd(nil, "HSET", orderKey, "Status", newStatus))
}

func makeKey(order *Order) string {
	hashInput := []byte(order.Customer + order.Pastry + order.OrderTime)
	hash := md5.Sum(hashInput)

	return fmt.Sprintf("%x", hash)
}
