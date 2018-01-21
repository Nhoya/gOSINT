# gOSINT [![Build Status](https://travis-ci.org/Nhoya/gOSINT.svg?branch=master)](https://travis-ci.org/Nhoya/gOSINT) [![GitHub stars](https://img.shields.io/github/stars/Nhoya/gOSINT.svg)](https://github.com/Nhoya/gOSINT/stargazers) [![GitHub forks](https://img.shields.io/github/forks/Nhoya/gOSINT.svg)](https://github.com/Nhoya/gOSINT/network) [![Twitter](https://img.shields.io/twitter/url/https/github.com/Nhoya/gOSINT.svg?style=social&style=plastic)](https://twitter.com/intent/tweet?text=Wow:&url=https%3A%2F%2Fgithub.com%2FNhoya%2FgOSINT) [![Go Report Card](https://goreportcard.com/badge/github.com/Nhoya/gOSINT)](https://goreportcard.com/report/github.com/Nhoya/gOSINT) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/76673062a30e48bd99d499d32c0c6af0)](https://www.codacy.com/app/Nhoya/gOSINT?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Nhoya/gOSINT&amp;utm_campaign=Badge_Grade)
OSINT framework in Go

Take a look at the [develop branch](https://github.com/Nhoya/gOSINT/tree/develop) for more updates.

## Introduction

gOSINT is a small OSINT framework in Golang. If you want, feel free to contribute and/or leave a feedback!

## Like my project? Consider donation :)

[![Paypal Badge](https://img.shields.io/badge/Donate-PayPal-yellow.svg)](https://www.paypal.me/Nhoya) [![BTC Badge](https://img.shields.io/badge/Donate-BTC-yellow.svg)](https://pastebin.com/raw/nyDDPwaM) [![Monero Badge](https://img.shields.io/badge/Donate-XMR-yellow.svg)](https://pastebin.com/raw/dNUFqwuC)

## What gOSINT can do

- [x] Find mails from git repository
- [x] Find Dumps for mail address
- [x] Search for  mail address linked to domain/mail address in PGP keyring
- [x] Retrieve Info from domain whois (waiting to be implemented)
- [x] Search for mail address in source code
- [x] Retrieve Telegram Public Groups History

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
go get "gopkg.in/ns3777k/go-shodan.v2/shodan"
```

`git clone https://github.com/Nhoya/gOSINT && cd gOSINT && go build`

## Modules

Currently `gOSINT` is still an early version and few modules are supported

- [x] git support for mail retriving (using github API, bitbucket API or RAW clone and search)
- [x] Search for mails in PGP Server
- [x] [https://haveibeenpwned.com/](http://haveibeenpwned.com/) search for mail in databreach
- [x] Retrive Telegram Public Group Messages
- [x] Search for mail address in source
- [x] [https://shodan.io](https://shodan.io) search
- [ ] Social Media search
- [ ] Search Engine search

## Usage

```
gOSINT [OPTIONS]

Application Options:
  -m, --module=[pgp|pwnd|git|plainSearch|telegram|shodan] Specify module
  -v, --version                                           Print version
      --url=                                              Specify target URL
      --gitAPI=[github|bitbucket]                         Specify git website API to use (for git module,optional)
  -c, --clone                                             Enable clone function for plainSearch module (need to specify repo URL)
      --mail=                                             Specify mail target (for pgp and pwnd module)
      --grace=                                            Specify telegram messages grace period (default: 15)
  -g, --tgroup=                                           Specify Telegram group/channel name
  -s, --tgstart=                                          Specify first message to scrape
  -e, --tgend=                                            Specify last message to scrape
      --dumpfile                                          Create and resume messages from dumpfile
      --ask-confirmation                                  Ask confirmation before adding mail to set (for plainSearch module)
  -p, --path=                                             Specify target path (for plainSearch module)
  -t, --target=                                           Specify shodan target host
  -f, --full                                              Make deep search using linked modules

Help Options:
  -h, --help                                              Show this help message
```

## Configuration file

The configuration file is in `$HOME/.config/gOSINT.conf`

If some API Key is missing insert it there

## Examples

Currently `gOSINT` supports the following actions:


`gOSINT -m git --url=[RepoURL] --gitAPI [github|bitbucket] (optional)`

retrieve mail from git commits

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

`gOSINT -m telegram --tgroup | -g  [PublicGroupName]`

retrieve message history for telegram public group

`gOSINT -m telegram --tgroup | -g [PublicGroupName] --dumpfile`

the output will be stored in a file, if the file is already populated it will resume from the last ID

`gOSINT -m telegram --tgroup | -g [PublicGroupName] --dumpfile -s [masageID] -e [messageID]`

Set start and end messages for scraping

`gOSINT -m shodan -t [HOST IP]`

Start Shodan Scan for Host

## PGP module Demo
[![asciicast](https://asciinema.org/a/21PCpbgFqyHiTbPINexHKEywj.png)](https://asciinema.org/a/21PCpbgFqyHiTbPINexHKEywj)

## Pwnd module Demo
[![asciicast](https://asciinema.org/a/x9Ap0IRcNNcLfriVujkNUhFSF.png)](https://asciinema.org/a/x9Ap0IRcNNcLfriVujkNUhFSF)

## Telegram Crawler Demo
[![asciicast](https://asciinema.org/a/nbRO9FNpjiYXAKeI87xn29j9z.png)](https://asciinema.org/a/nbRO9FNpjiYXAKeI87xn29j9z)
