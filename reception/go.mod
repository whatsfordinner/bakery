module github.com/whatsfordinner/bakery/reception

go 1.14

require (
	github.com/gorilla/mux v1.8.0
	github.com/mediocregopher/radix/v3 v3.5.2
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20200926130708-109ad85a4e9c
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.13.0
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/exporters/stdout v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)
