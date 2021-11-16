package feeds

import (
	"encoding/xml"

	"github.com/araddon/dateparse"
)

// type MediumPostCreator struct {
// 	XMLName xml.Name
// 	Url     string `xml:"url"`
// }

type MediumPost struct {
	Title   string `xml:"title"`
	Desc    string `xml:"description"`
	Link    string `xml:"link"`
	Guid    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
	Updated string `xml:"updated"`
	Creator string `xml:"creator"`
}

type MediumImage struct {
	Url string `xml:"url"`
}
type MediumChannel struct {
	Title         string       `xml:"title"`
	Link          string       `xml:"link"`
	Desc          string       `xml:"description"`
	Image         MediumImage  `xml:"image"`
	Generator     string       `xml:"generator"`
	LastBuildDate string       `xml:"lastBuildDate"`
	Creator       string       `xml:"creator"`
	Items         []MediumPost `xml:"item"`
}

type MediumFeed struct {
	Channel MediumChannel `xml:"channel"`
}

func ParseMediumFeed(url string) (*Feed, error) {
	bod, err := httpget(url)
	if err != nil {
		return nil, err
	}
	var f MediumFeed
	if err := xml.Unmarshal(bod, &f); err != nil {
		return nil, err
	}
	genericFeed, err := MediumFeedToGeneric(url, f)
	if err != nil {
		return nil, err
	}
	return &genericFeed, nil
}

func MediumFeedToGeneric(url string, mf MediumFeed) (Feed, error) {
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
		Generator:   c.Generator,
		Author:      c.Creator,
		DateUpdated: lbd.Unix(),
	}, nil
}
