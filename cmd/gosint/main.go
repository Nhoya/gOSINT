package main

import (
	"os"

	"github.com/Nhoya/gOSINT/internal/axfr"
	"github.com/Nhoya/gOSINT/internal/git"
	"github.com/Nhoya/gOSINT/internal/hibp"
	"github.com/Nhoya/gOSINT/internal/pgp"
	"github.com/Nhoya/gOSINT/internal/shodan"
	"github.com/Nhoya/gOSINT/internal/telegram"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	//gOSINT generic options
	app      = kingpin.New("gOSINT", "An Open Source INTelligence Swiss Army Knife")
	jsonFlag = app.Flag("json", "Enable JSON Output").Bool()

	//git module
	gitMod       = app.Command("git", "Get Emails and Usernames from repositories")
	gitRepo      = gitMod.Arg("url", "Repository URL").Required().URL()
	gitMethod    = gitMod.Flag("method", "Specify the API to use or plain clone").PlaceHolder("[clone|gh]").String()
	gitRecursive = gitMod.Flag("recursive", "Search for each repository of the user").Bool()

	//haveibeenpwned.com module
	hibpMod   = app.Command("pwd", "Check dumps for Email address using haveibeenpwned.com")
	hibpMails = hibpMod.Arg("mail", "Mail Address").Required().Strings()

	//pgp.mit module
	pgpMod     = app.Command("pgp", "Get Emails, KeyID and Aliases from PGP Keyring")
	pgpTargets = pgpMod.Arg("mail", "Mail or Domain").Required().Strings()

	//shodan.io scan module
	shodanMod      = app.Command("shodan", "Get info on host using shodan.io")
	shodanHosts    = shodanMod.Arg("host", "Remote Host IP").Required().Strings()
	shodanNewScan  = shodanMod.Flag("new-scan", "Schedule a new shodan scan (1 Shodan Credit will be deducted)").Bool()
	shodanHoneyPot = shodanMod.Flag("honeypot", "Get honeypot probability").Bool()
	//shodan.io query module
	shodanQMod  = app.Command("shodan-query", "Send a query to shodan.io")
	shodanQuery = shodanQMod.Arg("query", "Shodan query").Required().String()

	//crt.sh module (subdomain enumeration)
	axfrMod       = app.Command("axfr", "Subdomain enumeration using crt.sh")
	axfrURLs      = axfrMod.Arg("url", "Domain URL").Required().Strings()
	axfrURLStatus = axfrMod.Flag("verify", "Verify URL Status Code").Bool()

	//telegram.org module
	telegramMod         = app.Command("telegram", "Telegram public groups and channels scraper")
	telegramGroup       = telegramMod.Arg("group", "Telegram group or channel name").Required().String()
	telegramStart       = telegramMod.Flag("start", "Start message #").Int()
	telegramEnd         = telegramMod.Flag("end", "End message #").Int()
	telegramGracePeriod = telegramMod.Flag("grace", "The number of messages that will be considered deleted before the scraper stops").Default("15").Int()
	telegramDumpFlag    = telegramMod.Flag("dump", "Creates and resume messages from dumpfile").Bool()
)

func main() {
	app.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)
	app.Version("0.5dev")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// gosint git
	case gitMod.FullCommand():
		opts := new(git.Options)
		opts.Method = *gitMethod
		opts.Repo = *gitRepo
		opts.Recursive = *gitRecursive
		opts.JSONFlag = *jsonFlag
		opts.StartGit()
	// gosint pwnd
	case hibpMod.FullCommand():
		opts := new(hibp.Options)
		opts.Mails = *hibpMails
		opts.JSONFlag = *jsonFlag
		opts.StartHIBP()
	// gosint pgp
	case pgpMod.FullCommand():
		opts := new(pgp.Options)
		opts.JSONFlag = *jsonFlag
		opts.Targets = *pgpTargets
		opts.StartPGP()
		//gosint shodan
	case shodanMod.FullCommand():
		opts := new(shodan.Options)
		opts.Hosts = *shodanHosts
		opts.NewScan = *shodanNewScan
		opts.HoneyPotFlag = *shodanHoneyPot
		opts.StartShodanScan()
	//gosint shodan-query
	case shodanMod.FullCommand():
		opts := new(shodan.QueryOptions)
		opts.Query = *shodanQuery
		opts.StartShodanQuery()
	//gosint axfr
	case axfrMod.FullCommand():
		opts := new(axfr.Options)
		opts.URLs = *axfrURLs
		opts.JSONFlag = *jsonFlag
		opts.VerifyURLStatus = *axfrURLStatus
		opts.StartAXFR()
		//gosint telegram
	case telegramMod.FullCommand():
		opts := new(telegram.Options)
		opts.Group = *telegramGroup
		opts.Start = *telegramStart
		opts.End = *telegramEnd
		opts.GracePeriod = *telegramGracePeriod
		opts.DumpFlag = *telegramDumpFlag
		opts.StartTelegram()
	}
}
