package main

import (
	"fmt"
	"gopkg.in/ns3777k/go-shodan.v2/shodan"
	"os"
	"strconv"
)

func initShodan() {
	if opts.ShodanTarget == "" {
		fmt.Println("[-] You need to specify the target")
		os.Exit(1)
	}
	APIKey := getConfigFile().ShodanAPIKey
	if APIKey == "" {
		fmt.Println("[-] Unable to retrive Shodan API Key from config file")
		os.Exit(1)
	}
	getShodanHostInfo(opts.ShodanTarget, APIKey)
}

func getShodanHostInfo(target string, APIKey string) {
	client := shodan.NewClient(nil, APIKey)
	report, err := client.GetServicesForHost(target, &shodan.HostServicesOptions{false, false})
	if err != nil {
		fmt.Println("[-] Unable to get Report")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("==== REPORT FOR " + target + " ====")
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
	}
}
