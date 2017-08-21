package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/jessevdk/go-flags"
	"os"
)

const ver = "v0.2"

var opts struct {
	Module     string `short:"m" long:"module" description:"Specify module"  choice:"pgp" choice:"pwnd"  choice:"git"`
	Url        string `long:"url" default:"" description:"Specify target URL"`
	GitAPIType string `long:"gitAPI" default:"" description:"Specify git website API to use (optional)" choice:"github" choice:"bitbucket"`
	Mail       string `long:"mail" default:"" description:"Specify mail target"`
	Mode       bool   `short:"f" long:"full" description:"Make deep search using linked modules"`
	Version    bool   `short:"v" long:"version" description:"Print version"`
}

func mailCheck(mailSet mapset.Set) {
	if opts.Mail == "" {
		fmt.Println("You must specify target mail")
		os.Exit(1)
	}
	mailSet.Add(opts.Mail)
}

func main() {

	mailSet := mapset.NewSet()
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println("gOSINT " + ver)
		os.Exit(0)
	}

	switch mod := opts.Module; mod {
	case "pwnd":
		mailCheck(mailSet)
		pwnd(mailSet)
	case "pgp":
		mailCheck(mailSet)
		mailSet = pgpSearch(mailSet)
		if opts.Mode {
			pwnd(mailSet)
		}
	case "git":
		if opts.Url == "" {
			fmt.Println("You must specify target URL")
			os.Exit(1)
		}
		mailSet = gitSearch(opts.Url, opts.GitAPIType, mailSet)
		if opts.Mode {
			mailSet = pgpSearch(mailSet)
			pwnd(mailSet)
		}
	}
}
