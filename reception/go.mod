module github.com/whatsfordinner/bakery/reception

go 1.14

require (
	github.com/gorilla/mux v1.8.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mediocregopher/radix/v3 v3.7.0
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20210606114944-15426c2bf092
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20210605113819-219fb02ae5ad
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.20.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/stdout v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
