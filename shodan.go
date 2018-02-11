package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"gopkg.in/ns3777k/go-shodan.v2/shodan"
)

func initShodan() {
	if opts.ShodanTarget != nil {
		initShodanHost()
	} else if opts.ShodanQuery != "" {
		initShodanQuery()
	} else {
		fmt.Println("[-] Invalid action only --query and --target are available")
	}
}

func initShodanHost() {
	if len(opts.ShodanTarget) < 1 {
		fmt.Println("[-] You need to specify the target")
		os.Exit(1)
	}
	APIKey := getShodanAPIKey()
	client := shodan.NewClient(nil, APIKey)

	if opts.ShodanScan {
		newShodanScan(client, opts.ShodanTarget)
	} else {
		for _, host := range opts.ShodanTarget {
			getShodanHostInfo(host, client)
		}
	}
}

func initShodanQuery() {
	if len(opts.ShodanQuery) < 1 {
		fmt.Println("[-] You need to specify the target")
		os.Exit(1)
	}
	APIKey := getShodanAPIKey()
	client := shodan.NewClient(nil, APIKey)
	getQueryInfo(client, opts.ShodanQuery)
}

func getQueryInfo(client *shodan.Client, queryTarget string) {
	options := shodan.HostQueryOptions{}
	options.Query = queryTarget
	query, err := client.GetHostsForQuery(&options)
	if err != nil {
		fmt.Println(err)
	}
	if query.Total != 0 {
		fmt.Println("==== Query result for \"" + queryTarget + "\" ====")
		for _, host := range query.Matches {
			fmt.Println("Host:", host.IP, host.Hostnames)
			if host.OS != "" {
				fmt.Println("\tOS: ", host.OS)
			}
			fmt.Println("\tLocation: " + host.Location.Country + ", " + host.Location.City)
		}
	} else {
		fmt.Println("[-] No results found")
	}
}

func getShodanAPIKey() string {
	APIKey := getConfigFile().ShodanAPIKey
	if APIKey == "" {
		fmt.Println("[-] Unable to retrive Shodan API Key from config file")
		os.Exit(1)
	}
	fmt.Println("[+] APIKey Found")
	return APIKey
}

func getShodanHostInfo(host string, client *shodan.Client) {
	fmt.Println("==== REPORT FOR " + host + " ====")
	report, err := client.GetServicesForHost(host, &shodan.HostServicesOptions{false, false})
	if err != nil {
		fmt.Println("[-] Unable to get Report")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("ISP: " + report.ISP)
	fmt.Println("Organization: " + report.Organization)
	if report.OS != "" {
		fmt.Println("OS: " + report.OS)
	}
	fmt.Println("Ports:", report.Ports)
	fmt.Println("Hostnames:", report.Hostnames)
	if len(report.Vulnerabilities) > 0 {
		fmt.Println("Vulnerabilities:", report.Vulnerabilities)
	}
	fmt.Println("Country:", report.HostLocation.Country)
	fmt.Println("City:", report.HostLocation.City)
	fmt.Println("Last Update: " + report.LastUpdate)
	getShodanServicesData(report.Data)

	if opts.ShodanHoneyPotFlag {
		checkHoneyPotProbability(client, host)
	}
}

func getShodanServicesData(services []*shodan.HostData) {
	for _, service := range services {
		if service.Product == "" {
			service.Product = "Unknown"
		}
		fmt.Println("Service on port " + strconv.Itoa(service.Port) + ": " + service.Product + " " + string(service.Version))
		if service.Title != "" {
			fmt.Println("\tTitle: " + service.Title)
		}
		if service.OS != "" {
			fmt.Println("\tOS " + service.OS)
		}
		if service.Data != "" {
			getServiceFingerprint(service.Data)
		}
	}
}

func newShodanScan(client *shodan.Client, hosts []string) {
	info, _ := client.GetAPIInfo()
	if info.ScanCredits < len(hosts) {
		fmt.Println("[-] Insufficient credits")
		os.Exit(1)
	}
	fmt.Println("[+] Current Scan credits:", info.ScanCredits)
	fmt.Println("[+] Requesting new scan")
	scan, err := client.Scan(hosts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("[+] Scan request ID: %s (1 credit will be deducted)\n", scan.ID)
	for {
		status, _ := client.GetScanStatus(scan.ID)
		if status.Status == shodan.ScanStatusDone {
			fmt.Println("[+] Scan started, the new result will be available in ~30 minutes")
			break
		}
	}
}

func checkHoneyPotProbability(client *shodan.Client, host string) {
	honeyscore, err := client.CalcHoneyScore(host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Honeypot Score (0-1):", honeyscore)
}

func getServiceFingerprint(serviceData string) {
	re := regexp.MustCompile(`Fingerprint\:\s(.*)`)
	match := re.FindStringSubmatch(serviceData)
	if len(match) == 2 {
		fmt.Println("\tFingerprint: " + match[1])
	}
}
