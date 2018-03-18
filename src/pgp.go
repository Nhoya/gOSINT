package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/deckarep/golang-set"
)

func initPGP(mailSet mapset.Set) {
	mailCheck(mailSet)
	mailSet = pgpSearch(mailSet)
	if opts.Mode {
		pwnd(mailSet)
	}
}

//PGPAlias identiies a user with Name and email
type PGPAlias struct {
	Name  string
	email string
}

//PGPEntity identifies a PGP address on pgp.mit.edu
type PGPEntity struct {
	Person  PGPAlias
	KeyID   string
	Aliases []PGPAlias
}

//NewAlias adds a new alias to a PGP identity
func (identity *PGPEntity) NewAlias(a *PGPAlias) {
	identity.Aliases = append(identity.Aliases, *a)
}

func pgpSearch(mailSet mapset.Set) mapset.Set {
	fmt.Println("==== PGP SEARCH ====")
	mailIterator := mailSet.Iterator()
	for mail := range mailIterator.C {
		pgpSet := mapset.NewSet()
		fmt.Println("[+] pgp search for " + mail.(string))
		domain := "http://pgp.mit.edu/pks/lookup?search=" + mail.(string)
		body := retrieveRequestBody(domain)
		extractPGPIdentities(body)
		//findMailInText(body, pgpSet)
		if pgpSet != nil {
			pgpIterator := pgpSet.Iterator()
			for email := range pgpIterator.C {
				fmt.Println(email)
			}
			mailSet = mailSet.Union(pgpSet)
		} else {
			fmt.Println("No result found")
		}
	}
	return mailSet
}

func extractPGPIdentities(body string) {
	re := regexp.MustCompile(`search=(0x[0-9A-F]+)[^>]*>\s*([^&]+)&lt;([^&]+)[^<]+</a>\s*([^<]+)?`)
	matches := re.FindAllStringSubmatch(body, -1)
	for _, id := range matches {
		identity := new(PGPEntity)
		identity.KeyID = strings.TrimSpace(id[1])
		identity.Person.Name = strings.TrimSpace(id[2])
		identity.Person.email = strings.TrimSpace(id[3])
		if id[4] != "" {
			extractAliases(id[4], identity)
		}
		fmt.Println(identity.KeyID, identity.Person.Name, identity.Person.email, identity.Aliases)
	}

	os.Exit(1)
}

func extractAliases(aliases string, identity *PGPEntity) {
	re := regexp.MustCompile(`(?m)^\s*(.*) +&lt;([^&]+)`)
	matches := re.FindAllStringSubmatch(aliases, -1)
	for _, j := range matches {
		alias := new(PGPAlias)
		alias.Name = strings.TrimSpace(j[1])
		alias.email = strings.TrimSpace(j[2])
		identity.NewAlias(alias)
	}
}
