package main

import (
	"encoding/json"
	"github.com/stakwork/sphinx-tribes/mqtt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HandelTribeMessageBundleFromRelay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)
	tribeBatch := TribeBatch{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	err = r.Body.Close()
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(body, &tribeBatch)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if tribeBatch.ChatUUID == "" {
		log.Println("supplied 'chat_uuid' is empty")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	tribe := DB.getTribe(tribeBatch.ChatUUID)
	if tribe.UUID == "" {
		log.Println("problem finding a tribe with the supplied 'chat_uuid'")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if tribe.OwnerPubKey == pubKeyFromAuth {
		c := mqtt.CLIENT
		defer c.Disconnect()

		for _, bc := range tribeBatch.BatchContents {
			c.Publish(bc.MQTTTopic, bc.Data, false)
		}

		now := time.Now().Unix()
		DB.updateTribe(tribeBatch.ChatUUID, map[string]interface{}{
			"last_active": now,
		})
	} else {
		log.Println("keys dont match - unauthorised")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(true)
	if err != nil {
		log.Println(err)
	}
}
