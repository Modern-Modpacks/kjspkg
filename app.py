#!/usr/bin/python3

# IMPORTS

# Built-in modules
from os import path, remove, getcwd, makedirs, walk, chmod, system, environ, listdir # Working with files
from shutil import rmtree, move, copy # More file stuff
from pathlib import Path # EVEN MORE FILE STUFF
from tempfile import gettempdir # Get tmp dir of current os
from json import dump, load, dumps, loads # Json
from zipfile import ZipFile # Working with .jars

from http import server # Discord login stuff
from urllib.parse import urlparse, parse_qs # Parse url path

from stat import S_IWRITE # Windows stuff
from warnings import filterwarnings # Disable the dumb fuzz warning

from random import choice # Random splash

# External libraries
from fire import Fire # CLI tool

from requests import get, post, exceptions # Requests
from git import Repo, GitCommandNotFound, GitCommandError # Git cloning

from toml import loads as tomlload # Read .tomls

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
CONFIG = { # Default config
    "_": "Please, do not delete this file or any other file/directory labeled \".kjspkg\", even if they might seem empty and useless, they are still required. Thanks for your understanding!",

    "installed": {},
    "trustgithub": False
}

# VARIABLES
kjspkgfile = {} # .kjspkg file

# CLASSES
# class HTTPDiscordLoginRequestHandler(server.BaseHTTPRequestHandler):
#     def log_message(self, format, *args): return
#     def do_GET(self):
#         self.send_response(200)
#         self.send_header("Content-type", "text/html")
#         self.end_headers()
#         self.wfile.write(
# b"""<html>
#     <head>
#         <script>
#             close();
#         </script>
#     </head>
# </html>"""
#         )

#         token = post("https://discord.com/api/v10/oauth2/token", data={
#             "client_id": "1108295881247166496",
#             "client_secret": "yopSwUCpisfib6RR3aqHyp293vOyG4F5",
#             "grant_type": "authorization_code",
#             "code": parse_qs(urlparse(self.path).query)["code"][0],
#             "redirect_uri": "http://localhost:1337"
#         }).json()["access_token"]

#         print(post("https://discord.com/api/v10/channels/303440391124942858/messages", headers={
#             "Authorization": token
#         }).text)

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
def _update_manifest(): # Update .kjspkg file
    global kjspkgfile

    # Update the config by adding keys that don't exist
    for k, v in CONFIG.items():
        if k not in kjspkgfile.keys(): kjspkgfile[k] = v
def _enable_reflection(): # Enable reflection on 1.16
    with open(path.join("config", "common.properties"), "a+") as f:
        if ("invertClassLoader=true" not in f.read().splitlines()): f.write("invertClassLoader=true")
def _get_modid(modpath:str) -> str: # Get mod id from a mod file
    modfile = ZipFile(path.join(getcwd(), "..", "mods", modpath))

    if kjspkgfile["modloader"]=="forge": return tomlload(modfile.open(path.join("META-INF", "mods.toml")).read().decode("utf-8"))["mods"][0]["modId"]
    else: return loads(modfile.open(modfile.open("fabric.mod.json").read().decode("utf-8")))["id"]
def _get_modids() -> list: # Get all mod ids
    modids = []
    for i in listdir(path.join(getcwd(), "..", "mods")):
        if i.endswith(".jar"): modids.append(_get_modid(i)) # For each mod file, get the mod id and append

    return modids # Return the list of modids

# def _discord_login(): # Login with discord for discord prefixes
#     server.HTTPServer(("", 1337), HTTPDiscordLoginRequestHandler).handle_request()

# PKG HELPER FUNCTIONS
def _reload_pkgs(): # Reload package registry cache
    with open(path.join(gettempdir(), "kjspkgs.json"), "w") as tmp: tmp.write(get(f"https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/pkgs.json").text)
def _pkgs_json() -> dict: # Get the pkgs.json file
    pkgspath = path.join(gettempdir(), "kjspkgs.json")
    if not path.exists(pkgspath): _reload_pkgs()
    return load(open(pkgspath))

def _pkg_info(pkg:str, getlicense:bool=False, refresh:bool=True) -> dict: # Get info about the pkg
    if refresh: _reload_pkgs() # Refresh pkgs

    prefix = pkg.split(":")[0] if len(pkg.split(":"))>1 else "kjspkg" # kjspkg: - default prefix
    packagename = _remove_prefix(pkg) # Get package name without the prefix

    # Call correct function based on prefix
    if prefix=="kjspkg": info = _kjspkginfo(packagename)
    elif prefix in ("carbon", "carbonjs"): info = _carbonpkginfo(packagename)
    elif prefix in ("github", "external"): info = _githubpkginfo(packagename)
    # elif prefix=="discord": 
    #     _discord_login()
    #     exit()
    else: _err("Unknown prefix: "+_bold(prefix))

    if info and getlicense: # If the license is requested
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
    modids = []
    if (("dependencies" in package.keys() and any([i.startswith("mod:") for i in package["dependencies"]])) or ("incompatibilities" in package.keys() and any([i.startswith("mod:") for i in package["incompatibilities"]]))): modids = _get_modids() # Get a list of all mod ids

    if "dependencies" in package.keys():
        for dep in package["dependencies"]: 
            if dep.lower().startswith("mod:") and _remove_prefix(dep.lower()) not in modids: _err(f"Mod \"{_remove_prefix(dep.replace('_', ' ').replace('-', ' ')).title()}\" not found.") # Check for mod dependency
            elif not dep.lower().startswith("mod:"): _install_pkg(dep.lower(), False, skipmissing, noreload) # Install package dependency
    if "incompatibilities" in package.keys(): 
        for i in package["incompatibilities"]:
            if i.lower().startswith("mod:") and _remove_prefix(i.lower()) in modids: _err(f"Incompatible mod: "+_remove_prefix(i.replace('_', ' ').replace('-', ' ')).title()) # Check for mod incompats
            elif i in kjspkgfile["installed"].keys(): _err(f"Incompatible package: "+i) # Throw err if incompats detected

    tmpdir = _create_tmp(pkg) # Create a temp dir
    try: Repo.clone_from(f"https://github.com/{package['repo']}.git", tmpdir, branch=package["branch"]) # Install the repo into the tmp dir
    except GitCommandError: Repo.clone_from(f"https://github.com/{package['repo']}.git", tmpdir) # If the branch is not found, try to install from the default one

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

        if (pkg.startswith("github:")): # If the package is external
            if (
                ("trustgithub" in kjspkgfile.keys() and kjspkgfile["trustgithub"]) or # And external packages are trusted
                (quiet or input(_bold(f"Package \"{_format_github(pkg)}\" uses the \"github:\" prefix, external packages coming from github are not tested and are not guaranteed to work, do you want to disable external packages in this project? (Y/n) ")).lower()!="n") # Or quiet mode is set/prompt is answered with y
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

    nl = "\n" # Bruh

    # Print it (pretty)
    print(f"""
{_bold(_remove_prefix(pkg).replace("-", " ").title())} by {_bold(info["author"])}

{info["description"]}

{_bold("Dependencies")}: {", ".join([_remove_prefix(i).title().replace("-", " ").replace("_", " ")+(" ("+i.split(":")[0].title()+")" if ":" in i else "") for i in info["dependencies"]]) if "dependencies" in info.keys() and len(info["dependencies"])>0 else "*nothing here*"}
{_bold("Incompatibilities")}: {", ".join([_remove_prefix(i).title().replace("-", " ").replace("_", " ")+(" ("+i.split(":")[0].title()+")" if ":" in i else "") for i in info["incompatibilities"]]) if "incompatibilities" in info.keys() and len(info["incompatibilities"])>0 else "*compatible with everything!*"}

{_bold("License")}: {info["license"]}
{_bold("GitHub")}: https://github.com/{info["repo"]}/tree/{info["branch"]}

{_bold("Versions")}: {", ".join([f"1.{10+i}" for i in info["versions"]])}
{_bold("Modloaders")}: {", ".join([i.title() for i in info["modloaders"]])}
{
    nl+_bold("KJSPKG Lookup")+": https://kjspkglookup.modernmodpacks.site/#amogus"+nl if ":" not in pkg or pkg.split(":")[0]=="kjspkg" else ""
}""")
def listall(*, count:bool=False, search:str="", reload:bool=True, carbon:bool=False): # List all pkgs
    if not carbon:
        if reload: _reload_pkgs() # Reload pkgs
        allpkgs = list(_pkgs_json().keys()) # All package names
    else: allpkgs = [i["name"] for i in get("https://api.github.com/orgs/carbon-kjs/repos").json()]

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
def init(*, quiet:bool=False, override:bool=False, cancreate:str=None, **configargs): # Init project
    global kjspkgfile

    params = [ # Changable parameters
        "version",
        "modloader",
        "trustgithub"
    ]
    for arg in configargs.keys(): # If a parameter is not found, raise key err
        if arg not in params: raise TypeError()

    if cancreate: # Scriptable cancreate option
        print(_check_project())
        return

    if not _check_project(): _err("Hmm... This directory doesn't look like a kubejs directory") # Wrong dir err

    if _project_exists() and not override: # Override
        if not quiet and input("\u001b[31;1mA PROJECT ALREADY EXISTS IN THIS REPOSITORY, CREATING A NEW ONE OVERRIDES THE PREVIOUS ONE, ARE YOU SURE YOU WANT TO PROCEED? (y/N): \u001b[0m").lower()=="y" or override: _delete_project()
        else: exit(0)

    # Ask for missing params
    if "version" not in configargs.keys(): configargs["version"] = input(_bold("Input your minecraft version (1.12/1.16/1.18/1.19): ")) # Version
    if configargs["version"] not in VERSIONS.keys(): _err("Unknown or unsupported version: "+str(configargs["version"]))
    configargs["version"] = VERSIONS[configargs["version"]]

    if "modloader" not in configargs.keys(): configargs["modloader"] = input(_bold("Input your modloader (forge/fabric/quilt): ")) # Modloader
    configargs["modloader"] = configargs["modloader"].lower()
    if configargs["modloader"] not in ("forge", "fabric", "quilt"): _err("Unknown or unsupported modloader: "+configargs["modloader"].title())

    _create_project_directories() # Create .kjspkg directories
    if configargs["version"]==6: _enable_reflection() # Enable reflection in the config for 1.16.5

    kjspkgfile = CONFIG # Set .kjspkg to default config
    for k, v in configargs.items(): kjspkgfile[k] = v # Change config as needed

    with open(".kjspkg", "w+") as f: dump(kjspkgfile, f) # Create .kjspkg file
    if not quiet: print(_bold("Project created!")) # Woo!
def uninit(*, confirm:bool=False): # Remove the project
    if confirm or input("\u001b[31;1mDOING THIS WILL REMOVE ALL PACKAGES AND UNINSTALL KJSPKG COMPLETELY, ARE YOU SURE YOU WANT TO PROCEED? (y/N): \u001b[0m").lower()=="y": _delete_project() 
def info(): # Print the help page
    SPLASHES = [ # Splash list
        "You should run `kjspkg uninit`, NOW!",
        "Run `kjspkg mold` to brew kombucha",
        "Thanks Lat 👍",
        "Help, I'm locked in a basement packaging scripts!",
        "kjspkg rm -rf / --no-preserve-root",
        "Made in Python 3.whatever!",
        "Also try CarbonJS!",
        "https://modernmodpacks.site",
        "Made by Modern Modpacks!",
        "gimme gimme gimme",
        "`amogus` is a real package!",
        "Supports 1.12!",
        "Procrastinating doing one project by doing another project, genius!",
        "Also try Magna!"
    ]

    # Info string
    INFO = f"""
{_bold("KJSPKG")}, a package manager for KubeJS.
{choice(SPLASHES)}

{_bold("Commands:")}

kjspkg install/download [pkgname1] [pkgname2] [--quiet/--skipmissing] [--update] [--noreload] - installs packages
kjspkg remove/uninstall [pkgname1] [pkgname2] [--quiet/--skipmissing] - removes packages
kjspkg update [pkgname1/*] [pkgname2] [--quiet/--skipmissing] - updates packages
kjspkg updateall [--quiet/--skipmissing] [--carbon] - updates all packages

kjspkg install [pkgname] - installs packages from kjspkg's repo
kjspkg install kjspkg:[pkgname] - installs packages from kjspkg's repo
kjspkg install carbon:[pkgname] - installs packages from carbonjs' repo (https://github.com/carbon-kjs)
kjspkg install github:[author]/[name] - installs external packages from github

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

Modern Modpacks - Owner
G_cat101 - Coder
malezjaa - Creator of CarbonJS
Juh9870 - Wanted to be here
    """

    print(INFO)
def guiinfo(): # Print info about the GUI app
    print(f"{_bold('Did you know there is a GUI app for KJSPKG?')} Check it out at https://github.com/Modern-Modpacks/kjspkg-gui!")
def kombucha(): # Kombucha easter egg
    RECIPE = f"""
{_bold("Ingredients")}

* 2 organic green teabags (or 2 tsp loose leaf)
* 2 organic black teabags bags (or 2 tsp loose leaf)
* 100-200g granulated sugar, to taste
* 1 medium scoby, plus 100-200ml starter liquid

{_bold("Method")}

STEP 1
For essential information on brewing safely, our top recipe tips and fun flavours to try, read our guide on how to make kombucha. Pour 1.8 litres boiled water into a saucepan, add the teabags and sugar (depending on how sweet you like it or the bitterness of your tea), stir to dissolve the sugar and leave for 6-10 mins to infuse.

STEP 2
Remove and discard the teabags without squeezing them. Leave the tea to cool completely before pouring into a large 2.5- to 3-litre glass jar. Add the scoby and its starter liquid, leaving a minimum of 5cm space at the top of the jar.

STEP 3
Cover the jar with a thin tea towel or muslin cloth so the scoby can 'breathe'. Secure with an elastic band and label the jar with the date and its contents.

STEP 4
Leave to ferment for one to two weeks at room temperature and away from radiators, the oven or direct sunlight. Do not put the jar in a cupboard, as air circulation is important.

STEP 5
After the first week, taste the kombucha daily – the longer you leave it, the more acidic the flavour will become. When ready, pour the kombucha into bottles, making sure to reserve the scoby and 100-200ml of starter fluid for the next batch.

STEP 6
The kombucha is ready to drink immediately, or you can start a ‘secondary fermentation’ by adding flavours such as fruit, herbs and spices to the drawn-off liquid and leaving it bottled for a few more days before drinking. Will keep in the fridge for up to three months.
    """
    print(RECIPE)
 
# PARSER FUNCTION
def _parser(func:str="help", *args, help:bool=False, **kwargs):
    global kjspkgfile

    system("") # Enable color codes on windows, don't ask me why

    if help: func="help"

    FUNCTIONS = { # Command mappings
        "install": install,
        "download": install,
        "add": install,
        "remove": removepkg,
        "uninstall": removepkg,
        "rm": removepkg,
        "update": update,
        "updateall": updateall,
        "list": listpkgs,
        "pkg": pkginfo,
        "pkginfo": pkginfo,
        "listall": listall,
        "all": listall,
        "search": search,
        "reload": reload,
        "refresh": reload,
        "init": init,
        "uninit": uninit,
        "help": info,
        "info": info,
        "gui": guiinfo,
        "mold": kombucha
    }

    if func not in FUNCTIONS.keys(): _err("Command \""+func+"\" is not found. Run \"kjspkg help\" to see all of the available commands") # Wrong command err
    
    helperfuncs = (info, guiinfo, init, pkginfo, listall, search, kombucha) # Helper commands that don't require a project
    if FUNCTIONS[func] not in helperfuncs and not _project_exists(): # If a project is not found, call init
        print(_bold("Project not found, a new one will be created.\n"))
        init()
    if (FUNCTIONS[func]==init or FUNCTIONS[func] not in helperfuncs) and path.exists(".kjspkg"): kjspkgfile = load(open(".kjspkg")) # Open .kjspkg

    FUNCTIONS[func](*args, **kwargs) # Run the command

    # Clean up
    if path.exists(".kjspkg") and FUNCTIONS[func] not in helperfuncs: # If uninit wasn't called and the command isn't a helper command
        _update_manifest() # Update .kjspkg
        with open(".kjspkg", "w") as f: dump(kjspkgfile, f) # Save .kjspkg

# RUN
if __name__=="__main__": # If not imported
    _clear_tmp() # Remove tmp

    try: Fire(_parser) # Run parser with fire
    except (KeyboardInterrupt, EOFError): exit(0) # Ignore some exceptions
    # except TypeError: _err("Wrong syntax") # Wrong syntax err
    except GitCommandNotFound: _err("Git not found. Install it here: https://git-scm.com/downloads") # Git not found err
    except (exceptions.ConnectionError, exceptions.ReadTimeout): _err("Low internet connection") # Low internet connection err

    _clear_tmp() # Remove tmp again

# Ok that's it bye