package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		fmt.Println("no .env file")
	}

	db.InitDB()
	db.InitCache()
	// Config has to be inited before JWT, if not it will lead to NO JWT error
	config.InitConfig()
	auth.InitJwt()

	skipLoops := os.Getenv("SKIP_LOOPS")
	if skipLoops != "true" {
		go handlers.ProcessTwitterConfirmationsLoop()
		go handlers.ProcessGithubIssuesLoop()
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
