package main

import (
	"context"

	"github.com/whatsfordinner/bakery/pkg/config"
)

func main() {
	a := new(app)
	c := config.GetConfig(context.Background())
	a.init(c)
	a.run()
}
