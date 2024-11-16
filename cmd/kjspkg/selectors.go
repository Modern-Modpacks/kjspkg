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
	description string,
) error {
	options := []huh.Option[T]{}
	getOptions(&options)

	for _, opt := range options {
		if *selected == opt.Value {
			return nil
		}
	}

	sel := huh.NewSelect[T]().
		Title(title).
		Options(options...).
		Value(selected)
	if description != "" {
		sel.Description(description)
	}

	return sel.WithTheme(huh.ThemeBase16()).Run()
}

// NewMultiSelect is a wrapper function that displays a selection prompt only if
// the provided variable doesn't already match one of the available options.
func NewMultiSelect[T comparable](
	title string,
	getOptions func(opts *[]huh.Option[T]),
	selected *[]T,
	description string,
	preselectedEmpty T,
	validate func(input []T) error,
) error {
	options := []huh.Option[T]{}
	getOptions(&options)

	if len(*selected) == 1 {
		sel := *selected
		if sel[0] == preselectedEmpty {
			return nil
		}
	}

	if len(*selected) > 0 {
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
	}

	sel := huh.NewMultiSelect[T]().
		Title(title).
		Options(options...).
		Validate(validate).
		Height(9). // 8
		Value(selected)
	if description != "" {
		sel.Description(description)
	}

	return sel.WithTheme(huh.ThemeBase16()).Run()
}

// NewInput is a wrapper function that displays a selection prompt only if
// the provided variable isn't empty and can be validated.
func NewInput(
	title string,
	validate func(input string) error,
	selected *string,
) error {
	if *selected != "" && validate(*selected) == nil {
		return nil
	}

	sel := huh.NewInput().
		Title(title).
		Value(selected).
		Validate(validate)
	return sel.WithTheme(huh.ThemeBase16()).Run()
}
