package main

import (
	"fmt"

	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
	"github.com/charmbracelet/huh"
)

type CInit struct {
	Version   int              `help:"Version to use (as number, 1.16 -> 16 - 10 -> 6)"`
	Modloader kjspkg.ModLoader `help:"Mod loader (forge, fabric, quilt, neoforge)"`
}

func (c *CInit) Run(ctx *Context) error {
	if !kjspkg.IsKube(ctx.Path) {
		return fmt.Errorf("are you sure this is the kubejs directory?")
	}
	if kjspkg.HasConfig(ctx.Path) {
		return fmt.Errorf("config already exists, use 'uninit' instead!")
	}

	cfg := kjspkg.DefaultConfig()

	err := NewSelect("Pick a game version", func(opts *[]huh.Option[int]) {
		for display, id := range kjspkg.Versions {
			*opts = append(*opts, huh.NewOption(display, id))
		}
	}, &c.Version, "")
	if err != nil {
		return err
	}
	info("Game version: %s", kjspkg.GetVersionString(c.Version))
	cfg.Version = c.Version

	err = NewSelect("Pick a mod loader", func(opts *[]huh.Option[kjspkg.ModLoader]) {
		for _, loader := range kjspkg.ModLoaders {
			*opts = append(*opts, huh.NewOption(loader.String(), loader))
		}
	}, &c.Modloader, "")
	if err != nil {
		return err
	}
	info("Mod loader: %s", c.Modloader.String())
	cfg.ModLoader = c.Modloader

	err = kjspkg.SetConfig(ctx.Path, cfg)
	if err != nil {
		return err
	}

	info("Done! Use 'install' to install your first package.")
	return nil
}
