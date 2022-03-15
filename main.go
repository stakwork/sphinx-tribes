package main

import (
	"context"
	"fmt"
	"github.com/stakwork/sphinx-tribes/mqtt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	var err error

	err = godotenv.Load()
	if err != nil {
		fmt.Println("no .env file here")
	}

	fmt.Println("Init DB...")
	initDB()
	fmt.Println("Done.")

	fmt.Println("Init Cache...")
	initCache()
	fmt.Println("Done")

	fmt.Println("Init MQTT Client...")
	mqtt.Init()
	fmt.Println("Done")

	skipLoops := os.Getenv("SKIP_LOOPS")
	if skipLoops != "true" {
		go processTwitterConfirmationsLoop()
		go processGithubIssuesLoop()
	}

	run()
}

// Start the MQTT plugin
func run() {

	router := NewRouter()

	shutdownSignal := make(chan os.Signal)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownSignal

	// shutdown web server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := router.Shutdown(ctx); err != nil {
		fmt.Printf("error shutting down server: %s", err.Error())
	}
}
