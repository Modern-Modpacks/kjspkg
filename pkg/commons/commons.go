package commons

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
