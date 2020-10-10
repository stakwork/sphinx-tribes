package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const apiKey = "BVJTWLERYJXHWA7MYWXV"

const baseURL = "https://api.podcastindex.org/api/1.0/"

func unix() string {
	return strconv.Itoa(int(int32(time.Now().Unix())))
}

func makeHeaders() map[string]string {
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

func getFeed(feedURL string) (*Podcast, error) {
	client := &http.Client{}
	url := baseURL + "podcasts/byfeedurl?url=" + feedURL
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

	return &r.Feed, nil
}
func getEpisodes(feedURL string) ([]Episode, error) {

	client := &http.Client{}
	url := baseURL + "episodes/byfeedurl?url=" + feedURL
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
	Address string `json:"address"`
	Split   uint   `json:"split"`
	Type    string `json:"type"`
}
