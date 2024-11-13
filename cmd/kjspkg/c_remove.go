package main

import (
	"context"
	"fmt"

	"g.tizu.dev/colr"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
	"golang.org/x/sync/errgroup"
)

type CRemove struct {
	Packages []string `arg:"" help:"The packages to remove"`
}

func (c *CRemove) Run(ctx *Context) error {
	cfg, err := kjspkg.GetConfig(ctx.Path, false)
	if err != nil {
		return err
	}

	info("Removing")
	for _, id := range c.Packages {
		_, ok := cfg.Installed[id]
		if !ok {
			return fmt.Errorf("package not installed: %s", id)
		}
	}

	errs, _ := errgroup.WithContext(context.Background())
	for _, id := range c.Packages {
		errs.Go(func() error {
			err := kjspkg.Remove(ctx.Path, id, cfg)
			fmt.Printf(colr.Red(" -")+" %s\n", id)
			delete(cfg.Installed, id)
			return err
		})
	}
	if err := errs.Wait(); err != nil {
		return err
	}

	info("Successfully removed %d package(s)!", len(c.Packages))
	kjspkg.SetConfig(ctx.Path, cfg)
	return nil
}
