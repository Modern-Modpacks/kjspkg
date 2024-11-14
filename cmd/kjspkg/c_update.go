package main

import (
	"fmt"
	"slices"

	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
)

type CUpdate struct {
	Packages      []string `optional:"" arg:"" help:"The packages to update ('github:author/repo$path@branch' syntax supported)"`
	TrustExternal bool     `help:"If GitHub packages should be trusted"`
	NoModCheck    bool     `help:"If mod dependency check should be skipped (experimental)"`
	All           bool     `help:"Update all packages ('update *'/'updateall' use this)"`
}

func (c *CUpdate) Run(ctx *Context) error {
	packages := []string{}
	if slices.Contains(c.Packages, "*") {
		if len(c.Packages) > 1 {
			return fmt.Errorf("cannot combine * with others")
		}

		cfg, err := kjspkg.GetConfig(ctx.Path, true)
		if err != nil {
			return err
		}

		for pkg := range cfg.Installed {
			packages = append(packages, pkg)
		}
	} else {
		packages = c.Packages
	}

	cmd := CInstall{
		Packages:      packages,
		TrustExternal: c.TrustExternal,
		NoModCheck:    c.NoModCheck,
		Update:        true,
	}
	return cmd.Run(ctx)
}

func (c *CUpdate) AfterApply() error {
	if !c.All && len(c.Packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	if c.All && len(c.Packages) > 0 {
		return fmt.Errorf("cannot upgrade specific packages if --all specified")
	}

	if c.All {
		c.Packages = []string{"*"}
	}

	return nil
}

type CUpdateAll struct {
	TrustExternal bool `help:"If GitHub packages should be trusted"`
	NoModCheck    bool `help:"If mod dependency check should be skipped (experimental)"`
}

func (c *CUpdateAll) Run(ctx *Context) error {
	cmd := CUpdate{
		Packages:      []string{"*"},
		TrustExternal: c.TrustExternal,
		NoModCheck:    c.NoModCheck,
	}
	return cmd.Run(ctx)
}
