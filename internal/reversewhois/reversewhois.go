package reversewhois

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type Options struct {
	Target   string
	JSONFlag bool
}

type report struct {
	Target  string    `json:"target"`
	Domains []*Domain `json:"domains"`
}

type Domain struct {
	DomainName   string `json:"domainName"`
	CreationDate string `json:"creationDate"`
	Registrar    string `json:"registrar"`
}

func (opts *Options) StartReverseWhois() {
	if opts.Target == "" {
		fmt.Println("You must inster a valid target")
		os.Exit(1)
	}
	report := new(report)
	report.Target = opts.Target

	report.Domains = getDomains(opts.Target)
	if len(report.Domains) == 0 {
		fmt.Println("Unable to find results")
		os.Exit(0)
	}
	report.printReport(opts.JSONFlag)
}

func getDomains(target string) []*Domain {
	resp, err := http.Get("https://viewdns.info/reversewhois/?q=" + url.QueryEscape(target))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var domains []*Domain
	for _, i := range extractValues(body) {
		dom := new(Domain)
		dom.DomainName = i[1]
		dom.CreationDate = i[2]
		dom.Registrar = i[3]
		domains = append(domains, dom)
	}
	return domains
}

func extractValues(body []byte) [][]string {
	re := regexp.MustCompile(`(?mU)<\/tr><tr><td>([^\s]+)<\/td><td>(\d{4}-\d{2}-\d{2})<\/td><td>(.*)<\/td>`)
	match := re.FindAllStringSubmatch(string(body), -1)
	return match
}

func (report *report) printReport(jsonFlag bool) {
	if jsonFlag {
		jsonreport, _ := json.Marshal(&report)
		fmt.Println(string(jsonreport))
	} else {
		fmt.Println("==== Reverse Whois for " + report.Target + " ====")
		fmt.Println("Domain\t\tDate\t\tRegistar")
		for _, i := range report.Domains {
			fmt.Println(i.DomainName + "\t" + i.CreationDate + "\t" + i.Registrar)
		}
	}
}
