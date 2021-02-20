package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://api.podcastindex.org/api/1.0/"

func unix() string {
	return strconv.Itoa(int(int32(time.Now().Unix())))
}

func makeHeaders() map[string]string {
	apiKey := os.Getenv("PODCAST_INDEX_KEY")
	apiSecret := os.Getenv("PODCAST_INDEX_SECRET")
	ts := unix()
	s := apiKey + apiSecret + ts
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return map[string]string{
		"User-Agent":    "golang",
		"X-Auth-Date":   ts,
		"X-Auth-Key":    apiKey,
		"Authorization": fmt.Sprintf("%x", bs),
	}
}

func getFeed(feedURL string, feedID string) (*Podcast, error) {
	client := &http.Client{}

	url := ""
	if feedURL != "" {
		url = baseURL + "podcasts/byfeedurl?url=" + feedURL
	} else if feedID != "" {
		url = baseURL + "podcasts/byfeedid?id=" + feedID
	}
	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	req, err := http.NewRequest("GET", url, nil)

	headers := makeHeaders()
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("GET error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r PodcastResponse
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshall error", err)
		return nil, err
	}

	feed := addToFeed(r.Feed)

	return &feed, nil
}
func getEpisodes(feedURL string, feedID string) ([]Episode, error) {
	client := &http.Client{}

	url := ""
	if feedURL != "" {
		url = baseURL + "episodes/byfeedurl?url=" + feedURL
	} else if feedID != "" {
		url = baseURL + "episodes/byfeedid?id=" + feedID
	}
	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	req, err := http.NewRequest("GET", url, nil)

	headers := makeHeaders()
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("GET error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r EpisodeResponse
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshall error", err)
		return nil, err
	}

	return r.Items, nil
}

func addToFeed(feed Podcast) Podcast {
	if feed.URL == "" {
		return feed
	}
	tribe := DB.getTribeByFeedURL(feed.URL)
	if tribe.OwnerPubKey == "" {
		return feed
	}
	if feed.Value != nil {
		if feed.Value.Destinations != nil {
			if len(feed.Value.Destinations) == 1 {
				first := feed.Value.Destinations[0]
				firstSplit, _ := first.Split.Int64()
				if firstSplit == 1 {
					// this is the auto
					tribe := DB.getTribeByFeedURL(feed.URL)
					if tribe.OwnerPubKey != first.Address {
						feed.Value.Destinations = append(feed.Value.Destinations, Destination{
							Address: tribe.OwnerPubKey,
							Split:   json.Number(99),
							Type:    "node",
						})
					}
				}
			}
		}
	} else {
		feed.Value = &Value{
			Model: Model{
				Type:      "lightning",
				Suggested: "0.00000015000",
			},
			Destinations: []Destination{
				Destination{
					Address: tribe.OwnerPubKey,
					Type:    "node",
					Split:   json.Number(100),
				},
			},
		}
	}
	return feed
}

type PodcastResponse struct {
	Feed Podcast `json:"feed"`
}
type Podcast struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	URL            string    `json:"url"`
	Description    string    `json:"description"`
	Author         string    `json:"author"`
	Image          string    `json:"image"`
	Link           string    `json:"link"`
	LastUpdateTime int32     `json:"lastUpdateTime"`
	ContentType    string    `json:"contentType"`
	Language       string    `json:"language"`
	Episodes       []Episode `json:"episodes"`
	Value          *Value    `json:"value"`
}
type EpisodeResponse struct {
	Items []Episode `json:"items"`
	Count uint      `json:"count"`
}
type Episode struct {
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	DatePublished   int32  `json:"datePublished"`
	EnclosureURL    string `json:"enclosureUrl"`
	EnclosureType   string `json:"enclosureType"`
	EnclosureLength int32  `json:"enclosureLength"`
	Image           string `json:"image"`
	Link            string `json:"link"`
}
type Value struct {
	Model        Model         `json:"model"`
	Destinations []Destination `json:"destinations"`
}
type Model struct {
	Type      string `json:"type"`
	Suggested string `json:"suggested"`
}
type Destination struct {
	Address string      `json:"address"`
	Split   json.Number `json:"split"`
	Type    string      `json:"type"`
}
