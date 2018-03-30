package telegram

import "os"

var (
	// TelegramDumpPath defines the path for delegram dumps
	TelegramDumpPath = os.Getenv("HOME") + "/.local/share/gOSINT/tgdumps/"
)
