package feeds

import (
	"encoding/xml"

	"github.com/araddon/dateparse"
)

type SubstackPost struct {
	Title   string `xml:"title"`
	Desc    string `xml:"description"`
	Link    string `xml:"link"`
	Guid    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
	Updated string `xml:"updated"`
	Creator string `xml:"creator"`
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
	Creator       string         `xml:"creator"`
	Items         []SubstackPost `xml:"item"`
}

type SubstackFeed struct {
	Channel SubstackChannel `xml:"channel"`
}

func ParseSubstackFeed(url string, bod []byte) (*Feed, error) {
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
		t, _ := dateparse.ParseAny(post.PubDate)
		tu, _ := dateparse.ParseAny(post.Updated)
		items = append(items, Item{
			Id:            post.Guid,
			Title:         post.Title,
			Link:          post.Link,
			EnclosureURL:  post.Link,
			Description:   post.Desc,
			Author:        post.Creator,
			DatePublished: t.Unix(),
			DateUpdated:   tu.Unix(),
		})
	}
	lbd, _ := dateparse.ParseAny(c.LastBuildDate)
	return Feed{
		ID:          url,
		FeedType:    FeedTypeBlog,
		Title:       c.Title,
		Url:         url,
		Link:        c.Link,
		Description: c.Desc,
		Items:       items,
		ImageUrl:    c.Image.Url,
		DateUpdated: lbd.Unix(),
		Generator:   c.Generator,
		Author:      c.Creator,
	}, nil
}
