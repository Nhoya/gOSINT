package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/deckarep/golang-set"
	"gopkg.in/src-d/go-git.v4"
)

func initPlainSearch(mailSet mapset.Set) {
	if opts.Clone {
		if opts.Url == "" {
			fmt.Println("You must specify target URL")
			os.Exit(1)
		}
		mailSet = cloneAndSearch(opts.Url, mailSet, opts.Confirm)
	} else {
		if opts.Path == "" {
			fmt.Println("You must specify Path")
			os.Exit(1)
		}
		mailSet = plainMailSearch(opts.Path, mailSet, opts.Confirm)
	}
	if opts.Mode {
		mailSet = pgpSearch(mailSet)
		pwnd(mailSet)
	}
}
func checkFile(mailSet mapset.Set) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file, _ := ioutil.ReadFile(path)
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
		for tmpMail := range tmpIt.C {
			resp := simpleQuestion("Found " + tmpMail.(string) + " remove it?")
			if resp {
				diffSet.Add(tmpMail)
			}
		}
	}
	tmpSet = tmpSet.Difference(diffSet)
	mailSet = mailSet.Union(tmpSet)
	readFromSet(mailSet)
	return mailSet
}

func cloneRepo(repo string) string {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Println("Unable to create tmp directory")
		os.Exit(1)
	}
	fmt.Println("[+] Cloning Repo")
	_, err = git.PlainClone(tmpdir, false, &git.CloneOptions{
		URL:      repo,
		Progress: os.Stdout,
	})

	if err != nil {
		fmt.Println("[-] Unable to clone clone the repo")
		os.Exit(1)
	}
	return tmpdir
}

func cloneAndSearch(repo string, mailSet mapset.Set, confirm bool) mapset.Set {
	tmpdir := cloneRepo(repo)
	defer func() {
		fmt.Println("[+] Deleting Repo")
		os.RemoveAll(tmpdir)
		fmt.Println("[+] Done")
	}()

	mailSet = plainMailSearch(tmpdir, mailSet, confirm)
	return mailSet
}
