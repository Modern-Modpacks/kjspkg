package main

type CListall struct {
	Search string `help:"The search query"`
}

func (c *CListall) Run(ctx *Context) error {
	cmd := CList{
		All:    true,
		Search: c.Search,
	}
	return cmd.Run(ctx)
}

type CSearch struct {
	Query string `arg:"" help:"The search query"`
}

func (c *CSearch) Run(ctx *Context) error {
	cmd := CListall{
		Search: c.Query,
	}
	return cmd.Run(ctx)
}
