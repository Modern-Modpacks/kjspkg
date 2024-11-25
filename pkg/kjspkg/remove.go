package kjspkg

import (
	"os"
	"path/filepath"
)

// This is the recommended way to remove packages. This will NOT error if the
// package wasn't installed. Note that this will also remove assets that were
// added by the package.
func Remove(path string, id string, cfg *Config) error {
	for _, name := range ScriptDirs {
		err := os.RemoveAll(filepath.Join(path, name, ".kjspkg", id))
		if err != nil {
			return err
		}
	}
	if assets, ok := cfg.Installed[id]; ok {
		for _, name := range assets {
			err := os.RemoveAll(filepath.Join(path, name))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
