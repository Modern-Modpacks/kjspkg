#!/bin/bash

INFO="\033[34m::\033[0m"
WARN="\033[33m::\033[0m"

if [ "$(id -u)" != 0 ]; then
    echo -e "$INFO Escalation is required to install KJSPKG globally"
    sudo -v
    if [ $? -ne 0 ]; then
        echo -e "$WARN Failed to escalate privileges. Please run the script again."
        exit 1
    fi
fi

echo -e "$INFO Getting ready"

REPO="Modern-Modpacks/kjspkg"
LATEST_RELEASE_URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep -oP '"browser_download_url": "\K[^"]*kjspkg_linux_amd64')
if [ -z "$LATEST_RELEASE_URL" ]; then
    echo -e "$WARN Failed to fetch the latest release. Please check your internet connection!"
    exit 1
fi

echo -e "$INFO Downloading KJSPKG"
sudo curl --progress-bar -L -o /usr/local/bin/kjspkg "$LATEST_RELEASE_URL" | cat
if [ $? -ne 0 ]; then
    echo -e "$WARN Download failed. Please try again."
    exit 1
fi

sudo chmod +x /usr/local/bin/kjspkg
echo -e "$INFO Done! Run \`kjspkg\` to get started"
