package feeds

import (
	"context"
	"fmt"
	"os"

	"github.com/araddon/dateparse"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func YoutubeSearch(term string) ([]Feed, error) {
	apiKey := os.Getenv("YOUTUBE_KEY")
	ctx := context.Background()
	tube, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	call := tube.Search.List([]string{"snippet"})
	call.Q(term)
	call.MaxResults(50)
	call.Type("channel", "playlist")

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	fs := []Feed{}
	for _, r := range response.Items {
		tp, _ := dateparse.ParseAny(r.Snippet.PublishedAt)
		thumb := ""
		if r.Snippet.Thumbnails != nil {
			thumb = r.Snippet.Thumbnails.Default.Url
		}
		url := "https://www.youtube.com/feeds/videos.xml?"
		if r.Id.PlaylistId != "" && r.Id.ChannelId == "" {
			url = url + "playlist_id=" + r.Id.PlaylistId
		} else {
			url = url + "channel_id=" + r.Snippet.ChannelId
		}
		f := Feed{
			ID:            r.Snippet.ChannelId,
			FeedType:      FeedTypeVideo,
			Url:           url,
			Title:         r.Snippet.Title,
			Description:   r.Snippet.Description,
			DatePublished: tp.Unix(),
			ImageUrl:      thumb,
		}
		fs = append(fs, f)
	}

	return fs, err
}

func YoutubeVideosForChannel(channelId string) ([]Item, error) {
	apiKey := os.Getenv("YOUTUBE_KEY")
	ctx := context.Background()
	tube, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	call := tube.Search.List([]string{"snippet"})
	call.ChannelId(channelId)
	call.Type("video")
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	fs := []Item{}
	for _, r := range response.Items {
		tp, _ := dateparse.ParseAny(r.Snippet.PublishedAt)
		thumb := ""
		if r.Snippet.Thumbnails != nil {
			thumb = r.Snippet.Thumbnails.Default.Url
		}
		if r.Id == nil {
			continue
		}
		id := r.Id.VideoId
		link := "https://www.youtube.com/watch?v=" + r.Id.VideoId
		f := Item{
			Id:            id,
			Title:         r.Snippet.Title,
			Description:   r.Snippet.Description,
			DatePublished: tp.Unix(),
			ImageUrl:      thumb,
			EnclosureURL:  link,
			Link:          link,
		}
		fs = append(fs, f)
	}

	fmt.Printf("%+v\n", fs)

	return fs, err
}

func YoutubeVideoSearch(term string) ([]Item, error) {
	apiKey := os.Getenv("YOUTUBE_KEY")
	ctx := context.Background()
	tube, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	call := tube.Search.List([]string{"snippet"})
	call.Q(term)
	call.MaxResults(50)
	call.Type("video")

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	fs := []Item{}
	for _, r := range response.Items {
		tp, _ := dateparse.ParseAny(r.Snippet.PublishedAt)
		thumb := ""
		if r.Snippet.Thumbnails != nil {
			thumb = r.Snippet.Thumbnails.Default.Url
		}
		url := "https://www.youtube.com/feeds/videos.xml?"
		if r.Id.PlaylistId != "" && r.Id.ChannelId == "" {
			url = url + "playlist_id=" + r.Id.PlaylistId
		} else {
			url = url + "channel_id=" + r.Snippet.ChannelId
		}
		f := Item{
			Id:            r.Id.VideoId,
			FeedType:      FeedTypeVideo,
			Url:           url,
			Title:         r.Snippet.Title,
			Description:   r.Snippet.Description,
			DatePublished: tp.Unix(),
			ImageUrl:      thumb,
			FeedId:        r.Snippet.ChannelId,
		}
		fs = append(fs, f)
	}

	return fs, err
}
