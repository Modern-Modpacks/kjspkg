package kjspkg

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Modern-Modpacks/kjspkg/pkg/commons"
)

// *mods* may be nil, in which case mod check gets ignored
// skipMissing == true may return partial dependency trees or an empty array if they can't be resolved.
func CollectPackages(refs map[string]PackageLocator, cfg *Config, id string, trustExternal, update bool, mods map[string]string, skipMissing bool) (map[string]PackageLocator, error) {
	toInstall := map[string]PackageLocator{}

	ref, err := PackageLocatorFromPointer(id, refs, trustExternal)
	if err != nil {
		if skipMissing {
			return toInstall, nil
		}
		return nil, err
	}

	if _, ok := cfg.Installed[ref.Id]; ok {
		if !update {
			return toInstall, nil
		}
	}

	pkg, err := GetPackage(ref, false)
	if err != nil {
		if skipMissing {
			return toInstall, nil
		}
		return nil, err
	}

	if err := CollectEnsureVersion(pkg, cfg); err != nil {
		return nil, err
	}

	for _, dep := range pkg.Incompatibilities {
		if !strings.HasPrefix(dep, "mod:") {
			if _, ok := cfg.Installed[dep]; ok {
				return nil, fmt.Errorf("%s is incompatible with %s", ref.Id, dep)
			}
		} else if mods != nil {
			dep = strings.TrimPrefix(dep, "mod:")
			if _, ok := mods[dep]; ok {
				return nil, fmt.Errorf("%s is incompatible with %s", ref.Id, dep)
			}
		}
	}
	for _, dep := range pkg.Dependencies {
		if !strings.HasPrefix(dep, "mod:") {
			if _, ok := cfg.Installed[dep]; !ok {
				deps, err := CollectPackages(refs, cfg, dep, trustExternal, update, mods, skipMissing)
				if err != nil {
					return nil, err
				}
				for dep, loc := range deps {
					toInstall[dep] = loc
				}
			}
		} else if mods != nil {
			dep = strings.TrimPrefix(dep, "mod:")
			if _, ok := mods[dep]; ok {
				return nil, fmt.Errorf("%s requires mod %s, but not found", ref.Id, dep)
			}
		}
	}
	toInstall[ref.String()] = ref

	return toInstall, nil
}

// intended to be used by CollectPackages!
func CollectEnsureVersion(pkg Package, cfg *Config) error {
	if !slices.Contains(pkg.Versions, cfg.Version) {
		versions := []string{}
		for _, i := range pkg.Versions {
			versions = append(versions, GetVersionString(i))
		}
		return fmt.Errorf("not available for %s, only %s", GetVersionString(cfg.Version), strings.Join(versions, ", "))
	}
	if !slices.Contains(pkg.ModLoaders, cfg.ModLoader) {
		loaders := []string{}
		for _, l := range pkg.ModLoaders {
			loaders = append(loaders, commons.TitleCase(string(l)))
		}
		return fmt.Errorf("not available for %s, only %s", commons.TitleCase(string(cfg.ModLoader)), strings.Join(loaders, ", "))
	}
	return nil
}
