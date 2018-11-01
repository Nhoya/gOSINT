package pni

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/Nhoya/gOSINT/internal/utils"
	"github.com/otiai10/gosseract"
)

//SyncMeAnswer contains a (partial) representation of the SyncMe JSON answer
type SyncMeAnswer struct {
	ErrorCode   int    `json:"error_code"`
	PremiumType int    `json:"premium_type"`
	Name        string `json:"name"`
}

const (
	ua = "Mozilla/5.0 (Linux; Android 8.1.0; LG-D802 Build/OPM6.171019.030.K1) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/67.0.3396.87 Mobile Safari/537.36"
)

//Options contains the options for the PNI module
type Options struct {
	PhoneNumber []string
	JSONFlag    bool
}

//StartPNI is the init function of the module ¯\_(ツ)_/¯
func (opts *Options) StartPNI() {
	//TODO: add check on number lenght
	for _, num := range opts.PhoneNumber {
		if !strings.HasPrefix(num, "+") {
			fmt.Println(num + " is invalid, You must specify the country code with +")
			os.Exit(1)
		}
		retrievePhoneOwner(num[1:], opts.JSONFlag)
	}
}

func retrievePhoneOwner(target string, jsonFlag bool) {
	cj, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cj,
	}
	captchaID, err := sendGETRequest("https://sync.me/search/?number="+target, client)
	if err != nil {
		utils.Panic(err, "Unable to send GET captchaID")
	}

	//extract captcha ID
	re := regexp.MustCompile(`var\scaptchaURL\s=\s'(?:\/\.\.){2}/server/simple-php-captcha/simple-php-captcha\.php\?_CAPTCHA&amp;t=(.*)';`)
	match := re.FindAllStringSubmatch(string(captchaID), -1)

	fmt.Println("https://sync.me/server/simple-php-captcha/simple-php-captcha.php?_CAPTCHA&amp;t=" + match[0][1])
	os.Exit(1)
	//retrieve captcha
	challenge, err := sendGETRequest("https://sync.me/server/simple-php-captcha/simple-php-captcha.php?_CAPTCHA&amp;t="+match[0][1], client)
	if err != nil {
		utils.Panic(err, "Unable to get captcha")
	}
	//solve it
	solution := solveCaptcha(challenge)
	//send solution
	data := url.Values{}
	data.Set("action", "captcha")
	data.Add("data[g-recaptcha]", solution)
	data.Add("captchaResult", "")
	data.Add("isMobile", "true")

	_, err = sendPOSTRequest("https://sync.me/server/main.php", data, client, target)
	if err != nil {
		utils.Panic(err, "Unable to read response")
	}
	data = url.Values{}
	data.Set("action", "search")
	data.Add("number", target)
	report, err := sendPOSTRequest("https://sync.me/server/main.php", data, client, target)
	if err != nil {
		utils.Panic(err, "Unable to read response")
	}
	//TODO: move this on an indipended function
	if len(report) == 0 {
		fmt.Println("Unable to complete the challenge correctly. Please retry, if the error persist open an issue.")
	} else {
		if jsonFlag {
			fmt.Println(report)
		} else {
			var answer SyncMeAnswer
			err := json.Unmarshal(report, &answer)
			if err != nil {
				utils.Panic(err, "Unable to read JSON output")
			}
			if answer.ErrorCode == 8203 {
				fmt.Println("Unable to find owner")
			} else if answer.ErrorCode == 0 {
				fmt.Println(answer.Name)
			} else {
				utils.Panic(err, "Unexpected error occured, please retry using --debug flag and send the output as issue")
			}
		}
	}

}

func sendGETRequest(URL string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", ua)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func sendPOSTRequest(URL string, data url.Values, client *http.Client, target string) ([]byte, error) {
	req, _ := http.NewRequest("POST", URL, bytes.NewBufferString(data.Encode()))
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Referer", "https://sync.me/search/?number="+target)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}

func solveCaptcha(challenge []byte) string {
	//OCR init
	ocr := gosseract.NewClient()
	ocr.SetWhitelist("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	defer ocr.Close()
	ocr.SetImageFromBytes(challenge)
	text, _ := ocr.Text()
	return text
}
