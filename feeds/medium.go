package feeds

import (
	"encoding/xml"
)

type MediumPost struct {
	Title string `xml:"title"`
	Desc  string `xml:"description"`
	Link  string `xml:"link"`
	Guid  string `xml:"guid"`
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
		items = append(items, Item{
			Id:           post.Guid,
			Title:        post.Title,
			EnclosureURL: post.Link,
			Description:  post.Desc,
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
