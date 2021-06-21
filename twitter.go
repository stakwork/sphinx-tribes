package main

// https://github.com/twitterdev/Twitter-API-v2-sample-code/blob/master/User-Tweet-Timeline/user_tweets.js

import (
	"errors"
	"fmt"
	"os"

	"github.com/imroc/req"
)

func ConfirmIdentityTweet() (bool, error) {
	id, err := LookupUserID("TwitterDev")
	if err != nil {
		return false, err
	}
	LookupUserTweet(id)
	return false, nil
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

	fmt.Println(err)
	fmt.Printf("RES: %+v\n", re)

	return "", nil
}

func LookupUserTweet(userID string) error {

	twitterToken := os.Getenv("TWITTER_TOKEN")
	if twitterToken == "" {
		return errors.New("no twitter token")
	}

	// userID := "2244994945"

	url := "https://api.twitter.com/2/users/" + userID + "/tweets"
	url += "?max_results=100"
	url += "&tweet.fields=created_at"
	url += "&expansions=author_id"

	re, err := req.Get(url, req.Header{
		"User-Agent":    "sphinx_tribes",
		"authorization": "Bearer " + twitterToken,
	})

	fmt.Println(err)
	fmt.Printf("RES: %+v\n", re)

	return nil
}
