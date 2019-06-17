package axfr

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/Nhoya/gOSINT/internal/utils"
)

//Options contains the options for the axfr module
type Options struct {
	URLs            []string
	VerifyURLStatus bool
	JSONFlag        bool
}

//AXFRReport defines the report for the subdomains enumeration
type report struct {
	Domain     string      `json:"domain"`
	Subdomains []subdomain `json:"subdomains"`
}

//Subdomain are represented by their URL and HTTP status code
type subdomain struct {
	URL        string `json:"URL"`
	StatusCode int    `json:"statusCode"`
}

//StartAXFR is the init function of the module
func (opts *Options) StartAXFR() {
	report := new(report)
	for _, dom := range opts.URLs {
		if dom == "" {
			fmt.Println("You must specify domain target")
			os.Exit(1)
		}
		report.enumerateSubdomains(dom, opts.JSONFlag, opts.VerifyURLStatus)
	}
	report.printReport(opts.JSONFlag, opts.VerifyURLStatus)
}

func (report *report) enumerateSubdomains(domain string, jsonFlag bool, verifyURLStatus bool) {
	report.Domain = domain
	querystring := "https://crt.sh/?q=%." + report.Domain + "&output=json"
	body := utils.RetrieveRequestBody(querystring)

	match := getDomains(string(body))
	doms := make(map[string]bool)
	for _, i := range match {
		if _, ok := doms[i[1]]; !ok {
			doms[i[1]] = true
		}
	}

	for i := range doms {
		subdomain := new(subdomain)
		subdomain.URL = i
		if verifyURLStatus {
			log.Println("checking", i)
			client := http.Client{
				Timeout: time.Duration(2 * time.Second),
			}
			resp, err := client.Get("https://" + i)
			if err != nil {
				subdomain.StatusCode = -1
			} else {
				subdomain.StatusCode = resp.StatusCode
			}
		}
		report.Subdomains = append(report.Subdomains, *subdomain)
	}

}

//get subsomains, their json output is not standard and i'm too lazy to write a parser :)
func getDomains(body string) [][]string {
	re := regexp.MustCompile(`"name_value":"([^\*\.][^"]+)"`)
	match := re.FindAllStringSubmatch(body, -1)
	return match
}

func (report *report) printReport(jsonFlag bool, verifyURLStatus bool) {
	if jsonFlag {
		jsonreport, _ := json.Marshal(&report)
		fmt.Println(string(jsonreport))
	} else {
		fmt.Println("==== Report for " + report.Domain + " ====")
		if verifyURLStatus {
			for _, dom := range report.Subdomains {
				fmt.Println(dom.URL, dom.StatusCode)
			}
		} else {
			for _, dom := range report.Subdomains {
				fmt.Println(dom.URL)
			}
		}
	}
}
