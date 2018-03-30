package pgp

import "testing"

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
