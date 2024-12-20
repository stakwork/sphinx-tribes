package utils

// https://github.com/twitterdev/Twitter-API-v2-sample-code/blob/master/User-Tweet-Timeline/user_tweets.js

import (
	"errors"
	"os"
	"strings"

	"github.com/imroc/req"
	"github.com/stakwork/sphinx-tribes/auth"
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
	Log.Info("Twitter verification token: %s", token)
	pubkey, err := auth.VerifyArbitrary(token, "Sphinx Verification")
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
		"Sphinx verify: ",
		"Sphinx Verify: ",
	}

	// these are used in case the prefix is buried in the middle of a tweet
	split_prefixes := []string{
		"Verification:",
		"verification:",
		"verify:",
		"Verify:",
	}

	for _, tweet := range res.Data {
		for _, prefix := range prefixes {
			if strings.Contains(tweet.Text, prefix) {
				// yes the tweet contains to prefix
				if strings.HasPrefix(tweet.Text, prefix) {
					// if not preceded by text, trim the prefix to return the code
					return strings.TrimPrefix(tweet.Text, prefix), nil
				} else {
					// split the tweet by spaces, compare to split_prefixes
					tweet_words := strings.Split(tweet.Text, " ")
					// find the prefix index, and return the following index (the verification code)
					for i, word := range tweet_words {
						for _, s_prefix := range split_prefixes {
							if word == s_prefix {
								// Found it!
								// make sure that there is a word after this one
								if (len(tweet_words) - 1) > i {
									// the next index should be the verification code
									next_index := (i + 1)
									return tweet_words[next_index], nil
								}

							}
						}
					}

				}
			}

		}
	}

	return "", errors.New("did not find the tweet")
}
