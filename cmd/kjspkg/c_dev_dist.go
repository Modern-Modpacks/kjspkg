package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
	"github.com/charmbracelet/huh"
)

type CDevDist struct {
	kjspkg.Package
	Scripts []string `help:"Scripts to include in the package"`
	Assets  []string `help:"Assets to include in the package"`
	Target  string   `help:"Target directory name (will be placed in the kubejs dir)" default:"dist"`
}

func (c *CDevDist) Run(ctx *Context) error {
	cfg, err := kjspkg.GetConfig(ctx.Path, false)
	if err != nil {
		warn("dev dist migrates a kubejs modpack to a kjspkg package.")
		return err
	}

	targetDir := filepath.Join(ctx.Path, c.Target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	scriptFiles := kjspkg.GetStandaloneScripts(ctx.Path)
	if err := NewMultiSelect("Select scripts to include in your package", func(opts *[]huh.Option[string]) {
		for _, file := range scriptFiles {
			path, err := filepath.Rel(ctx.Path, file)
			if err != nil {
				path = file
			}
			*opts = append(*opts, huh.NewOption(path, path))
		}
	}, &c.Scripts, "You may scroll to discover more scripts!", "NONE", func(input []string) error {
		return nil
	}); err != nil {
		return err
	}

	assetFiles := kjspkg.GetAssets(ctx.Path)
	if err := NewMultiSelect("Select assets to include in your package", func(opts *[]huh.Option[string]) {
		for _, file := range assetFiles {
			path, err := filepath.Rel(ctx.Path, file)
			if err != nil {
				path = file
			}
			*opts = append(*opts, huh.NewOption(path, path))
		}
	}, &c.Assets, "You may scroll to discover more assets!", "NONE", func(input []string) error {
		return nil
	}); err != nil {
		return err
	}

	installed := []string{}
	for dep := range cfg.Installed {
		installed = append(installed, dep)
	}
	init := CDevInit{
		Package:          c.Package,
		DependencyFilter: installed,
		UseFilter:        true,
	}
	if err := init.Run(&Context{Path: c.Target}); err != nil {
		return err
	}

	info("Just a moment, migrating your package...")

	types := [][]string{c.Scripts, c.Assets}
	for _, srcType := range types {
		for _, file := range srcType {
			srcPath := filepath.Join(ctx.Path, file)
			dstPath := filepath.Join(targetDir, file)

			if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
				return err
			}

			src, err := os.Open(srcPath)
			if err != nil {
				return err
			}
			defer src.Close()

			dst, err := os.Create(dstPath)
			if err != nil {
				return err
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				return err
			}
		}
	}

	info("Package created at %s", targetDir)
	return nil
}
