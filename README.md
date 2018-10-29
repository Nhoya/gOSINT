# gOSINT [![Build Status](https://travis-ci.org/Nhoya/gOSINT.svg?branch=master)](https://travis-ci.org/Nhoya/gOSINT) [![Build status](https://ci.appveyor.com/api/projects/status/9qn2y2f8t5up8ww2?svg=true)](https://ci.appveyor.com/project/Nhoya/gosint) [![GitHub stars](https://img.shields.io/github/stars/Nhoya/gOSINT.svg)](https://github.com/Nhoya/gOSINT/stargazers) [![GitHub forks](https://img.shields.io/github/forks/Nhoya/gOSINT.svg)](https://github.com/Nhoya/gOSINT/network) [![Twitter](https://img.shields.io/twitter/url/https/github.com/Nhoya/gOSINT.svg?style=social&style=plastic)](https://twitter.com/intent/tweet?text=Wow:&url=https%3A%2F%2Fgithub.com%2FNhoya%2FgOSINT) [![Go Report Card](https://goreportcard.com/badge/github.com/Nhoya/gOSINT)](https://goreportcard.com/report/github.com/Nhoya/gOSINT) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/76673062a30e48bd99d499d32c0c6af0)](https://www.codacy.com/app/Nhoya/gOSINT?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Nhoya/gOSINT&amp;utm_campaign=Badge_Grade) [![Mentioned in Awesome Pentest](https://awesome.re/mentioned-badge.svg)](https://github.com/enaqx/awesome-pentest)

OSINT Swiss Army Knife in Go

Take a look at the [develop branch](https://github.com/Nhoya/gOSINT/tree/develop) for more updates.

## Introduction

gOSINT is a multiplatform OSINT Swiss army knife in Golang. If you want, feel free to contribute and/or leave a feedback!

## Like my project? Please consider donation :)

[![Paypal Badge](https://img.shields.io/badge/Donate-PayPal-yellow.svg)](https://www.paypal.me/Nhoya) [![BTC Badge](https://img.shields.io/badge/Donate-BTC-yellow.svg)](https://pastebin.com/raw/nyDDPwaM) [![Monero Badge](https://img.shields.io/badge/Donate-XMR-yellow.svg)](https://pastebin.com/raw/dNUFqwuC) [![Ethereum Badge](https://img.shields.io/badge/Donate-Ethereum-yellow.svg)](https://pastebin.com/raw/S6XMmSiv)

## What gOSINT can do

Currently `gOSINT` has different modules:

- [x] git support for mail retriving (using github API, or plain clone and search)
- [x] Search for mails, aliases and KeyID in PGP Server
- [x] [haveibeenpwned.com/](http://haveibeenpwned.com/) search for mail in databreach
- [x] Retrieve Telegram Public Group Message History
- [x] Search for mail address in source
- [x] [shodan.io](https://shodan.io) search
- [x] Subdomain enumeration using [crt.sh](https://crt.sh)
- [x] Given a phone number, can retrieve the owner name
- [x] Search for password relatives to email address :P
- [x] Reverse Whois given Email Address or Name

A complete features list and roadmap is available under [Projects Tab](https://github.com/Nhoya/gOSINT/projects)

## Installation

### Dependencies

gOSINT currently depends from [tesseract-ocr](https://github.com/tesseract-ocr/) so you need to install on your system `tesseract-ocr`, `libtesseract-dev` and `libleptonica-dev`

### Install on a go-dependent way (is the easier and faster way)

You can install `gOSINT` using `go get` with a simple `go get github.com/Nhoya/gOSINT/cmd/gosint`

### Install On Windows

Check the AppVeyor Build page for builds

## Manual Building

### Building On Linux

Build gOSINT on linux is really easy, you just need to install [dep](https://github.com/golang/dep), clone the repository and `make` and `make install`

### Building On Windows

If you have `make` installed you can follow the Linux instructions (and skip `make install`) otherwise be sure to have [dep](https://github.com/golang/dep) installed, clone the directory and run

```bash
dep ensure
go build cmd/gosint.go
```

## Usage

```bash
usage: gOSINT [<flags>] <command> [<args> ...]

An Open Source INTelligence Swiss Army Knife

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --json     Enable JSON Output
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.


  git [<flags>] <url>
    Get Emails and Usernames from repositories

    --method=[clone|gh]  Specify the API to use or plain clone
    --recursive          Search for each repository of the user

  pwd [<flags>] <mail>...
    Check dumps for Email address using haveibeenpwned.com

    --get-passwords  Search passwords for mail

  pgp <mail>...
    Get Emails, KeyID and Aliases from PGP Keyring


  shodan [<flags>] <host>...
    Get info on host using shodan.io

    --new-scan  Schedule a new shodan scan (1 Shodan Credit will be deducted)
    --honeypot  Get honeypot probability

  shodan-query <query>
    Send a query to shodan.io


  axfr [<flags>] <url>...
    Subdomain enumeration using crt.sh

    --verify  Verify URL Status Code

  pni <number>...
    Retrieve info about a give phone number


  telegram [<flags>] <group>
    Telegram public groups and channels scraper

    --start=START  Start message #
    --end=END      End message #
    --grace=15     The number of messages that will be considered deleted before the scraper stops
    --dump         Creates and resume messages from dumpfile

  rev-whois <target>
    Find domains for name or email address

```

## Configuration file

The default configuration file is in `$HOME/.config/gosint.toml` for linux environment and `./config/toml` for windows env

You can place it in different paths, load prioriy is:

- `.`
- `./config/ or $HOME/.config/`
- `/etc/gosint/`

If some API Keys are missing insert it there

## PGP module Demo (**OUTDATED**)

[![asciicast](https://asciinema.org/a/21PCpbgFqyHiTbPINexHKEywj.png)](https://asciinema.org/a/21PCpbgFqyHiTbPINexHKEywj)

## Pwnd module Demo (**OUTDATED**)

[![asciicast](https://asciinema.org/a/x9Ap0IRcNNcLfriVujkNUhFSF.png)](https://asciinema.org/a/x9Ap0IRcNNcLfriVujkNUhFSF)

## Telegram Crawler Demo (**OUTDATED**)

[![asciicast](https://asciinema.org/a/nbRO9FNpjiYXAKeI87xn29j9z.png)](https://asciinema.org/a/nbRO9FNpjiYXAKeI87xn29j9z)

## Shodan module Demo (**OUTDATED**)

[![asciicast](https://asciinema.org/a/9lfzAZ65n9MJCkrUrxoHZQYwP.png)](https://asciinema.org/a/9lfzAZ65n9MJCkrUrxoHZQYwP)
