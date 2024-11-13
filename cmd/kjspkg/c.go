package main

import (
	"fmt"

	"g.tizu.dev/colr"
)

type CNotImplemented struct{}

func (c *CNotImplemented) Run(ctx *Context) error {
	fmt.Printf(colr.Red("not implemented\n"))

	return nil
}
