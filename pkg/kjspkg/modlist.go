package kjspkg

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Returns mod list (K/V id/version) from KUBEJS path (so ../mods will be applied!).
// modlessOk == false causes it to error if mods dir can't be found.
// corruptedOk == false causes it to error if a mod can't be parsed (recommended true).
func GetMods(path string, modlessOk bool, corruptedOk bool) (map[string]string, error) {
	modDir := filepath.Join(path, "..", "mods")
	files, err := os.ReadDir(modDir)
	if err != nil {
		if !modlessOk {
			return nil, fmt.Errorf("mods directory not found: %w", err)
		}
		return nil, nil
	}

	modVersions := make(map[string]string)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jar") {
			zipPath := filepath.Join(modDir, file.Name())
			r, err := zip.OpenReader(zipPath)
			if err != nil {
				if !corruptedOk {
					return nil, fmt.Errorf("failed to open mod file %s: %w", file.Name(), err)
				}
				continue
			}
			defer r.Close()

			var modID, modVersion string
			var foundManifest bool

			for _, f := range r.File {
				forgeLike, fabricLike := strings.HasSuffix(f.Name, "mods.toml"), strings.HasSuffix(f.Name, ".mod.json")
				if forgeLike || fabricLike {
					file, err := f.Open()
					if err != nil {
						if !corruptedOk {
							return nil, fmt.Errorf("failed to read manifest from mod file %s: %w", f.Name, err)
						}
						continue
					}
					defer file.Close()

					if forgeLike {
						var manifest struct {
							Mods []struct {
								ModId   string `toml:"modId"`
								Version string `toml:"version"`
							} `toml:"mods"`
						}
						if _, err := toml.NewDecoder(file).Decode(&manifest); err != nil {
							if !corruptedOk {
								return nil, fmt.Errorf("failed to parse mods.toml for mod %s: %w", f.Name, err)
							}
							continue
						}
						if len(manifest.Mods) > 0 {
							modID = manifest.Mods[0].ModId
							modVersion = manifest.Mods[0].Version
							foundManifest = true
						}
					} else if fabricLike {
						var manifest struct {
							ID      string `json:"id"`
							Version string `json:"version"`
						}
						decoder := json.NewDecoder(file)
						if err := decoder.Decode(&manifest); err != nil {
							if !corruptedOk {
								return nil, fmt.Errorf("failed to parse fabric.mod.json for mod %s: %w", f.Name, err)
							}
							continue
						}
						modID = manifest.ID
						modVersion = manifest.Version
						foundManifest = true
					}
				}
			}

			if foundManifest {
				modVersions[modID] = modVersion
			} else if !corruptedOk {
				return nil, fmt.Errorf("manifest not found in mod file %s", file.Name())
			}
		}
	}

	return modVersions, nil
}
