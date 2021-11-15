package feeds

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

const PodcastIndexBaseURL = "https://api.podcastindex.org/api/1.0/"

func unix() string {
	return strconv.Itoa(int(int32(time.Now().Unix())))
}

func PodcastToGeneric(url string, p *Podcast) (Feed, error) {
	items := []Item{}
	fmt.Println("P EPISODES", len(p.Episodes), p)
	for _, ep := range p.Episodes {
		items = append(items, Item{
			Id:              strconv.Itoa(int(ep.ID)),
			Link:            ep.Link,
			Description:     ep.Description,
			Title:           ep.Title,
			ImageUrl:        ep.Image,
			EnclosureURL:    ep.Link,
			EnclosureType:   ep.EnclosureType,
			EnclosureLength: ep.EnclosureLength,
			DatePublished:   int64(ep.DatePublished),
		})
	}
	return Feed{
		ID:          strconv.Itoa(int(p.ID)),
		FeedType:    FeedTypePodcast,
		Title:       p.Title,
		Url:         url,
		Description: p.Description,
		Author:      p.Author,
		Generator:   p.Generator,
		Items:       items,
		ImageUrl:    p.Image,
		Link:        p.Link,
		DateUpdated: int64(p.LastUpdateTime),
		ContentType: p.ContentType,
		Language:    p.Language,
		Value:       p.Value,
	}, nil
}

func PodcastIndexHeaders() map[string]string {
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

func ParsePodcastFeed(url string) (*Feed, error) {
	pod, err := PodcastFeed(url)
	fmt.Println("GOT A POD!", pod)
	if err != nil || pod == nil {
		return nil, err
	}
	eps, err := PodcastEpisodes(url)
	if err != nil {
		return nil, err
	}
	pod.Episodes = eps
	feed, err := PodcastToGeneric(url, pod)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func PodcastFeed(url string) (*Podcast, error) {
	client := &http.Client{}

	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	requrl := PodcastIndexBaseURL + "podcasts/byfeedurl?url=" + url
	req, err := http.NewRequest("GET", requrl, nil)

	headers := PodcastIndexHeaders()
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
		fmt.Println("=> json unmarshall error", err)
		return nil, err
	}

	return &r.Feed, nil
}

func PodcastEpisodes(url string) ([]Episode, error) {
	client := &http.Client{}

	requrl := PodcastIndexBaseURL + "episodes/byfeedurl?url=" + url
	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	req, err := http.NewRequest("GET", requrl, nil)

	headers := PodcastIndexHeaders()
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
	Generator      string    `json:"generator"`
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
