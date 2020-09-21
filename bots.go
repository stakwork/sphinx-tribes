package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

func createOrEditBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	bot := Bot{}
	body, err := ioutil.ReadAll(r.Body)
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

	extractedPubkey, err := VerifyTribeUUID(bot.UUID)
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

	// existing := DB.getTribe(tribe.UUID)
	// if existing.UUID != "" {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	bot.OwnerPubKey = extractedPubkey
	bot.Updated = &now
	bot.UniqueName, _ = botUniqueNameFromName(bot.Name)

	DB.createOrEditBot(bot)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bot)
}

func getListedBots(w http.ResponseWriter, r *http.Request) {
	bots := DB.getListedBots()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bots)
}

func getBot(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	bot := DB.getBot(uuid)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bot)
}

func getBotByUniqueName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	bot := DB.getBotByUniqueName(name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bot)
}

func searchBots(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitString)
	offset, _ := strconv.Atoi(offsetString)
	if limit == 0 {
		limit = 10
	}
	bots := DB.searchBots(query, limit, offset)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bots)
}

func botUniqueNameFromName(name string) (string, error) {
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
		existing := DB.getBotByUniqueName(uniquepath)
		if existing.UUID != "" {
			n = n + 1
		} else {
			path = uniquepath
			break
		}
	}
	return path, nil
}

func deleteBot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")

	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := VerifyTribeUUID(uuid)
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

	DB.updateBot(uuid, map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}
