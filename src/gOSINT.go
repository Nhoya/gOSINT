package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/deckarep/golang-set"
	"github.com/jessevdk/go-flags"
)

const ver = "v0.5dev"

var opts struct {
	Module  string `short:"m" long:"module" description:"Specify module"  choice:"pgp" choice:"pwnd" choice:"git" choice:"plainSearch" choice:"telegram" choice:"shodan"`
	Version bool   `short:"v" long:"version" description:"Print version"`
	// git module
	URL        string `long:"url" default:"" description:"Specify target URL"`
	GitAPIType string `long:"gitAPI" default:"" description:"Specify git website API to use (for git module,optional)" choice:"github" choice:"bitbucket" choice:"clone"`
	Clone      bool   `short:"c" long:"clone" description:"Enable clone function for plainSearch module (need to specify repo URL)"`
	// pwn and pgp module
	Mail string `long:"mail" default:"" description:"Specify mail target (for pgp and pwnd module)"`
	// telegram module
	TgGrace  int    `long:"grace" default:"15" description:"Specify telegram messages grace period"`
	TgGroup  string `short:"g" long:"tgroup" default:"" description:"Specify Telegram group/channel name"`
	TgStart  int    `short:"s" long:"tgstart" default:"1" default-mask:"-" description:"Specify first message to scrape"`
	TgEnd    int    `short:"e" long:"tgend" description:"Specify last message to scrape"`
	DumpFile bool   `long:"dumpfile" description:"Create and resume messages from dumpfile"`
	// plainSearch module
	Confirm bool   `long:"ask-confirmation" description:"Ask confirmation before adding mail to set (for plainSearch module)"`
	Path    string `short:"p" long:"path" description:"Specify target path (for plainSearch module)"`
	// shodan module
	ShodanTarget       []string `short:"t" long:"target" description:"Specify shodan target host"`
	ShodanQuery        string   `short:"q" long:"query"  description:"Specify shodan query"`
	ShodanScan         bool     `long:"newscan" description:"Ask shodan for a new scan (-1 Scan credit)"`
	ShodanHoneyPotFlag bool     `long:"honeypot" description:"Check Honeypot probability"`
	// generic
	Mode bool `short:"f" long:"full" description:"Make deep search using linked modules"`
	JSON bool `long:"json" description:"Print JSON formatted output"`
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

	if opts.Module == "" {
		fmt.Println("You need to specify the module you want to use, -h for more info")
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println("gOSINT " + ver)
		os.Exit(0)
	}
	if opts.URL != "" {
		_, err := url.ParseRequestURI(opts.URL)
		if err != nil {
			fmt.Println("Invalid URL")
			os.Exit(1)
		}
	}

	switch mod := opts.Module; mod {
	case "pwnd":
		initPwnd(mailSet)
	case "pgp":
		initPGP(mailSet)
	case "git":
		initGit(mailSet)
	case "plainSearch":
		initPlainSearch(mailSet)
	case "telegram":
		initTelegram()
	case "shodan":
		initShodan()
	}
}
