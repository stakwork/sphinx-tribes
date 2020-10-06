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

const apiKey = "BVJTWLERYJXHWA7MYWXV"

const baseURL = "https://api.podcastindex.org/api/1.0/"

func unix() string {
	return strconv.Itoa(int(int32(time.Now().Unix())))
}

func makeHeaders() map[string]string {
	apiSecret := os.Getenv("PODCAST_INDEX_SECRET")
	fmt.Println("SERCRET", apiSecret)
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
func getLatestEpisode(feedURL string) (*PodcastEpisode, error) {

	client := &http.Client{}
	url := baseURL + "episodes/byfeedurl?url=" + feedURL
	fmt.Println("URL", url)
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

	if r.Items != nil && len(r.Items) > 0 {
		latest := r.Items[0]
		fmt.Printf("%+v\n", latest)
		return &latest, nil
	}
	return nil, errors.New("no items")
}

type PodcastEpisode struct {
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
type PodcastResponse struct {
	Items []PodcastEpisode `json:"items"`
	Count uint             `json:"count"`
}
