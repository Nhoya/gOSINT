package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"os"
	"strconv"
	"strings"
)

func gitSearch(target string, WebsiteAPI string, mailSet mapset.Set) mapset.Set {
	domain := ""
	targetSplit := strings.Split(target, "/")
	commits := ""

	fmt.Println("==== GIT SEARCH FOR " + target + " ==== ")

	//If using GitHub API
	if strings.Contains(target, "https://github.com") || WebsiteAPI == "github" {
		fmt.Println("[+] Using github API")
		domain = targetSplit[0] + "//api." + targetSplit[2] + "/repos/" + targetSplit[3] + "/" + targetSplit[4] + "/commits?per_page=100"
		//GitHub Pagination
		lastPage := retriveLastPage(domain)
		fmt.Println("[+] Looping through pages.This MAY take a while...")
		for i := 1; i < lastPage+1; i++ {
			commits = retriveRequestBody(domain + "&page=" + strconv.Itoa(i))
			findMailInText(commits, mailSet)
		}
	} else if strings.Contains(target, "https://bitbucket.org") || WebsiteAPI == "bitbucket" {
		// If using BitBucket API
		fmt.Println("[+] Using bitbucket API")
		domain = targetSplit[0] + "//api." + targetSplit[2] + "/2.0/repositories/" + targetSplit[3] + "/" + targetSplit[4] + "/commits?per_page=5000"
		//TODO: add BitBucket Pagination: https://developer.atlassian.com/bitbucket/api/2/reference/meta/pagination
		commits = retriveRequestBody(domain)
		findMailInText(commits, mailSet)
	} else {
		commits = cloneAndSearchCommit(target)
		findMailInText(commits, mailSet)
	}

	//Check if the mailset has been populated (this avoids problems with mispelled repositories too)
	if mailSet == nil {
		fmt.Println("[-] Nothing Found")
		os.Exit(1)
	}
	fmt.Println("[+] Mails Found")
	readFromSet(mailSet)
	return mailSet
}

func cloneAndSearchCommit(Url string) string {
	fmt.Println("[+] Cloning Repo")
	r, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: Url,
	})
	ref, _ := r.Head()
	cIter, _ := r.Log(&git.LogOptions{From: ref.Hash()})

	commits := ""
	_ = cIter.ForEach(func(c *object.Commit) error {
		commits += fmt.Sprintf("%s", c.Author)
		return nil
	})

	return commits
}
