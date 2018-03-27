package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	//	"github.com/deckarep/golang-set"
)

func initPGP() {
	for _, mail := range opts.Mail {
		pgpSearch(mail)
	}
}

//PGPReport
type PGPReport struct {
	Target   string       `json:"Target"`
	Entities []*PGPEntity `json:"Entities"`
}

//PGPEntity identifies a PGP address on pgp.mit.edu
type PGPEntity struct {
	Person  PGPAlias   `json:"Person"`
	KeyID   string     `json:"KeyID"`
	Aliases []PGPAlias `json:"Aliases"`
}

//PGPAlias identiies a user with Name and email
type PGPAlias struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

//NewAlias adds a new alias to a PGP identity
func (identity *PGPEntity) newAlias(a *PGPAlias) {
	identity.Aliases = append(identity.Aliases, *a)
}

//send the request to pgp.mit.edu for each mail passed as argument
func pgpSearch(mail string) {
	domain := "http://pgp.mit.edu/pks/lookup?search=" + mail
	body := retrieveRequestBody(domain)
	//create a new Report for each search term
	report := new(PGPReport)
	report.Target = mail
	report.extractPGPIdentities(body)
}

func (report *PGPReport) extractPGPIdentities(body string) {
	re := regexp.MustCompile(`search=(0x[0-9A-F]+)[^>]*>\s*([^&]+)&lt;([^&]+)[^<]+</a>\s*([^<]+)?`)
	matches := re.FindAllStringSubmatch(body, -1)
	for _, id := range matches {
		identity := new(PGPEntity)
		identity.KeyID = strings.TrimSpace(id[1])
		identity.Person.Name = strings.TrimSpace(id[2])
		identity.Person.Email = strings.TrimSpace(id[3])
		if id[4] != "" {
			identity.extractAliases(id[4])
		}
		report.Entities = append(report.Entities, identity)
	}
	report.printPGPReport()
}

func (identity *PGPEntity) extractAliases(aliases string) {
	re := regexp.MustCompile(`(?m)^\s*(.*) +&lt;([^&]+)`)
	matches := re.FindAllStringSubmatch(aliases, -1)
	for _, j := range matches {
		alias := new(PGPAlias)
		alias.Name = strings.TrimSpace(j[1])
		alias.Email = strings.TrimSpace(j[2])
		identity.newAlias(alias)
	}
}

func (report *PGPReport) printPGPReport() {
	if opts.JSON {
		jsonreport, _ := json.MarshalIndent(&report, "", " ")
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
				fmt.Println(j)
			}
		}
	}
}
