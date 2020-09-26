package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/whatsfordinner/bakery/pkg/orders"
)

type app struct {
	Router *mux.Router
	DB     *orders.OrderDB
}

func main() {
	app := new(app)
	app.init()
	defer app.DB.Disconnect()
	runServer(app.Router)
	os.Exit(0)
}

func (a *app) init() {
	a.DB = new(orders.OrderDB)
	a.DB.Connect("127.0.0.1:6379")
	a.buildRouter()
}

func (a *app) buildRouter() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/", a.homeHandler).Methods("GET")
	a.Router.HandleFunc("/orders", a.newOrderHandler).Methods("POST")
	a.Router.HandleFunc("/orders/{key}", a.orderStatusHandler).Methods("GET")
}

func runServer(router *mux.Router) {
	// This is taken wholesale from the gorilla/mux README
	server := &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	log.Printf("Starting HTTP server on %s", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("Shutting down")

}
