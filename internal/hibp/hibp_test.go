package hibp

import "testing"

func generateReport() *report {
	e := new(PwnedEntity)
	e.Breaches = []string{"google.com", "facebook.com"}
	e.Email = ("mail@domain.tld")

	r := new(report)
	r.Pwnd = append(r.Pwnd, e)
	return r
}

func Example_report_printReport() {
	r := generateReport()
	r.printHIBPReport(false)
	//Output:
	//Mail: mail@domain.tld
	//Breaches [google.com facebook.com]
}

func Example_report_printReport_second() {
	r := generateReport()
	r.printHIBPReport(true)
	//Output:
	//{"pwnd":[{"email":"mail@domain.tld","breaches":["google.com","facebook.com"]}]}
}

func TestStartHIBP(t *testing.T) {
	opts := new(Options)
	opts.JSONFlag = true
	opts.Mails = []string{"me@gmail.com", "hi@gmail.com"}
	opts.StartHIBP()
}

func Test_report_getBreachesForMail(t *testing.T) {
	r := new(report)
	r.getBreachesForMail("hiall@gmail.com")
	if r.Pwnd == nil {
		t.Errorf("The report should contains breaches, it is empty!\n")
	}
}
