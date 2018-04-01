package pgp

import (
	"os"
	"os/exec"
	"testing"
)

func TestExtractPGPIdentities(t *testing.T) {
	body := `"search=0x31337">John Doe &lt;testingmail@mail.it&gt;</a>
    Super Secret Agent &lt;007@NSA.gov&gt;"`

	r := new(report)
	r.extractPGPIdentities(body)
	for _, e := range r.Entities {
		keyTest := "0x31337"
		if e.KeyID != keyTest {
			t.Errorf("Incorrect Key ID. Is %s, should be %s", e.KeyID, keyTest)
		}
		nameTest := "John Doe"
		if e.Person.Name != nameTest {
			t.Errorf("Incorrect Name. Is %s, should be %s", e.Person.Name, nameTest)
		}
		mailTest := "testingmail@mail.it"
		if e.Person.Email != mailTest {
			t.Errorf("Incorrect Email. Is %s, should be %s", e.Person.Email, mailTest)
		}
		for _, a := range e.Aliases {
			aliasNameTest := "Super Secret Agent"
			if a.Name != aliasNameTest {
				t.Errorf("Incorrect Alias Name. Is %s, should be %s", a.Name, aliasNameTest)
			}
			aliasMailTest := "007@NSA.gov"
			if a.Email != aliasMailTest {
				t.Errorf("Incorrect Alias Email. Is %s, should be %s", a.Email, aliasMailTest)
			}
		}
	}

}

func generateReport() *report {
	a := new(alias)
	a.Email = "mail@domain.tld"
	a.Name = "John Doe"

	b := new(alias)
	b.Email = "anothermail@domain.tld"
	b.Name = "Secret Agent"

	e := new(entity)
	e.Person = *a
	e.KeyID = "0x31337"
	e.Aliases = append(e.Aliases, *a)

	r := new(report)
	r.Target = "mail@domain.tld"
	r.Entities = append(r.Entities, e)

	return r
}

func generateEmptyReport() *report {
	r := new(report)
	r.Target = "mail@domain.tld"
	r.Entities = nil

	return r
}

func Example_report_printReport_first() {
	r := generateReport()
	r.printReport(false)
	//Output:
	//==== PGP SEARCH FOR: mail@domain.tld====
	//0x31337 {John Doe mail@domain.tld}
	//	Alias: {John Doe mail@domain.tld}
}

func Example_report_printReport_second() {
	r := generateReport()
	r.printReport(true)
	//Output:
	//{"Target":"mail@domain.tld","Entities":[{"Person":{"Name":"John Doe","Email":"mail@domain.tld"},"KeyID":"0x31337","Aliases":[{"Name":"John Doe","Email":"mail@domain.tld"}]}]}
}

func Test_report_printReport(t *testing.T) {
	if os.Getenv("TEST_CRASH") == "1" {
		r := generateEmptyReport()
		r.printReport(false)
	}

	cmd := exec.Command(os.Args[0], "-test.run=Test_report_printReport")
	cmd.Env = append(os.Environ(), "TEST_CRASH=1")

	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Fatalf("Process ran with err %v, want exit status 1", err)
	}
}

func TestOptionsStartPGP(t *testing.T) {
	opts := new(Options)
	opts.JSONFlag = true
	opts.Targets = []string{"hackingteam.it"}

	opts.StartPGP()

}
