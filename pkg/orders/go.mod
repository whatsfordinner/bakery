module github.com/whatsfordinner/bakery/pkg/orders

go 1.14

require (
	github.com/mediocregopher/radix/v3 v3.5.2
	github.com/streadway/amqp v1.0.0
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20210619103517-be0208586703
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20210619113101-7ecc63acc32a
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
)

replace github.com/whatsfordinner/bakery/pkg/config => ../config
replace github.com/whatsfordinner/bakery/pkg/trace => ../trace
