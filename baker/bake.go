package main

import (
	"hash/fnv"
	"log"
	"time"

	"github.com/whatsfordinner/bakery/pkg/orders"
)

func (a *app) bakeOrder(orderMessage *orders.OrderMessage) error {
	log.Printf("Received order: %+v", *orderMessage)
	err := a.DB.UpdateOrder(orderMessage.OrderKey, "baking")

	if err != nil {
		return err
	}

	bakePastry(orderMessage.Pastry)
	err = a.DB.UpdateOrder(orderMessage.OrderKey, "finished")

	if err != nil {
		return err
	}

	return nil
}

func (a *app) rejectOrder(err error) {
	log.Print(err.Error())
}

func bakePastry(pastry string) {
	h := fnv.New32a()
	h.Write([]byte(pastry))
	bakingTime := time.Duration(h.Sum32() % 5000)
	log.Printf("%s will take %dms to bake", pastry, bakingTime)
	time.Sleep(bakingTime * time.Millisecond)
}
