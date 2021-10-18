module github.com/whatsfordinner/bakery/baker

go 1.14

require (
	github.com/mediocregopher/radix/v3 v3.8.0 // indirect
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20211016120743-e056df21de46
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20211016120743-e056df21de46
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20211016120743-e056df21de46
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/trace v1.0.1
)

replace (
	github.com/whatsfordinner/bakery/pkg/config => ../pkg/config
	github.com/whatsfordinner/bakery/pkg/orders => ../pkg/orders
	github.com/whatsfordinner/bakery/pkg/trace => ../pkg/trace
)
