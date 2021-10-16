package trace

import (
	"context"

	"github.com/whatsfordinner/bakery/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

// InitTracer initialises a new OTLP trace provider and adds trace providers
// according to environment variables from a config object
func InitTracer(ctx context.Context, c *config.Config) (func(), error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(c.ServiceName),
		),
	)

	if err != nil {
		return nil, err
	}

	otlpExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(otlpExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
		),
	)

	return func() {
		_ = tp.Shutdown(ctx)
	}, nil
}
