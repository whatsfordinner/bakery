module github.com/whatsfordinner/bakery/reception

go 1.14

require (
	github.com/gorilla/mux v1.8.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mediocregopher/radix/v3 v3.7.0
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20210619103517-be0208586703
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20210616124049-f531dac6597d
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20210616124049-f531dac6597d
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.20.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	github.com/whatsfordinner/bakery/pkg/config => ../pkg/config
	github.com/whatsfordinner/bakery/pkg/orders => ../pkg/orders
	github.com/whatsfordinner/bakery/pkg/trace => ../pkg/trace
)
