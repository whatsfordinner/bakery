package main

import (
	"context"
	"hash/fnv"
	"log"
	"time"

	"github.com/whatsfordinner/bakery/pkg/orders"
)

func (a *app) bakeOrder(ctx context.Context, orderMessage *orders.OrderMessage) error {
	log.Printf("Received order: %+v", *orderMessage)
	err := a.DB.UpdateOrder(ctx, orderMessage.OrderKey, "baking")

	if err != nil {
		return err
	}

	a.bakePastry(ctx, orderMessage.Pastry)
	err = a.DB.UpdateOrder(ctx, orderMessage.OrderKey, "finished")

	if err != nil {
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
	bakingTime := time.Duration(h.Sum32() % 5000)
	log.Printf("%s will take %dms to bake", pastry, bakingTime)
	time.Sleep(bakingTime * time.Millisecond)
}
