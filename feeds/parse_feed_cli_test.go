package feeds

import (
    "flag"
    "fmt"
    "os"
    "testing"
)

var urlFlag = flag.String("url", "", "URL of the feed to parse")

func TestParseFeedFromCLI(t *testing.T) {
    // Parse command-line flags
    flag.Parse()

    // Check if URL is provided
    if *urlFlag == "" {
        fmt.Println("Please provide a URL using the -url flag")
        os.Exit(1)
    }

    // Parse the feed
    feed, err := ParseFeed(*urlFlag, false)
    if err != nil {
        t.Fatalf("Error parsing feed: %v", err)
    }

    // Print feed information
    t.Logf("Feed URL: %s", *urlFlag)
    t.Logf("Feed Title: %s", feed.Title)
    t.Logf("Feed Type: %d", feed.FeedType)
    t.Logf("Number of items: %d", len(feed.Items))
    if len(feed.Items) > 0 {
        t.Logf("First item title: %s", feed.Items[0].Title)
    }

    // Additional feed information
    t.Logf("Feed Description: %s", feed.Description)
    t.Logf("Feed Author: %s", feed.Author)
    t.Logf("Feed Image URL: %s", feed.ImageUrl)
    t.Logf("Feed Link: %s", feed.Link)
    t.Logf("Feed Language: %s", feed.Language)
}
