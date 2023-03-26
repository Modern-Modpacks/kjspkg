#!/usr/bin/python3

# IMPORTS

# Built-in modules
from os import path, remove, mkdir, makedirs
from shutil import rmtree
from json import dump, load

# External libraries
from fire import Fire # CLI tool
from requests import get
from git import Repo

# CONSTANTS
VERSIONS = {
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
SCRIPT_DIRS = ("server_scripts", "client_scripts", "startup_scripts")

# VARIABLES
kjspkgfile = {}

def _bold(s:str) -> str: return "\u001b[1m"+s+"\u001b[0m"
def _err(err:str):
    print("\u001b[31;1m"+err+"\u001b[0m")
    exit(1)
def _create_tmp(path:str) -> str: 
    makedirs("tmp/"+path)
    return "tmp/"+path

def _check_project() -> bool:
    if (path.basename(path.curdir)=="kubejs"): return True

    for dir in SCRIPT_DIRS:
        if path.exists(dir): return True
    return False
def _create_project_directories():
    for dir in SCRIPT_DIRS:
        if not path.exists(dir): mkdir(dir)
        if not path.exists(dir+"/.kjspkg"): mkdir(dir+"/.kjspkg")
def _project_exists() -> bool: return path.exists(".kjspkg")
def _delete_project():
    remove(".kjspkg")
    for dir in SCRIPT_DIRS: rmtree(dir+"/.kjspkg")

def install(pkg:str):
    package = get()

    Repo.clone_from(url, to_path)

    kjspkgfile["installed"].append(pkg)
def removepkg(pkg:str):
    kjspkgfile["installed"].remove(pkg)
def listmods(*, count:bool=False):
    if count:
        print(len(kjspkgfile["installed"]))
        return

    print("\n".join(kjspkgfile["installed"]))
def init(*, version:str=None, modloader:str=None, quiet:bool=False, override:bool=False):
    if _project_exists():
        if not quiet and input("\u001b[31;1mA PROJECT ALREADY EXISTS IN THIS REPOSITORY, CREATING A NEW ONE OVERRIDES THE PREVIOUS ONE, ARE YOU SURE YOU WANT TO PROCEED? (y/N): \u001b[0m").lower()=="y" or override: _delete_project()
        else: exit(0)

    if not version: version = input(_bold("Input your minecraft version (1.12/1.16/1.18/1.19): "))
    if version not in VERSIONS.keys(): _err("Unknown or unsupported version: "+version)

    if not modloader: modloader = input(_bold("Input your modloader (forge/fabric/quilt): "))
    modloader = modloader.lower()
    if modloader not in ("forge", "fabric", "quilt"): _err("Unknown or unsupported modloader: "+modloader.title())

    _create_project_directories()
    with open(".kjspkg", "w+") as f: dump({
        "version": VERSIONS[version],
        "modloader": modloader if modloader!="quilt" else "fabric",
        "installed": []
    }, f)
    if not quiet: print(_bold("Project created!"))
def info():
    INFO = f"""
Hello world!
    """

    print(INFO)
    
def _parser(func:str="help", *args, **kwargs):
    global kjspkgfile

    FUNCTIONS = {
        "install": install,
        "remove": removepkg,
        "list": listmods,
        "init": init,
        "help": info,
        "info": info
    }

    if func not in FUNCTIONS.keys(): _err("Command \""+func+"\" is not found. Run \"kjspkg help\" to see all of the available commands")

    if not _check_project(): _err("Hmm... This directory doesn't look like a kubejs directory")

    if func not in ("init", "help") and not _project_exists(): 
        print(_bold("Project not found, a new one will be created.\n"))
        init()
    else:
        kjspkgfile = load(open(".kjspkg"))

        FUNCTIONS[func](*args, **kwargs)

        with open(".kjspkg", "w") as f: dump(kjspkgfile, f)
        if path.exists("tmp"): rmtree("tmp")

if __name__=="__main__": 
    try: Fire(_parser)
    except (KeyboardInterrupt, EOFError): exit(0)