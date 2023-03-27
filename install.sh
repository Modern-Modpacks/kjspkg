if [ "$EUID" -ne 0 ] then
    sudo su
fi

rm /usr/local/bin/kjspkg

curl -s https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/app.py > /usr/local/bin/kjspkg
pip -q install $(curl -s https://raw.githubusercontent.com/Modern-Modpacks/kjspkg/main/requirements.txt)
chmod +x /usr/local/bin/kjspkg

echo "Done!"