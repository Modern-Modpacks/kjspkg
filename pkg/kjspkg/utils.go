package kjspkg

import (
	"strings"

	"github.com/Modern-Modpacks/kjspkg/pkg/commons"
)

var Logo = "    ⢀⣤⣶⣿⣿⣶⣤⡀    \n  ⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣷⣦  \n⢠⣄⡀⠉⠻⢿⣿⣿⣿⣿⡿⠟⠉⢀⣠⡄\n⢸⣿⣿⣷⣦⣀⠈⠙⠋⠁⣀⣴⣾⣿⣿⡇\n⢸⣿⣿⣿⣿⣿⣿  ⣿⣿⣿⣿⣿⣿⡇\n⢸⣿⣿⣿⣿⣿⣿  ⣿⣿⣿⣿⣿⣿⡇\n ⠙⠻⣿⣿⣿⣿  ⣿⣿⣿⣿⠟⠋ \n    ⠉⠻⢿  ⡿⠟⠉    "

func DepsJoin(deps []string) string {
	out := ""
	for _, dep := range deps {
		if out != "" {
			out = out + ", "
		}
		out = out + commons.TitleCase(strings.TrimPrefix(dep, "mod:"))
		if strings.HasPrefix(dep, "mod:") {
			out = out + " (mod)"
		}
	}
	if out == "" {
		return "none!"
	}
	return out
}
