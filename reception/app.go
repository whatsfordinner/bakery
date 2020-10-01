package main

import (
	"github.com/gorilla/mux"
	"github.com/whatsfordinner/bakery/pkg/orders"
)

type app struct {
	Router *mux.Router
	DB     *orders.OrderDB
}

func (a *app) init(c *config) {
	a.DB = new(orders.OrderDB)
	a.DB.Connect(*c.DBHost)
	a.buildRouter()
}

func (a *app) buildRouter() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/", a.homeHandler).Methods("GET")
	a.Router.HandleFunc("/orders", a.newOrderHandler).Methods("POST")
	a.Router.HandleFunc("/orders/{key}", a.orderStatusHandler).Methods("GET")
}
