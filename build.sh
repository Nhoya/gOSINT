#!/bin/bash
readonly RED="\033[01;31m"
readonly GREEN="\033[01;32m"
readonly BLUE="\033[01;34m"
readonly YELLOW="\033[00;33m"
readonly BOLD="\033[01m"
readonly END="\033[0m"

dependencies=( github.com/deckarep/golang-set github.com/nhoya/goPwned github.com/jessevdk/go-flags gopkg.in/src-d/go-git.v4)


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
