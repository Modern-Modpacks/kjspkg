package kjspkg

import (
	"fmt"
	"slices"
	"strings"
)

// TODO: refactor this; gcast, trust me, this was worse
// *mods* may be nil, in which case mod check gets ignored
func CollectPackages(refs map[string]PackageLocator, cfg *Config, id string, trustExternal, update bool, mods map[string]string) (map[string]PackageLocator, error) {
	toInstall := map[string]PackageLocator{}

	ref := PackageLocator{}
	err := ref.FromPointer(id, refs, trustExternal)
	if err != nil {
		return nil, err
	}

	// TODO: make this better
	if _, ok := cfg.Installed[ref.Id]; ok {
		if !update {
			return toInstall, nil
		}
	}

	pkg, err := GetPackage(ref, false)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(pkg.Versions, cfg.Version) {
		versions := []string{}
		for _, i := range pkg.Versions {
			versions = append(versions, GetVersionString(i))
		}
		return nil, fmt.Errorf("not available for %s, only %s", GetVersionString(cfg.Version), strings.Join(versions, ", "))
	}
	if !slices.Contains(pkg.ModLoaders, cfg.ModLoader) {
		loaders := []string{}
		for _, l := range pkg.ModLoaders {
			loaders = append(loaders, TitleCase(string(l)))
		}
		return nil, fmt.Errorf("not available for %s, only %s", TitleCase(string(cfg.ModLoader)), strings.Join(loaders, ", "))
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
				deps, err := CollectPackages(refs, cfg, dep, trustExternal, update, mods)
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

	// TODO: make this better
	if !update {
		for _, ref := range toInstall {
			if _, ok := cfg.Installed[ref.Id]; ok {
				delete(toInstall, ref.Id)
			}
		}
	}

	return toInstall, nil
}
