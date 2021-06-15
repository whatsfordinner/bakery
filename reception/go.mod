module github.com/whatsfordinner/bakery/reception

go 1.14

require (
	github.com/gorilla/mux v1.8.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mediocregopher/radix/v3 v3.7.0
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20210615113430-c552aa96f02b
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20210615134525-eb257147d41e
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20210615113430-c552aa96f02b
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.20.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
