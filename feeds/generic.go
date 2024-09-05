package feeds

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const (
	FeedTypePodcast = 0
	FeedTypeVideo   = 1
	FeedTypeBlog    = 2
)

func ParseFeed(url string, fulltext bool) (*Feed, error) {

	gen, bod, err := FindGenerator(url)
	if err != nil {
		return nil, err
	}

	if strings.Contains(url, "https://medium.com/") || gen == GeneratorWordpress {
		f, err := ParseMediumFeed(url, bod)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	if strings.Contains(url, ".substack.com/feed") {
		f, err := ParseSubstackFeed(url, bod)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	if strings.Contains(url, "youtube.com/feeds/videos.xml") {
		f, err := ParseYoutubeFeed(url, bod)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	if strings.Contains(url, "bitcointv.com/feeds/videos.xml") {
		f, err := ParseBitcoinTVFeed(url, bod)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	if gen == FeedTypePodcast {
		f, err := ParseSubstackFeed(url, bod)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	f, err := ParsePodcastFeed(url, fulltext)

	if err != nil {
		f, err = ParseSubstackFeed(url, bod) // this one is quite generic
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

func AddedValue(value *Value, tribeOwnerPubkey string) *Value {
	if tribeOwnerPubkey == "" {
		return value
	}
	if value != nil {
		if value.Destinations != nil {
			if len(value.Destinations) == 1 {
				first := value.Destinations[0]
				firstSplit, _ := first.Split.(json.Number).Int64()
				if firstSplit == 1 {
					// this is the auto
					if tribeOwnerPubkey != first.Address {
						value.Destinations = append(value.Destinations, Destination{
							Address: tribeOwnerPubkey,
							Split:   99,
							Type:    "node",
						})
					}
				}
			}
		}
	} else {
		value = &Value{
			Model: Model{
				Type:      "lightning",
				Suggested: "0.00000015000",
			},
			Destinations: []Destination{
				{
					Address: tribeOwnerPubkey,
					Type:    "node",
					Split:   100,
				},
			},
		}
	}
	return value
}

type Feed struct {
	ID            string `json:"id"`
	FeedType      int    `json:"feedType"` // podcast, video, blog
	Title         string `json:"title"`
	Url           string `json:"url"`
	Description   string `json:"description"`
	Author        string `json:"author"`
	Generator     string `json:"generator"`
	ImageUrl      string `json:"imageUrl"`
	OwnerUrl      string `json:"ownerUrl"`
	Link          string `json:"link"`
	DatePublished int64  `json:"datePublished"`
	DateUpdated   int64  `json:"dateUpdated"`
	ContentType   string `json:"contentType"`
	Language      string `json:"language"`
	Items         []Item `json:"items"`
	Value         *Value `json:"value"`
	ItemId        string `json:"itemId"`
}
type Item struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	DatePublished int64  `json:"datePublished"`
	DateUpdated   int64  `json:"dateUpdated"`
	Author        string `json:"author"`
	EnclosureURL  string `json:"enclosureUrl"`
	EnclosureType string `json:"enclosureType"`
	Duration      int32  `json:"duration"`
	ImageUrl      string `json:"imageUrl"`
	ThumbnailUrl  string `json:"thumbnailUrl"`
	Link          string `json:"link"`
	// for search
	FeedId   string `json:"feedId"`
	FeedType int    `json:"feedType"`
	Url      string `json:"url"`
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
	Address     string      `json:"address"`
	Split       interface{} `json:"split"`
	Type        string      `json:"type"`
	CustomKey   string      `json:"customKey"`
	CustomValue string      `json:"customValue"`
}

func httpget(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	return body, nil
}
