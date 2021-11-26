package feeds

import (
	"encoding/xml"
	"strconv"

	"github.com/araddon/dateparse"
)

type BitcoinTVEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	Type    string   `xml:"type,attr"`
	Url     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
}

type BitcoinTVMediaThumbnail struct {
	XMLName xml.Name
	Url     string `xml:"url,attr"`
}
type BitcoinTVMediaGroupContent struct {
	Url      string `xml:"url,attr"`
	Duration string `xml:"duration,attr"`
	Type     string `xml:"type,attr"`
}
type BitcoinTVMediaGroup struct {
	XMLName xml.Name
	Content []BitcoinTVMediaGroupContent `xml:"content"`
}

type BitcoinTVVideo struct {
	Title      string                  `xml:"title"`
	Desc       string                  `xml:"description"`
	Link       string                  `xml:"link"`
	Guid       string                  `xml:"guid"`
	PubDate    string                  `xml:"pubDate"`
	Enclosure  BitcoinTVEnclosure      `xml:"enclosure"`
	MediaGroup BitcoinTVMediaGroup     `xml:"group"`
	Thumbnail  BitcoinTVMediaThumbnail `xml:"thumbnail"`
}

type BitcoinTVImage struct {
	Url string `xml:"url"`
}

type BitcoinTVChannel struct {
	Title         string           `xml:"title"`
	Link          string           `xml:"link"`
	Desc          string           `xml:"description"`
	Image         BitcoinTVImage   `xml:"image"`
	Generator     string           `xml:"generator"`
	LastBuildDate string           `xml:"lastBuildDate"`
	Copyright     string           `xml:"copyright"`
	Items         []BitcoinTVVideo `xml:"item"`
}

type BitcoinTVFeed struct {
	Channel BitcoinTVChannel `xml:"channel"`
}

func ParseBitcoinTVFeed(url string) (*Feed, error) {
	bod, err := httpget(url)
	if err != nil {
		return nil, err
	}
	var f BitcoinTVFeed
	if err := xml.Unmarshal(bod, &f); err != nil {
		return nil, err
	}
	genericFeed, err := BitcoinTVToGeneric(url, f)
	if err != nil {
		return nil, err
	}
	return &genericFeed, nil
}

func BitcoinTVToGeneric(url string, mf BitcoinTVFeed) (Feed, error) {
	c := mf.Channel
	items := []Item{}
	for _, post := range c.Items {
		t, _ := dateparse.ParseAny(post.PubDate)
		if len(post.MediaGroup.Content) > 0 {
			content := post.MediaGroup.Content[0]
			dur, _ := strconv.Atoi(content.Duration)
			items = append(items, Item{
				Id:            post.Guid,
				Title:         post.Title,
				Link:          post.Link,
				EnclosureURL:  content.Url,
				EnclosureType: content.Type,
				Duration:      int32(dur),
				Description:   post.Desc,
				DatePublished: t.Unix(),
				ImageUrl:      post.Thumbnail.Url,
				ThumbnailUrl:  post.Thumbnail.Url,
			})
		}
	}
	tu, _ := dateparse.ParseAny(c.LastBuildDate)
	return Feed{
		ID:          url,
		FeedType:    FeedTypeVideo,
		Title:       c.Title,
		Url:         url,
		Description: c.Desc,
		Items:       items,
		ImageUrl:    c.Image.Url,
		Generator:   c.Generator,
		DateUpdated: tu.Unix(),
	}, nil
}
