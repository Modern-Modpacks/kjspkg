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
func (s ModLoader) StringLong() string {
	return ModLoadersLong[s]
}
func (s ModLoader) Identifier() string {
	return string(s)
}

var ModLoaders = []ModLoader{MLForge, MLFabric}
var ModLoadersLong = map[ModLoader]string{
	MLForge:  "Forge/NeoForge",
	MLFabric: "Fabric/Quilt",
}

const (
	MLForge  ModLoader = "forge"
	MLFabric ModLoader = "fabric"
)

func DefaultConfig() *Config {
	return &Config{
		Installed: map[string][]string{},
	}
}
