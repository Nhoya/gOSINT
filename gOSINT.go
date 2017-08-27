package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/jessevdk/go-flags"
	"os"
)

const ver = "v0.3"

var opts struct {
	Module     string `short:"m" long:"module" description:"Specify module"  choice:"pgp" choice:"pwnd"  choice:"git" choice:"plainSearch"`
	Url        string `long:"url" default:"" description:"Specify target URL"`
	GitAPIType string `long:"gitAPI" default:"" description:"Specify git website API to use (for git module,optional)" choice:"github" choice:"bitbucket"`
	Mail       string `long:"mail" default:"" description:"Specify mail target (for pgp and pwnd module)"`
	Path       string `short:"p" long:"path" description:"Specify target path (for plainSearch module)"`
	Mode       bool   `short:"f" long:"full" description:"Make deep search using linked modules"`
	Confirm    bool   `long:"ask-confirmation" description:"Ask confirmation before adding mail to set (for plainSearch module)"`
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
	case "plainSearch":
		if opts.Path == "" {
			fmt.Println("You must specify target Path")
			os.Exit(1)
		}
		mailSet = plainMailSearch(opts.Path, mailSet, opts.Confirm)
		if opts.Mode {
			mailSet = pgpSearch(mailSet)
			pwnd(mailSet)
		}
	}
}
