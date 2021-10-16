module github.com/whatsfordinner/bakery/pkg/orders

go 1.14

require (
	github.com/mediocregopher/radix/v3 v3.5.2
	github.com/streadway/amqp v1.0.0
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20211016094353-651c27ce8445
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20211016094353-651c27ce8445
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/trace v1.0.1
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211016002631-37fc39342514 // indirect
)

replace github.com/whatsfordinner/bakery/pkg/config => ../config

replace github.com/whatsfordinner/bakery/pkg/trace => ../trace
