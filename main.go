package main

import (
    "fmt"
    "os"
    "github.com/joho/godotenv"
    "github.com/stakwork/sphinx-tribes/simulation"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("no .env file")
    }

    // Keep only the necessary initializations
    // Remove or comment out the parts that are causing import cycles

    url := os.Getenv("FEED_URL")
    if url == "" {
        url = "https://fixthefood.substack.com/feed" // Default URL if not provided
    }

    fmt.Printf("Simulating GET /feeds?url=%s\n", url)
    simulation.SimulateGetGenericFeed(url)
}