package main

import (
	"bufio"
	"fmt"
	"github.com/deckarep/golang-set"
	"io/ioutil"
	"os"
	"path/filepath"
)

func checkFile(mailSet mapset.Set) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file, _ := ioutil.ReadFile(path)
			//fmt.Println(path)
			findMailInText(string(file), mailSet)

		}
		return nil
	}
}

func plainMailSearch(path string, mailSet mapset.Set, confirm bool) mapset.Set {
	tmpSet := mapset.NewSet()
	diffSet := mapset.NewSet()
	fmt.Println("[+] Searching for mail in " + path)
	filepath.Walk(path, checkFile(tmpSet))
	if confirm {
		fmt.Println("confirm?")
		tmpIt := tmpSet.Iterator()
		scanner := bufio.NewScanner(os.Stdin)
		for tmpMail := range tmpIt.C {
			fmt.Println("[?] Found " + tmpMail.(string) + " remove it?[Y/n]")
			scanner.Scan()
			text := scanner.Text()
			if text == "y" || text == "Y" {
				diffSet.Add(tmpMail)
			}
		}
	}
	tmpSet = tmpSet.Difference(diffSet)
	mailSet = mailSet.Union(tmpSet)
	readFromSet(mailSet)
	return mailSet
}
