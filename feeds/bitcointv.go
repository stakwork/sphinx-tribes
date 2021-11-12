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

type BitcoinTVMediaGroupThumbnail struct {
	XMLName xml.Name
	Url     string `xml:"url,attr"`
}
type BitcoinTVMediaGroup struct {
	XMLName   xml.Name
	Thumbnail BitcoinTVMediaGroupThumbnail `xml:"thumbnail"`
}

type BitcoinTVVideo struct {
	Title      string              `xml:"title"`
	Desc       string              `xml:"description"`
	Link       string              `xml:"link"`
	Guid       string              `xml:"guid"`
	PubDate    string              `xml:"pubDate"`
	Enclosure  BitcoinTVEnclosure  `xml:"enclosure"`
	MediaGroup BitcoinTVMediaGroup `xml:"group"`
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
		l, _ := strconv.Atoi(post.Enclosure.Length)
		items = append(items, Item{
			Id:              post.Guid,
			Title:           post.Title,
			Link:            post.Link,
			EnclosureURL:    post.Enclosure.Url,
			EnclosureType:   post.Enclosure.Type,
			EnclosureLength: int32(l),
			Description:     post.Desc,
			DatePublished:   t.Unix(),
			ImageUrl:        post.MediaGroup.Thumbnail.Url,
			ThumbnailUrl:    post.MediaGroup.Thumbnail.Url,
		})
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
		DateUpdated: int32(tu.Unix()),
	}, nil
}
