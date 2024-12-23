package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/go-github/v39/github"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
	"golang.org/x/oauth2"
)

func GetGithubIssue(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	issueString := chi.URLParam(r, "issue")
	issueNum, err := strconv.Atoi(issueString)
	if err != nil || issueNum < 1 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	issue, err := GetIssue(owner, repo, issueNum)
	if err != nil {
		utils.Log.Error("Github error: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(issue)
}

func GetOpenGithubIssues(w http.ResponseWriter, r *http.Request) {
	issue_count, err := db.DB.GetOpenGithubIssues(r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	json.NewEncoder(w).Encode(issue_count)
}

func githubClient() *github.Client {
	gh_token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh_token},
	)
	tc := oauth2.NewClient(ctx, ts)
	gc := github.NewClient(tc)
	return gc
}

func GetRepoIssues(owner string, repo string) ([]db.GithubIssue, error) {
	client := githubClient()
	issues, _, err := client.Issues.ListByRepo(context.Background(), owner, repo, nil)
	ret := []db.GithubIssue{}
	if err == nil {
		for _, iss := range issues {
			assignee := ""
			if iss.Assignee != nil {
				assignee = *iss.Assignee.Login
			}
			ret = append(ret, db.GithubIssue{
				Title:    *iss.Title,
				Status:   *iss.State,
				Assignee: assignee,
			})
		}
	}
	return ret, err
}

func GetIssue(owner string, repo string, id int) (db.GithubIssue, error) {
	client := githubClient()
	iss, _, err := client.Issues.Get(context.Background(), owner, repo, id)
	issue := db.GithubIssue{}
	if err == nil && iss != nil {
		assignee := ""
		if iss.Assignee != nil {
			assignee = *iss.Assignee.Login
		}
		issue = db.GithubIssue{
			Title:       *iss.Title,
			Status:      *iss.State,
			Assignee:    assignee,
			Description: *iss.Body,
		}
	}
	return issue, err
}

func PubkeyForGithubUser(owner string) (string, error) {
	client := githubClient()
	gs, _, err := client.Gists.List(context.Background(), owner, nil)
	if err == nil && gs != nil {
		for _, g := range gs {
			if g.Files != nil {
				for k := range g.Files {
					if strings.Contains(string(k), "Sphinx Verification") {
						// get the actual gist
						gist, _, err := client.Gists.Get(context.Background(), *g.ID)
						gistFile := gist.Files[k]
						pubkey, err := auth.VerifyArbitrary(*gistFile.Content, "Sphinx Verification")
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
