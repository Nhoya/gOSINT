package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/jessevdk/go-flags"
	"os"
)

var opts struct {
	Module     string `short:"m" long:"module" description:"Specify module" required:"true" choice:"pgp" choice:"pwnd"  choice:"git"`
	Url        string `long:"url" default:"" description:"Specify target URL"`
	gitAPIType string `long:"gitAPI" default:"" description:"Specify git website API to use" choice:"github" choice"bitbucket"`
	Mail       string `long:"mail" default:"" description:"Specify mail target"`
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
		fmt.Println(err)
		os.Exit(1)
	}

	switch mod := opts.Module; mod {
	case "pwnd":
		mailCheck(mailSet)
		pwnd(mailSet)
	case "pgp":
		mailCheck(mailSet)
		pgpSearch(mailSet)
	case "git":
		if opts.Url == "" {
			fmt.Println("You must specify target URL")
			os.Exit(1)
		}
		mailSet = gitSearch(opts.Url, opts.gitAPIType, mailSet)
		mailSet = pgpSearch(mailSet)
		pwnd(mailSet)
	}
}
