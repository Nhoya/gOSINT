package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func getTelegramGroupHistory(group string, grace int) {
	i := 1
	graceCounter := 0
	ret := ""
	for i != 0 {
		messageid := strconv.Itoa(i)
		body := retriveRequestBody("https://t.me/" + group + "/" + messageid + "?embed=1")
		message := getTelegramMessage(body)
		if message != "" {
			username, nickname := getTelegramUsername(body)
			date, time := getTelegramMessageDateTime(body)
			if username == "" {
				ret = "[" + date + " " + time + "]" + nickname + ": " + message
			} else {
				ret = "[" + date + " " + time + "]" + nickname + "(" + username + "): " + message
			}
			fmt.Println(ret)
		} else {
			graceCounter++
			if graceCounter == grace {
				break
			}
		}
		i++
		time.Sleep(time.Millisecond * 500)
	}
}

func getTelegramMessage(body string) string {
	re := regexp.MustCompile(`class=\"tgme_widget_message_text\" dir=\"auto\">(.*)<\/div>\n`)
	//match := re.FindStringSubmatch(body)
	//if len(match) == 2 {
	//	return match[1]
	//} else {
	//	return ""
	//}
	match := re.FindAllStringSubmatch(body, -1)
	if len(match) == 1 {
		return match[0][1]
	} else if len(match) == 2 {
		return ">" + match[0][1] + "\n" + match[1][1]
	}
	return ""

}

func getTelegramUsername(body string) (string, string) {
	re := regexp.MustCompile(`class=\"tgme_widget_message_author_name\" (?:href="https://t\.me/(\w+)")? dir=\"auto\">(.*)</a>&nbsp;in&nbsp;<a`)
	match := re.FindStringSubmatch(body)
	return match[1], match[2]
}

func getTelegramMessageDateTime(body string) (string, string) {
	re := regexp.MustCompile(`<time datetime="(\d+-\d+-\d+)T(\d+:\d+:\d+)\+\d+:\d+\">`)
	match := re.FindStringSubmatch(body)
	return match[1], match[2]
}
