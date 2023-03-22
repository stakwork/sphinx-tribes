package main

import "os"

var host string

func initConfig() {
	host = os.Getenv("LN_SERVER_BASE_URL")

	if host == "" {
		host = "people.sphinx.chat"
	}
}
