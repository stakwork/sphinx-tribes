package main

import "testing"

func TestInitConfig(t *testing.T) {
	initConfig()

	if host != "https://people.sphinx.chat" {
		t.Error("Could not load default host")
	}

	if jwtKey == "" {
		t.Error("Could not load random jwtKey")
	}
}
