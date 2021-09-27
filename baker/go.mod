module github.com/whatsfordinner/bakery/baker

go 1.14

require (
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20210619103517-be0208586703
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20210619113101-7ecc63acc32a
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20210619113101-7ecc63acc32a
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
)

replace (
	github.com/whatsfordinner/bakery/pkg/config => ../pkg/config
	github.com/whatsfordinner/bakery/pkg/orders => ../pkg/orders
	github.com/whatsfordinner/bakery/pkg/trace => ../pkg/trace
)
