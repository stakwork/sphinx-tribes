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
		r.Get("/tickets", frontend.IndexRoute)
	})

	r.Group(func(r chi.Router) {
		r.Get("/tribes", getListedTribes)
		r.Get("/tribes/{uuid}", getTribe)
		r.Get("/tribe_by_feed", getFirstTribeByFeed)
		r.Get("/tribes/total", getTotalribes)
		r.Get("/tribe_by_un/{un}", getTribeByUniqueName)

		r.Get("/leaderboard/{tribe_uuid}", getLeaderBoard)

		r.Get("/tribes_by_owner/{pubkey}", getTribesByOwner)
		r.Post("/tribes", createOrEditTribe)
		r.Post("/bots", createOrEditBot)
		r.Get("/bots", getListedBots)
		r.Get("/bots/owner/{pubkey}", getBotsByOwner)
		r.Get("/bots/{uuid}", getBot)

		r.Get("/bot/{name}", getBotByUniqueName)
		r.Get("/search/bots/{query}", searchBots)
		r.Get("/podcast", getPodcast)
		r.Get("/feed", getGenericFeed)
		r.Get("/search_podcasts", searchPodcasts)
		r.Get("/search_youtube", searchYoutube)
		r.Get("/youtube_videos", youtubeVideosForChannel)
		r.Get("/people", getListedPeople)
		r.Get("/people/search", getPeopleBySearch)
		r.Get("/people/posts", getListedPosts)
		r.Get("/people/wanteds", getListedWanteds)
		r.Get("/people/wanteds/assigned/{pubkey}", getPersonAssignedWanteds)
		r.Get("/people/wanteds/header", getWantedsHeader)
		r.Get("/people/short", getPeopleShortList)
		r.Get("/people/offers", getListedOffers)
		r.Get("/admin_pubkeys", getAdminPubkeys)

		r.Get("/ask", ask)
		r.Get("/poll/{challenge}", poll)
		r.Get("/person/{pubkey}", getPersonByPubkey)
		r.Get("/person/uuid/{uuid}", getPersonByUuid)
		r.Get("/person/uuid/{uuid}/assets", getPersonAssetsByUuid)
		r.Get("/person/githubname/{github}", getPersonByGithubName)

		r.Get("/github_issue/{owner}/{repo}/{issue}", getGithubIssue)
		r.Get("/github_issue/status/open", getOpenGithubIssues)
		r.Post("/save", postSave)
		r.Get("/save/{key}", pollSave)
	})

	r.Group(func(r chi.Router) {
		r.Use(PubKeyContext)
		r.Post("/channel", createChannel)
		r.Post("/leaderboard/{tribe_uuid}", createLeaderBoard)
		r.Put("/leaderboard/{tribe_uuid}", updateLeaderBoard)
		r.Put("/tribe", createOrEditTribe)
		r.Put("/tribestats", putTribeStats)
		r.Delete("/tribe/{uuid}", deleteTribe)
		r.Put("/tribeactivity/{uuid}", putTribeActivity)
		r.Put("/tribepreview/{uuid}", setTribePreview)
		r.Delete("/bot/{uuid}", deleteBot)
		r.Post("/person", createOrEditPerson)
		r.Post("/verify/{challenge}", verify)
		r.Post("/badges", addOrRemoveBadge)
		r.Put("/bot", createOrEditBot)
		r.Delete("/person/{id}", deletePerson)
		r.Delete("/channel/{id}", deleteChannel)
		r.Delete("/ticket/{pubKey}/{created}", deleteTicketByAdmin)
	})

	r.Group(func(r chi.Router) {
		r.Post("/connectioncodes", createConnectionCode)
		r.Get("/connectioncodes", getConnectionCode)
	})

	r.Group(func(r chi.Router) {
		r.Get("/lnauth_login", receiveLnAuthData)
		r.Get("/lnauth", getLnurlAuth)
		r.Get("/lnauth_poll", pollLnurlAuth)
		r.Get("/refresh_jwt", refreshToken)
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
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-User", "authorization", "x-jwt", "Referer", "User-Agent"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		// Debug:            true,
	})
	r.Use(cors.Handler)
	r.Use(middleware.Timeout(60 * time.Second))
	return r
}
