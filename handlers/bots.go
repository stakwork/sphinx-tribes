package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
)

type botHandler struct {
	db              db.Database
	verifyTribeUUID func(uuid string, checkTimestamp bool) (string, error)
}

func NewBotHandler(db db.Database) *botHandler {
	return &botHandler{
		db:              db,
		verifyTribeUUID: auth.VerifyTribeUUID,
	}
}

// CreateOrEditBot godoc
//
//	@Summary		Create or edit a bot
//	@Description	Create or edit a bot
//	@Tags			Bots
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			bot	body		db.Bot	true	"Bot object"
//	@Success		200	{object}	db.Bot
//	@Router			/bots [put]
func (bt *botHandler) CreateOrEditBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	bot := db.Bot{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &bot)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if bot.UUID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now()

	extractedPubkey, err := bt.verifyTribeUUID(bot.UUID, false)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pubKeyFromAuth == "" {
		bot.Created = &now
	} else { // IF PUBKEY IN CONTEXT, MUST AUTH!
		if pubKeyFromAuth != extractedPubkey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	bot.OwnerPubKey = extractedPubkey
	bot.Updated = &now
	bot.UniqueName, _ = bt.BotUniqueNameFromName(bot.Name)

	_, err = bt.db.CreateOrEditBot(bot)
	if err != nil {
		logger.Log.Error("=> ERR createOrEditBot: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bot)
}

// GetListedBots godoc
//
//	@Summary		Get listed bots
//	@Description	Get a list of listed bots
//	@Tags			Bots
//	@Success		200	{array}	db.Bot
//	@Router			/bots [get]
func (bt *botHandler) GetListedBots(w http.ResponseWriter, r *http.Request) {
	bots := bt.db.GetListedBots(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bots)
}

// GetBot godoc
//
//	@Summary		Get a bot
//	@Description	Get a bot by UUID
//	@Tags			Bots
//	@Param			uuid	path		string	true	"Bot UUID"
//	@Success		200		{object}	db.Bot
//	@Router			/bots/{uuid} [get]
func (bt *botHandler) GetBot(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	bot := bt.db.GetBot(uuid)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bot)
}

// GetBotByUniqueName godoc
//
//	@Summary		Get a bot by unique name
//	@Description	Get a bot by unique name
//	@Tags			Bots
//	@Param			unique_name	path		string	true	"Unique name"
//	@Success		200			{object}	db.Bot
//	@Router			/bot/{name} [get]
func (bt *botHandler) GetBotByUniqueName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	bot := bt.db.GetBotByUniqueName(name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bot)
}

// GetBotsByOwner godoc
//
//	@Summary		Get bots by owner
//	@Description	Get a list of bots by owner public key
//	@Tags			Bots
//	@Param			pubkey	path	string	true	"Owner public key"
//	@Success		200		{array}	db.Bot
//	@Router			/bots/owner/{pubkey} [get]
func (bt *botHandler) GetBotsByOwner(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "pubkey")
	bots := bt.db.GetBotsByOwner(name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bots)
}

// SearchBots godoc
//
//	@Summary		Search bots
//	@Description	Search for bots
//	@Tags			Bots
//	@Param			query	query	string	true	"Search query"
//	@Success		200		{array}	db.Bot
//	@Router			/search/bots/{query} [get]
func (bt *botHandler) SearchBots(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitString)
	offset, _ := strconv.Atoi(offsetString)
	if limit == 0 {
		limit = 10
	}
	bots := bt.db.SearchBots(query, limit, offset)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bots)
}

// DeleteBot godoc
//
//	@Summary		Delete a bot
//	@Description	Delete a bot by UUID
//	@Tags			Bots
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string	true	"Bot UUID"
//	@Success		200		{object}	bool
//	@Router			/bot/{uuid} [delete]
func (bt *botHandler) DeleteBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")

	logger.Log.Info("uuid: %s", uuid)

	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := bt.verifyTribeUUID(uuid, false)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bt.db.UpdateBot(uuid, map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

// BotUniqueNameFromName godoc
//
//	@Summary		Get unique name from bot name
//	@Description	Get unique name from bot name
//	@Tags			Bots
//	@Param			name	query		string	true	"Bot name"
//	@Success		200		{object}	string
//	@Router			/bot/unique_name [get]
func (h *botHandler) BotUniqueNameFromName(name string) (string, error) {
	pathOne := strings.ToLower(strings.Join(strings.Fields(name), ""))
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	path := reg.ReplaceAllString(pathOne, "")
	n := 0
	for {
		uniquepath := path
		if n > 0 {
			uniquepath = path + strconv.Itoa(n)
		}
		existing := h.db.GetBotByUniqueName(uniquepath)
		if existing.UUID != "" {
			n = n + 1
		} else {
			path = uniquepath
			break
		}
	}
	return path, nil
}

// TribeUniqueNameFromName godoc
//
//	@Summary		Get unique name from tribe name
//	@Description	Get unique name from tribe name
//	@Tags			Tribes
//	@Param			name	query		string	true	"Tribe name"
//	@Success		200		{object}	string
//	@Router			/tribes/unique_name [get]
func TribeUniqueNameFromName(name string) (string, error) {
	pathOne := strings.ToLower(strings.Join(strings.Fields(name), ""))
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	path := reg.ReplaceAllString(pathOne, "")
	n := 0
	for {
		uniquepath := path
		if n > 0 {
			uniquepath = path + strconv.Itoa(n)
		}
		existing := db.DB.GetTribeByUniqueName(uniquepath)
		if existing.UUID != "" {
			n = n + 1
		} else {
			path = uniquepath
			break
		}
	}
	return path, nil
}
