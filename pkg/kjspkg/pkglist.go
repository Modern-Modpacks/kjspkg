package kjspkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var PackageList = "https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/pkgs.json"
var httpClient = http.Client{Timeout: time.Second * 10}
var packageListCache map[string]PackageLocator

func GetPackageList() (map[string]PackageLocator, error) {
	if packageListCache != nil {
		return packageListCache, nil
	}

	locators, err := plFetch()
	if err != nil {
		return nil, err
	}

	list, err := plTransform(locators)
	if err != nil {
		return nil, err
	}

	packageListCache = list
	return list, nil
}

func plFetch() (map[string]string, error) {
	r, err := httpClient.Get(PackageList)
	if err != nil || r.StatusCode != 200 {
		return nil, EN200(err, PackageList)
	}
	defer r.Body.Close()

	var locators map[string]string
	err = json.NewDecoder(r.Body).Decode(&locators)
	return locators, err
}

func plTransform(list map[string]string) (map[string]PackageLocator, error) {
	var result = map[string]PackageLocator{}
	for id, input := range list {
		loc := PackageLocator{}
		err := loc.FromString(id, input)
		if err != nil {
			return nil, err
		}
		result[id] = loc
	}

	return result, nil
}

type PackageLocator struct {
	Id         string
	User       string
	Repository string
	Branch     *string
	Path       *string
}

func (p *PackageLocator) FromString(id, input string) error {
	// I stole this from https://github.com/Modern-Modpacks/kjspkg-lookup/blob/main/src/lib/consts.ts#L7
	regex := regexp.MustCompile(`([^/@$]*)\/([^/@$]*)(\$[^@$]*)?(@[^/@$]*)?`)

	match := regex.FindStringSubmatch(input)
	if match == nil {
		return fmt.Errorf("input string does not match the expected format: %s", input)
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

	return nil
}

// Obtains a package id from a pointer (that is to say, a string like 'kjspkg:amogus' or 'github:...')
func (p *PackageLocator) FromPointer(id string, refs map[string]PackageLocator, trustExternal bool) error {
	if strings.HasPrefix(id, "github:") {
		l := PackageLocator{}
		err := l.FromString("", strings.TrimPrefix(id, "github:"))
		if err != nil {
			return err
		}
		if !trustExternal {
			return fmt.Errorf("you may not use external packages unless you trust them")
		}
		*p = l
	} else {
		id = strings.TrimPrefix(id, "kjspkg:")
		r, ok := refs[id]
		if !ok {
			return fmt.Errorf("cannot find package: %s", id)
		}
		*p = r
	}
	return nil
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
