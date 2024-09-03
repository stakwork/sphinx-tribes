package feeds

import (
    "testing"
)

func TestParseFeed(t *testing.T) {
    testCases := []struct {
        name     string
        url      string
        fulltext bool
        wantType int
    }{
        // {"Medium Blog", "https://medium.com/feed/@yourmediumusername", false, FeedTypeBlog},
        {"Substack Newsletter", "https://patrickholland.substack.com/feed", false, FeedTypeBlog},
        // {"YouTube Channel", "https://www.youtube.com/feeds/videos.xml?channel_id=YOUR_CHANNEL_ID", false, FeedTypeVideo},
        // {"Bitcoin TV Channel", "https://bitcointv.com/feeds/videos.xml?videoChannelId=YOUR_CHANNEL_ID", false, FeedTypeVideo},
        // {"Podcast Feed", "https://feeds.simplecast.com/YOUR_PODCAST_ID", true, FeedTypePodcast},
        // Add more test cases as needed
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            feed, err := ParseFeed(tc.url, tc.fulltext)
            if err != nil {
                t.Fatalf("ParseFeed(%s) returned error: %v", tc.url, err)
            }
            if feed.FeedType != tc.wantType {
                t.Errorf("ParseFeed(%s) got feed type %d, want %d", tc.url, feed.FeedType, tc.wantType)
            }
            // Add more assertions here to check other feed properties
        })
    }
}