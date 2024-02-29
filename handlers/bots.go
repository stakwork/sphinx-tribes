package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
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

func CreateOrEditBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	bot := db.Bot{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &bot)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if bot.UUID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now()

	extractedPubkey, err := auth.VerifyTribeUUID(bot.UUID, false)
	if err != nil {
		fmt.Println(err)
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
	bot.UniqueName, _ = BotUniqueNameFromName(bot.Name)

	_, err = db.DB.CreateOrEditBot(bot)
	if err != nil {
		fmt.Println("=> ERR createOrEditBot", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bot)
}

func GetListedBots(w http.ResponseWriter, r *http.Request) {
	bots := db.DB.GetListedBots(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bots)
}

func GetBot(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	bot := db.DB.GetBot(uuid)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bot)
}

func (bt *botHandler) GetBotByUniqueName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	bot := bt.db.GetBotByUniqueName(name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bot)
}

func (bt *botHandler) GetBotsByOwner(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "pubkey")
	bots := bt.db.GetBotsByOwner(name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bots)
}

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

func (bt *botHandler) DeleteBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")

	fmt.Println("uuid: ", uuid)

	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := bt.verifyTribeUUID(uuid, false)
	if err != nil {
		fmt.Println(err)
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

func BotUniqueNameFromName(name string) (string, error) {
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
		existing := db.DB.GetBotByUniqueName(uniquepath)
		if existing.UUID != "" {
			n = n + 1
		} else {
			path = uniquepath
			break
		}
	}
	return path, nil
}

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
