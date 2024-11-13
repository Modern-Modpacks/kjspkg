package main

import (
	"fmt"

	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
	"github.com/charmbracelet/huh"
)

type CInit struct{}

// TODO: god
func (c *CInit) Run(ctx *Context) error {
	if !kjspkg.IsKube(ctx.Path) {
		return fmt.Errorf("are you sure this is the kubejs directory?")
	}
	if kjspkg.HasConfig(ctx.Path) {
		return fmt.Errorf("config already exists, use 'uninit' instead!")
	}

	cfg := kjspkg.DefaultConfig()

	var gamever int
	err := huh.NewSelect[int]().
		Title("Pick a game version").Options((func() []huh.Option[int] {
		opts := []huh.Option[int]{}
		for display, id := range kjspkg.Versions {
			opts = append(opts, huh.NewOption(display, id))
		}
		return opts
	}())...).Value(&gamever).WithTheme(huh.ThemeBase16()).Run()
	if err != nil {
		return err
	}
	info("Game version: %s", kjspkg.GetVersionString(gamever))
	cfg.Version = gamever

	var modloader kjspkg.ModLoader
	err = huh.NewSelect[kjspkg.ModLoader]().
		Title("Pick a mod loader").Options((func() []huh.Option[kjspkg.ModLoader] {
		opts := []huh.Option[kjspkg.ModLoader]{}
		for _, loader := range kjspkg.ModLoaders {
			opts = append(opts, huh.NewOption(loader.String(), loader))
		}
		return opts
	}())...).Value(&modloader).WithTheme(huh.ThemeBase16()).Run()
	if err != nil {
		return err
	}
	info("Mod loader: %s", modloader.String())
	cfg.ModLoader = modloader

	err = kjspkg.SetConfig(ctx.Path, cfg)
	if err != nil {
		return err
	}

	info("Done! Use 'install' to install your first package.")
	return nil
}
