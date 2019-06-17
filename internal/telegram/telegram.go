package telegram

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Nhoya/gOSINT/internal/utils"

	"github.com/jaytaylor/html2text"
)

//Options contains the options for the telegram module
type Options struct {
	Group       string
	Start       int
	End         int
	GracePeriod int
	DumpFlag    bool
}

//StartTelegram is the init function for the Telegram module
func (opts *Options) StartTelegram() {
	getTelegramGroupHistory(opts.Group, opts.GracePeriod, opts.DumpFlag, (opts.Start - 1), opts.End)
}

func getTelegramGroupHistory(group string, grace int, dumpFlag bool, startMessage int, endMessage int) {
	checkGroupName(group)
	// check if -e option is set
	if endMessage != 0 {
		// end can't be less than start
		if endMessage <= startMessage {
			fmt.Println("[-] The final message number (-e)  must be >= than start message number (-s)")
			os.Exit(1)
		}

		fmt.Println("[?] End  message set, grace time will be ignored")
	}

	//dump file
	dumpfile := TelegramDumpPath + group + ".dump"
	//counter for deleted messages
	dmCounter := 0
	//set messageCounter as startMessage, is -e is not used the default value of startMessage is 0 (Note: the first message on group is id:1)
	messageCounter := startMessage
	readFromTelegramDump(&startMessage, dumpfile, dumpFlag, &messageCounter)
	//add a counter to remember the first message
	firstMessageCounter := messageCounter - startMessage
	//this is needed because if a file is available it will start from the next to the last found
	messageCounter++
	//if -e or - s is set but on the dumpfile the message is already scraped
	if dumpFlag && ((endMessage != 0 && messageCounter >= endMessage) || (startMessage != 0 && messageCounter >= startMessage)) {
		fmt.Println("[-] You already have this messages")
		os.Exit(1)
	}
	startTime := time.Now()
	fmt.Println("==== [" + startTime.Format(time.RFC3339) + "] Dumping messages for " + group + " ====")

	//we don't know how many first how many messages the group has
	for {
		messageid := strconv.Itoa(messageCounter)
		body := string(utils.RetrieveRequestBody("https://t.me/" + group + "/" + messageid + "?embed=1"))
		message := getTelegramMessage(body)

		if message != "" && dmCounter > 0 {
			//this is to avoid to write on file the last n empty messages
			for j := 0; j < dmCounter; j++ {
				msg := "[MESSAGE REMOVED]"
				writeTelegramLogs(messageCounter, msg, dumpFlag, dumpfile)
			}
			dmCounter = 0
		} else if message != "" {
			//retrive the message the message message
			msg := createMessage(body, message)
			writeTelegramLogs(messageCounter, msg, dumpFlag, dumpfile)
		} else if messageCounter == 1 {
			//the first message is always a service message, if doesn't exist the groups is not valid
			fmt.Println("[!!] Invalid group")
			break
		} else if endMessage == 0 && dmCounter == grace {
			dmCounter++
			messageCounter = messageCounter - dmCounter
			break
		} else if endMessage == 0 {
			//if -e is not set and is not the last message increase the counter
			dmCounter++
		} else {
			//if -e is set and the message is empty dmCounter is 0 and grace is 0 so print the message
			msg := "[DELETED MESSAGE]"
			writeTelegramLogs(messageCounter, msg, dumpFlag, dumpfile)
		}
		// if this is the last message (defined  with -e) quit
		if endMessage != 0 && messageCounter == endMessage {
			break
		}
		messageCounter++
		time.Sleep(time.Millisecond * 100)
	}
	fmt.Println("==== [" + time.Now().Format(time.RFC3339) + " (elapsed:" + time.Since(startTime).String() + ")] End of history | " + strconv.Itoa(messageCounter-startMessage-firstMessageCounter) + " messages scraped  ==== ")
	if endMessage == 0 && messageCounter > 0 {
		fmt.Println("[=] If you think there are more messages try to increase the grace period (--grace [INT])")
	}
}

func getTelegramMessage(body string) string {
	re := regexp.MustCompile(`class=\"tgme_widget_message_text\" dir=\"auto\">(.*)<\/div>\n`)
	match := re.FindAllStringSubmatch(body, -1)
	messageBody := ""
	if len(match) == 1 {
		messageBody = match[0][1]
	} else if len(match) == 2 {
		quotedUser := getMessageRepliedAuthor(body)
		messageBody = "{ " + quotedUser + match[0][1] + " } " + match[1][1]
	}
	messageBody = messageBody + getTelegramMedia(body)
	return messageBody
}

func getTelegramMedia(body string) string {
	messageBody := getTelegramVideo(body) + getTelegramPhoto(body) + getTelegramVoice(body) + getTelegramServiceMessage(body) + getTelegramDocument(body)
	return messageBody
}

func getTelegramPhoto(body string) string {
	re := regexp.MustCompile(`image:url\('(https:\/\/[\w+.\/-]+)'`)
	match := re.FindStringSubmatch(body)
	if len(match) == 2 {
		return "Image: " + match[1]
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
	re := regexp.MustCompile(`video"?\ssrc="(https:\/\/[\w.\/-]+)"`)
	match := re.FindStringSubmatch(body)
	if len(match) == 2 {
		return "Video: " + match[1]
	}
	return ""
}

func getTelegramDocument(body string) string {
	re := regexp.MustCompile(`document_title" dir="auto">(.*)</div>`)
	match := re.FindStringSubmatch(body)
	if len(match) == 2 {
		return "Document: " + match[1]
	}
	return ""
}

func getTelegramServiceMessage(body string) string {
	re := regexp.MustCompile(`<div\sclass="message_media_not_supported_label">Service\smessage<\/div>`)
	if re.MatchString(body) {
		return "[SERVICE MESSAGE]"
	}
	return ""
}

func getMessageRepliedAuthor(body string) string {
	re := regexp.MustCompile(`reply"\shref="https:\/\/t.me\/[\w/]+">[\n\s]+<div\sclass="tgme_widget_message_author">[\n\s]+<span\sclass="tgme_widget_message_author_name"\s?dir="auto">(.*)</span>`)
	match := re.FindStringSubmatch(body)
	if len(match) == 2 {
		return " " + match[1] + ": "
	}
	return ""
}

func getTelegramUsername(body string) (string, string) {
	re := regexp.MustCompile(`class=\"tgme_widget_message_author_name\"\s?(?:href="https://t\.me/(\w+)")? dir=\"auto\">(.*)<\/(?:span>)?(?:a>)?&nbsp;in&nbsp;<a`)
	match := re.FindStringSubmatch(body)
	if len(match) == 3 {
		return match[1], match[2]
	}
	return "", ""
}

func getTelegramMessageDateTime(body string) (string, string) {
	re := regexp.MustCompile(`<time datetime="(\d+-\d+-\d+)T(\d+:\d+:\d+)\+\d+:\d+\">`)
	match := re.FindStringSubmatch(body)
	return match[1], match[2]
}

func checkGroupName(group string) {
	re := regexp.MustCompile(`^[[:alpha:]](?:\-?[[:alnum:]]){3,31}$`)
	if !re.MatchString(group) {
		fmt.Println("Invalid Group name, valid chars alphanum and -")
		os.Exit(1)
	}
}

func writeTelegramLogs(messageCounter int, msg string, dumpFlag bool, dumpfile string) {
	if dumpFlag {
		utils.WriteOnFile(dumpfile, fmt.Sprintf("[%d] %s\n", messageCounter, strings.Replace(msg, "\n", " ⏎ ", -1)))
	}
	fmt.Println(msg)
}

func createMessage(body string, message string) string {
	username, nickname := getTelegramUsername(body)
	date, time := getTelegramMessageDateTime(body)
	msgtxt := fmt.Sprintf("[%s %s] ", date, time)
	//channels don't have username and nickname
	if nickname == "" {
		msgtxt += message
	} else if username == "" {
		msgtxt += fmt.Sprintf("%s: %s", nickname, message)
	} else {
		msgtxt += fmt.Sprintf("%s(%s): %s", nickname, username, message)
	}
	//html format the message before printing it
	msg, _ := html2text.FromString(msgtxt)
	return msg
}

func readFromTelegramDump(startMessage *int, dumpfile string, dumpFlag bool, messageCounter *int) {
	if dumpFlag {
		utils.CreateDirectory(TelegramDumpPath)
		fmt.Println("[=] --dumpfile used, ignoring --startpoint")
		*startMessage = 0
		if utils.FileExists(dumpfile) {
			fmt.Println("[=] The dump will be saved in " + dumpfile)
			resp := utils.SimpleQuestion("Print the existing dump before resuming it?")
			fmt.Println("[+] Calculating the last message")
			file, _ := os.Open(dumpfile)
			scan := bufio.NewScanner(file)
			for scan.Scan() {
				messageSlice := strings.Split(scan.Text(), " ")
				if resp {
					fmt.Println(strings.Join(messageSlice[1:], " "))
				}
				*messageCounter, _ = strconv.Atoi(strings.Trim(messageSlice[0], "[]"))
			}
			fmt.Println("[=] Starting from message n.", *messageCounter)
		}
	}
}
