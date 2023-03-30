# KJSPKG

A simple package manager for KubeJS.

![badge](https://img.shields.io/badge/contributions-welcome-c374e4?style=for-the-badge&labelColor=480066&logo=hackthebox&logoColor=white)

![logo](https://user-images.githubusercontent.com/79367505/227798123-5454e9b1-b39b-4c45-9e02-e18f2e807585.png)

## Installation & Update

### Requirements

* [Python 3.8](https://www.python.org/) (or higher)
* Pip
* [Git](https://git-scm.com/)
* [Curl](https://curl.se/) (probably pre-installed)

### Linux

```sh
curl -s https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/install.sh | sudo sh
```

### Windows

Download [this bat file](https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/install.bat) and run it

or use [WSL](https://learn.microsoft.com/en-us/windows/wsl/install)

## Usage

Installing packages:

```sh
kjspkg install [package] [package]
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
2. Add a file to your repo named `.kjspkg` and format it like this:

    ```json
    {
        "author": "<your_name>",
        "description": "<description>",
        
        "versions": [<Version key numbers (see the numbers in parentheses in the table below). Can contain multiple numbers>],
        "modloaders": [<Modloaders ("fabric"/"forge", "fabric" will for for quilt as well)>. Can contain multiple modloaders],
        "dependencies": [<package names that your package depends on, blank if none>],
        "incompatibilities": [<package names that your package is incompatible with, blank if none>]
    }
    ```

3. Fork this repo
4. Clone it
5. Add your package to `pkgs.json` file. Format it like this: `"your_package_id": "your_github_name/your_repo_name",`
    * You can specify the branch by adding `@branch_name` at the end of the string, otherwise it will automatically use the `main` branch
6. Create a pull request
7. Wait for it to be accepted

## Supported versions

| |Forge|Fabric/Quilt|
|-|-----|------------|
|1.19 (9)|![tested_partial](https://img.shields.io/badge/Tested-Partial%20Support-green)|![nottested_partial](https://img.shields.io/badge/Not%20tested-Partial%20Support-green)|
|1.18 (8)|![nottested_partial](https://img.shields.io/badge/Not%20Tested-Partial%20Support-green)|![nottested_partial](https://img.shields.io/badge/Not%20tested-Partial%20Support-green)|
|1.16 (6)|![tested_full](https://img.shields.io/badge/Tested-Full%20Support-brightgreen)|![nottested_partial](https://img.shields.io/badge/Not%20tested-Partial%20Support-green)|
|1.12 (2)|![nottested_no](https://img.shields.io/badge/Not%20Tested-No%20Support-yellow)||
|1.7 (🧌)|![borked](https://img.shields.io/badge/Borked-red)||

Tested means that the version is confirmed to be working;

Not tested means that the version should work, but hasn't been tested. Feel free to test it yourself and let us know so we'll update the readme.

Full support means that we focus on that version;

Partial support means that the version is supported, but not as much as the fully supported ones;

No support means that the version works, but any issues that you have with it won't be fixed.

Borked means it doesn't work lmao.
