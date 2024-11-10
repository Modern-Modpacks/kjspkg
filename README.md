# KJSPKG

A simple package manager for KubeJS.

[![contributions](https://github.com/Modern-Modpacks/kjspkg/assets/79367505/d2519e70-ce96-4bbc-b35b-af3e674bf421)](https://github.com/Modern-Modpacks/kjspkg#adding-your-own-package)
[![lat](https://img.shields.io/badge/approved%20by-lat-c374e4?style=for-the-badge&labelColor=480066)](https://github.com/user-attachments/assets/0df64919-6333-447e-9869-c270138941bd)

![logo](https://user-images.githubusercontent.com/79367505/227798123-5454e9b1-b39b-4c45-9e02-e18f2e807585.png)

_Work in progress_

<!-- ## Installation & Update

### Requirements

* [Python 3.8](https://www.python.org/) (or higher)
* Pip
* [Git](https://git-scm.com/)
* [Curl](https://curl.se/) (probably pre-installed)

### Linux

```sh
curl -s https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/install.sh | sh
```

### Windows

Download [this bat file](https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/install.bat) and run it

or use [WSL](https://learn.microsoft.com/en-us/windows/wsl/install)

## Usage

Installing packages:

```sh
kjspkg install [package] [package]
kjspkg install kjspkg:[package]
kjspkg install carbon:[package] # CarbonJS compatibility (https://github.com/carbon-kjs)
kjspkg install github:[author]/[package] # External packages
```

Removing packages:

```sh
kjspkg remove [package] [package]
```

Updating packages:

```sh
kjspkg update [package] [package]
```

More info in the help page:

```sh
kjspkg help
```

## Adding your own package

1. Create a repository containing your scripts and assets
2. [Don't forget to license your code](https://choosealicense.com/)
3. Add a file to your repo named `.kjspkg` and format it like this:

    ```json
    {
        "author": "<your_name>",
        "description": "<description>",

        "versions": [<Version key numbers (see the numbers in parentheses in the table below, or use this formula: "version title = 1.(version key + 10).whatever"). Can contain multiple numbers>],
        "modloaders": [<Modloaders ("fabric"/"forge", "fabric" will for quilt as well)>. Can contain multiple modloaders],
        "dependencies": [<Package names that your package depends on, blank if none. To depend on mods add "mod:" before the mod id>],
        "incompatibilities": [<Package names that your package is incompatible with, blank if none. Incompatible mods are also supported (use the same syntax)>]
    }
    ```

4. Fork this repo
5. Clone it
6. Add your package to `pkgs.json` file. Format it like this: `"your_package_id": "your_github_name/your_repo_name[$path/to/your/package/directory][@branch_name]",`
    * Things in [] are optional
    * Only specify the path if you have multiple packages in one repository. If you do, specify the path where the .kjspkg file is located at
    * Branch is `main` by default
7. Create a pull request
8. Wait for it to be accepted

### KJSPKG badges

[![kjspkg-available](https://github-production-user-asset-6210df.s3.amazonaws.com/79367505/250114674-fb848719-d52e-471b-a6cf-2c0ea6729f1c.svg)](https://kjspkglookup.modernmodpacks.site/#)

```md
[![kjspkg-available](https://github-production-user-asset-6210df.s3.amazonaws.com/79367505/250114674-fb848719-d52e-471b-a6cf-2c0ea6729f1c.svg)](https://kjspkglookup.modernmodpacks.site/#{packagename})
```

## Supported versions

![Version list](https://github.com/user-attachments/assets/5a3b8e3a-bd91-456e-8443-bbffa894a38f)

(Thanks tizu.dev on discord for the figma template)

Tested means that the version is confirmed to be working;

Not tested means that the version should work, but hasn't been tested. Feel free to test it yourself and let us know so we'll update the readme.

Full support means that we focus on that version;

Partial support means that the version is supported, but not as much as the fully supported ones;

No support means that the version works, but any issues that you have with it won't be fixed.

Borked means it doesn't work lmao. -->
