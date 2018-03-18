package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"os"

	"github.com/google/go-github/github"
)

func checkGHLogin(ctx context.Context) *http.Client {
	GHToken := getConfigFile().GHToken
	if GHToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: GHToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		return tc
	}
	return nil
}

func getUsersFromGitHub(user string, repository string) [][]string {
	ctx := context.Background()
	client := github.NewClient(checkGHLogin(ctx))

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
