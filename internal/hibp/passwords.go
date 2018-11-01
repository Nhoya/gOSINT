package hibp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

func getPasswords(mail string) ([]string, error) {
	client := &http.Client{}
	data := url.Values{}
	data.Set("query", mail)
	data.Add("fuck", "off")
	data.Add("search", "")

	req, _ := http.NewRequest("POST", "http://dumpedlqezarfife.onion.ws", bytes.NewBufferString(data.Encode()))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pwds []string
	re := regexp.MustCompile(mail + ".*:(.*)")
	match := re.FindAllStringSubmatch(string(body), -1)
	for _, i := range match {
		pwds = append(pwds, i[1])
	}

	return pwds, nil

}
