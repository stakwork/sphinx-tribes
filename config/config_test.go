package config

import "testing"

func TestInitConfig(t *testing.T) {
	InitConfig()

	if Host != "https://people.sphinx.chat" {
		t.Error("Could not load default host")
	}

	if MemeUrl != "https://memes.sphinx.chat" {
		t.Error("Could not load default meme url")
	}

	if JwtKey == "" {
		t.Error("Could not load random jwtKey")
	}
}
