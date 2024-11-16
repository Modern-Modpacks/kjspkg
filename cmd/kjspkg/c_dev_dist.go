package main

import (
	"encoding/json"
	"fmt"
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
	loc, err := LoadLocators()
	if err != nil {
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

	if err := NewInput("What does this package do?", func(input string) error {
		if input == "" {
			return fmt.Errorf("description cannot be empty")
		}
		return nil
	}, &c.Description); err != nil {
		return err
	}

	if err := NewInput("What's your name?", func(input string) error {
		if input == "" {
			return fmt.Errorf("no authors?")
		}
		return nil
	}, &c.Author); err != nil {
		return err
	}

	if err := NewMultiSelect("What versions of the game does this support?", func(opts *[]huh.Option[int]) {
		for k, v := range kjspkg.Versions {
			*opts = append(*opts, huh.NewOption(k, v))
		}
	}, &c.Versions, "", -1, func(input []int) error {
		if len(input) == 0 {
			return fmt.Errorf("at least one version is required")
		}
		return nil
	}); err != nil {
		return err
	}

	if err := NewMultiSelect("What modloaders does this support?", func(opts *[]huh.Option[kjspkg.ModLoader]) {
		for _, l := range kjspkg.ModLoaders {
			*opts = append(*opts, huh.NewOption(l.String(), l))
		}
	}, &c.ModLoaders, "", "NONE", func(input []kjspkg.ModLoader) error {
		if len(input) == 0 {
			return fmt.Errorf("at least one modloader is required")
		}
		return nil
	}); err != nil {
		return err
	}

	if err := NewMultiSelect("What does this depend on?", func(opts *[]huh.Option[string]) {
		for dep := range cfg.Installed {
			*opts = append(*opts, huh.NewOption(dep, dep))
		}
	}, &c.Dependencies, "You may scroll to discover more!", "NONE", func(input []string) error {
		return nil
	}); err != nil {
		return err
	}

	if err := NewMultiSelect("What does this conflict with?", func(opts *[]huh.Option[string]) {
		for dep := range loc {
			_, ok := cfg.Installed[dep]
			if !ok { // given you can only dep on installed, this is fine without: && slices.Contains(c.Dependencies, dep)
				*opts = append(*opts, huh.NewOption(dep, dep))
			}
		}
	}, &c.Incompatibilities, "You may scroll to discover more!", "NONE", func(input []string) error {
		return nil
	}); err != nil {
		return err
	}

	info("Just a moment, creating your package...")

	targetDir := filepath.Join(ctx.Path, c.Target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

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

	pkgBytes, err := json.MarshalIndent(c.Package, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(targetDir, ".kjspkg"), pkgBytes, 0644); err != nil {
		return err
	}

	info("Package created at %s", targetDir)
	return nil
}
