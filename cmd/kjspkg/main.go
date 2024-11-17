package main

import (
	"fmt"
	"os"

	"g.tizu.dev/colr"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
	"github.com/alecthomas/kong"
)

type Context struct {
	Verbose bool
	Path    string
}

var cli struct {
	Verbose bool   `help:"Print verbose"`
	Quiet   bool   `help:"No non-error output" short:"q" hidden:""`
	Path    string `help:"Path to KubeJS directory (defaults to current)" default:"." type:"existingdir"`
	Source  string `help:"URL source to package list" type:"url"`

	Install   CInstall        `cmd:"" help:"Installs packages" aliases:"download"`
	Remove    CRemove         `cmd:"" help:"Removes packages" aliases:"uninstall"`
	Update    CUpdate         `cmd:"" help:"Update packages (same as 'install --update')"`
	Updateall CUpdateAll      `cmd:"" help:"Update all packages (same as 'update *')"`
	List      CList           `cmd:"" help:"Lists packages (and the count of them)"`
	Fetch     CFetch          `cmd:"" help:"neofetch but different"`
	Pkg       CPkg            `cmd:"" help:"Shows info about the package"`
	Listall   CListall        `cmd:"" help:"Lists online packages (same as 'list --all')" aliases:"all"`
	Search    CSearch         `cmd:"" help:"Search online packages (same as 'listall --search')"`
	Init      CInit           `cmd:"" help:"Initialize a KJSPKG env"`
	Uninit    CUninit         `cmd:"" help:"Remove all KJSPKG-related things in your project"`
	Reload    CNotImplemented `cmd:"" hidden:"" aliases:"refresh"`
	Dev       struct {
		Init CDevInit        `cmd:"" help:"Initialize a new KJSPKG package"`
		Run  CDevRun         `cmd:"" help:"Runs your package in a test Minecraft instance"`
		Dist CDevDist        `cmd:"" help:"Creates a package from your packs' KubeJS folder"`
		Test CNotImplemented `cmd:"" hidden:""`
	} `cmd:"" help:"Helper functions for developing KJSPKG packages"`
	Gui CNotImplemented `cmd:"" hidden:""`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("kjspkg"),
		kong.Description("A KubeJS package manager"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}))

	if cli.Quiet && cli.Verbose {
		fmt.Printf("   "+colr.BlackOnRed(" :( ")+" %v\n", "Cannot use -q in combination with --verbose")
		os.Exit(1)
	}
	if cli.Quiet {
		os.Stdout.Close()
	}
	if cli.Source != "" {
		kjspkg.PackageList = cli.Source
	}

	err := ctx.Run(&Context{
		Path:    cli.Path,
		Verbose: cli.Verbose,
	})
	if err != nil {
		fmt.Printf("   "+colr.BlackOnRed(" :( ")+" %v\n", err)
		os.Exit(1)
	}
}
