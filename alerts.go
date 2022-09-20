package main

import "os"

type Action struct {
	Action    string `json:"action"`    // "dm"
	Chat_uuid string `json:"chat_uuid"` // tribe uuid
	Pubkey    string `json:"pubkey"`
	Content   string `json:"content"`
}

func processAlerts(p Person) {
	// get the existing person's github_issues
	// compare to see if there are new ones
	// check the new ones for coding languages like "#Javascript"

	// people need a new "extras"->"alerts" that lets them toggle on and off alerts
	// pull people who have alerts on
	// of those people, who have "Javascript" in their "coding_languages"?
	// if they match, build an Action with their pubkey
	// post all the Actions you have build to relay with the HMAC header

	relay_url := os.Getenv("ALERT_URL")
	alert_secret := os.Getenv("ALERT_SECERET")
	alert_tribe_uuid := os.Getenv("ALERT_TRIBE_UUID")
	if relay_url == "" || alert_secret == "" || alert_tribe_uuid == "" {
		return
	}
}
