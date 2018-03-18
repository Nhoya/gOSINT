package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/deckarep/golang-set"
)

//Configuration struct will contain the configuration parameters
//the config file is available in $HOME/.config/gOSINT.conf
type Configuration struct {
	ShodanAPIKey string `json:"ShodanAPIKey"`
	GHToken      string `json:"GHToken"`
}

func retrieveRequestBody(domain string) string {
	resp, err := http.Get(domain)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func findMailInText(body string, mailSet mapset.Set) {
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

func isURL(URL string) {
	validURL, _ := regexp.MatchString(`(?i)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:'".,<>?«»“”‘’]))`, URL)
	if !validURL {
		fmt.Println("[-] " + URL + " is not a valid URL")
		os.Exit(1)
	}
}

func writeOnFile(filename string, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Unabale to open file")
		fmt.Println(err)
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
func createDirectory(dirname string) {
	if !fileExists(dirname) {
		fmt.Println("[+] Creating directory " + dirname)
		os.MkdirAll(dirname, os.ModePerm)
	}
}

func simpleQuestion(question string) bool {
	fmt.Println("[?] " + question + " [Y/N]")
	var resp string
	_, err := fmt.Scanln(&resp)
	if err != nil {
		fmt.Println("[-] Unable to read answer")
		os.Exit(1)
	}
	if resp == "y" || resp == "Y" {
		return true
	}
	return false
}

func getConfigFile() Configuration {
	file, err := os.Open(ConfigFilePath)
	if err != nil {
		fmt.Println("[-] Unable to open config file, be sure it exists")
		os.Exit(1)
	}
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[-] Unable to read config file")
		os.Exit(1)
	}
	return config
}
