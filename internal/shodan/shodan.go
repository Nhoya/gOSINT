package shodan

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Nhoya/gOSINT/internal/utils"
	"gopkg.in/ns3777k/go-shodan.v2/shodan"
)

//Options contains the options for the shodan scan module
type Options struct {
	Hosts        []string
	NewScan      bool
	HoneyPotFlag bool
}

//QueryOptions contains the options for the shodan query module
type QueryOptions struct {
	Query string
}

//StartShodanScan is the init function of the shodan scan module
func (opts *Options) StartShodanScan() {
	//init the configuration file
	utils.WriteConfigFile("shodanApiKey", "")
	//get the API Key from the configuration file
	APIKey := getShodanAPIKey()

	client := shodan.NewClient(nil, APIKey)

	if opts.NewScan {
		newShodanScan(client, opts.Hosts)
	} else {
		for _, host := range opts.Hosts {
			getShodanHostInfo(host, client, opts.HoneyPotFlag)
		}
	}
}

//StartShodanQuery is the init function of the shodan query module
func (opts *QueryOptions) StartShodanQuery() {
	//init the configuration file
	utils.WriteConfigFile("shodanApiKey", "")
	//get the API Key from the configuration file

	if len(opts.Query) < 1 {
		fmt.Println("[-] You need to specify the target")
		os.Exit(1)
	}
	APIKey := getShodanAPIKey()
	client := shodan.NewClient(nil, APIKey)
	getQueryInfo(client, opts.Query)
}

func getQueryInfo(client *shodan.Client, queryTarget string) {
	options := shodan.HostQueryOptions{}
	options.Query = queryTarget
	query, err := client.GetHostsForQuery(&options)
	if err != nil {
		utils.Panic(err, "Unable to get hosts for query")
	}
	if query.Total != 0 {
		fmt.Printf("==== Query result for \"%s\" ====\n", queryTarget)
		for _, host := range query.Matches {
			fmt.Println("Host:", host.IP, host.Hostnames)
			if host.OS != "" {
				fmt.Println("\tOS: ", host.OS)
			}
			fmt.Printf("\tLocation: %s, %s\n", host.Location.Country, host.Location.City)
		}
	} else {
		fmt.Println("[-] No results found")
	}
}

func getShodanAPIKey() string {
	APIKey := utils.GetConfigValue("shodanapikey")
	if APIKey == "" {
		fmt.Println("[-] API KEY Can't be empty")
		os.Exit(1)
	}
	fmt.Println("[+] APIKey Found")
	return APIKey
}

func getShodanHostInfo(host string, client *shodan.Client, honeyPotFlag bool) {
	fmt.Println("==== REPORT FOR " + host + " ====")
	report, err := client.GetServicesForHost(host, new(shodan.HostServicesOptions))
	if err != nil {
		utils.Panic(err, "Unable to get report")
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

	if honeyPotFlag {
		checkHoneyPotProbability(client, host)
	}
}

func getShodanServicesData(services []*shodan.HostData) {
	for _, service := range services {
		if service.Product == "" {
			service.Product = "Unknown"
		}
		fmt.Printf("Service on port %d: %s %s\n", service.Port, service.Product, service.Version)
		if service.Title != "" {
			fmt.Printf("\tTitle: %s\n", service.Title)
		}
		if service.OS != "" {
			fmt.Printf("\tOS %s\n", service.OS)
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
		utils.Panic(err, "Unable to schedule scan")
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
		utils.Panic(err, "Unable to get Honeypot Score")
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
