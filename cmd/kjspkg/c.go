package main

import (
	"fmt"

	"g.tizu.dev/colr"
)

type CNotImplemented struct{}

func (c *CNotImplemented) Run(ctx *Context) error {
	fmt.Print(colr.Red("not implemented\n"))

	return nil
}
