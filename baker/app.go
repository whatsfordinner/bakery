package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/whatsfordinner/bakery/pkg/config"
	"github.com/whatsfordinner/bakery/pkg/orders"
)

type app struct {
	DB    *orders.OrderDB
	Queue *orders.OrderQueue
}

func (a *app) init(c *config.Config) {
	a.DB = new(orders.OrderDB)
	err := a.DB.Connect(c.DBHost)

	if err != nil {
		log.Fatal(err.Error())
	}

	a.Queue = new(orders.OrderQueue)
	err = a.Queue.Connect(c.RabbitHost, c.RabbitUsername, c.RabbitPassword)

	if err != nil {
		log.Fatal(err.Error())
	}
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
