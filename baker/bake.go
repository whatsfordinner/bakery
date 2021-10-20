package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"time"

	"github.com/whatsfordinner/bakery/pkg/orders"
	"go.opentelemetry.io/otel/attribute"
)

func (a *app) bakeOrder(ctx context.Context, orderMessage *orders.OrderMessage) error {
	_, span := a.tracer.Start(ctx, "processing-order")
	defer span.End()

	span.SetAttributes(
		attribute.String("baker.order_key", orderMessage.OrderKey),
	)

	order, err := a.DB.ReadOrder(ctx, orderMessage.OrderKey)

	if err != nil {
		span.SetAttributes(
			attribute.Bool("baker.error", true),
		)
		span.AddEvent(
			fmt.Sprintf("Failed to read order from DB: %s", err.Error()),
		)

		return err
	}

	err = a.DB.UpdateOrder(ctx, orderMessage.OrderKey, "baking")

	if err != nil {
		span.SetAttributes(
			attribute.Bool("baker.error", true),
		)
		span.AddEvent(
			fmt.Sprintf("Failed to update order: %s", err.Error()),
		)

		return err
	}

	a.bakePastry(ctx, order.Pastry)

	timeOrdered, err := time.Parse(time.RFC3339, order.OrderTime)

	if err != nil {
		span.AddEvent(
			fmt.Sprintf("Failed to parse order time: %s", err.Error()),
		)
	} else {
		timeToCompletion := time.Since(timeOrdered)
		span.SetAttributes(
			attribute.Int64("baker.time_to_completion_ms", timeToCompletion.Milliseconds()),
		)
	}

	err = a.DB.UpdateOrder(ctx, orderMessage.OrderKey, "finished")

	if err != nil {
		span.SetAttributes(
			attribute.Bool("baker.error", true),
		)
		span.AddEvent(
			fmt.Sprintf("Failed to update order: %s", err.Error()),
		)
		return err
	}

	return nil
}

func (a *app) rejectOrder(ctx context.Context, err error) {
	log.Print(err.Error())
}

func (a *app) bakePastry(ctx context.Context, pastry string) {
	_, span := a.tracer.Start(ctx, "baking")
	defer span.End()
	h := fnv.New32a()
	h.Write([]byte(pastry))
	bakingTime := time.Duration(h.Sum32() % 100)
	time.Sleep(bakingTime * time.Millisecond)

	if rand.Float64() < 0.1 {
		span.AddEvent("baker made mistake, redoing")
		span.SetAttributes(
			attribute.Bool("baker.mistake", true),
		)
		time.Sleep(bakingTime * time.Millisecond)
	}

	span.SetAttributes(
		attribute.String("baker.pastry", pastry),
	)
}
