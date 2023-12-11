package config

import (
	"fmt"
	"github.com/h2non/gock"
	"os"
	"testing"
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

	if RelayAuthKey != "" {
		t.Error("Could not load RelayAuthKey")
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

	//response := map[string]string{"identity_pubkey": "1234"}
	//success := map[string]bool{"success": true}
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
