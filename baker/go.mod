module github.com/whatsfordinner/bakery/baker

go 1.14

require (
	github.com/mediocregopher/radix/v3 v3.8.0 // indirect
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20211018102008-efcf29d76d3d
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20211018102008-efcf29d76d3d
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20211018102008-efcf29d76d3d
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/trace v1.0.1
	google.golang.org/genproto v0.0.0-20211018162055-cf77aa76bad2 // indirect
)

replace (
	github.com/whatsfordinner/bakery/pkg/config => ../pkg/config
	github.com/whatsfordinner/bakery/pkg/orders => ../pkg/orders
	github.com/whatsfordinner/bakery/pkg/trace => ../pkg/trace
)
