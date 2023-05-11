#!/usr/bin/python3

# IMPORTS

# Built-in modules
from os import path, remove, getcwd, makedirs, walk, chmod, system, environ # Working with files
from shutil import rmtree, move, copy # More file stuff
from pathlib import Path # EVEN MORE FILE STUFF
from tempfile import gettempdir # Get tmp dir of current os
from json import dump, load, dumps # Json
from stat import S_IWRITE # Windows stuff
from warnings import filterwarnings # Disable the dumb fuzz warning

# External libraries
from fire import Fire # CLI tool

from requests import get, exceptions # Requests
from git import Repo, GitCommandNotFound, GitCommandError # Git cloning

# CONSTANTS
VERSIONS = { # Version and version keys
    "1.12.2": 2,
    "1.12": 2,

    "1.16.5": 6,
    "1.16": 6,

    "1.18.2": 8,
    "1.18": 8,

    "1.19.2": 9,
    "1.19.3": 9,
    "1.19.4": 9,
    "1.19": 9
}
SCRIPT_DIRS = ("server_scripts", "client_scripts", "startup_scripts") # Script directories
ASSET_DIRS = ("data", "assets") # Asset directories

# VARIABLES
kjspkgfile = {} # .kjspkg file

# HELPER FUNCTIONS
def _bold(s:str) -> str: return "\u001b[1m"+s+"\u001b[0m" # Make the text bold
def _err(err:str): # Handle errors
    print("\u001b[31;1m"+err+"\u001b[0m") # Print error
    exit(1) # Quit
def _remove_prefix(pkgname:str) -> str: return pkgname.split(":")[-1] # Remove prefix
def _format_github(pkgname:str) -> str: return _remove_prefix(pkgname).split('/')[-1].split("@")[0] # Remove github author and branch
def _dumbass_windows_path_error(f, p:str, e): chmod(p, S_IWRITE) # Dumbass windows path error

# TMP HELPER FUNCTIONS
def _create_tmp(pathintmp:str) -> str: # Create a temp directory and return its path
    tmppath = path.join("tmp", pathintmp)
    makedirs(tmppath, exist_ok=True)
    return tmppath
def _clear_tmp(): # Clear tmp directory
    if path.exists("tmp"): rmtree("tmp", onerror=_dumbass_windows_path_error) # Delete if exists

# PROJECT HELPER FUNCTIONS
def _check_project() -> bool: # Check if the current directory is a kubejs directory
    for dir in SCRIPT_DIRS:
        if path.exists(dir) and path.basename(getcwd())=="kubejs": return True
    return False
def _create_project_directories():
    for dir in SCRIPT_DIRS+ASSET_DIRS: makedirs(dir, exist_ok=True)  # Create asset and script directories
    for dir in SCRIPT_DIRS: makedirs(path.join(dir, ".kjspkg"), exist_ok=True)  # Create .kjspkg directories
def _project_exists() -> bool: return path.exists(".kjspkg") # Check if a kjspkg project exists
def _delete_project(): # Delete the project and all of the files
    for pkg in list(kjspkgfile["installed"].keys()): _remove_pkg(pkg, True) # Remove all packages
    for dir in SCRIPT_DIRS: rmtree(path.join(dir, ".kjspkg"), onerror=_dumbass_windows_path_error) # Remove .kjspkg dirs
    remove(".kjspkg") # Remove .kjspkg file
def _enable_reflection(): # Enable reflection on 1.16
    with open(path.join("config", "common.properties"), "a+") as f:
        if ("invertClassLoader=true" not in f.read().splitlines()): f.write("invertClassLoader=true")

# PKG HELPER FUNCTIONS
def _reload_pkgs(): # Reload package registry cache
    with open(path.join(gettempdir(), "kjspkgs.json"), "w") as tmp: tmp.write(get(f"https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/pkgs.json").text)
def _pkgs_json() -> dict: # Get the pkgs.json file
    pkgspath = path.join(gettempdir(), "kjspkgs.json")
    if not path.exists(pkgspath): _reload_pkgs()
    return load(open(pkgspath))

def _pkg_info(pkg:str, getlicense:bool=False) -> dict: # Get info about the pkg
    prefix = pkg.split(":")[0] if len(pkg.split(":"))>1 else "kjspkg" # kjspkg: - default prefix
    packagename = _remove_prefix(pkg) # Get package name without the prefix

    # Call correct function based on prefix
    if prefix=="kjspkg": info = _kjspkginfo(packagename)
    elif prefix in ("carbon", "carbonjs"): info = _carbonpkginfo(packagename)
    elif prefix=="github": info = _githubpkginfo(packagename)
    else: _err("Unknown prefix: "+_bold(prefix))

    if getlicense: # If the license is requested
        pkglicense = get(f"https://api.github.com/repos/{info['repo']}/license?ref={info['branch']}") # Get license
        if pkglicense.status_code!=403: info["license"] = pkglicense.json()["license"]["spdx_id"] if pkglicense.status_code!=404 else "All Rights Reserved" # Add the license to info
        if ("license" not in info.keys() or info["license"]=="NOASSERTION"): info["license"] = f"Other ({pkglicense.json()['html_url'] if pkglicense.status_code!=403 else 'https://github.com/Modern-Modpacks/kjspkg/blob/'+info['branch']+'/LICENSE'})" # Custom licenses or api rate exceeded

    return info
def _kjspkginfo(pkg:str) -> dict: # Get info about a default kjspkg pkg
    pkgregistry = _pkgs_json() # Get the pkgs.json file
    if pkg not in pkgregistry.keys(): return # Return nothing if the pkg doesn't exist

    repo = pkgregistry[pkg] # Get the repo with branch
    branch = "main" # Set the branch to main by default
    if "@" in repo: # If the branch is specifed
        branch = repo.split("@")[-1] # Set the branch
        repo = repo.split("@")[0] # Remove the branch from the repo

    package = get(f"https://raw.githubusercontent.com/{repo}/{branch}/.kjspkg").json() # Get package info

    package["repo"] = repo # Add the repo to info
    package["branch"] = branch # Add the branch to info

    return package # Return the json object
def _carbonpkginfo(pkg:str) -> dict: # Get info about a carbonjs pkg (https://github.com/carbon-kjs)
    request = get(f"https://raw.githubusercontent.com/carbon-kjs/{pkg}/main/carbon.config.json") # Request info about the package
    if request.status_code==404: return # Return nothing if the pkg doesn't exist

    json = request.json() # Get info in json

    return { # Return formatted info
        "author": json["author"],
        "description": json["description"],
        
        "versions": [VERSIONS[json["minecraftVersion"]]],
        "modloaders": json["modloaders"],
        "dependencies": [],
        "incompatibilities": [],

        "repo": f"carbon-kjs/{pkg}",
        "branch": "main"
    }
def _githubpkginfo(pkg:str) -> dict: # Get dummy info about an external pkg
    return {
        "author": pkg.split("/")[0],
        "description": "",
        
        "versions": [kjspkgfile["version"]],
        "modloaders": [kjspkgfile["modloader"]],
        "dependencies": [],
        "incompatibilities": [],

        "repo": pkg,
        "branch": "main" if "@" not in pkg else pkg.split("@")[-1]
    }
def _install_pkg(pkg:str, update:bool, skipmissing:bool, noreload:bool): # Install the pkg
    if not update and _format_github(pkg) in kjspkgfile["installed"]: return # If the pkg is already installed and the update parameter is false, do nothing
    if update: 
        if pkg=="*": # If updating all packages
            for p in list(kjspkgfile["installed"].keys()): _install_pkg(p, True, skipmissing) # Update all packages
            return

        _remove_pkg(pkg, False) # If update is true, remove the previous version of the pkg

    package = _pkg_info(pkg) # Get pkg
    if not package and reload: 
        _reload_pkgs() # Reload if not found
        package = _pkg_info(pkg) # Try to get the pkg again

    if not package: # If pkg doesn't exist after reload
        if not skipmissing: _err(f"Package \"{pkg}\" does not exist") # Err
        else: return # Or just ingore

    pkg = _remove_prefix(pkg) # Remove pkg prefix

    # Unsupported version/modloader errs
    if kjspkgfile["version"] not in package["versions"]: _err(f"Unsupported version 1.{10+kjspkgfile['version']} for package \"{pkg}\"")
    if kjspkgfile["modloader"] not in package["modloaders"]: _err(f"Unsupported modloader \"{kjspkgfile['modloader'].title()}\" for package \"{pkg}\"")

    # Install dependencies & check for incompats
    if "dependencies" in package.keys():
        for dep in package["dependencies"]: _install_pkg(dep.lower(), False)
    if "incompatibilities" in package.keys(): 
        for i in kjspkgfile["installed"].keys(): 
            if i in package["incompatibilities"]: _err(f"Incompatible package: "+i) # Throw err if incompats detected

    tmpdir = _create_tmp(pkg) # Create a temp dir
    try: Repo.clone_from(f"https://github.com/{package['repo']}.git", tmpdir, branch=package["branch"]) # Install the repo into the tmp dir
    except GitCommandError: tryRepo.clone_from(f"https://github.com/{package['repo']}.git", tmpdir) # If the branch is not found, try to install from the default one

    pkg = _format_github(pkg) # Remove github author and branch if present

    licensefile = path.join(tmpdir, "LICENSE") 
    if not path.exists(licensefile): licensefile = path.join(tmpdir, "LICENSE.txt")
    for dir in SCRIPT_DIRS: # Clone scripts & licenses into the main kjs folders
        tmppkgpath = path.join(tmpdir, dir)
        finalpkgpath = path.join(dir, ".kjspkg", pkg)
        if path.exists(tmppkgpath):
            move(tmppkgpath, finalpkgpath) # Files
            if path.exists(licensefile): copy(licensefile, finalpkgpath) # License

    assetfiles = [] # Pkg's asset files
    for dir in ASSET_DIRS: # Clone assets
        tmppkgpath = path.join(tmpdir, dir) # Get asset path
        if not path.exists(tmppkgpath): continue

        for dirpath, _, files in walk(tmppkgpath): # For each file in assets/data
            for name in files:
                tmppath = path.join(dirpath, name)
                finalpath = path.sep.join(tmppath.split(path.sep)[2:])

                makedirs(path.sep.join(finalpath.split(path.sep)[:-1]), exist_ok=True) # Create parent dirs
                move(tmppath, finalpath) # Move it to the permanent dir
                assetfiles.append(finalpath) # Add it to assetfiles

    kjspkgfile["installed"][pkg] = assetfiles # Add the pkg to installed
def _remove_pkg(pkg:str, skipmissing:bool): # Remove the pkg
    if pkg not in kjspkgfile["installed"].keys():
        if not skipmissing: _err(f"Package \"{pkg}\" is not installed") # If the pkg is not installed, err
        else: return # Or just ignore

    for dir in SCRIPT_DIRS: # Remove script files
        scriptpath = path.join(dir, ".kjspkg", pkg)
        if path.exists(scriptpath): rmtree(scriptpath, onerror=_dumbass_windows_path_error)
    for file in kjspkgfile["installed"][pkg]: # Remove asset files
        if path.exists(file): remove(file)

    del kjspkgfile["installed"][pkg] # Remove the pkg from installed

# COMMAND FUNCTIONS
def install(*pkgs:str, update:bool=False, quiet:bool=False, skipmissing:bool=False, reload:bool=True): # Install pkgs
    if update and not pkgs: pkgs = ("*",) # If update is passed and not pkgs are, set them to all pkgs

    # Install each given package
    for pkg in pkgs:
        pkg = pkg.lower() # Format

        if (
            pkg.startswith("github:") and ( # If the package is external
                ("trustgithub" in kjspkgfile.keys() and kjspkgfile["trustgithub"]) or # And external packages are trusted
                (quiet or input(_bold(f"Package \"{_format_github(pkg)}\" uses the \"github:\" prefix, external packages coming from github are not tested and are not guaranteed to work, do you want to disable external packages in this project? (Y/n) ")).lower()!="n") # Or quiet mode is set/prompt is answered with y
            )
        ): kjspkgfile["trustgithub"]=True # Trust external packages
        else: # Otherwise
            kjspkgfile["trustgithub"]=False # Don't trust them
            if skipmissing: continue # Skip if skipmissing
            else: return # Stop if no skipmissing

        if update and pkg not in kjspkgfile["installed"].keys() and not skipmissing and pkg!="*": _err(f"Package \"{_format_github(pkg)}\" not found") # Err if package not found during update
        _install_pkg(pkg, update, skipmissing, reload) # Install package
        if not quiet and pkg=="*": print(_bold(f"All packages updated succesfully!")) # Show message if all packages are updated
        elif not quiet: print(_bold(f"Package \"{_format_github(pkg)}\" {'installed' if not update else 'updated'} succesfully!")) # Show message if one package is installed/updated
def removepkg(*pkgs:str, quiet:bool=False, skipmissing:bool=False): # Remove pkgs
    for pkg in pkgs:
        pkg = pkg.lower()
        _remove_pkg(pkg, skipmissing)
        if not quiet: print(_bold(f"Package \"{pkg}\" removed succesfully!"))
def update(*pkgs:str, **kwargs): # Update pkgs
    install(*pkgs, update=True, **kwargs)
def updateall(**kwargs): # Update all pkgs
    update("*", **kwargs)
def listpkgs(*, count:bool=False): # List pkgs
    if count: # Only show the pkg count if the "count" option is passed
        print(len(kjspkgfile["installed"].keys()))
        return

    if (len(kjspkgfile["installed"].keys())==0): # Easter egg if 0 pkg installed
        print(_bold("*nothing here, noone around to help*"))
        return

    print("\n".join(kjspkgfile["installed"].keys())) # Print the list
def pkginfo(pkg:str, *, script:bool=False): # Print info about a pkg
    if(pkg.startswith("github:")): _err(f"Can't show info about an external package \"{_format_github(pkg)}\"") # Err if pkg is external

    info = _pkg_info(pkg, True) # Get the info
    if not info: _err(f"Package {pkg} not found") # Err if pkg not found

    # Print it (scripty)
    if script:
        print(dumps(info))
        return

    # Print it (pretty)
    print(f"""
{_bold(_remove_prefix(pkg).replace("-", " ").title())} by {_bold(info["author"])}

{info["description"]}

{_bold("Dependencies")}: {", ".join([i.title() for i in info["dependencies"]]) if "dependencies" in info.keys() and len(info["dependencies"])>0 else "*nothing here*"}
{_bold("Incompatibilities")}: {", ".join([i.title() for i in info["incompatibilities"]]) if "incompatibilities" in info.keys() and len(info["incompatibilities"])>0 else "*compatible with everything!*"}

{_bold("License")}: {info["license"]}
{_bold("GitHub")}: https://github.com/{info["repo"]}/tree/{info["branch"]}

{_bold("Versions")}: {", ".join([f"1.{10+i}" for i in info["versions"]])}
{_bold("Modloaders")}: {", ".join([i.title() for i in info["modloaders"]])}
    """)
def listall(*, count:bool=False, search:str="", reload:bool=True): # List all pkgs
    if reload: _reload_pkgs() # Reload pkgs
    allpkgs = list(_pkgs_json().keys()) # All package names

    if count: # If count is true
        print(len(allpkgs)) # Print the pkg count
        return
    if not search: # If no search query
        print("\n".join(allpkgs)) # Print all pkg names
        return

    filterwarnings("ignore") # Ignore the warning thefuzz produces
    from thefuzz import process # Fuzzy search

    # Get results and print the best ones
    results = process.extract(search, allpkgs)
    for result, ratio in results:
        if ratio>75: print(result)
def search(*query:str, **kwags): # Search for pkgs
    listall(search="-".join(query), **kwags) # Call listall with joined spaces
def reload(): _reload_pkgs() # Reload packages
def init(*, version:str=None, modloader:str=None, trustgithub:bool=False, quiet:bool=False, override:bool=False, cancreate:str=None): # Init project
    global kjspkgfile

    if cancreate: # Scriptable cancreate option
        print(_check_project())
        return

    if not _check_project(): _err("Hmm... This directory doesn't look like a kubejs directory") # Wrong dir err

    if _project_exists(): # Override
        if not quiet and input("\u001b[31;1mA PROJECT ALREADY EXISTS IN THIS REPOSITORY, CREATING A NEW ONE OVERRIDES THE PREVIOUS ONE, ARE YOU SURE YOU WANT TO PROCEED? (y/N): \u001b[0m").lower()=="y" or override: _delete_project()
        else: exit(0)

    # Ask for missing params
    if not version: version = input(_bold("Input your minecraft version (1.12/1.16/1.18/1.19): ")) # Version
    if version not in VERSIONS.keys(): _err("Unknown or unsupported version: "+str(version))
    version = VERSIONS[version]

    if not modloader: modloader = input(_bold("Input your modloader (forge/fabric/quilt): ")) # Modloader
    modloader = modloader.lower()
    if modloader not in ("forge", "fabric", "quilt"): _err("Unknown or unsupported modloader: "+modloader.title())

    _create_project_directories() # Create .kjspkg directories
    if version==6: _enable_reflection() # Enable reflection in the config for 1.16.5

    kjspkgfile = {
        "version": version,
        "modloader": modloader if modloader!="quilt" else "fabric",
        "installed": {},
        "trustgithub": trustgithub
    }
    with open(".kjspkg", "w+") as f: dump(kjspkgfile, f) # Create .kjspkg file
    if not quiet: print(_bold("Project created!")) # Woo!
def uninit(*, confirm:bool=False): # Remove the project
    if confirm or input("\u001b[31;1mDOING THIS WILL REMOVE ALL PACKAGES AND UNINSTALL KJSPKG COMPLETELY, ARE YOU SURE YOU WANT TO PROCEED? (y/N): \u001b[0m").lower()=="y": _delete_project() 
def info(): # Print the help page
    INFO = f"""
{_bold("Commands:")}

kjspkg install/download [pkgname1] [pkgname2] [--quiet/--skipmissing] [--update] [--noreload] - installs packages
kjspkg remove/uninstall [pkgname1] [pkgname2] [--quiet/--skipmissing] - removes packages
kjspkg update [pkgname1/*] [pkgname2] [--quiet/--skipmissing] - updates packages
kjspkg updateall [--quiet/--skipmissing] - updates all packages

kjspkg install [pkgname] - installs packages from kjspkg's repo
kjspkg install kjspkg:[pkgname] - installs packages from kjspkg's repo
kjspkg install carbon:[pkgname] - installs packages from carbonjs' repo (https://github.com/carbon-kjs)

kjspkg list [--count] - lists packages (or outputs the count of them)
kjspkg pkg [package] [--script] - shows info about the package
kjspkg listall/all [--count] [--search "<query>"] [--noreload] - lists all packages
kjspkg search [query] - searches for packages with a similar name
kjspkg reload/refresh - reloads the cached package registry

kjspkg init [--override/--quiet] [--version "<version>"] [--modloader "<modloader>"] [--cancreate "<path>"] - inits a new project (will be run by default)
kjspkg uninit [--confirm] - removes all packages and the project

kjspkg help/info - shows this message (default behavior)
kjspkg gui - shows info about the GUI app

{_bold("Credits:")}

Modern Modpacks - owner
G_cat101 - coder
malezjaa - creator of CarbonJS
Juh9870 - Wanted to be here
    """

    print(INFO)
def guiinfo(): # Print info about the GUI app
    print(f"{_bold('Did you know there is a GUI app for KJSPKG?')} Check it out at https://github.com/Modern-Modpacks/kjspkg-gui!")
    
# PARSER FUNCTION
def _parser(func:str="help", *args, help:bool=False, **kwargs):
    global kjspkgfile

    system("") # Enable color codes on windows, don't ask me why

    if help: func="help"

    FUNCTIONS = { # Command mappings
        "install": install,
        "download": install,
        "remove": removepkg,
        "uninstall": removepkg,
        "rm": removepkg,
        "update": update,
        "updateall": updateall,
        "list": listpkgs,
        "pkg": pkginfo,
        "listall": listall,
        "all": listall,
        "search": search,
        "reload": reload,
        "refresh": reload,
        "init": init,
        "uninit": uninit,
        "help": info,
        "info": info,
        "gui": guiinfo
    }

    if func not in FUNCTIONS.keys(): _err("Command \""+func+"\" is not found. Run \"kjspkg help\" to see all of the available commands") # Wrong command err
    
    helperfuncs = (info, guiinfo, init, pkginfo, listall, search) # Helper commands that don't require a project
    if FUNCTIONS[func] not in helperfuncs and not _project_exists(): # If a project is not found, call init
        print(_bold("Project not found, a new one will be created.\n"))
        init()
    if (FUNCTIONS[func]==init or FUNCTIONS[func] not in helperfuncs) and path.exists(".kjspkg"): kjspkgfile = load(open(".kjspkg")) # Open .kjspkg

    FUNCTIONS[func](*args, **kwargs) # Run the command

    # Clean up
    if path.exists(".kjspkg") and FUNCTIONS[func] not in helperfuncs: # If uninit wasn't called and the command isn't a helper command
        with open(".kjspkg", "w") as f: dump(kjspkgfile, f) # Save .kjspkg

# RUN
if __name__=="__main__": # If not imported
    _clear_tmp() # Remove tmp

    try: Fire(_parser) # Run parser with fire
    except (KeyboardInterrupt, EOFError): exit(0) # Ignore some exceptions
    except TypeError: _err("Wrong syntax") # Wrong syntax err
    except GitCommandNotFound: _err("Git not found. Install it here: https://git-scm.com/downloads") # Git not found err
    except (exceptions.ConnectionError, exceptions.ReadTimeout): _err("Low internet connection") # Low internet connection err

    _clear_tmp() # Remove tmp again

# Ok that's it bye