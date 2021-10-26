package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"

	"github.com/stakwork/sphinx-tribes/frontend"
)

// NewRouter creates a chi router
func NewRouter() *http.Server {
	r := initChi()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("pong")
	})

	r.Group(func(r chi.Router) {
		r.Get("/", frontend.IndexRoute)
		r.Get("/static/*", frontend.StaticRoute)
		r.Get("/manifest.json", frontend.ManifestRoute)
		r.Get("/favicon.ico", frontend.FaviconRoute)
	})
	r.Group(func(r chi.Router) {
		r.Get("/t/static/*", frontend.StaticRoute)
		r.Get("/t/manifest.json", frontend.ManifestRoute)
		r.Get("/t/favicon.ico", frontend.FaviconRoute)
		r.Get("/t/{unique_name}", frontend.IndexRoute)
		r.Get("/t", frontend.IndexRoute)
	})
	r.Group(func(r chi.Router) {
		r.Get("/p/static/*", frontend.StaticRoute)
		r.Get("/p/manifest.json", frontend.ManifestRoute)
		r.Get("/p/favicon.ico", frontend.FaviconRoute)
		r.Get("/p/{pubkey}", frontend.IndexRoute)
		r.Get("/p", frontend.IndexRoute)
		r.Get("/b", frontend.IndexRoute)
	})

	r.Group(func(r chi.Router) {
		r.Get("/tribes", getListedTribes)
		r.Get("/tribes/{uuid}", getTribe)
		r.Get("/tribes_by_owner/{pubkey}", getTribesByOwner)
		r.Post("/tribes", createOrEditTribe)
		r.Post("/bots", createOrEditBot)
		r.Get("/bots", getListedBots)
		r.Get("/bots/{uuid}", getBot)
		r.Get("/bot/{name}", getBotByUniqueName)
		r.Get("/search/bots/{query}", searchBots)
		r.Get("/podcast", getPodcast)
		r.Get("/people", getListedPeople)

		r.Get("/ask", ask)
		r.Get("/poll/{challenge}", poll)
		r.Get("/person/{pubkey}", getPersonByPubkey)
		r.Get("/github_issue/{owner}/{repo}/{issue}", getGithubIssue)
		r.Post("/save", postSave)
		r.Get("/save/{key}", pollSave)
	})

	r.Group(func(r chi.Router) {
		r.Use(PubKeyContext)
		r.Put("/tribe", createOrEditTribe)
		r.Put("/tribestats", putTribeStats)
		r.Delete("/tribe/{uuid}", deleteTribe)
		r.Put("/tribeactivity/{uuid}", putTribeActivity)
		r.Delete("/bot/{uuid}", deleteBot)
		r.Post("/person", createOrEditPerson)
		r.Post("/verify/{challenge}", verify)
		r.Put("/bot", createOrEditBot)
		r.Delete("/person/{id}", deletePerson)
	})

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5002"
	}

	server := &http.Server{Addr: ":" + PORT, Handler: r}
	go func() {
		fmt.Println("Listening on port " + PORT)
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("server err:", err.Error())
		}
	}()

	return server
}

func getPodcast(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	feedid := r.URL.Query().Get("id")
	podcast, err := getFeed(url, feedid)
	episodes, err := getEpisodes(url, feedid)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	podcast.Episodes = episodes

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(podcast)
	if err != nil {
		fmt.Println(err)
	}
}

func getAllTribes(w http.ResponseWriter, r *http.Request) {
	tribes := DB.getAllTribes()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

func getListedTribes(w http.ResponseWriter, r *http.Request) {
	tribes := DB.getListedTribes()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

func getTribesByOwner(w http.ResponseWriter, r *http.Request) {
	all := r.URL.Query().Get("all")
	tribes := []Tribe{}
	pubkey := chi.URLParam(r, "pubkey")
	if all == "true" {
		tribes = DB.getAllTribesByOwner(pubkey)
	} else {
		tribes = DB.getTribesByOwner(pubkey)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

func getListedPeople(w http.ResponseWriter, r *http.Request) {
	people := DB.getListedPeople()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}

func getTribe(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	tribe := DB.getTribe(uuid)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribe)
}

func createOrEditTribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	tribe := Tribe{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &tribe)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if tribe.UUID == "" {
		fmt.Println("createOrEditTribe no uuid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now() //.Format(time.RFC3339)

	extractedPubkey, err := VerifyTribeUUID(tribe.UUID, false)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pubKeyFromAuth == "" {
		tribe.Created = &now
	} else { // IF PUBKEY IN CONTEXT, MUST AUTH!
		if pubKeyFromAuth != extractedPubkey {
			fmt.Println("createOrEditTribe pubkeys dont match")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	existing := DB.getTribe(tribe.UUID)
	if existing.UUID == "" { // doesnt exist already, create unique name
		tribe.UniqueName, _ = tribeUniqueNameFromName(tribe.Name)
	} else { // already exists! make sure its owned
		if existing.OwnerPubKey != extractedPubkey {
			fmt.Println("createOrEditTribe tribe.ownerPubKey not match")
			fmt.Println(existing.OwnerPubKey)
			fmt.Println(extractedPubkey)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	tribe.OwnerPubKey = extractedPubkey
	tribe.Updated = &now
	tribe.LastActive = now.Unix()

	_, err = DB.createOrEditTribe(tribe)
	if err != nil {
		fmt.Println("=> ERR createOrEditTribe", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribe)
}

func putTribeActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := VerifyTribeUUID(uuid, false)
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

	now := time.Now().Unix()
	DB.updateTribe(uuid, map[string]interface{}{
		"last_active": now,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func putTribeStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	tribe := Tribe{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &tribe)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if tribe.UUID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := VerifyTribeUUID(tribe.UUID, false)
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

	now := time.Now()
	tribe.Updated = &now
	DB.updateTribe(tribe.UUID, map[string]interface{}{
		"member_count": tribe.MemberCount,
		"updated":      &now,
		"bots":         tribe.Bots,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func deleteTribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")

	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := VerifyTribeUUID(uuid, false)
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

	DB.updateTribe(uuid, map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

type extractResponse struct {
	Pubkey string `json:"pubkey"`
	Valid  bool   `json:"valid"`
}

func getFromAuth(path string) (*extractResponse, error) {

	authURL := "http://auth:9090"
	resp, err := http.Get(authURL + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var inter map[string]interface{}
	err = json.Unmarshal(body2, &inter)
	if err != nil {
		return nil, err
	}
	pubkey, _ := inter["pubkey"].(string)
	valid, _ := inter["valid"].(bool)
	return &extractResponse{
		Pubkey: pubkey,
		Valid:  valid,
	}, nil
}

func initChi() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-User", "authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		//Debug:            true,
	})
	r.Use(cors.Handler)
	r.Use(middleware.Timeout(60 * time.Second))
	return r
}
