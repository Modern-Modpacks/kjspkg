package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"g.tizu.dev/colr"
	"github.com/Modern-Modpacks/kjspkg/pkg/commons"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
)

type CPkg struct {
	Package string `arg:"" help:"The package to inspect"`
	Script  bool   `help:"Return JSON instead?"`
}

func (c *CPkg) Run(ctx *Context) error {
	if c.Script {
		doLog = false
	}

	pkg, loc, err := LoadPackageById(c.Package, true)
	if err != nil {
		return err
	}

	loaders := []string{}
	for _, l := range pkg.ModLoaders {
		loaders = append(loaders, l.String())
	}
	versions := []string{}
	for _, i := range pkg.Versions {
		versions = append(versions, kjspkg.GetVersionString(i))
	}

	if c.Script {
		b, err := json.Marshal(pkg)
		if err != nil {
			return err
		}
		fmt.Printf("%s", string(b))
		return nil
	}

	fmt.Printf(colr.Bold(colr.Blue("%s"))+" by "+colr.Blue("%s")+"\n", commons.TitleCase(c.Package), pkg.Author)
	fmt.Printf("%s\n", pkg.Description)
	fmt.Printf("\n")
	fmt.Printf(colr.Blue("Lookup:")+" https://kjspkglookup.modernmodpacks.site/#%s\n", c.Package)
	fmt.Printf(colr.Blue("GitHub:")+" %s\n", loc.URLFrontend())
	fmt.Printf(colr.Blue("Views:")+" %-6d  "+colr.Blue("Downloads:")+" %-6d\n", pkg.Views, pkg.Downloads)
	fmt.Printf("\n")
	fmt.Printf(colr.Blue("Available for:") + "\n")
	fmt.Printf(colr.Blue("  Modloaders:")+" %s\n", strings.Join(loaders, ", "))
	fmt.Printf(colr.Blue("  Versions:")+" %s\n", strings.Join(versions, ", "))
	fmt.Printf("\n")
	fmt.Printf(colr.Blue("Dependencies:")+" %s\n", kjspkg.DepsJoin(pkg.Dependencies))
	fmt.Printf(colr.Blue("Incompatibilities:")+" %s\n", kjspkg.DepsJoin(pkg.Incompatibilities))

	return nil
}
