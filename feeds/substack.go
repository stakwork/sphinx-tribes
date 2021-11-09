package feeds

import (
	"encoding/xml"
	"time"
)

const SubstackTimeFormat = time.RFC1123

type SubstackPost struct {
	Title   string `xml:"title"`
	Desc    string `xml:"description"`
	Link    string `xml:"link"`
	Guid    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
}

type SubstackImage struct {
	Url string `xml:"url"`
}

type SubstackChannel struct {
	Title         string         `xml:"title"`
	Link          string         `xml:"link"`
	Desc          string         `xml:"description"`
	Image         SubstackImage  `xml:"image"`
	Generator     string         `xml:"generator"`
	LastBuildDate string         `xml:"lastBuildDate"`
	Copyright     string         `xml:"copyright"`
	Language      string         `xml:"language"`
	Items         []SubstackPost `xml:"item"`
}

type SubstackFeed struct {
	Channel SubstackChannel `xml:"channel"`
}

func ParseSubstackFeed(url string) (*Feed, error) {
	bod, err := httpget(url)
	if err != nil {
		return nil, err
	}
	var f SubstackFeed
	if err := xml.Unmarshal(bod, &f); err != nil {
		return nil, err
	}
	genericFeed, err := SubstackFeedToGeneric(url, f)
	if err != nil {
		return nil, err
	}
	return &genericFeed, nil
}

func SubstackFeedToGeneric(url string, mf SubstackFeed) (Feed, error) {
	c := mf.Channel
	items := []Item{}
	for _, post := range c.Items {
		t, _ := time.Parse(SubstackTimeFormat, post.PubDate)
		items = append(items, Item{
			Id:            post.Guid,
			Title:         post.Title,
			EnclosureURL:  post.Link,
			Description:   post.Desc,
			DatePublished: t.Unix(),
		})
	}
	return Feed{
		ID:          url,
		FeedType:    FeedTypeBlog,
		Title:       c.Title,
		Url:         url,
		Description: c.Desc,
		Items:       items,
		ImageUrl:    c.Image.Url,
	}, nil
}
