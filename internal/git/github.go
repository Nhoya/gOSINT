package git

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Nhoya/gOSINT/internal/utils"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func ghLogin() (*github.Client, context.Context) {
	var tc *http.Client
	ctx := context.Background()
	//load config values
	GHToken := utils.GetConfigValue("GHToken")
	if GHToken != "" {
		fmt.Println("[+] Found GitHub Token")
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: GHToken},
		)
		tc = oauth2.NewClient(ctx, ts)
	}
	client := github.NewClient(tc)
	return client, ctx
}

func getUsersFromGitHub(user string, repository string) [][]string {
	//init configuration file
	utils.WriteConfigFile("GHToken", "")
	client, ctx := ghLogin()
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var extractedValues [][]string
	//GitHub Pagination
	for {
		commits, resp, err := client.Repositories.ListCommits(ctx, user, repository, opt)
		if err != nil {
			fmt.Println("Unable to reach the repository")
			fmt.Println(err)
			os.Exit(1)
		}

		//Extract Author Name and Email
		for _, commit := range commits {
			var userTuple []string
			userTuple = append(userTuple, *commit.Commit.Author.Name)
			userTuple = append(userTuple, *commit.Commit.Author.Email)
			//add usersdata to extractedValues slice
			extractedValues = append(extractedValues, userTuple)
		}
		//check if this is the last page, if is exit
		if resp.NextPage == 0 {
			break
		}
		//Change page
		opt.Page = resp.NextPage
	}
	return extractedValues
}

func getGHUserRepositories(user string) []*github.Repository {
	opt := new(github.RepositoryListOptions)
	client, ctx := ghLogin()
	repos, _, err := client.Repositories.List(ctx, user, opt)
	if err != nil {
		fmt.Println(err)
	}
	return repos
}
