module github.com/whatsfordinner/bakery/baker

go 1.14

require (
	github.com/whatsfordinner/bakery/pkg/config v0.0.0-20210619103517-be0208586703
	github.com/whatsfordinner/bakery/pkg/orders v0.0.0-20210616124049-f531dac6597d
	github.com/whatsfordinner/bakery/pkg/trace v0.0.0-20210616124049-f531dac6597d
)

replace (
	github.com/whatsfordinner/bakery/pkg/config => ../pkg/config
	github.com/whatsfordinner/bakery/pkg/orders => ../pkg/orders
	github.com/whatsfordinner/bakery/pkg/trace => ../pkg/trace
)
