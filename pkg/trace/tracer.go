package trace

import (
	"context"

	"github.com/whatsfordinner/bakery/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"google.golang.org/grpc"
)

// InitTracer initialises a new OTLP trace provider and adds trace providers
// according to environment variables from a config object
func InitTracer(ctx context.Context, c *config.Config) (func(), error) {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource.NewWithAttributes(semconv.ServiceNameKey.String(c.ServiceName))),
	)

	if c.JaegerEndpoint != "" {
		jaegerExporter, err := otlp.NewExporter(
			ctx,
			otlpgrpc.NewDriver(
				otlpgrpc.WithEndpoint(c.JaegerEndpoint),
				otlpgrpc.WithDialOption(grpc.WithInsecure()),
			),
		)

		if err != nil {
			return nil, err
		}

		tp.RegisterSpanProcessor(sdktrace.NewBatchSpanProcessor(jaegerExporter))
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return func() {
		_ = tp.Shutdown(ctx)
	}, nil
}
