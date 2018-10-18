package axfr

import (
	"os"
	"os/exec"
	"testing"
)

func TestGetDomains(t *testing.T) {
	body := `{"issuer_ca_id":31337,"issuer_name":"C=US, O=Let's Encrypt, CN=Let's Encrypt Authority X3","name_value":"subdomain.domain.tld","min_cert_id":31337,"min_entry_timestamp":"2018-03-04T21:24:23.145","not_before":"2018-03-04T20:24:23","not_after":"2018-06-02T20:24:23"}`
	match := getDomains(body)
	for _, i := range match {
		subdomainTest := "subdomain.domain.tld"
		if i[1] != subdomainTest {
			t.Errorf("Incorrect Subdomain Match. Is %s, should be %s", i[1], subdomainTest)
		}
	}
}

func generateReport() *report {

	r := new(report)
	r.Domain = "google.com"
	domains := []string{"mail.google.com", "secret.google.com", "priv8.google.com"}
	for _, i := range domains {
		s := new(subdomain)
		s.URL = i
		r.Subdomains = append(r.Subdomains, *s)
	}
	return r
}

func Example_report_printReport() {
	r := generateReport()
	r.printReport(false, false)
	//Output:
	//==== Report for google.com ====
	//mail.google.com
	//secret.google.com
	//priv8.google.com
}

func Example_report_printReport_second() {
	r := generateReport()
	r.printReport(true, false)
	//Output:
	//{"domain":"google.com","subdomains":[{"URL":"mail.google.com","statusCode":0},{"URL":"secret.google.com","statusCode":0},{"URL":"priv8.google.com","statusCode":0}]}
}
func TestStartAXFR(t *testing.T) {
	opts := new(Options)
	opts.JSONFlag = (false)
	opts.URLs = []string{"twitter.com", "haveibeenpwned.com"}
	opts.StartAXFR()

}

func TestStartAXFR_second(t *testing.T) {
	opts := new(Options)
	opts.JSONFlag = (true)
	opts.URLs = []string{"twitter.com", "haveibeenpwned.com"}
	opts.StartAXFR()
}

//This funcion is really tricky, I hope there is a better way to get the exit status
func TestStartAXFR_third(t *testing.T) {
	if os.Getenv("TEST_CRASH") == "1" {
		opts := new(Options)
		opts.JSONFlag = (true)
		opts.URLs = []string{""}
		opts.StartAXFR()
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestStartAXFR_third")
	cmd.Env = append(os.Environ(), "TEST_CRASH=1")

	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Fatalf("Process ran with err %v, want exit status 1", err)
	}
}
