package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type config struct {
	DBHost *string
}

func main() {
	shutdownTracer := initTracer()
	defer shutdownTracer()
	c := getConfig()
	app := new(app)
	app.init(c)
	defer app.DB.Disconnect()
	runServer(app.Router)

	os.Exit(0)
}

func getConfig() *config {
	c := new(config)
	c.DBHost = flag.String("dbhost", "127.0.0.1:6379", "connection string for Redis DB")
	flag.Parse()

	return c
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
