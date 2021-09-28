package main

import (
	"context"

	"github.com/google/go-github/v39/github"
)

type Issue struct {
	Title    string
	Status   string
	Assignee string
}

func GetRepoIssues(owner string, repo string) ([]Issue, error) {
	client := github.NewClient(nil)
	issues, _, err := client.Issues.ListByRepo(context.Background(), owner, repo, nil)
	ret := []Issue{}
	if err == nil {
		for _, iss := range issues {
			assignee := ""
			if iss.Assignee != nil {
				assignee = *iss.Assignee.Login
			}
			ret = append(ret, Issue{
				Title:    *iss.Title,
				Status:   *iss.State,
				Assignee: assignee,
			})
		}
	}
	return ret, err
}

func GetIssue(owner string, repo string, id int) (Issue, error) {
	client := github.NewClient(nil)
	iss, _, err := client.Issues.Get(context.Background(), owner, repo, id)
	issue := Issue{}
	if err == nil && iss != nil {
		assignee := ""
		if iss.Assignee != nil {
			assignee = *iss.Assignee.Login
		}
		issue = Issue{
			Title:    *iss.Title,
			Status:   *iss.State,
			Assignee: assignee,
		}
	}
	return issue, err
}
