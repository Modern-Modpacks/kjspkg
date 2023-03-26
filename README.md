# KJSPKG

A simple package manager for KubeJS.

![badge](https://img.shields.io/badge/contributions-welcome-green)

![logo](https://user-images.githubusercontent.com/79367505/227797641-6ce7d3d6-235b-4668-adfb-bef3c02a91d9.png)

## Installation

### Linux

```sh
curl -s https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/install.sh | sh
```

### Windows

Use [WSL](https://learn.microsoft.com/en-us/windows/wsl/install)

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
        
        "versions": [<Version key numbers (2 for 1.12, 6 for 1.16, 8 for 1.18, 9 for 1.19). Can contain multiple numbers>],
        "modloaders": [<Modloaders ("fabric"/"forge", "fabric" will for for quilt as well)>. Can contain multiple modloaders]
    }
    ```

6. Create a pull request
7. Wait for it to be accepted
