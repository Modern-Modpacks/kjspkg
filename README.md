# <img alt="icon" height="32" src="https://user-images.githubusercontent.com/79367505/227798123-5454e9b1-b39b-4c45-9e02-e18f2e807585.png"> KJSPKG

A simple package manager for KubeJS written in Go.

[![Contributions welcome](https://github.com/Modern-Modpacks/kjspkg/assets/79367505/d2519e70-ce96-4bbc-b35b-af3e674bf421)](https://github.com/Modern-Modpacks/kjspkg#adding-your-own-package)
[![Approved by latvian.dev](https://img.shields.io/badge/approved%20by-lat-c374e4?style=for-the-badge&labelColor=480066)](https://github.com/user-attachments/assets/0df64919-6333-447e-9869-c270138941bd)

**The Go API exposed at github.com/Modern-Modpacks/kjspkg/pkg/kjspkg is not
stable. It is not recommended to use it just yet.**

## Installation & Update

### Install script

This script will install KJSPKG to your system and add it to your PATH.

```sh
# On Linux
curl -fsSL https://g.tizu.dev/mm.kj/install.sh?r | bash
# On Windows
powershell -c "irm https://g.tizu.dev/mm.kj/install.ps1?r | iex"
```

<!-- ### Arch Linux

KJSPKG is available in the AUR as `kjspkg-git`.

```sh
paru -S kjspkg-git
``` -->

<!-- ### winget

KJSPKG is available in the winget repository as `modernmodpacks.kjspkg`.

```sh
winget install modernmodpacks.kjspkg
``` -->

### Using Go

This requires a working Go installation of at least 1.23.2.

```sh
go install github.com/Modern-Modpacks/kjspkg/cmd/kjspkg@latest
```

## Usage

KJSPKG comes with extensive help text, so you can just run `kjspkg` to see all
the commands and options available. You may also use `--help` after any command
to get more information about it.

```sh
kjspkg install [package] [package]
kjspkg remove [package] [package]
kjspkg update [package] [package]
```

## Adding your own package

1. Create a repository containing your scripts and assets
2. [Don't forget to license your code](https://choosealicense.com/)
3. Create an empty directory and run `kjspkg dev init`
4. Do your thing and create a repository with the code
4. Fork this repo
5. Clone it
6. Add your package to `pkgs.json` file. Format it like this: `"your_package_id": "your_github_name/your_repo_name[$path/to/your/package/directory][@branch_name]",`
    * Things in [] are optional
    * Only specify the path if you have multiple packages in one repository. If you do, specify the path where the .kjspkg file is located at
    * Branch is `main` by default
7. Create a pull request
8. Wait for it to be accepted
9. profit <!-- im not even kidding, copilot wrote this -->

### KJSPKG badges

[![kjspkg-available](https://github-production-user-asset-6210df.s3.amazonaws.com/79367505/250114674-fb848719-d52e-471b-a6cf-2c0ea6729f1c.svg)](https://kjspkglookup.modernmodpacks.site/#)

```md
[![kjspkg-available](https://github-production-user-asset-6210df.s3.amazonaws.com/79367505/250114674-fb848719-d52e-471b-a6cf-2c0ea6729f1c.svg)](https://kjspkglookup.modernmodpacks.site/#{packagename})
```

## Supported versions

![Version list](https://github.com/user-attachments/assets/5a3b8e3a-bd91-456e-8443-bbffa894a38f)

Tested means that the version is confirmed to be working;

Not tested means that the version should work, but hasn't been tested. Feel free to test it yourself and let us know so we'll update the readme.

Full support means that we focus on that version;

Partial support means that the version is supported, but not as much as the fully supported ones;

No support means that the version works, but any issues that you have with it won't be fixed.

Borked means it doesn't work lmao.
