package feeds

import (
	"encoding/xml"

	"github.com/araddon/dateparse"
)

type MediaGroupContent struct {
	XMLName xml.Name
	Url     string `xml:"url,attr"`
	Type    string `xml:"type,attr"`
}
type MediaGroupThumbnail struct {
	XMLName xml.Name
	Url     string `xml:"url,attr"`
}
type MediaGroup struct {
	XMLName     xml.Name
	Content     MediaGroupContent   `xml:"content"`
	Thumbnail   MediaGroupThumbnail `xml:"thumbnail"`
	Description string              `xml:"description"`
}

type YoutubeLink struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
}

type YoutubeEntry struct {
	ID         string        `xml:"id"`
	Title      string        `xml:"title"`
	Link       YoutubeLink   `xml:"link"`
	Published  string        `xml:"published"`
	Updated    string        `xml:"updated"`
	Author     YoutubeAuthor `xml:"author"`
	MediaGroup MediaGroup    `xml:"group"`
}

type YoutubeAuthor struct {
	Name string `xml:"name"`
	Uri  string `xml:"uri"`
}

type YoutubeFeed struct {
	ID        string         `xml:"id"`
	Title     string         `xml:"title"`
	Link      YoutubeLink    `xml:"link"`
	Published string         `xml:"published"`
	Author    YoutubeAuthor  `xml:"author"`
	Items     []YoutubeEntry `xml:"entry"`
}

func ParseYoutubeFeed(url string) (*Feed, error) {
	bod, err := httpget(url)
	if err != nil {
		return nil, err
	}
	var f YoutubeFeed
	if err := xml.Unmarshal(bod, &f); err != nil {
		return nil, err
	}
	genericFeed, err := YoutubeFeedToGeneric(url, f)
	if err != nil {
		return nil, err
	}
	return &genericFeed, nil
}

func YoutubeFeedToGeneric(url string, f YoutubeFeed) (Feed, error) {
	items := []Item{}
	for _, item := range f.Items {
		tp, _ := dateparse.ParseAny(item.Published)
		tu, _ := dateparse.ParseAny(item.Updated)
		id := item.ID
		if id == "" {
			id = item.Link.Href
		}
		items = append(items, Item{
			Id:            id,
			Title:         item.Title,
			Link:          item.Link.Href,
			EnclosureURL:  item.MediaGroup.Content.Url,
			EnclosureType: item.MediaGroup.Content.Type,
			Description:   item.MediaGroup.Description,
			Author:        item.Author.Name,
			ImageUrl:      item.MediaGroup.Thumbnail.Url,
			ThumbnailUrl:  item.MediaGroup.Thumbnail.Url,
			DatePublished: tp.Unix(),
			DateUpdated:   tu.Unix(),
		})
	}
	id := f.ID
	if id == "" {
		id = url
	}
	pd, _ := dateparse.ParseAny(f.Published)
	return Feed{
		ID:            id,
		FeedType:      FeedTypeVideo,
		Title:         f.Title,
		Url:           url,
		Link:          f.Link.Href,
		Items:         items,
		Author:        f.Author.Name,
		DatePublished: pd.Unix(),
	}, nil
}
