# gOSINT [![Build Status](https://travis-ci.org/Nhoya/gOSINT.svg?branch=master)](https://travis-ci.org/Nhoya/gOSINT) [![GitHub stars](https://img.shields.io/github/stars/Nhoya/gOSINT.svg)](https://github.com/Nhoya/gOSINT/stargazers) [![GitHub forks](https://img.shields.io/github/forks/Nhoya/gOSINT.svg)](https://github.com/Nhoya/gOSINT/network) [![Twitter](https://img.shields.io/twitter/url/https/github.com/Nhoya/gOSINT.svg?style=social&style=plastic)](https://twitter.com/intent/tweet?text=Wow:&url=https%3A%2F%2Fgithub.com%2FNhoya%2FgOSINT)
OSINT framework in Go

you probably want to take a look at the develop branch for more updates.

## Introduction

gOSINT is a small OSINT framework in Golang, it's currently in development and still not ready for production if you want, feel free to contribute!


## What gOSINT can do

- [x] Find mails from git repository
- [x] Find Dumps for mail address
- [x] Search for  mail address linked to domain/mail address in PGP keyring
- [x] Retrive Info from domain whois (waiting to be implemented)
- [x] Search for mail address in source code
- [x] Retrive Telegram Public Groups History

## Building

You can use the building script, just clone the directory and execute it

```
git clone https://github.com/Nhoya/gOSINT
./build.sh
```

The package will be installed in `/usr/local/bin`

You can then call `gOSINT` from command line

`$ gOSINT --help`
 

## Manual Building 

#### Dependecies
Before building `gOSINT` manually you need to solve the dependencies:

```
go get "github.com/deckarep/golang-set"
go get "github.com/nhoya/goPwned"
go get "github.com/jessevdk/go-flags"
go get "gopkg.in/src-d/go-git.v4"
go get "github.com/jaytaylor/html2text"
```

`git clone https://github.com/Nhoya/gOSINT && cd gOSINT && go build`

## Modules

Currently `gOSINT` is still an early version and few modules are supported

- [x] git support for mail retriving (using github API, bitbucket API or RAW clone and search) *Now with Pagination*!
- [x] Search for mails in PGP Server
- [x] [https://haveibeenpwned.com/](http://haveibeenpwned.com/) search for mail in databreach
- [x] Retrive Telegram Public Group Messages
- [ ] WHOIS support (the module is ready but has to be integrated)
- [x] Search for mail address in source
- [ ] [https://shodan.io](https://shodan.io) search
- [ ] Social Media search
- [ ] Search Engine search

## Usage

```
Usage:
  gOSINT [OPTIONS]

Application Options:
  -m, --module=[pgp|pwnd|git|plainSearch] Specify module
      --url=                              Specify target URL
      --gitAPI=[github|bitbucket]         Specify git website API to use (for git module,optional)
      --mail=                             Specify mail target (for pgp and pwnd module)
  -p, --path=                             Specify target path (for plainSearch module)
  -f, --full                              Make deep search using linked modules
  -c, --clone                             Enable clone function for plainSearch module (need to specify repo URL)
      --ask-confirmation                  Ask confirmation before adding mail to set (for plainSearch module)
  -v, --version                           Print version

Help Options:
  -h, --help                              Show this help message
```

## Examples

Currently `gOSINT` supports the following actions


`gOSINT -m git --url=[RepoURL] --gitAPI [github|bitbucket] (optional)`

retrive mail from git commits

`gOSINT -m git --url [RepoURL] --gitAPI [github|bitbucket] (optional) -f`

pass the result to pgp search and pwnd module

`gOSINT -m pwnd --mail [targetMail]`

search for breaches where targetMail is preset

`gOSINT -m pgp --mail [targetMail]`

search for others mail in PGP Server

`gOSINT -m pgp --mail [targetMail] -f`

pass the result to haveibeenpwn module

`gOSINT -m sourceSerch --path [targetDirectory]`

search for mails in source code (recursively)

`gOSINT -m sourceSearh --path [targetDirectory] --ask-confirmation`

ask confirmation before adding  mail to search results

`gOSINT -m sourceSearch --path [targetDirectory] -f`

pass the result to pgp search and haveibeenpwnd modules

`gOSINT -m sourceSearch --clone --url [targetRepository]`

clone and search mail in repository source

`gOSINT -m sourceSearch --clone --url [targetRepository] -f`

pass the resoult to pgp search and haveibeenpwnd modules

`gOSINT -m sourceSearch --clone --url [targetRepository] --ask-confirmation`

ask confirmation before adding mail to search results

`gOSINT -m telegram --target [PublicGroupName]`

retrive message history for telegram public group

`gOSINT -m telegram --target [PublicGroupName] --dumpfile`

the output will be stored in a file, if the file is already populated it will resume from the last ID
