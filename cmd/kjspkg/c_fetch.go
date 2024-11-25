package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"g.tizu.dev/colr"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
)

type CFetch struct{}

func (c *CFetch) Run(ctx *Context) error {
	cfg, err := kjspkg.GetConfig(ctx.Path, false)
	if err != nil {
		return err
	}

	mods, err := kjspkg.GetMods(ctx.Path, true, true)
	if err != nil {
		return err
	}

	logoLines := strings.Split(kjspkg.Logo, "\n")
	type Line struct{ K, V string }
	infoLines := []Line{
		{K: "version ", V: kjspkg.GetVersionString(cfg.Version)},
		{K: "loader  ", V: cfg.ModLoader.String()},
		{K: "pkgs    ", V: fmt.Sprint(len(cfg.Installed))},
		{K: "kube    ", V: fmt.Sprint(mods["kubejs"])},
		{K: "rhino   ", V: fmt.Sprint(mods["rhino"])},
		{K: "arch    ", V: fmt.Sprint(mods["architectury"])},
	}

	for i := 0; i < max(len(logoLines), len(infoLines)+1); i++ {
		if len(logoLines) > i {
			fmt.Print(colr.Purple(logoLines[i] + "    "))
		} else {
			fmt.Print(colr.Purple(strings.Repeat(" ", utf8.RuneCountInString(logoLines[1])+2)))
		}
		if i == 0 {
			fmt.Printf(colr.Bold("KJSPKG")+colr.Dim(" @ %s"), ctx.Path)
		} else if len(infoLines) >= i {
			info := infoLines[i-1]
			fmt.Printf(colr.Purple("%s")+"%s", info.K, info.V)
		}
		fmt.Printf("\n")
	}
	return nil
}
