package utils

import "os"

var (
	// ConfigFilePath defines the path where the standard configuration file will be written
	ConfigFilePath = os.Getenv("HOME") + "/Library/Application Support/gOSINT/config/"
)
