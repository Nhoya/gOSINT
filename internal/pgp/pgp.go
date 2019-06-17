package pgp

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Nhoya/gOSINT/internal/utils"
)

//Options contains the options needed for the PGP modules
type Options struct {
	Targets  []string
	JSONFlag bool
}

//report contains the report for pgp search
type report struct {
	Target   string    `json:"Target"`
	Entities []*entity `json:"Entities"`
}

//PGPEntity identifies a PGP address on pgp.mit.edu
type entity struct {
	Person  alias   `json:"Person"`
	KeyID   string  `json:"KeyID"`
	Aliases []alias `json:"Aliases"`
}

//PGPAlias identiies a user with Name and email
type alias struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

//StartPGP is the init function of the PGP Module
func (opts *Options) StartPGP() {
	report := new(report)
	for _, mail := range opts.Targets {
		report.pgpSearch(mail)
	}
	report.printReport(opts.JSONFlag)
}

//send the request to pgp.mit.edu for each mail passed as argument
func (report *report) pgpSearch(mail string) {
	domain := "https://pgp.mit.edu/pks/lookup?search=" + mail
	body := string(utils.RetrieveRequestBody(domain))
	//create a new Report for each search term
	report.Target = mail
	report.extractPGPIdentities(body)
}

func (report *report) extractPGPIdentities(body string) {
	re := regexp.MustCompile(`search=(0x[0-9A-F]+)[^>]*>\s*([^&]+)&lt;([^&]+)[^<]+</a>\s*([^<]+)?`)
	matches := re.FindAllStringSubmatch(body, -1)
	for _, id := range matches {
		identity := new(entity)
		identity.KeyID = strings.TrimSpace(id[1])
		identity.Person.Name = strings.TrimSpace(id[2])
		identity.Person.Email = strings.TrimSpace(id[3])
		if id[4] != "" {
			identity.extractAliases(id[4])
		}
		report.Entities = append(report.Entities, identity)
	}
}

func (identity *entity) extractAliases(aliases string) {
	re := regexp.MustCompile(`(?m)^\s*(.*) +&lt;([^&]+)`)
	matches := re.FindAllStringSubmatch(aliases, -1)
	for _, j := range matches {
		alias := new(alias)
		alias.Name = strings.TrimSpace(j[1])
		alias.Email = strings.TrimSpace(j[2])
		identity.newAlias(alias)
	}
}

//NewAlias adds a new alias to a PGP identity
func (identity *entity) newAlias(a *alias) {
	identity.Aliases = append(identity.Aliases, *a)
}

func (report *report) printReport(jsonFlag bool) {
	if jsonFlag {
		jsonreport, _ := json.Marshal(&report)
		fmt.Println(string(jsonreport))
	} else {
		fmt.Println("==== PGP SEARCH FOR: " + report.Target + "====")
		if report.Entities == nil {
			fmt.Println("No results found")
			os.Exit(1)
		}
		for _, k := range report.Entities {
			fmt.Println(k.KeyID, k.Person)
			for _, j := range k.Aliases {
				fmt.Println("\tAlias:", j)
			}
		}
	}
}
