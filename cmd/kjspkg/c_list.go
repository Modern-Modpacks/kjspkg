package main

import (
	"fmt"

	"g.tizu.dev/colr"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type CList struct {
	// Count bool `optional:"" help:"Show package count instead?"`
	All    bool   `help:"Also show uninstalled packages?"`
	Search string `help:"The search query"`
}

func (c *CList) Run(ctx *Context) error {
	filtered := []string{}
	if !c.All {
		cfg, err := kjspkg.GetConfig(ctx.Path, false)
		if err != nil {
			return err
		}

		info("Installed: %d", len(cfg.Installed))
		for id := range cfg.Installed {
			if fuzzy.MatchNormalizedFold(c.Search, id) {
				filtered = append(filtered, id)
			}
		}
	} else {
		loc, err := LoadLocators()
		if err != nil {
			return err
		}

		info("Available: %d", len(loc))
		for id := range loc {
			if fuzzy.MatchNormalizedFold(c.Search, id) {
				filtered = append(filtered, id)
			}
		}
	}
	for _, id := range filtered {
		fmt.Printf(colr.Dim(" -")+" %s\n", id)
	}
	return nil
}

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
