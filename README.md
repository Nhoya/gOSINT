[![Build Status](https://travis-ci.org/Nhoya/gOSINT.svg?branch=master)](https://travis-ci.org/Nhoya/gOSINT)
# gOSINT
OSINT framework in Go

## Introduction
gOSINT is a little OSINT framework in golang, it's actually in development and still not ready for production if you want, feel free to contribute!


## Dependecies
Before building `gOSINT` you need to solve the dependencies

```
"github.com/deckarep/golang-set"
"github.com/nhoya/goPwned"
"github.com/jessevdk/go-flags"
"gopkg.in/src-d/go-git.v4"
```

## Building

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

Currently `gOSINT` supports the following actions


`./gOSINT -m git --url=[RepoURL] --gitAPI [github|bitbucket] (optional)`

retrive mail from git commits

`./gOSINT -m git --url [RepoURL] --gitAPI [github|bitbucket] (optional) -f`

pass the result to pgp search ad pwnd module

`./gOSINT -m pwnd --mail [targetMail]`

search for breaches where targetMail is preset

`./gOSINT -m pgp --mail [targetMail]`

search for others mail in PGP Server

`./gOSINT -m pgp --mail [targetMail] -f`

pass the result to haveibeenpwn module
