package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/go-github/v39/github"
)

func getGithubIssue(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	issueString := chi.URLParam(r, "issue")
	issueNum, err := strconv.Atoi(issueString)
	if err != nil || issueNum < 1 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	issue, err := GetIssue(owner, repo, issueNum)
	json.NewEncoder(w).Encode(issue)
}

func getOpenGithubIssues(w http.ResponseWriter, r *http.Request) {
	issue_count, err := DB.getOpenGithubIssues(r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	json.NewEncoder(w).Encode(issue_count)
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
			Title:       *iss.Title,
			Status:      *iss.State,
			Assignee:    assignee,
			Description: *iss.Body,
		}
	}
	return issue, err
}

func PubkeyForGithubUser(owner string) (string, error) {
	client := github.NewClient(nil)
	gs, _, err := client.Gists.List(context.Background(), owner, nil)
	if err == nil && gs != nil {
		for _, g := range gs {
			if g.Files != nil {
				for k := range g.Files {
					if strings.Contains(string(k), "Sphinx Verification") {
						// get the actual gist
						gist, _, err := client.Gists.Get(context.Background(), *g.ID)
						gistFile := gist.Files[k]
						pubkey, err := VerifyArbitrary(*gistFile.Content, "Sphinx Verification")
						if err != nil {
							return "", err
						}
						return pubkey, nil
					}
				}
			}
		}
	}
	return "", errors.New("nope")
}
