package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/deckarep/golang-set"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type GitReport struct {
	Repository   string         `json:"repository"`
	Data         []*GitUserData `json:"data,omitempty"`
	InvalidUSers mapset.Set     `json:"invalidUsers,omitempty"`
}

type GitUserData struct {
	Mails   mapset.Set `json:"mails,omitempty"`
	Aliases mapset.Set `json:"aliases,omitempty"`
}

func initGit() {
	if opts.URL == "" {
		fmt.Println("You must specify target URL")
		os.Exit(1)
	}
	gitSearch()
}

func gitSearch() {
	domain, _ := url.Parse(opts.URL)
	targetSplit := strings.Split(domain.Path, "/")
	user := targetSplit[1]
	repository := targetSplit[2]
	var extractedValues [][]string
	//If using GitHub API
	if (domain.Host == "github.com" && opts.GitAPIType != "clone") || opts.GitAPIType == "github" {
		fmt.Println("[+] Using github API")
		extractedValues = getUsersFromGitHub(user, repository)
	} else {
		extractedValues = cloneAndSearchCommit(opts.URL)
	}

	rawUserMap, invalidUsers := parseGitEntities(extractedValues)
	refinedUserMap := removeGitDuplicates(rawUserMap)
	report := generateGitReport(refinedUserMap, invalidUsers)
	report.printGitReport()
}

func cloneAndSearchCommit(URL string) [][]string {
	//create temporary file
	tmpdir, err := ioutil.TempDir(".", ".gOSINT")
	defer os.Remove(tmpdir)
	if err != nil {
		fmt.Println("Unable create temporary directory")
		os.Exit(1)
	}

	fmt.Println("[+] Cloning Repo")
	r, _ := git.PlainClone(tmpdir, false, &git.CloneOptions{
		URL:      URL,
		Progress: os.Stdout,
		//don't waste space on disk
		NoCheckout: true,
	})

	ref, err := r.Head()
	if err != nil {
		fmt.Println("Unable to clone Repository")
		os.Exit(1)
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	var extractedValues [][]string

	_ = cIter.ForEach(func(c *object.Commit) error {
		var userdata []string
		userdata = append(userdata, c.Author.Name)
		userdata = append(userdata, c.Author.Email)
		extractedValues = append(extractedValues, userdata)
		return nil
	})
	//delete the directory
	return extractedValues
}

func parseGitEntities(extractedValues [][]string) (map[string]*GitUserData, mapset.Set) {
	rawUserMap := make(map[string]*GitUserData)
	invalidUsers := mapset.NewSet()
	for _, commit := range extractedValues {
		name, mail := commit[0], commit[1]
		//if is a new username we need to create a new entity
		if isValidMail(mail) {
			if _, present := rawUserMap[name]; !present {
				newData := buildNewGitEntity()
				rawUserMap[name] = newData
			}
			rawUserMap[name].Mails.Add(mail)
			rawUserMap[name].Aliases.Add(name)
		} else {
			invalidUsers.Add(name)
		}
	}
	//removing valid users form invalidUSersSet
	diffSet := mapset.NewSet()
	it := invalidUsers.Iterator()
	for name := range it.C {
		if _, present := rawUserMap[name.(string)]; present {
			diffSet.Add(name)
		}
	}

	return rawUserMap, invalidUsers.Difference(diffSet)
}

func removeGitDuplicates(rawUserMap map[string]*GitUserData) mapset.Set {
	refinedUserMap := make(map[string]*GitUserData)
	for _, user := range rawUserMap {
		it1 := user.Mails.Iterator()
		for mail := range it1.C {
			if user2, ok := refinedUserMap[mail.(string)]; ok {
				user2.Mails = user2.Mails.Union(user.Mails)
				user2.Aliases = user2.Aliases.Union(user.Aliases)
				//copy aliases for struct with the same mail inside
				it2 := user.Aliases.Iterator()
				for alias := range it2.C {
					refinedUserMap[alias.(string)] = user2
				}
				//put mail on the same struct
				it3 := user.Mails.Iterator()
				for mail := range it3.C {
					refinedUserMap[mail.(string)] = user2
				}
			} else {
				refinedUserMap[mail.(string)] = user
			}
		}
	}
	//we need this set to get rid of the latest duplicates struct
	//now they are ordered and uniqe
	gitEntitiesSet := mapset.NewSet()
	for _, data := range refinedUserMap {
		gitEntitiesSet.Add(data)
	}
	return gitEntitiesSet

}

func generateGitReport(gitEntitiesSet mapset.Set, invalidUsers mapset.Set) *GitReport {
	it := gitEntitiesSet.Iterator()
	gitReport := new(GitReport)
	gitReport.Repository = opts.URL
	gitReport.InvalidUSers = invalidUsers
	for gitStruct := range it.C {
		gitReport.Data = append(gitReport.Data, gitStruct.(*GitUserData))
	}
	return gitReport
}

func (report *GitReport) printGitReport() {
	if opts.JSON {
		jsonreport, _ := json.MarshalIndent(&report, "", " ")
		fmt.Println(string(jsonreport))
	} else {
		fmt.Println("==== GIT SEARCH FOR " + report.Repository + " ====")
		fmt.Println("Users found:", len(report.Data))
		for _, userdata := range report.Data {
			fmt.Printf("User Alias: %v Email Address: %v\n", userdata.Aliases.ToSlice(), userdata.Mails.ToSlice())
		}
		if report.InvalidUSers.Cardinality() != 0 {
			fmt.Printf("Users with invalid mail: %v\n", report.InvalidUSers.ToSlice())
		}
	}
}

func buildNewGitEntity() *GitUserData {
	m := new(GitUserData)
	a := mapset.NewSet()
	b := mapset.NewSet()
	m.Mails = a
	m.Aliases = b
	return m
}

func isValidMail(mail string) bool {
	return !(strings.HasSuffix(mail, "@users.noreply.github.com") || strings.HasSuffix(mail, ".local"))
}
