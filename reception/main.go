package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/whatsfordinner/bakery/pkg/config"
)

func main() {
	shutdownTracer := initTracer()
	defer shutdownTracer()
	c := config.GetConfig(context.Background())
	app := new(app)
	app.init(c)
	defer app.DB.Disconnect()
	runServer(app.Router)

	os.Exit(0)
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
