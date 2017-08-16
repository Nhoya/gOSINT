package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"strings"
)

func gitSearch(target string, WebsiteAPI string, mailSet mapset.Set) mapset.Set {
	domain := ""
	rawParse := false
	targetSplit := strings.Split(target, "/")

	//TODO: work on pagination but i'm too lazy :(
	fmt.Println("==== GIT SEARCH FOR " + target + " ==== ")
	if strings.Contains(target, "https://github.com") || WebsiteAPI == "github" {
		fmt.Println("[+] Using github API")
		domain = targetSplit[0] + "//api." + targetSplit[2] + "/repos/" + targetSplit[3] + "/" + targetSplit[4] + "/commits?per_page=5000"
	} else if strings.Contains(target, "https://bitbucket.org") || WebsiteAPI == "bitbucket" {
		fmt.Println("[+] Using bitbucket API")
		domain = targetSplit[0] + "//api." + targetSplit[2] + "/2.0/repositories/" + targetSplit[3] + "/" + targetSplit[4] + "/commits?per_page=5000"
	} else {
		rawParse = true
	}

	commits := ""
	if rawParse {
		commits = cloneAndSearchCommit(target)
	} else {
		commits = retriveRequestBody(domain)
	}

	mailSet = findMailInText(commits, mailSet)
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
