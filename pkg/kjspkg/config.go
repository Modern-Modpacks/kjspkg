package kjspkg

import "github.com/Modern-Modpacks/kjspkg/pkg/commons"

type Config struct {
	Installed map[string][]string `json:"installed"`

	Version   int       `json:"version"`
	ModLoader ModLoader `json:"modloader"`
}

type ModLoader string

func (s ModLoader) String() string {
	return commons.TitleCase(string(s))
}
func (s ModLoader) Identifier() string {
	return string(s)
}

var ModLoaders = []ModLoader{MLForge, MLFabric}

const (
	MLForge  ModLoader = "forge"
	MLFabric ModLoader = "fabric"
)

func DefaultConfig() *Config {
	return &Config{
		Installed: map[string][]string{},
	}
}
