package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func retriveRequestBody(domain string) string {
	req, err := http.Get(domain)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	return string(body)
}

func findMailInText(body string, mailSet mapset.Set) mapset.Set {

	re := regexp.MustCompile(`[\w\-\.]+\@[\w \.\-]+\.[\w]+`)
	mails := re.FindAllString(body, -1)
	if len(mails) == 0 {
		return nil
	}

	for _, mail := range mails {
		if !strings.Contains(mail, "noreply") {
			mailSet.Add(mail)
		}
	}

	return (mailSet)
}

func readFromSet(mailSet mapset.Set) {
	mailIterator := mailSet.Iterator()
	if mailIterator != nil {
		for addr := range mailIterator.C {
			fmt.Println(addr)
		}
	}
}
