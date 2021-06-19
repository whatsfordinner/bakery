module github.com/whatsfordinner/bakery/pkg/trace

go 1.14

require (
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20210619103517-be0208586703
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/otlp v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	google.golang.org/grpc v1.37.0
)

replace github.com/whatsfordinner/bakery/pkg/config => ../config
