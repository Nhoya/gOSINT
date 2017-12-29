package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"io/ioutil"
	"net/http"
	//"os"
	"regexp"
	"strconv"
	"strings"
)

func retriveLastPage(domain string) int {
	req, err := http.Get(domain)
	if err != nil {
		panic(err)
	}
	pagInfo := req.Header.Get("Link")
	if pagInfo != "" {
		re := regexp.MustCompile(`page=(\d+)>;\srel="last"`)
		match := re.FindStringSubmatch(pagInfo)
		lastPage, _ := strconv.Atoi(match[1])
		return lastPage
	}
	return 1
}

func retriveRequestBody(domain string) string {
	req, err := http.Get(domain)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	return string(body)
}

func findMailInText(body string, mailSet mapset.Set) {

	//re := regexp.MustCompile(`[\w\-\.]+\@[\w \.\-]+\.[\w]+`)
	re := regexp.MustCompile(`(?:![\n|\s])*(?:[\w\d\.\w\d]|(?:[\w\d]+[\-]+[\w\d]+))+[\@]+[\w]+[\.]+[\w]+`)
	mails := re.FindAllString(body, -1)
	if len(mails) == 0 {
		return
	}
	for _, mail := range mails {
		if !strings.Contains(mail, "noreply") {
			mailSet.Add(mail)
		}
	}

}

func readFromSet(mailSet mapset.Set) {
	mailIterator := mailSet.Iterator()
	if mailIterator != nil {
		for addr := range mailIterator.C {
			fmt.Println(addr)
		}
	}
}
