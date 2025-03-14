#!/bin/bash

export DEV_HOME=$(pwd)

mkdir -p $DEV_HOME/kubejs/server_scripts
mkdir -p $DEV_HOME/kubejs/client_scripts
mkdir -p $DEV_HOME/kubejs/startup_scripts
mkdir -p $DEV_HOME/mods

echo "done! => cd $DEV_HOME/kubejs"
echo "      => go run ../cmd/kjspkg init"
