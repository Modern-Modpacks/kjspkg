package main

import (
	"slices"

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
		Value(selected).
		WithTheme(huh.ThemeBase16()).
		Run()
	return err
}

// NewMultiSelect is a wrapper function that displays a selection prompt only if
// the provided variable doesn't already match one of the available options.
func NewMultiSelect[T comparable](
	title string,
	getOptions func(opts *[]huh.Option[T]),
	selected *[]T,
) error {
	options := []huh.Option[T]{}
	getOptions(&options)

	allValid := true
	for _, opt := range *selected {
		if !slices.ContainsFunc(options, func(e huh.Option[T]) bool { return e.Value == opt }) {
			allValid = false
			break
		}
	}
	if allValid {
		return nil
	}

	err := huh.NewMultiSelect[T]().
		Title(title).
		Options(options...).
		Value(selected).
		WithTheme(huh.ThemeBase16()).
		Run()
	return err
}
