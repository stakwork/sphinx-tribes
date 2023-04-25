package config

import "testing"

func TestInitConfig(t *testing.T) {
	InitConfig()

	if Host != "https://people.sphinx.chat" {
		t.Error("Could not load default host")
	}

	if JwtKey == "" {
		t.Error("Could not load random jwtKey")
	}
}
