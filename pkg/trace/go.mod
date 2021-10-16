module github.com/whatsfordinner/bakery/pkg/trace

go 1.14

require (
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20210619103517-be0208586703
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/exporters/otlp v0.20.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.1
	go.opentelemetry.io/otel/sdk v1.0.1
	go.opentelemetry.io/otel/sdk/metric v0.24.0 // indirect
	golang.org/x/sys v0.0.0-20211015200801-69063c4bb744 // indirect
	google.golang.org/grpc v1.41.0
)

replace github.com/whatsfordinner/bakery/pkg/config => ../config
