package main

import (
	"fmt"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/nhoya/goPwned"
)

func initPwnd(mailSet mapset.Set) {
	mailCheck(mailSet)
	pwnd(mailSet)
}

func pwnd(mailSet mapset.Set) {
	fmt.Println("==== HAVEIBEENPWND SEARCH ====")
	mailIterator := mailSet.Iterator()
	for mail := range mailIterator.C {
		fmt.Println("[+] Dump for " + mail.(string))
		stuff, err := gopwned.GetAllBreachesForAccount(mail.(string), "", "true")
		if err == nil {
			for _, data := range stuff {
				fmt.Println(data.Name)
			}
		}
		time.Sleep(time.Second * 2)
	}
}
