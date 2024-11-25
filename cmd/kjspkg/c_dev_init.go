package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/Modern-Modpacks/kjspkg/pkg/commons"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
	"github.com/charmbracelet/huh"
)

type CDevInit struct {
	kjspkg.Package
	DependencyFilter []string `help:"Filter dependencies by ID"`
	UseFilter        bool     `help:"Use dependency filter to filter listed dependencies?"`
}

func (c *CDevInit) Run(ctx *Context) error {
	info("Let's get you started with a new KJSPKG package!")
	if empty, err := commons.IsEmpty(ctx.Path); err != nil || !empty {
		return fmt.Errorf("the provided directory (current or --path: %s) is not empty", ctx.Path)
	}

	loc, err := LoadLocators()
	if err != nil {
		return err
	}
	if !c.UseFilter {
		c.DependencyFilter = []string{}
		for dep := range loc {
			c.DependencyFilter = append(c.DependencyFilter, dep)
		}
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
		for _, dep := range c.DependencyFilter {
			*opts = append(*opts, huh.NewOption(dep, dep))
		}
	}, &c.Dependencies, "You may scroll to discover more!", "NONE", func(input []string) error {
		return nil
	}); err != nil {
		return err
	}

	if err := NewMultiSelect("What does this conflict with?", func(opts *[]huh.Option[string]) {
		for dep := range loc {
			if !slices.Contains(c.Dependencies, dep) {
				*opts = append(*opts, huh.NewOption(dep, dep))
			}
		}
	}, &c.Incompatibilities, "You may scroll to discover more!", "NONE", func(input []string) error {
		return nil
	}); err != nil {
		return err
	}

	pkgBytes, err := json.MarshalIndent(c.Package, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(ctx.Path, ".kjspkg"), pkgBytes, 0644); err != nil {
		return err
	}
	info("Package manifest created at %s", ctx.Path)
	return nil
}
