package main

import (
	"context"

	"github.com/google/go-github/v39/github"
)

type GithubIssue struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Assignee string `json:"assignee"`
}

func GetRepoIssues(owner string, repo string) ([]GithubIssue, error) {
	client := github.NewClient(nil)
	issues, _, err := client.Issues.ListByRepo(context.Background(), owner, repo, nil)
	ret := []GithubIssue{}
	if err == nil {
		for _, iss := range issues {
			assignee := ""
			if iss.Assignee != nil {
				assignee = *iss.Assignee.Login
			}
			ret = append(ret, GithubIssue{
				Title:    *iss.Title,
				Status:   *iss.State,
				Assignee: assignee,
			})
		}
	}
	return ret, err
}

func GetIssue(owner string, repo string, id int) (GithubIssue, error) {
	client := github.NewClient(nil)
	iss, _, err := client.Issues.Get(context.Background(), owner, repo, id)
	issue := GithubIssue{}
	if err == nil && iss != nil {
		assignee := ""
		if iss.Assignee != nil {
			assignee = *iss.Assignee.Login
		}
		issue = GithubIssue{
			Title:    *iss.Title,
			Status:   *iss.State,
			Assignee: assignee,
		}
	}
	return issue, err
}
