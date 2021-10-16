module github.com/whatsfordinner/bakery/baker

go 1.14

require (
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20211016094353-651c27ce8445
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20210619113101-7ecc63acc32a
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20211016094353-651c27ce8445
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/trace v1.0.1
)

replace (
	github.com/whatsfordinner/bakery/pkg/config => ../pkg/config
	github.com/whatsfordinner/bakery/pkg/orders => ../pkg/orders
	github.com/whatsfordinner/bakery/pkg/trace => ../pkg/trace
)
