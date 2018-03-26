package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/deckarep/golang-set"
)

type AxfrReport struct {
	Domain     string     `json:'domain'`
	Subdomains mapset.Set `json:'subdomains'`
}

func initAXFR() {
	if opts.Target == "" {
		fmt.Println("You must specify domain target")
		os.Exit(1)
	}
	enumerateSubdomains(opts.Target)
}
func enumerateSubdomains(domain string) {
	report := new(AxfrReport)
	report.Domain = domain
	querystring := "https://crt.sh/?q=%." + report.Domain + "&output=json"
	body := retrieveRequestBody(querystring)
	subdomains := mapset.NewSet()
	//get subsomains, their json output is not standard and i'm too lazy to write a parser :)
	re := regexp.MustCompile(`"name_value":"([^"]+)"`)
	match := re.FindAllStringSubmatch(body, -1)
	for _, i := range match {
		subdomains.Add(i[1])
	}
	report.Subdomains = subdomains
	report.printAXFRReport()
}

func (report *AxfrReport) printAXFRReport() {
	if opts.JSON {
		jsonreport, _ := json.MarshalIndent(&report, "", " ")
		fmt.Println(string(jsonreport))
	} else {
		fmt.Println("==== Report for " + report.Domain + " ====")
		di := report.Subdomains.Iterator()
		for dom := range di.C {
			fmt.Println(dom)
		}
	}
}
