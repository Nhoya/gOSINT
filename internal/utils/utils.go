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

//DEBUG_FLAG if true will trigger debug actions like coredump printing and more detailed errors
var DebugFlag bool

//RetrieveRequestBody send a GET Request to a domain and return the body casted to string
func RetrieveRequestBody(domain string) []byte {
	resp, err := http.Get(domain)
	if err != nil {
		Panic(err, "Unable to send request")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

//WriteOnFile open a file with Append and write on it, if the file doesn't exist will create it
func WriteOnFile(filename string, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		Panic(err, "Unable to Open file")
	}
	_, err = f.WriteString(text)
	if err != nil {
		Panic(err, "Unable to Write file")
	}
}

//FileExists return true if the path exists
func FileExists(file string) bool {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
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
	fmt.Printf("[?] %s [Y/N]\n", question)
	var resp string
	_, err := fmt.Scanln(&resp)
	if err != nil {
		Panic(err, "Unable to read answer")
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
	v.AddConfigPath("/etc/gosint")
	v.AddConfigPath(os.Getenv("HOME") + "/.config/")

	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("[-] Unable to find config file")
		fmt.Println("[+] Creating config file in "+ConfigFilePath + "" + ConfigFile + ".toml")
		
		f,err := os.Create(ConfigFilePath + "" + ConfigFile + ".toml")
		if err != nil {
			fmt.Println("Unable to generate config file")
		}
		defer f.Close()
		v.ReadInConfig()	
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

//AddToSMap put a string in a map if is not already present
func AddToSMap(elem string, m map[string]bool) map[string]bool {
	if _, present := m[elem]; !present {
		m[elem] = true
	}
	return m
}

//MapSDifference takes two matps [string]bool as agument and return a third map containing m1 - m2
func MapSDifference(m1 map[string]bool, m2 map[string]bool) map[string]bool {
	diffMap := make(map[string]bool)
	for key := range m1 {
		if _, present := m2[key]; !present {
			diffMap[key] = true
		} else {
			delete(m2, key)
		}
	}
	return diffMap
}

//Panic is a wrapper function on top of builtin panic,
//if DEBUG_FLAG is true it will print the core dump,
//otherwise it will print a message and exit with exitcode 1
func Panic(err error, msg string) {
	fmt.Println(msg)
	if DebugFlag {
		panic(err)
	} else {
		os.Exit(1)
	}
}
