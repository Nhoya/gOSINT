package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/deckarep/golang-set"
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

func isUrl(url string) {
	validUrl, _ := regexp.MatchString(`(?i)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:'".,<>?«»“”‘’]))`, url)
	if !validUrl {
		fmt.Println("[-] " + url + " is not a valid URL")
		os.Exit(1)
	}
}

func writeOnFile(filename string, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Unabale to open file")
		os.Exit(1)
	}
	_, err = f.WriteString(text)
	if err != nil {
		fmt.Println("Unable to wite on file")
	}
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}
