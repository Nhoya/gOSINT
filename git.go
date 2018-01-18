package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/deckarep/golang-set"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func initGit(mailSet mapset.Set) {
	if opts.Url == "" {
		fmt.Println("You must specify target URL")
		os.Exit(1)
	}
	mailSet = gitSearch(opts.Url, opts.GitAPIType, mailSet)
	if opts.Mode {
		mailSet = pgpSearch(mailSet)
		pwnd(mailSet)
	}
}

func gitSearch(target string, WebsiteAPI string, mailSet mapset.Set) mapset.Set {
	// TODO: add worker for pagination
	domain := ""
	targetSplit := strings.Split(target, "/")
	commits := ""

	fmt.Println("==== GIT SEARCH FOR " + target + " ==== ")

	//If using GitHub API
	if strings.Contains(target, "https://github.com") || WebsiteAPI == "github" {
		fmt.Println("[+] Using github API")
		domain = targetSplit[0] + "//api." + targetSplit[2] + "/repos/" + targetSplit[3] + "/" + targetSplit[4] + "/commits?per_page=100"
		//GitHub Pagination
		lastPage := retrieveLastGHPage(domain)
		fmt.Println("[+] Looping through pages.This MAY take a while...")
		for page := 1; page < lastPage+1; page++ {
			fmt.Println("[+] Analyzing commits page: " + strconv.Itoa(page))
			commits = retrieveRequestBody(domain + "&page=" + strconv.Itoa(page))
			findMailInText(commits, mailSet)
		}
	} else if strings.Contains(target, "https://bitbucket.org") || WebsiteAPI == "bitbucket" {
		// If using BitBucket API
		fmt.Println("[+] Using bitbucket API")
		domain = targetSplit[0] + "//api." + targetSplit[2] + "/2.0/repositories/" + targetSplit[3] + "/" + targetSplit[4] + "/commits?pagelen=100"
		//BitBucket Pagination
		page := 1
		fmt.Println("[+] Looping through pages.This MAY take a while...")
		for page != 0 {
			fmt.Println("[+] Analyzing commits page: " + strconv.Itoa(page))
			pageDom := domain + "&page=" + strconv.Itoa(page)
			//This is needed because we can't unluckily retrieve max_page from one single request
			pageContent := retrieveRequestBody(pageDom)
			nextPage := "\"next\": \"" + domain + "&page="

			findMailInText(pageContent, mailSet)
			if strings.Contains(pageContent, nextPage) {
				page++
			} else {
				page = 0
			}
		}
	} else {
		commits = cloneAndSearchCommit(target)
		findMailInText(commits, mailSet)
	}

	//Check if the mailset has been populated (this avoids problems with misspelled repositories too)
	if mailSet == nil {
		fmt.Println("[-] Nothing Found")
		os.Exit(1)
	}
	fmt.Println("[+] Mails Found")
	readFromSet(mailSet)
	return mailSet
}

func retrieveLastGHPage(domain string) int {
	req, err := http.Get(domain)
	if err != nil {
		panic(err)
	}
	pagInfo := req.Header.Get("Link")
	if pagInfo != "" {
		re := regexp.MustCompile(`page=(\d+)>;\srel="last"`)
		match := re.FindStringSubmatch(pagInfo)
		lastPage, _ := strconv.Atoi(match[1])
		return lastPage
	}
	return 1
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
