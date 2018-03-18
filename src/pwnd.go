package main

import (
	//	"encoding/json"
	"fmt"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/nhoya/goPwned"
)

type HIBPReport struct {
	Pwnd []PwnedEntity
}
type PwnedEntity struct {
	Email    string   `json:email`
	Breaches []string `json:breaches`
}

func initPwnd(mailSet mapset.Set) {
	mailCheck(mailSet)
	pwnd(mailSet)
}

func pwnd(mailSet mapset.Set) {
	mailIterator := mailSet.Iterator()
	for mail := range mailIterator.C {
		fmt.Println("[+] Dump for " + mail.(string))
		stuff, err := gopwned.GetAllBreachesForAccount(mail.(string), "", "true")
		pwnd := new(PwnedEntity)
		pwnd.Email = mail.(string)
		if err == nil {
			for _, data := range stuff {
				pwnd.Breaches = append(pwnd.Breaches, data.Name)
			}
			fmt.Println(pwnd)
		}
		time.Sleep(time.Second * 2)
	}
}
