package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"g.tizu.dev/colr"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
	"golang.org/x/sync/errgroup"
)

type CInstall struct {
	Packages      []string `arg:"" help:"The packages to install ('github:author/repo$path@branch' syntax supported)"`
	TrustExternal bool     `help:"If GitHub packages should be trusted (experimental! updates are only partially supported)"`
	NoModCheck    bool     `help:"If mod dependency check should be skipped"`
	Skipmissing   bool     `help:"Skips dependencies that can't be found"`
	Update        bool     `help:"If packages already downloaded should be updated (same as 'update')"`
	NoInstall     bool     `help:"If new packages should not be installed (requires --update)" hidden:""`
}

func (c *CInstall) Run(ctx *Context) error {
	cfg, err := kjspkg.GetConfig(ctx.Path, false)
	if err != nil {
		return err
	}

	if c.NoInstall && !c.Update {
		return fmt.Errorf("--no-install requires --update")
	}

	refs, err := LoadLocators() // this caches them
	if err != nil {
		return err
	}

	info("Loading installed mods")
	mods, err := kjspkg.GetMods(ctx.Path, false, true)
	if err != nil {
		return err
	}
	if c.NoModCheck {
		mods = nil
	}

	toInstall := map[string]kjspkg.PackageLocator{}
	info("Resolving dependencies")
	for _, id := range c.Packages {
		list, err := kjspkg.CollectPackages(refs, cfg, id, c.TrustExternal, c.Update, mods, c.Skipmissing)
		if err != nil {
			return err
		}
		for dep, loc := range list {
			if c.NoInstall {
				if _, ok := cfg.Installed[dep]; !ok {
					return fmt.Errorf("package %s is not installed", dep)
				}
			}
			toInstall[dep] = loc
			fmt.Printf(colr.Dim(" >")+" %s\n", loc.Id)
		}
	}

	errs, _ := errgroup.WithContext(context.Background())
	info(If(c.NoInstall, "Updating", "Installing"))
	os.RemoveAll(filepath.Join(ctx.Path, "tmp"))
	for _, ref := range toInstall {
		errs.Go(func() error {
			startTime := time.Now()
			assets, err := kjspkg.Install(ctx.Path, ref, cfg, true, If(ctx.Verbose, os.Stdout, nil))
			tookTime := time.Since(startTime).Milliseconds()
			fmt.Printf(colr.Green(" +")+" %s "+colr.Dim("(took %dms)\n"), ref.Id, tookTime)
			cfg.Installed[ref.Id] = assets
			return err
		})
	}
	if err := errs.Wait(); err != nil {
		return err
	}

	os.RemoveAll(filepath.Join(ctx.Path, "tmp"))
	info("Successfully %s %d package(s)!", If(c.NoInstall, "updated", "installed"), len(toInstall))
	kjspkg.SetConfig(ctx.Path, cfg)
	return nil
}
