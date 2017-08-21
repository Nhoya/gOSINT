[![Build Status](https://travis-ci.org/Nhoya/gOSINT.svg?branch=master)](https://travis-ci.org/Nhoya/gOSINT) [![Twitter](https://img.shields.io/twitter/url/https/github.com/Nhoya/gOSINT.svg?style=social)](https://twitter.com/intent/tweet?text=Wow:&url=%5Bobject%20Object%5D)
# gOSINT
OSINT framework in Go

## Introduction
gOSINT is a little OSINT framework in golang, it's actually in development and still not ready for production if you want, feel free to contribute!

## What gOSINT can do

- [x] Find mails from git repository
- [x] Find Dumps for mail address
- [x] Search for  mail address linked to domain/mail address in PGP keyring
- [x] Retrive Info from domain whois (waiting to be implemented)
- [x] Search for mail address in source code (waiting to be implemented)

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
```

`git clone https://github.com/Nhoya/gOSINT && cd gOSINT && go build`

## Modules

Currently `gOSINT` is still an early version and few modules are supported

- [x] git support for mail retriving (using github API, bitbucket API or RAW clone and search)
- [x] Search for mails in PGP Server
- [x] [https://haveibeenpwned.com/](http://haveibeenpwned.com/) search for mail in databreach
- [ ] WHOIS support (the module is ready but has to be integrated)
- [ ] Search for mail address in source (module ready, needs to be integrated)
- [ ] [https://shodan.io](https://shodan.io) search
- [ ] Social Media search
- [ ] Search Engine search

## Usage

```
Usage:
  gOSINT [OPTIONS]

Application Options:
  -m, --module=[pgp|pwnd|git]     Specify module
      --url=                      Specify target URL
      --gitAPI=[github|bitbucket] Specify git website API to use (optional)
      --mail=                     Specify mail target
  -f, --full                      Make deep search using linked modules
  -v, --version                   Print version

Help Options:
  -h, --help                      Show this help message
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
