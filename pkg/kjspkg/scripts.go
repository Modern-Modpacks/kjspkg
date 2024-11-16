package kjspkg

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

var ScriptFileTypes = []string{".js", ".ts", ".tsx"}

func GetStandaloneScripts(path string) []string {
	files := []string{}
	for _, dir := range ScriptDirs {
		_ = filepath.Walk(filepath.Join(path, dir), func(p string, info fs.FileInfo, err error) error {
			if info.IsDir() || strings.HasPrefix(p, filepath.Join(path, dir, ".kjspkg")) {
				return nil
			}
			if slices.Contains(ScriptFileTypes, filepath.Ext(p)) {
				files = append(files, p)
			}
			return nil
		})
	}
	return files
}

func GetAssets(path string) []string {
	files := []string{}
	for _, dir := range AssetDirs {
		_ = filepath.Walk(filepath.Join(path, dir), func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	}
	return files
}
