package main

import (
	"log"

	"github.com/whatsfordinner/bakery/pkg/orders"
)

func bakeOrder(orderMessage *orders.OrderMessage) error {
	log.Printf("Received order: %+v", *orderMessage)
	return nil
}

func rejectOrder(err error) {
	log.Printf("Error")
}
