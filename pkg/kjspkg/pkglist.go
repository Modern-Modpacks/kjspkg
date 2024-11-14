package kjspkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Modern-Modpacks/kjspkg/pkg/commons"
)

var PackageList = "https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/pkgs.json"
var httpClient = http.Client{Timeout: time.Second * 10}
var packageListCache map[string]PackageLocator

func GetPackageList() (map[string]PackageLocator, error) {
	if packageListCache != nil {
		return packageListCache, nil
	}

	r, err := httpClient.Get(PackageList)
	if err != nil || r.StatusCode != 200 {
		return nil, commons.EN200(err, PackageList)
	}
	defer r.Body.Close()

	var locators map[string]string
	err = json.NewDecoder(r.Body).Decode(&locators)
	if err != nil {
		return nil, err
	}

	var list = map[string]PackageLocator{}
	for id, input := range locators {
		loc, err := PackageLocatorFromString(id, input)
		if err != nil {
			return nil, err
		}
		list[id] = loc
	}

	packageListCache = list
	return list, nil
}

type PackageLocator struct {
	Id         string
	User       string
	Repository string
	Branch     *string
	Path       *string
}

func PackageLocatorFromString(id, input string) (PackageLocator, error) {
	p := PackageLocator{}

	// I stole this from https://github.com/Modern-Modpacks/kjspkg-lookup/blob/main/src/lib/consts.ts#L7
	regex := regexp.MustCompile(`([^/@$]*)\/([^/@$]*)(\$[^@$]*)?(@[^/@$]*)?`)
	match := regex.FindStringSubmatch(input)
	if match == nil {
		return p, fmt.Errorf("input string does not match the expected format: %s", input)
	}

	p.User = match[1]
	p.Repository = match[2]
	if len(match) > 3 && match[3] != "" {
		trimmed := strings.TrimPrefix(match[3], "$")
		p.Path = &trimmed
	}
	if len(match) > 4 && match[4] != "" {
		trimmed := strings.TrimPrefix(match[4], "@")
		p.Branch = &trimmed
	}

	if id != "" {
		p.Id = id
	} else {
		p.Id = p.Repository
	}

	return p, nil
}

// Obtains a package id from a pointer (that is to say, a string like 'kjspkg:amogus' or 'github:...')
func PackageLocatorFromPointer(id string, refs map[string]PackageLocator, trustExternal bool) (PackageLocator, error) {
	if strings.HasPrefix(id, "github:") {
		l, err := PackageLocatorFromString("", strings.TrimPrefix(id, "github:"))
		if err != nil {
			return l, err
		}
		if !trustExternal {
			return l, fmt.Errorf("you may not use external packages unless you trust them")
		}
		return l, nil
	} else {
		id = strings.TrimPrefix(id, "kjspkg:")
		r, ok := refs[id]
		if !ok {
			return PackageLocator{}, fmt.Errorf("cannot find package: %s", id)
		}
		return r, nil
	}
}

func (p *PackageLocator) String() string {
	str := p.User + "/" + p.Repository
	if p.Path != nil {
		str = str + "$" + *p.Path
	}
	if p.Branch != nil {
		str = str + "@" + *p.Branch
	}
	return str
}

func (p *PackageLocator) URL() string {
	return p.URLPath() + "/.kjspkg"
}

func (p *PackageLocator) URLPath() string {
	// TODO: what the fuck
	// possibly make this smaller with tl.tizu.dev?
	str := "https://raw.githubusercontent.com/" + p.User + "/" + p.Repository
	if p.Branch != nil {
		str = str + "/" + *p.Branch
	} else {
		str = str + "/main"
	}
	if p.Path != nil {
		str = str + "/" + *p.Path
	}
	return str
}

func (p *PackageLocator) URLFrontend() string {
	str := "https://github.com/" + p.User + "/" + p.Repository + "/tree/"
	if p.Branch != nil {
		str = str + *p.Branch
	} else {
		str = str + "main"
	}
	return str
}

func (p *PackageLocator) URLBase() string {
	return "https://github.com/" + p.User + "/" + p.Repository
}
