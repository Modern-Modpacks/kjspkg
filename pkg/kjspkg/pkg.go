package kjspkg

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Modern-Modpacks/kjspkg/pkg/commons"
)

var StatsViews = "https://tizudev.vercel.app/automatin/api/1025316079226064966/kjspkg?stat=views"
var StatsDownloads = "https://tizudev.vercel.app/automatin/api/1025316079226064966/kjspkg?stat=downloads"

// Note: withStats will slow down load, and MAY not result in stats being propagated!
func GetPackage(ref PackageLocator, withStats bool) (Package, error) {
	r, err := httpClient.Get(ref.URL())
	if err != nil || r.StatusCode != 200 {
		return Package{}, commons.EN200(err, ref.URL())
	}
	defer r.Body.Close()

	var pkg Package
	err = json.NewDecoder(r.Body).Decode(&pkg)
	if err != nil {
		return Package{}, err
	}

	/* func() {
		r, err := httpClient.Get(ref.URLPath() + "/README.md")
		if err != nil || r.StatusCode != 200 {
			return
		}
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}

		pkg.Readme = string(body)
	}() */

	if withStats {
		getStat := func(url string) int {
			r, err := httpClient.Get(url)
			if err != nil || r.StatusCode != 200 {
				return 0
			}
			defer r.Body.Close()

			bbody, err := io.ReadAll(r.Body)
			if err != nil {
				return 0
			}
			// ffs
			body := strings.ReplaceAll(string(bbody), "\"@ZeldaLord\": \"Activate Windows\"", "\"@ZeldaLord\": 0")

			var counts map[string]int
			if err := json.Unmarshal([]byte(body), &counts); err != nil {
				return 0
			}
			return counts[ref.Id]
		}
		pkg.Views, pkg.Downloads = getStat(StatsViews), getStat(StatsDownloads)
	}

	return pkg, nil
}

type Package struct {
	Author      string `json:"author" help:"Author(s) of the package"`
	Description string `json:"description" help:"Description of the package"`
	// Readme      string `json:"readme"`

	Versions          []int       `json:"versions" help:"Game versions the package supports"`
	ModLoaders        []ModLoader `json:"modloaders" help:"Mod loaders the package supports"`
	Dependencies      []string    `json:"dependencies" help:"Which dependencies to include in the package"`
	Incompatibilities []string    `json:"incompatabilities" help:"Which dependencies to include in the package"`

	Views     int `json:"-" hidden:""`
	Downloads int `json:"-" hidden:""`
}

func PackageFromPath(path string) (Package, error) {
	pkg := Package{}
	filePath := filepath.Join(path, ".kjspkg")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return pkg, fmt.Errorf(".kjspkg pkg file not found")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return pkg, fmt.Errorf("failed to read pkg file: %w", err)
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return pkg, fmt.Errorf("failed to parse pkg file as JSON: %w", err)
	}

	if pkg.Description == "" {
		return pkg, fmt.Errorf("package is invalid (is this an instance?)")
	}

	return pkg, nil
}
