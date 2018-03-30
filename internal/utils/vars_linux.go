package utils

import "os"

var (
	// ConfigFilePath defines the path of the config file
	ConfigFilePath = os.Getenv("HOME") + "/.config/gOSINT.conf"
)
