package kjspkg

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var Logo = "    ⢀⣤⣶⣿⣿⣶⣤⡀    \n  ⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣷⣦  \n⢠⣄⡀⠉⠻⢿⣿⣿⣿⣿⡿⠟⠉⢀⣠⡄\n⢸⣿⣿⣷⣦⣀⠈⠙⠋⠁⣀⣴⣾⣿⣿⡇\n⢸⣿⣿⣿⣿⣿⣿  ⣿⣿⣿⣿⣿⣿⡇\n⢸⣿⣿⣿⣿⣿⣿  ⣿⣿⣿⣿⣿⣿⡇\n ⠙⠻⣿⣿⣿⣿  ⣿⣿⣿⣿⠟⠋ \n    ⠉⠻⢿  ⡿⠟⠉    "

func DepsJoin(deps []string) string {
	out := ""
	for _, dep := range deps {
		if out != "" {
			out = out + ", "
		}
		out = out + TitleCase(strings.TrimPrefix(dep, "mod:"))
		if strings.HasPrefix(dep, "mod:") {
			out = out + " (mod)"
		}
	}
	if out == "" {
		return "none!"
	}
	return out
}

func TitleCase(input string) string {
	input = strings.ReplaceAll(input, "-", " ")
	input = strings.ReplaceAll(input, "_", " ")
	titleCase := cases.Title(language.English)
	return titleCase.String(input)
}

func EN200(err error, ctx string) error {
	if err != nil {
		return err
	}
	return fmt.Errorf("request did not return 200: %s", ctx)
}
