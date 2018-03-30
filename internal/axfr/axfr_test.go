package axfr

import (
	"fmt"
	"testing"
)

func TestGetDomains(t *testing.T) {
	body := `{"issuer_ca_id":31337,"issuer_name":"C=US, O=Let's Encrypt, CN=Let's Encrypt Authority X3","name_value":"subdomain.domain.tld","min_cert_id":31337,"min_entry_timestamp":"2018-03-04T21:24:23.145","not_before":"2018-03-04T20:24:23","not_after":"2018-06-02T20:24:23"}`
	match := getDomains(body)
	for _, i := range match {
		fmt.Println("in")
		subdomainTest := "subdomain.domain.tld"
		if i[1] != subdomainTest {
			t.Errorf("Incorrect Subdomain Match. Is %s, should be %s", i[1], subdomainTest)
		}
	}
}
