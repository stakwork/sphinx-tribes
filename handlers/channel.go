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

// DeleteChannel godoc
//
//	@Summary		Delete a channel
//	@Description	Delete a channel by marking it as deleted. Only the tribe owner can perform this action.
//	@Tags			Channel
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		int		true	"Channel ID"
//	@Success		200	{object}	bool	"Channel deleted successfully"
//	@Failure		400	{object}	nil		"Bad request: Invalid channel ID"
//	@Failure		401	{object}	nil		"Unauthorized: User is not the tribe owner or invalid credentials"
//	@Failure		404	{object}	nil		"Not found: Channel does not exist"
//	@Router			/channel/{id} [delete]
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

// CreateChannel godoc
//
//	@Summary		Create a channel
//	@Description	Create a new channel within a tribe. Only the tribe owner can perform this action.
//	@Tags			Channel
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			request	body		db.Channel	true	"Request body containing tribe UUID and channel name"
//	@Success		200		{object}	db.Channel				"Channel created successfully"
//	@Failure		400		{object}	nil			"Bad request: Invalid request body"
//	@Failure		401		{object}	nil			"Unauthorized: User is not the tribe owner"
//	@Failure		406		{object}	nil			"Not acceptable: Channel name already in use or invalid data"
//	@Router			/channel [post]
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
