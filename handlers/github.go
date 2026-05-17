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
	"github.com/stakwork/sphinx-tribes/logger"
	"golang.org/x/oauth2"
)

// GithubIssue godoc
//
//	@Summary		Get Github Issue
//	@Description	Get a Github issue by owner, repo, and issue number
//	@Tags			Github
//	@Accept			json
//	@Produce		json
//	@Param			owner	path		string	true	"Owner"
//	@Param			repo	path		string	true	"Repository"
//	@Param			issue	path		int		true	"Issue Number"
//	@Success		200		{object}	db.GithubIssue
//	@Router			/github_issues/{owner}/{repo}/{issue} [get]
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
		logger.Log.Error("Github error: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(issue)
}

// GetOpenGithubIssues godoc
//
//	@Summary		Get Open Github Issues
//	@Description	Get the count of open Github issues
//	@Tags			Github
//	@Accept			json
//	@Produce		json
//	@Success		200	{int}	int
//	@Router			/github_issues/status/open [get]
func GetOpenGithubIssues(w http.ResponseWriter, r *http.Request) {
	issue_count, err := db.DB.GetOpenGithubIssues(r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	json.NewEncoder(w).Encode(issue_count)
}

func githubClient() *github.Client {
	return githubClientWithToken(os.Getenv("GITHUB_TOKEN"))
}

func githubClientWithToken(gh_token string) *github.Client {
	gh_token = strings.TrimSpace(gh_token)
	if gh_token == "" {
		return github.NewClient(nil)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh_token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func GetRepoIssues(owner string, repo string) ([]db.GithubIssue, error) {
	client := githubClient()
	issues, _, err := client.Issues.ListByRepo(context.Background(), owner, repo, nil)
	ret := []db.GithubIssue{}
	if err == nil {
		for _, iss := range issues {
			ret = append(ret, githubIssueToDBIssue(iss))
		}
	}
	return ret, err
}

func GetIssue(owner string, repo string, id int) (db.GithubIssue, error) {
	client := githubClient()
	iss, _, err := client.Issues.Get(context.Background(), owner, repo, id)
	issue := db.GithubIssue{}
	if err == nil && iss != nil {
		issue = githubIssueToDBIssue(iss)
	}
	return issue, err
}

func githubIssueToDBIssue(issue *github.Issue) db.GithubIssue {
	if issue == nil {
		return db.GithubIssue{}
	}

	assignee := ""
	if issue.Assignee != nil {
		assignee = githubString(issue.Assignee.Login)
	}

	return db.GithubIssue{
		Title:       githubString(issue.Title),
		Status:      githubString(issue.State),
		Assignee:    assignee,
		Description: githubString(issue.Body),
	}
}

func githubString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
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
