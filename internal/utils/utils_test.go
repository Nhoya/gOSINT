package utils

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/deckarep/golang-set"
)

func TestRetrieveRequestBody(t *testing.T) {
	body := RetrieveRequestBody("https://pastebin.com/raw/hbT8ATxJ")
	bodyTest := "Working"
	if body != bodyTest {
		t.Errorf("Incorrect request body. Is %s, should be %s", body, bodyTest)
	}
}

func TestFileExists(t *testing.T) {
	tmpfile, _ := ioutil.TempFile(".", "test")
	if !FileExists(tmpfile.Name()) {
		t.Errorf("File Exists functions is broken")
	}
	os.RemoveAll(tmpfile.Name())
}

func TestCreateDirectory(t *testing.T) {
	dirTest := strconv.Itoa(rand.Int()) //totally not secure, but is just a test :)
	CreateDirectory(dirTest)
	if _, err := os.Stat(dirTest); os.IsNotExist(err) {
		t.Errorf("%s has not been created", dirTest)
	} else {
		os.RemoveAll(dirTest)
	}
}

func TestSetToSlice(t *testing.T) {
	newSet := mapset.NewSet()
	newSet.Add("Aab")
	newSet.Add("bbbA")
	newSet.Add("213494")
	newSlice := SetToSlice(newSet)
	testSlice := []string{"213494", "Aab", "bbbA"}
	for i, e := range newSlice {
		if e != testSlice[i] {
			t.Errorf("Slices are different. Is %s, should be %s", newSlice, testSlice)
			break
		}
	}
}

func TestWriteOnFile(t *testing.T) {
	tmpFile, _ := ioutil.TempFile(".", "test")
	stringTest := "Working"
	WriteOnFile(tmpFile.Name(), stringTest)
	b := make([]byte, 7)
	tmpFile.Read(b)
	if string(b) != stringTest {
		t.Errorf("Wrong file content. Is %s, should be %s", string(b), stringTest)
	}
	os.Remove(tmpFile.Name())

}

func Test_readConfigFile(t *testing.T) {
	os.Create("gosint.toml")
	configTest := "testconfig=1"
	b := []byte(configTest + "\n")
	ioutil.WriteFile("gosint.toml", b, 0644)
	v := readConfigFile()
	a := v.GetInt("testconfig")
	if a != 1 {
		t.Errorf("Wrong configuration key. Is %d, should be 1", a)
	}
	os.Remove("gosint.toml")
}

func TestWriteConfig(t *testing.T) {
	os.Create("gosint.toml")
	ConfigFilePath = "./"
	WriteConfigFile("test", "1")
	dat, _ := ioutil.ReadFile("gosint.toml")
	if string(dat) != "test = \"1\"\n" {
		t.Errorf("Wrong configuration. Is %s, should be test = \"1\"", dat)
	}
	os.Remove("gosint.toml")
}

//todo: Write SimpleQuestion test
