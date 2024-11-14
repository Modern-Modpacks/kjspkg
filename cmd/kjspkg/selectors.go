package main

import (
	"github.com/charmbracelet/huh"
)

// NewSelect is a wrapper function that displays a selection prompt only if
// the provided variable doesn't already match one of the available options.
func NewSelect[T comparable](
	title string,
	getOptions func(opts *[]huh.Option[T]),
	selected *T,
) error {
	options := []huh.Option[T]{}
	getOptions(&options)

	for _, opt := range options {
		if *selected == opt.Value {
			return nil
		}
	}

	err := huh.NewSelect[T]().
		Title(title).
		Options(options...).
		WithTheme(huh.ThemeBase16()).
		Run()
	return err
}
