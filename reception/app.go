package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/whatsfordinner/bakery/pkg/orders"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	otelglobal "go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/propagators"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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
	a.Router.Use(otelmux.Middleware("bakery-reception"))
	a.Router.HandleFunc("/", a.homeHandler).Methods("GET")
	a.Router.HandleFunc("/orders", a.newOrderHandler).Methods("POST")
	a.Router.HandleFunc("/orders/{key}", a.orderStatusHandler).Methods("GET")
}

func initTracer() func() {
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithSyncer(exporter),
	)
	if err != nil {
		log.Fatal(err)
	}
	otelglobal.SetTracerProvider(tp)
	otelglobal.SetTextMapPropagator(otel.NewCompositeTextMapPropagator(propagators.TraceContext{}, propagators.Baggage{}))

	// Some exporters have shutdown methods which need to be invoked before the program quits
	return func() {}
}
