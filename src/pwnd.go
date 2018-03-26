package main

import (
	"encoding/json"
	//	"github.com/deckarep/golang-set"
	"fmt"
	"time"

	"github.com/nhoya/goPwned"
)

type HIBPReport struct {
	Pwnd []*PwnedEntity
}
type PwnedEntity struct {
	Email    string   `json:"email"`
	Breaches []string `json:"breaches"`
}

func initPwnd() {
	mailCheck()
	report := new(HIBPReport)
	for n, mail := range opts.Mail {
		report.getBreachesForMail(mail)
		if n != 0 {
			time.Sleep(time.Second * 2)
		}
	}
	report.printHIBPReport()
}

func (report *HIBPReport) getBreachesForMail(mail string) {
	fmt.Println("[+] Dump for " + mail)
	stuff, err := gopwned.GetAllBreachesForAccount(mail, "", "true")
	if err == nil {
		pwnd := new(PwnedEntity)
		pwnd.Email = mail
		for _, data := range stuff {
			pwnd.Breaches = append(pwnd.Breaches, data.Name)
		}
		report.Pwnd = append(report.Pwnd, pwnd)
	}
}

func (report *HIBPReport) printHIBPReport() {
	if opts.JSON {
		jsonreport, _ := json.MarshalIndent(&report, "", " ")
		fmt.Println(string(jsonreport))
	} else {
		for _, k := range report.Pwnd {
			fmt.Println(k)
		}
	}
}
