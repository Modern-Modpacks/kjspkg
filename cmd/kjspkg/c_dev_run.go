package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/Modern-Modpacks/kjspkg/assets"
	"github.com/Modern-Modpacks/kjspkg/pkg/kjspkg"
)

type CDevRun struct {
	Loader   kjspkg.ModLoader `arg:"" help:"Mod loader"`
	Version  string           `arg:"" help:"Game version"`
	Launcher string           `help:"CLI name of launcher to use" default:"prismlauncher"`
}

func (c *CDevRun) Run(ctx *Context) error {
	warn("This command is experimental and may not work as expected.")

	ver, ok := kjspkg.Versions[c.Version]
	if !ok {
		return fmt.Errorf("not a game version, expected version like '1.18'")
	}

	pkg, err := kjspkg.PackageFromPath(ctx.Path)
	if err != nil {
		return err
	}

	if !slices.Contains(pkg.Versions, ver) {
		return fmt.Errorf("package does not support version %s", c.Version)
	}
	if !slices.Contains(pkg.ModLoaders, c.Loader) {
		return fmt.Errorf("package does not support loader %s", c.Loader)
	}

	launcherPath, err := getLauncherPath(c.Launcher)
	if err != nil {
		return err
	}

	instanceName := fmt.Sprintf("kjspkg%d%s", ver, c.Loader.Identifier())
	instancePath := getInstancePath(c.Launcher, instanceName)

	data, err := assets.Instances.ReadFile("instances/" + instanceName + ".zip")
	if err != nil {
		return fmt.Errorf("not supported for this version/loader combination: %s", instanceName)
	}

	putPath := filepath.Join(ctx.Path, instanceName+".zip")
	if _, err := os.Stat(instancePath); os.IsNotExist(err) {
		info("Creating new test instance...")
		info("Please name your instance '%s'!!!", instanceName)
		info("Then configure the instance (if required) and CLOSE the launcher.")
		info("Rerun the same command after.")

		if err := os.WriteFile(putPath, data, 0666); err != nil {
			return err
		}

		cmd := exec.Command(launcherPath, "-I", putPath)
		return cmd.Run()
	} else {
		os.Remove(putPath)
	}

	kubejsPath := filepath.Join(instancePath, ".minecraft", "kubejs")
	if err := os.MkdirAll(kubejsPath, 0755); err != nil {
		return err
	}

	// we have to assume this pack contains other kjspkg packages
	// this is really hacky, but it works :)
	cu := CUninit{Confirm: true}
	if err := cu.Run(&Context{Path: kubejsPath}); err != nil {
		return err
	}
	ci := CInit{Version: ver, Modloader: c.Loader}
	if err := ci.Run(&Context{Path: kubejsPath}); err != nil {
		return err
	}
	if err := kjspkg.InstallEnsureKube(kubejsPath); err != nil {
		return err
	}

	cin := CInstall{Packages: pkg.Dependencies, TrustExternal: true, NoModCheck: true}
	if err := cin.Run(&Context{Path: kubejsPath}); err != nil {
		return err
	}

	cfg, err := kjspkg.GetConfig(kubejsPath, false)
	if err != nil {
		return err
	}

	pkgAssets, err := kjspkg.InstallCopy(kubejsPath, ctx.Path)
	if err != nil {
		return err
	}

	cfg.Installed["kjspkg-dev-test"] = pkgAssets
	if err := kjspkg.SetConfig(kubejsPath, cfg); err != nil {
		return err
	}

	// if dependencies contain mods, warn the user
	for _, dep := range pkg.Dependencies {
		if strings.HasPrefix(dep, "mod:") {
			warn("This package depends on mods mods, which may have to be manually installed: %s", dep)
			break
		}
	}

	info("Launching test instance...")
	cmd := exec.Command(launcherPath, "-l", instanceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func getLauncherPath(launcher string) (string, error) {
	launcher = strings.ToLower(launcher)

	if runtime.GOOS == "windows" {
		switch launcher {
		case "prismlauncher":
			return filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "PrismLauncher", "prismlauncher.exe"), nil
		case "multimc":
			return filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "MultiMC", "MultiMC.exe"), nil // wtf gcat
		}
	} else {
		switch launcher {
		case "prismlauncher":
			return "/usr/bin/prismlauncher", nil
		case "multimc":
			return "/opt/multimc/run.sh", nil
		}
	}

	return "", fmt.Errorf("unsupported launcher: %s", launcher)
}

func getInstancePath(launcher, instanceName string) string {
	if runtime.GOOS == "windows" {
		switch strings.ToLower(launcher) {
		case "prismlauncher":
			return filepath.Join(os.Getenv("APPDATA"), "PrismLauncher", "instances", instanceName)
		case "multimc":
			return filepath.Join(os.Getenv("USERPROFILE"), "Downloads", "MultiMC", "instances", instanceName) // wtf gcat
		}
	} else {
		switch strings.ToLower(launcher) {
		case "prismlauncher":
			return filepath.Join(os.Getenv("HOME"), ".local", "share", "PrismLauncher", "instances", instanceName)
		case "multimc":
			return filepath.Join(os.Getenv("HOME"), ".local", "share", "MultiMC", "instances", instanceName)
		}
	}
	return ""
}
