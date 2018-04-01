package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/deckarep/golang-set"
	"github.com/spf13/viper"
)

//ConfigFile is the name of the configuration file
const ConfigFile = "gosint"

//RetrieveRequestBody send a Get Request to a domain and return the body casted to string
func RetrieveRequestBody(domain string) string {
	resp, err := http.Get(domain)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

//WriteOnFile open a file with Append and write on it, if the file doesn't exist will create it
func WriteOnFile(filename string, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Unabale to open file")
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = f.WriteString(text)
	if err != nil {
		fmt.Println("Unable to wite on file")
	}
}

//FileExists return true if the path exists
func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

//CreateDirectory creates a directory in the path passed as argument
func CreateDirectory(dirname string) {
	if !FileExists(dirname) {
		fmt.Println("[+] Creating directory " + dirname)
		os.MkdirAll(dirname, os.ModePerm)
	}
}

//SimpleQuestion prompt a simple Y/N question
func SimpleQuestion(question string) bool {
	fmt.Println("[?] " + question + " [Y/N]")
	var resp string
	_, err := fmt.Scanln(&resp)
	if err != nil {
		fmt.Println("[-] Unable to read answer")
		os.Exit(1)
	}
	if resp == "y" || resp == "Y" {
		return true
	}
	return false
}

func readConfigFile() *viper.Viper {
	v := viper.New()
	v.SetConfigName(ConfigFile)
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("$HOME/.config")
	v.AddConfigPath("/etc/gosint")

	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return v
}

//WriteConfigFile will write default values in the config file specifiend in var_os.go
func WriteConfigFile(key string, value string) {
	v := readConfigFile()
	v.SetDefault(key, value)
	v.WriteConfigAs(ConfigFilePath + "" + ConfigFile + ".toml")
}

//GetConfigValue take as input a key value and will return the relative value set in the configuration file
func GetConfigValue(key string) string {
	v := readConfigFile()
	a := v.GetString(key)
	return a
}

//SetToSlice convers a set in a slice of strings and then order them
func SetToSlice(oldset mapset.Set) []string {
	var newSlice []string
	setIt := oldset.Iterator()
	for elem := range setIt.C {
		newSlice = append(newSlice, elem.(string))
	}
	sort.Strings(newSlice)
	return newSlice
}
