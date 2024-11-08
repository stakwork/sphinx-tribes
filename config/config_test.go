package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	InitConfig()

	if Host == "" {
		t.Error("Could not load default host")
	}

	if MemeUrl == "" {
		t.Error("Could not load default meme url")
	}

	if JwtKey == "" {
		t.Error("Could not load random jwtKey")
	}
}

func TestGenerateRandomString(t *testing.T) {

	testRandString := GenerateRandomString()

	if testRandString == "" {
		t.Error("randstring should not be empty")
	}

	if len(testRandString) < 24 {
		t.Error("randstring cannot be less than length 24")
	}
	if len(testRandString) > 24 {
		t.Error("randstring cannot be greater than lenght 24")
	}
}

func TestGetNodePubKey(t *testing.T) {
	defer gock.Off()

	response := NodeGetInfoResponse{IdentityPubkey: "1234"}
	nodeGetInfo := NodeGetInfo{Success: true, Response: response}

	gock.New("https://relay.com").
		Get("/getinfo").
		Persist().
		Reply(200).
		JSON(nodeGetInfo)

	os.Setenv("RELAY_URL", "https://relay.com")
	InitConfig()
	nodePubKey := GetNodePubKey()
	fmt.Print(nodePubKey)
	if nodePubKey != "1234" {
		t.Error("Node pubkey is incorrect")
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
