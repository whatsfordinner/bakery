package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/whatsfordinner/bakery/pkg/config"
	"github.com/whatsfordinner/bakery/pkg/orders"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type app struct {
	Router *mux.Router
	DB     *orders.OrderDB
	Queue  *orders.OrderQueue
}

func (a *app) init(c *config.Config) {
	db, err := orders.NewDB(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	a.DB = db
	a.Queue = new(orders.OrderQueue)
	err = a.Queue.Connect(c.RabbitHost, c.RabbitUsername, c.RabbitUsername)
	if err != nil {
		log.Fatal(err.Error())
	}
	a.buildRouter()
}

func (a *app) buildRouter() {
	a.Router = mux.NewRouter()
	a.Router.Use(otelmux.Middleware("bakery-reception"))
	a.Router.HandleFunc("/", a.homeHandler).Methods("GET")
	a.Router.HandleFunc("/orders", a.newOrderHandler).Methods("POST")
	a.Router.HandleFunc("/orders/{key}", a.orderStatusHandler).Methods("GET")
}
