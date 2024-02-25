package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
)

func DeleteChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if id == 0 {
		fmt.Println("id is 0")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := db.DB.GetChannel(uint(id))
	existingTribe := db.DB.GetTribe(existing.TribeUUID)
	if existing.ID == 0 {
		fmt.Println("existing id is 0")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if existingTribe.OwnerPubKey != pubKeyFromAuth {
		fmt.Println("keys dont match")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db.DB.UpdateChannel(uint(id), map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func CreateChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	channel := db.Channel{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &channel)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	//check that the tribe has the same pubKeyFromAuth
	tribe := db.DB.GetTribe(channel.TribeUUID)
	if tribe.OwnerPubKey != pubKeyFromAuth {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	tribeChannels := db.DB.GetChannelsByTribe(channel.TribeUUID)
	for _, tribeChannel := range tribeChannels {
		if tribeChannel.Name == channel.Name {
			fmt.Println("Channel name already in use")
			w.WriteHeader(http.StatusNotAcceptable)
			return

		}
	}

	channel, err = db.DB.CreateChannel(channel)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(channel)
}
