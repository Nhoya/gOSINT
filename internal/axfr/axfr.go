package axfr

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/Nhoya/gOSINT/internal/utils"
	"github.com/deckarep/golang-set"
)

//Options contains the options for the axfr module
type Options struct {
	URLs     []string
	JSONFlag bool
}

//AXFRReport defines the report for the subdomains enumeration
type report struct {
	Domain     string     `json:"domain"`
	Subdomains mapset.Set `json:"subdomains"`
}

//StartAXFR is the init function of the module
func (opts *Options) StartAXFR() {
	report := new(report)
	for _, dom := range opts.URLs {
		if dom == "" {
			fmt.Println("You must specify domain target")
			os.Exit(1)
		}
		report.enumerateSubdomains(dom, opts.JSONFlag)
	}
	report.printReport(opts.JSONFlag)
}

func (report *report) enumerateSubdomains(domain string, jsonFlag bool) {
	report.Domain = domain
	querystring := "https://crt.sh/?q=%." + report.Domain + "&output=json"
	body := utils.RetrieveRequestBody(querystring)
	subdomains := mapset.NewSet()

	match := getDomains(body)
	for _, i := range match {
		//skip *.domain.tld
		if i[1] == "*."+domain {
			continue
		}
		subdomains.Add(i[1])
	}
	report.Subdomains = subdomains
}

//get subsomains, their json output is not standard and i'm too lazy to write a parser :)
func getDomains(body string) [][]string {
	re := regexp.MustCompile(`"name_value":"([^"]+)"`)
	match := re.FindAllStringSubmatch(body, -1)
	return match
}

func (report *report) printReport(jsonFlag bool) {
	if jsonFlag {
		jsonreport, _ := json.Marshal(&report)
		fmt.Println(string(jsonreport))
	} else {
		fmt.Println("==== Report for " + report.Domain + " ====")
		di := report.Subdomains.Iterator()
		for dom := range di.C {
			fmt.Println(dom)
		}
	}
}
