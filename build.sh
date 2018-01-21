#!/bin/bash
readonly RED="\033[01;31m"
readonly GREEN="\033[01;32m"
readonly BLUE="\033[01;34m"
readonly YELLOW="\033[00;33m"
readonly BOLD="\033[01m"
readonly END="\033[0m"

version=$(go version 2> /dev/null)
if [[ "$?" != 0 ]]; then
        echo "Unable to find go, you need go >= 1.8 to build gOSINT"
        exit 1
fi
go_version_regex="([0-9]).([0-9]).[0-9]"
if [[ "$version" =~ $go_version_regex ]]; then
        if [[ ${BASH_REMATCH[1]} -le 1 ]]; then
                if [[ ${BASH_REMATCH[2]} -lt 8 ]]; then
                        echo "This version of go is not supported, you need go >= 1.8"
                        echo "Current: $version"
                        exit 1
                fi
        fi
fi

dependencies=( github.com/deckarep/golang-set github.com/nhoya/goPwned github.com/jessevdk/go-flags gopkg.in/src-d/go-git.v4 github.com/jaytaylor/html2text gopkg.in/ns3777k/go-shodan.v2/shodan )


echo -e "${GREEN}[+] Installing dependencies${END}"
for i in "${dependencies[@]}"
do
        depname=$(echo $i |awk -F / '{print$3}') 
        echo -e "${BLUE}[+] Installing $depname${END}"
        go get -v "$i"
        if [ $? != 0 ]; then
            echo -e "${RED}[-]$i raised error during installation${END}"
            exit 2
        fi
done

echo -e "${GREEN}[+] Building gOSINT${END}"
go build
echo -e "${GREEN}[+] Installing gOSINT${END}"
sudo mv gOSINT /usr/local/bin
mv config/config.json $HOME/.config/gOSINT.conf
