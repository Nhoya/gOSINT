package main

import "os"

var (
	TelegramDumpPath = os.Getenv("HOME") + "/.local/share/gOSINT/tgdumps/"
	ConfigFilePath   = os.Getenv("HOME") + "/.config/gOSINT.conf"
)
