// TODO: hey gcast pls migrate this off of chmod 0744 :3
package kjspkg

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"
)

// This is the recommended way to install packages.
// This will also update/reinstall packages that have already been installed.
// If mass is provided, some actions won't be done that may cause other concurrent
// install calls to fail, like deleting the tmp directory.
func Install(path string, loc PackageLocator, cfg *Config, mass bool, stdout io.Writer) ([]string, error) {
	if !mass {
		os.RemoveAll(filepath.Join(path, "tmp"))
	}

	err := InstallEnsureKube(path)
	if err != nil {
		return nil, err
	}

	err = InstallDiscardExisting(path, loc, cfg)
	if err != nil {
		return nil, err
	}

	err = InstallClone(path, loc, stdout)
	if err != nil {
		return nil, err
	}

	if loc.Branch == nil {
		err = InstallBranch(path, loc, stdout)
		if err != nil {
			return nil, err
		}
	}

	assets, err := InstallCopy(path, loc)
	if err != nil {
		return nil, err
	}

	// if unsuccessful, this does not delete tmp. this is intentional.
	if !mass {
		os.RemoveAll(filepath.Join(path, "tmp"))
	}

	return assets, nil
}

func InstallEnsureKube(path string) error {
	if !IsKube(path) {
		return fmt.Errorf("are you sure this is the kubejs directory?")
	}

	for _, name := range ScriptDirs {
		err := os.MkdirAll(filepath.Join(path, name, ".kjspkg"), 0744)
		if err != nil {
			return err
		}
	}
	for _, name := range AssetDirs {
		err := os.MkdirAll(filepath.Join(path, name), 0744)
		if err != nil {
			return err
		}
	}

	return nil
}

func InstallDiscardExisting(path string, loc PackageLocator, cfg *Config) error {
	return Remove(path, loc.Id, cfg)
}

func InstallClone(path string, loc PackageLocator, stdout io.Writer) error {
	err := os.MkdirAll(filepath.Join(path, "tmp", loc.Id), 0744)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", loc.URLBase(), loc.Id)
	cmd.Dir = filepath.Join(path, "tmp")
	cmd.Stdout = stdout
	cmd.Stderr = stdout
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func InstallBranch(path string, loc PackageLocator, stdout io.Writer) error {
	if loc.Branch == nil {
		return nil
	}

	// TODO: migrate to 'git switch'
	cmd := exec.Command("git", "checkout", *loc.Branch)
	cmd.Dir = filepath.Join(path, "tmp", loc.Id)
	cmd.Stdout = stdout
	cmd.Stderr = stdout
	return cmd.Run()
}

type LocEsque interface {
	PackageLocator | string
}

func InstallCopy[L LocEsque](path string, loc L) ([]string, error) {
	id, repoPath := "", ""
	switch l := any(loc).(type) {
	case PackageLocator:
		id = l.Id
		repoPath = filepath.Join(path, "tmp", l.Id)
		if l.Path != nil && !strings.Contains(*l.Path, "..") {
			repoPath = filepath.Join(repoPath, *l.Path)
		}
	case string:
		id = filepath.Base(l)
		repoPath = filepath.Join(path, "tmp", id)
		os.MkdirAll(repoPath, 0744)
		err := cp.Copy(l, repoPath)
		if err != nil {
			return nil, err
		}
	}

	for _, name := range ScriptDirs {
		root, target := filepath.Join(repoPath, name), filepath.Join(path, name, ".kjspkg", id)
		os.MkdirAll(root, 0744) // TODO: creates dir so that I don't have to check if it exists or not
		err := os.Rename(root, target)
		if err != nil {
			return nil, err
		}
	}

	assets := []string{}
	for _, name := range AssetDirs {
		root := filepath.Join(repoPath, name)
		os.MkdirAll(root, 0744) // TODO: creates dir so that I don't have to check if it exists or not
		err := filepath.WalkDir(root, func(longpath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			relpath, err := filepath.Rel(root, longpath)
			if err != nil {
				return err
			}

			virtpath := filepath.Join(path, name, relpath)
			if _, err := os.Stat(virtpath); err == nil {
				return fmt.Errorf("asset file exists: %s", virtpath)
			}

			os.MkdirAll(filepath.Dir(virtpath), 0744) // TODO: idk why I need this here, but I do :)
			err = os.Rename(longpath, virtpath)
			if err != nil {
				return err
			}

			assets = append(assets, filepath.Join(name, relpath))
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return assets, nil
}
