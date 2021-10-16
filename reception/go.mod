module github.com/whatsfordinner/bakery/reception

go 1.14

require (
	github.com/gorilla/mux v1.8.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mediocregopher/radix/v3 v3.7.0
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20211016094353-651c27ce8445
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20210619113101-7ecc63acc32a
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20211016094353-651c27ce8445
	go.opentelemetry.io/contrib v0.20.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.25.0
	go.opentelemetry.io/contrib/propagators v0.20.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	github.com/whatsfordinner/bakery/pkg/config => ../pkg/config
	github.com/whatsfordinner/bakery/pkg/orders => ../pkg/orders
	github.com/whatsfordinner/bakery/pkg/trace => ../pkg/trace
)
