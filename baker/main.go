package main

import (
	"context"
	"log"

	"github.com/whatsfordinner/bakery/pkg/config"
	"github.com/whatsfordinner/bakery/pkg/trace"
)

func main() {
	ctx := context.Background()
	c := config.GetConfig(ctx)
	c.ServiceName = "baker"
	shutdownTracer, err := trace.InitTracer(ctx, c)
	if err != nil {
		log.Fatalf("Error initialising tracer: %v", err.Error())
	}
	defer shutdownTracer()
	app := new(app)
	app.init(c)
	app.run()
}
