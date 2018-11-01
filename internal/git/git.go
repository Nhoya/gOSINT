package git

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/Nhoya/gOSINT/internal/utils"
	"github.com/deckarep/golang-set"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

//Options contains the options needed for the git module
type Options struct {
	Repo      *url.URL
	Method    string
	Recursive bool
	JSONFlag  bool
}

// GitReport contains the struct of the git Report
type report struct {
	Repository   string      `json:"repository"`
	Data         []*userData `json:"data,omitempty"`
	InvalidUSers mapset.Set  `json:"invalidUsers,omitempty"`
}

//userData contains the user data like email address and aliases
type userData struct {
	Mails   mapset.Set `json:"mails,omitempty"`
	Aliases mapset.Set `json:"aliases,omitempty"`
}

//StartGit is the int module for the git Module
func (opts *Options) StartGit() {
	domainSplit := strings.Split(opts.Repo.Path, "/")
	if len(domainSplit) == 2 {
		user := domainSplit[1]
		host := opts.Repo.Host
		if host == "github.com" {
			fmt.Println("[+]Starting recursive git search")
			repos := getGHUserRepositories(user)
			opts.Method = "clone"
			for _, r := range repos {
				gitSearch(host, user, *r.Name, opts)
			}
		} else {
			fmt.Println("[-] Recursive search is only available on github.com repositories")
			os.Exit(1)
		}
	} else if len(domainSplit) == 3 {
		user := domainSplit[1]
		host := opts.Repo.Host
		repository := domainSplit[2]
		gitSearch(host, user, repository, opts)
	} else {
		fmt.Println("[-] Invalid URL")
		os.Exit(1)
	}
}

func gitSearch(host string, user string, repository string, opts *Options) {
	var extractedValues [][]string
	URL := "https://" + host + "/" + user + "/" + repository
	//If using GitHub API
	if (host == "github.com" && opts.Method != "clone") || opts.Method == "gh" {
		fmt.Println("[+] Using github API")
		extractedValues = getUsersFromGitHub(user, repository)
	} else {
		extractedValues = cloneAndSearchCommit(URL)
	}

	rawUserMap, invalidUsers := parseGitEntities(extractedValues)
	refinedUserMap := removeGitDuplicates(rawUserMap)
	report := generateGitReport(URL, refinedUserMap, invalidUsers)
	report.printReport(opts.JSONFlag)
}

func cloneAndSearchCommit(URL string) [][]string {
	fmt.Println(URL)
	//create temporary file
	tmpdir, err := ioutil.TempDir(".", ".gOSINT")
	if err != nil {
		utils.Panic(err, "Unable create temporary directory")
	}
	defer os.RemoveAll(tmpdir)

	fmt.Println("[+] Cloning Repo")
	r, _ := git.PlainClone(tmpdir, false, &git.CloneOptions{
		URL:      URL,
		Progress: os.Stdout,
		//don't waste space on disk
		NoCheckout: true,
	})

	ref, err := r.Head()
	if err != nil {
		utils.Panic(err, "Unable to clone Repository")
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

func parseGitEntities(extractedValues [][]string) (map[string]*userData, mapset.Set) {
	rawUserMap := make(map[string]*userData)
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

func removeGitDuplicates(rawUserMap map[string]*userData) mapset.Set {
	refinedUserMap := make(map[string]*userData)
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

func generateGitReport(URL string, gitEntitiesSet mapset.Set, invalidUsers mapset.Set) *report {
	it := gitEntitiesSet.Iterator()
	gitReport := new(report)
	gitReport.Repository = URL
	gitReport.InvalidUSers = invalidUsers
	for gitStruct := range it.C {
		gitReport.Data = append(gitReport.Data, gitStruct.(*userData))
	}

	return gitReport
}

func (report *report) printReport(jsonFlag bool) {
	if jsonFlag {
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

func buildNewGitEntity() *userData {
	m := new(userData)
	a := mapset.NewSet()
	b := mapset.NewSet()
	m.Mails = a
	m.Aliases = b
	return m
}

func isValidMail(mail string) bool {
	return !(strings.HasSuffix(mail, "@users.noreply.github.com") || strings.HasSuffix(mail, ".local"))
}
