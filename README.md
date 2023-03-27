# KJSPKG

A simple package manager for KubeJS.

![badge](https://img.shields.io/badge/contributions-welcome-c374e4?style=for-the-badge&labelColor=480066&logo=hackthebox&logoColor=white)

![logo](https://user-images.githubusercontent.com/79367505/227798123-5454e9b1-b39b-4c45-9e02-e18f2e807585.png)

## Installation & Update

### Linux

```sh
curl -s https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/install.sh | sh
```

### Windows

The only way to install this currently is through [WSL](https://learn.microsoft.com/en-us/windows/wsl/install)

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
    **Your "assets" folder and your "data" folder should not contain modid sub-directories, everything in them will be saved in assets/<package_name> and data/<package_name>. Please design your scripts with that in mind.**
2. Fork this repo
3. Clone it
4. Add a "<package_name>.json" file to the "pkgs" directory
5. Format it like this:

    ```json
    {
        "repo": "yourname/reponame",

        "author": "<your_name>",
        "description": "Description",
        
        "versions": [<Version key numbers (see the numbers in parentheses in the table below). Can contain multiple numbers>],
        "modloaders": [<Modloaders ("fabric"/"forge", "fabric" will for for quilt as well)>. Can contain multiple modloaders]
    }
    ```

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
