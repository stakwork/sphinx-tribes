package main

// https://github.com/twitterdev/Twitter-API-v2-sample-code/blob/master/User-Tweet-Timeline/user_tweets.js

import (
	"errors"
	"os"
	"strings"

	"github.com/imroc/req"
)

func ConfirmIdentityTweet(username string) (string, error) {
	id, err := LookupUserID(username)
	if err != nil {
		return "", err
	}
	token, err := LookupUserTweet(id)
	if err != nil {
		return "", err
	}
	pubkey, err := VerifyArbitrary(token, "Sphinx Verification")
	return pubkey, err
}

type UserResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
}
type LookupUserIDResponse struct {
	Data []UserResponse `json:"data"`
}

func LookupUserID(username string) (string, error) {

	twitterToken := os.Getenv("TWITTER_TOKEN")
	if twitterToken == "" {
		return "", errors.New("no twitter token")
	}

	url := "https://api.twitter.com/2/users/by?usernames=" + username
	url += "&user.fields=created_at,description"
	url += "&expansions=pinned_tweet_id"

	re, err := req.Get(url, req.Header{
		"User-Agent":    "sphinx_tribes",
		"authorization": "Bearer " + twitterToken,
	})

	if err != nil {
		return "", err
	}
	var res LookupUserIDResponse
	re.ToJSON(&res)
	if res.Data == nil {
		return "", errors.New("no data")
	}
	if len(res.Data) == 0 {
		return "", errors.New("no data")
	}
	if res.Data[0].ID == "" {
		return "", errors.New("no ID")
	}
	return res.Data[0].ID, nil
}

type TweetResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type LookupUserTweetResponse struct {
	Data []TweetResponse `json:"data"`
}

func LookupUserTweet(userID string) (string, error) {

	twitterToken := os.Getenv("TWITTER_TOKEN")
	if twitterToken == "" {
		return "", errors.New("no twitter token")
	}

	url := "https://api.twitter.com/2/users/" + userID + "/tweets"
	url += "?max_results=100"
	url += "&tweet.fields=created_at"
	url += "&expansions=author_id"

	re, err := req.Get(url, req.Header{
		"User-Agent":    "sphinx_tribes",
		"authorization": "Bearer " + twitterToken,
	})

	if err != nil {
		return "", err
	}
	var res LookupUserTweetResponse
	re.ToJSON(&res)
	if res.Data == nil {
		return "", errors.New("no data")
	}
	if len(res.Data) == 0 {
		return "", errors.New("no data")
	}

	prefixes := []string{
		"Sphinx Verification: ",
		"Sphinx verification: ",
		"sphinx verification: ",
	}

	for _, tweet := range res.Data {
		for _, prefix := range prefixes {
			if strings.HasPrefix(tweet.Text, prefix) {
				return strings.TrimPrefix(tweet.Text, prefix), nil
			}
		}
	}

	return "", errors.New("did not find the tweet")
}
