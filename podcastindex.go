package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stakwork/sphinx-tribes/feeds"
)

func getFeed(feedURL string, feedID string) (*feeds.Podcast, error) {
	client := &http.Client{}

	url := ""
	if feedURL != "" {
		url = feeds.PodcastIndexBaseURL + "podcasts/byfeedurl?url=" + feedURL
	} else if feedID != "" {
		url = feeds.PodcastIndexBaseURL + "podcasts/byfeedid?id=" + feedID
	}
	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	req, err := http.NewRequest("GET", url, nil)

	headers := feeds.PodcastIndexHeaders()
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("GET error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r feeds.PodcastResponse
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshall error", err)
		return nil, err
	}

	feed := r.Feed
	tribe := DB.getFirstTribeByFeedURL(r.Feed.URL)
	feed.Value = feeds.AddedValue(r.Feed.Value, tribe.OwnerPubKey)

	return &feed, nil
}
func getEpisodes(feedURL string, feedID string) ([]feeds.Episode, error) {
	client := &http.Client{}

	url := ""
	if feedURL != "" {
		url = feeds.PodcastIndexBaseURL + "episodes/byfeedurl?url=" + feedURL
	} else if feedID != "" {
		url = feeds.PodcastIndexBaseURL + "episodes/byfeedid?id=" + feedID
	}
	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	req, err := http.NewRequest("GET", url, nil)

	headers := feeds.PodcastIndexHeaders()
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("GET error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r feeds.EpisodeResponse
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshall error", err)
		return nil, err
	}

	return r.Items, nil
}

func searchPodcastIndex(term string) ([]feeds.Podcast, error) {
	client := &http.Client{}

	url := feeds.PodcastIndexBaseURL + "search/byterm?q=" + term

	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	req, err := http.NewRequest("GET", url, nil)

	headers := feeds.PodcastIndexHeaders()
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("GET error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r feeds.PodcastSearchResponse
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshall error", err)
		return nil, err
	}

	return r.Feeds, nil
}
