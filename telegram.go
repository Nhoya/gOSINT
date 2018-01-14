package main

import (
	"fmt"
	"github.com/jaytaylor/html2text"
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
			for j := 0; j < graceCounter; j++ {
				fmt.Println("[MESSAGE REMOVED]")
			}
			graceCounter = 0
			username, nickname := getTelegramUsername(body)
			date, time := getTelegramMessageDateTime(body)
			if username == "" {
				ret = "[" + date + " " + time + "] " + nickname + ": " + message
			} else {
				ret = "[" + date + " " + time + "] " + nickname + "(" + username + "): " + message
			}
			msg, _ := html2text.FromString(ret)
			fmt.Println(msg)
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
	match := re.FindAllStringSubmatch(body, -1)
	messageBody := ""
	if len(match) == 1 {
		messageBody = match[0][1]
	} else if len(match) == 2 {
		messageBody = "{" + match[0][1] + "}" + match[1][1]
	}
	messageBody = messageBody + getTelegramMedia(body)
	return messageBody

}

func getTelegramMedia(body string) string {
	messageBody := getTelegramVideo(body) + getTelegramPhoto(body) + getTelegramVoice(body)
	return messageBody
}

func getTelegramPhoto(body string) string {
	re := regexp.MustCompile(`image:url\('https:\/\/([\w+.\/-]+)'`)
	match := re.FindStringSubmatch(body)
	if len(match) == 2 {
		return "Photo: " + match[1]
	}
	return ""
}

func getTelegramVoice(body string) string {
	re := regexp.MustCompile(`voice"\ssrc="(https:\/\/[\w.\/-]+)"`)
	match := re.FindStringSubmatch(body)
	if len(match) == 2 {
		return "Voice: " + match[1]
	}
	return ""
}

func getTelegramVideo(body string) string {
	re := regexp.MustCompile(`video\ssrc="(https:\/\/[\w.\/-]+)"`)
	match := re.FindStringSubmatch(body)
	if len(match) == 2 {
		return "Video: " + match[1]
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
