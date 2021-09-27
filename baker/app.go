package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/whatsfordinner/bakery/pkg/config"
	"github.com/whatsfordinner/bakery/pkg/orders"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type app struct {
	DB     *orders.OrderDB
	Queue  *orders.OrderQueue
	tracer trace.Tracer
}

func (a *app) init(c *config.Config) {
	db, err := orders.NewDB(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	a.DB = db

	queue, err := orders.NewOrderQueue(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	a.Queue = queue

	a.tracer = otel.Tracer("")
}

func (a *app) run() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go a.Queue.ConsumeOrderQueue(ctx, a.bakeOrder, a.rejectOrder)
	<-c
	cancel()
	err := a.Queue.Disconnect()

	if err != nil {
		log.Fatal(err.Error())
	}

	err = a.DB.Disconnect()

	if err != nil {
		log.Fatal(err.Error())
	}
}
