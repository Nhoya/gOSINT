package main

import (
	"fmt"
	"os"

	"github.com/deckarep/golang-set"
	"github.com/jessevdk/go-flags"
)

const ver = "v0.4c"

var opts struct {
	Module     string `short:"m" long:"module" description:"Specify module"  choice:"pgp" choice:"pwnd" choice:"git" choice:"plainSearch" choice:"telegram"`
	Url        string `long:"url" default:"" description:"Specify target URL"`
	Target     string `short:"t" long:"target" default:"" description:"Specify target"`
	GitAPIType string `long:"gitAPI" default:"" description:"Specify git website API to use (for git module,optional)" choice:"github" choice:"bitbucket"`
	Mail       string `long:"mail" default:"" description:"Specify mail target (for pgp and pwnd module)"`
	Path       string `short:"p" long:"path" description:"Specify target path (for plainSearch module)"`
	TgGrace    int    `long:"grace" default:"15" description:"Specify telegram messages grace period"`
	DumpFile   bool   `long:"dumpfile" description:"Create and resume messages from dumpfile"`
	Mode       bool   `short:"f" long:"full" description:"Make deep search using linked modules"`
	Clone      bool   `short:"c" long:"clone" description:"Enable clone function for plainSearch module (need to specify repo URL)"`
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
	if opts.Url != "" {
		isUrl(opts.Url)
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
		if opts.Clone {
			if opts.Url == "" {
				fmt.Println("You must specify target URL")
				os.Exit(1)
			}
			mailSet = cloneAndSearch(opts.Url, mailSet, opts.Confirm)
		} else {
			if opts.Path == "" {
				fmt.Println("You must specify Path")
				os.Exit(1)
			}
			mailSet = plainMailSearch(opts.Path, mailSet, opts.Confirm)
		}
		if opts.Mode {
			mailSet = pgpSearch(mailSet)
			pwnd(mailSet)
		}
	case "telegram":
		if opts.Target == "" {
			fmt.Println("You must specify target")
			os.Exit(1)
		}
		getTelegramGroupHistory(opts.Target, opts.TgGrace, opts.DumpFile)
	}
}
