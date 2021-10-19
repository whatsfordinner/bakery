package orders

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/mediocregopher/radix/v3"
	"github.com/whatsfordinner/bakery/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// OrderDB manages the connection pool to Redis
type OrderDB struct {
	Pool   *radix.Pool
	tracer trace.Tracer
}

func NewDB(c *config.Config) (*OrderDB, error) {
	db := new(OrderDB)

	db.tracer = otel.Tracer("redis")
	if err := db.Connect(c.DBHost); err != nil {
		return nil, err
	}

	return db, nil
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
	if err := db.Pool.Close(); err != nil {
		return err
	}

	return nil
}

// CreateOrder writes a full hash to Redis and returns the key
func (db *OrderDB) CreateOrder(ctx context.Context, order *Order) (string, error) {
	_, span := db.tracer.Start(ctx, "create-order", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(
		attribute.String("db.redis_operation", "HSET"),
		attribute.Bool("db.success", false),
	)

	key := makeKey(order)
	span.SetAttributes(attribute.String("db.order_key", key))
	err := db.Pool.Do(radix.FlatCmd(nil, "HSET", key, *order))

	if err != nil {
		return "", err
	}

	span.SetAttributes(attribute.Bool("db.success", true))
	return key, nil
}

// ReadOrder returns the *Order associated with the provided orderID or an error
func (db *OrderDB) ReadOrder(ctx context.Context, orderKey string) (*Order, error) {
	_, span := db.tracer.Start(ctx, "read-order", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(
		attribute.String("db.redis_operation", "HGETALL"),
		attribute.String("db.order_key", orderKey),
		attribute.Bool("db.found_order", false),
	)

	order := new(Order)
	err := db.Pool.Do(radix.Cmd(order, "HGETALL", orderKey))

	if err != nil {
		return nil, err
	}

	if order.Customer == "" {
		return nil, nil
	}

	span.SetAttributes(attribute.Bool("found_order", true))
	return order, nil
}

// UpdateOrder updates the status of an order in Redis
func (db *OrderDB) UpdateOrder(ctx context.Context, orderKey string, newStatus string) error {
	ctx, span := db.tracer.Start(ctx, "update-order", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(
		attribute.String("db.redis_operation", "HSET"),
		attribute.String("db.order_key", orderKey),
		attribute.String("db.order_status", newStatus),
		attribute.Bool("db.success", false),
	)

	// Check that the order exists first
	order, err := db.ReadOrder(ctx, orderKey)

	if err != nil {
		return err
	}

	if order == nil {
		return fmt.Errorf("Order with key %s does not exist", orderKey)
	}

	if err = db.Pool.Do(radix.Cmd(nil, "HSET", orderKey, "Status", newStatus)); err != nil {
		return err
	}

	span.SetAttributes(attribute.Bool("db.success", true))
	return nil
}

func makeKey(order *Order) string {
	hashInput := []byte(order.Customer + order.Pastry + order.OrderTime)
	hash := md5.Sum(hashInput)

	return fmt.Sprintf("%x", hash)
}
