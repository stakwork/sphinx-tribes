package feeds

import (
    "testing"
)

func TestMultipleRealFeeds(t *testing.T) {
    urls := []struct {
        url      string
        expected int
    }{
        {"https://patrickholland.substack.com/feed", FeedTypeBlog},
        {"https://tftc.io/feed/", FeedTypeBlog},
        {"https://medium.com/@shosaski/feed", FeedTypeBlog},
        // Add more URLs here, each with its expected FeedType
    }

    for _, tc := range urls {
        t.Run(tc.url, func(t *testing.T) {
            feed, err := ParseFeed(tc.url, false)

            if err != nil {
                t.Fatalf("ParseFeed(%s) returned error: %v", tc.url, err)
            }

            if feed == nil {
                t.Fatalf("ParseFeed(%s) returned nil feed", tc.url)
            }

            // Basic checks
            if feed.Title == "" {
                t.Errorf("Feed title is empty")
            }

            if feed.FeedType != tc.expected {
                t.Errorf("Expected FeedType %d, got %d", tc.expected, feed.FeedType)
            }

            if len(feed.Items) == 0 {
                t.Errorf("No items found in the feed")
            }

            // Print some information about the feed
            t.Logf("Feed URL: %s", tc.url)
            t.Logf("Feed Title: %s", feed.Title)
            t.Logf("Feed Type: %d", feed.FeedType)
            t.Logf("Number of items: %d", len(feed.Items))
            if len(feed.Items) > 0 {
                t.Logf("First item title: %s", feed.Items[0].Title)
            }
            t.Logf("---")
        })
    }
}