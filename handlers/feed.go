package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/feeds"
	"github.com/stakwork/sphinx-tribes/logger"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetGenericFeed(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	feed, err := feeds.ParseFeed(url, false)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tribeUUID := r.URL.Query().Get("uuid")
	tribe := db.Tribe{}
	if tribeUUID != "" {
		tribe = db.DB.GetTribe(tribeUUID)
	} else {
		tribe = db.DB.GetFirstTribeByFeedURL(url)
	}

	feed.Value = feeds.AddedValue(feed.Value, tribe.OwnerPubKey)

	var data [][]string
	for z := 0; z < len(feed.Items); z++ {
		i := feed.Items[z]
		item := []string{i.Id, i.EnclosureURL}
		data = append(data, item)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feed)
}

func DownloadYoutubeFeed(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("YOUTUBE_KEY")
	ctx := context.Background()
	tube, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))

	youtube_download := db.YoutubeDownload{}
	body, err := io.ReadAll(r.Body)

	r.Body.Close()
	err = json.Unmarshal(body, &youtube_download)
	if err != nil {
		logger.Log.Error("[feed] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for i := 0; i < len(youtube_download.YoutubeUrls); i++ {
		url := youtube_download.YoutubeUrls[i]
		// Split URL to get video ID
		idSplit := strings.Split(url, "v=")
		id := idSplit[1]

		// Make a Youtube API call to check if the video exists
		call := tube.Videos.List([]string{"snippet,contentDetails,statistics"}).Id(id)
		response, _ := call.Do()

		// Add the Youtube results to the data result it should be one if the video exists
		if response.PageInfo.TotalResults < 1 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Could not process Youtube download, one of the youtueb videos does not exists")
			return
		}
	}

	processYoutubeDownload(youtube_download.YoutubeUrls)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Youtube download processed successfully")
}

func processYoutubeDownload(data []string) {
	stakworkKey := fmt.Sprintf("Token token=%s", os.Getenv("STAKWORK_KEY"))
	if stakworkKey == "" {
		logger.Log.Error("[feed] Youtube Download Error: Stakwork key not found")
	} else {
		type Vars struct {
			YoutubeContent []string `json:"youtube_content"`
		}

		type Attributes struct {
			Vars Vars `json:"vars"`
		}

		type SetVar struct {
			Attributes Attributes `json:"attributes"`
		}

		type WorkflowParams struct {
			SetVar SetVar `json:"set_var"`
		}

		workflows := WorkflowParams{
			SetVar: SetVar{
				Attributes: Attributes{
					Vars: Vars{YoutubeContent: data},
				},
			},
		}

		body := map[string]interface{}{
			"name":            "Sphinx Youtube Content Storage",
			"workflow_id":     "11848",
			"workflow_params": workflows,
		}

		buf, err := json.Marshal(body)
		if err != nil {
			logger.Log.Error("[feed] Youtube error: Unable to parse message into byte buffer: %v", err)
			return
		}

		requestUrl := "https://jobs.stakwork.com/api/v1/projects"
		request, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewBuffer(buf))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", stakworkKey)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			logger.Log.Error("[feed] Youtube Download Request Error: %v", err)
		}
		defer response.Body.Close()
		res, err := io.ReadAll(response.Body)
		if err != nil {
			logger.Log.Error("[feed] Youtube Download Request Error: %v", err)
		}
		logger.Log.Info("[feed] Youtube Download Success: %s", string(res))
	}
}

func GetPodcast(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	feedid := r.URL.Query().Get("id")
	podcast, err := getFeed(url, feedid)
	episodes, err := getEpisodes(url, feedid)

	if err != nil {
		logger.Log.Error("[feed] %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	podcast.Episodes = episodes

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(podcast)
	if err != nil {
		logger.Log.Error("[feed] %v", err)
	}
}

func SearchPodcasts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	podcasts, err := searchPodcastIndex(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fs := []feeds.Feed{}
	for _, pod := range podcasts {
		feed, err1 := feeds.PodcastToGeneric(pod.URL, &pod)
		if err1 == nil {
			fs = append(fs, feed)
		}
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(fs)
}

func SearchPodcastEpisodes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	eps, err := feeds.PodcastEpisodesByPerson(q, false)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fs := []feeds.Item{}
	for _, ep := range eps {
		episode := feeds.EpisodeToGeneric(ep, true)
		fs = append(fs, episode)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(fs)
}

func SearchYoutube(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	fs, err := feeds.YoutubeSearch(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(fs)
}

func SearchYoutubeVideos(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	fs, err := feeds.YoutubeVideoSearch(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(fs)
}

func YoutubeVideosForChannel(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("channelId")
	fs, err := feeds.YoutubeVideosForChannel(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(fs)
}

func getFeed(feedURL string, feedID string) (*feeds.Podcast, error) {
	client := &http.Client{}

	url := ""
	if feedURL != "" {
		url = feeds.PodcastIndexBaseURL + "podcasts/byfeedurl?url=" + feedURL
	} else if feedID != "" {
		url = feeds.PodcastIndexBaseURL + "podcasts/byfeedid?id=" + feedID
	}
	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	req, err := http.NewRequest("GET", url, nil)

	headers := feeds.PodcastIndexHeaders()
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		logger.Log.Error("[feed] GET error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r feeds.PodcastResponse
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		logger.Log.Error("[feed] json unmarshall error: %v", err)
		return nil, err
	}

	feed := r.Feed
	tribe := db.DB.GetFirstTribeByFeedURL(r.Feed.URL)
	feed.Value = feeds.AddedValue(r.Feed.Value, tribe.OwnerPubKey)

	return &feed, nil
}

func getEpisodes(feedURL string, feedID string) ([]feeds.Episode, error) {
	client := &http.Client{}

	url := ""
	if feedURL != "" {
		url = feeds.PodcastIndexBaseURL + "episodes/byfeedurl?url=" + feedURL
	} else if feedID != "" {
		url = feeds.PodcastIndexBaseURL + "episodes/byfeedid?id=" + feedID
	}
	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	req, err := http.NewRequest("GET", url, nil)

	headers := feeds.PodcastIndexHeaders()
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		logger.Log.Error("[feed] GET error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r feeds.EpisodeResponse
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		logger.Log.Error("[feed] json unmarshall error: %v", err)
		return nil, err
	}

	return r.Items, nil
}

func searchPodcastIndex(term string) ([]feeds.Podcast, error) {
	client := &http.Client{}

	url := feeds.PodcastIndexBaseURL + "search/byterm?q=" + term

	if url == "" {
		return nil, errors.New("no url or id supplied")
	}

	req, err := http.NewRequest("GET", url, nil)

	headers := feeds.PodcastIndexHeaders()
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)

	if err != nil {
		logger.Log.Error("[feed] GET error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r feeds.PodcastSearchResponse
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		logger.Log.Error("[feed] json unmarshall error: %v", err)
		return nil, err
	}

	return r.Feeds, nil
}
