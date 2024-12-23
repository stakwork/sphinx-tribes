package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
)

type channelHandler struct {
	db db.Database
}

func NewChannelHandler(db db.Database) *channelHandler {
	return &channelHandler{
		db: db,
	}
}

func (ch *channelHandler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if id == 0 {
		logger.Log.Info("id is 0")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := ch.db.GetChannel(uint(id))
	existingTribe := ch.db.GetTribe(existing.TribeUUID)
	if existing.ID == 0 {
		logger.Log.Info("existing id is 0")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if existingTribe.OwnerPubKey != pubKeyFromAuth {
		 logger.Log.Info("keys dont match")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ch.db.UpdateChannel(uint(id), map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func (ch *channelHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	channel := db.Channel{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &channel)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	//check that the tribe has the same pubKeyFromAuth
	tribe := ch.db.GetTribe(channel.TribeUUID)
	if tribe.OwnerPubKey != pubKeyFromAuth {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tribeChannels := ch.db.GetChannelsByTribe(channel.TribeUUID)
	for _, tribeChannel := range tribeChannels {
		if tribeChannel.Name == channel.Name {
			logger.Log.Info("Channel name already in use")
			w.WriteHeader(http.StatusNotAcceptable)
			return

		}
	}

	channel, err = ch.db.CreateChannel(channel)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(channel)
}
