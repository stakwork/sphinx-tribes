package simulation

import (
    "encoding/json"
    "fmt"
    "github.com/stakwork/sphinx-tribes/feeds"
)

func SimulateGetGenericFeed(url string) {
    feed, err := feeds.ParseFeed(url, false)
    if err != nil {
        fmt.Printf("Error parsing feed: %v\n", err)
        return
    }

    // Simulate the JSON encoding
    jsonData, err := json.MarshalIndent(feed, "", "  ")
    if err != nil {
        fmt.Printf("Error encoding JSON: %v\n", err)
        return
    }

    fmt.Println(string(jsonData))
}