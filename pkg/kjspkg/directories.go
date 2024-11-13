package kjspkg

import (
	"os"
	"path/filepath"
)

var ScriptDirs = []string{"server_scripts", "client_scripts", "startup_scripts"}
var AssetDirs = []string{"data", "assets"}

func IsKube(path string) bool {
	for _, dir := range ScriptDirs {
		dirPath := filepath.Join(path, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return false
		}
	}
	return true
}
