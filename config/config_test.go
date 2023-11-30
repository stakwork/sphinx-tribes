package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestStripSuperAdmins(t *testing.T) {
	testAdminList := "hello, hi, yes, now"
	admins := StripSuperAdmins(testAdminList)
	assert.Equal(t, len(admins), 4)

	testAdminNocomma := "hello"
	adminsNoComma := StripSuperAdmins(testAdminNocomma)
	assert.Equal(t, len(adminsNoComma), 1)

	testNoAdmins := ""
	noAdmins := StripSuperAdmins(testNoAdmins)
	assert.Equal(t, len(noAdmins), 0)

	test2Admins := "hello, hi"
	admins2 := StripSuperAdmins(test2Admins)
	assert.Equal(t, len(admins2), 2)
}
