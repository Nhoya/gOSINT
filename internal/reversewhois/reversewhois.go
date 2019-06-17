package reversewhois

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/Nhoya/gOSINT/internal/utils"
)

//Options contains the options for reversewhois module
type Options struct {
	Target   string
	JSONFlag bool
}

type report struct {
	Target  string    `json:"target"`
	Domains []*Domain `json:"domains"`
}

//Domain contains the filed relative to a domain
type Domain struct {
	DomainName   string `json:"domainName"`
	CreationDate string `json:"creationDate"`
	Registrar    string `json:"registrar"`
}

//StartReverseWhois will start the ReverseWhois module
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
	body := string(utils.RetrieveRequestBody("https://viewdns.info/reversewhois/?q=" + url.QueryEscape(target)))
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

func extractValues(body string) [][]string {
	re := regexp.MustCompile(`(?mU)<\/tr><tr><td>([^\s]+)<\/td><td>(\d{4}-\d{2}-\d{2})<\/td><td>(.*)<\/td>`)
	match := re.FindAllStringSubmatch(body, -1)
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
