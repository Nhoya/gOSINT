package telegram

import "os"

var (
	// TelegramDumpPath defines the path for delegram dumps
	TelegramDumpPath = os.Getenv("HOME") + "/Library/Application Support/gOSINT/tgdumps/"
)
