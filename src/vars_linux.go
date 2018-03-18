package main

import "os"

var (
	// TelegramDumpPath defines the path for delegram dumps
	TelegramDumpPath = os.Getenv("HOME") + "/.local/share/gOSINT/tgdumps/"
	// ConfigFilePath defines the path of the config file
	ConfigFilePath = os.Getenv("HOME") + "/.config/gOSINT.conf"
)
