package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"g.tizu.dev/colr"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
	"golang.org/x/sync/errgroup"
)

type CUninit struct {
	Confirm bool
}

func (c *CUninit) Run(ctx *Context) error {
	if !c.Confirm {
		warn(colr.Red("DOING THIS WILL REMOVE ALL PACKAGES AND UNINSTALL KJSPKG COMPLETELY!"))
		confirm := false
		if err := NewConfirm("Are you sure?", "You can't undo this, which may cause data loss!", &confirm, false); err != nil {
			return err
		}
		if !confirm {
			return fmt.Errorf("aborted")
		}
	}

	cfg, err := kjspkg.GetConfig(ctx.Path, true)
	if err != nil {
		warn("oopsie woopsie! %v", err)
		return nil
	}

	info("Removing packages")
	errs, _ := errgroup.WithContext(context.Background())
	for id := range cfg.Installed {
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

	info("Deleting stuff")
	for _, name := range kjspkg.ScriptDirs {
		os.RemoveAll(filepath.Join(ctx.Path, name, ".kjspkg"))
	}
	os.Remove(filepath.Join(ctx.Path, ".kjspkg"))

	info("Uninitialized successfully!")
	return nil
}
