package main

import (
	"fmt"

	"g.tizu.dev/colr"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
)

type CDevRun struct {
	Loader   string `arg:"" help:"Mod loader" enum:"forge,fabric"`
	Version  string `arg:"" help:"Game version"`
	Launcher string `help:"CLI name of launcher to use" default:"prismlauncher"`
}

func (c *CDevRun) Run(ctx *Context) error {
	ver, ok := kjspkg.Versions[c.Version]
	if !ok {
		return fmt.Errorf("not a game version, expected version like '1.18'")
	}

	pkg, err := kjspkg.PackageFromPath(ctx.Path)
	if err != nil {
		return err
	}

	_, _ = ver, pkg
	fmt.Print(colr.Red("not implemented\n"))

	return nil
}
