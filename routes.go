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

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/frontend"
	"github.com/stakwork/sphinx-tribes/handlers"
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
		r.Get("/tribes", handlers.GetListedTribes)
		r.Get("/tribes/{uuid}", handlers.GetTribe)
		r.Get("/tribe_by_feed", handlers.GetFirstTribeByFeed)
		r.Get("/tribes/total", handlers.GetTotalribes)
		r.Get("/tribe_by_un/{un}", handlers.GetTribeByUniqueName)

		r.Get("/leaderboard/{tribe_uuid}", handlers.GetLeaderBoard)

		r.Get("/tribes_by_owner/{pubkey}", handlers.GetTribesByOwner)
		r.Post("/tribes", handlers.CreateOrEditTribe)
		r.Post("/bots", handlers.CreateOrEditBot)
		r.Get("/bots", handlers.GetListedBots)
		r.Get("/bots/owner/{pubkey}", handlers.GetBotsByOwner)
		r.Get("/bots/{uuid}", handlers.GetBot)

		r.Get("/bot/{name}", handlers.GetBotByUniqueName)
		r.Get("/search/bots/{query}", handlers.SearchBots)
		r.Get("/podcast", handlers.GetPodcast)
		r.Get("/feed", handlers.GetGenericFeed)
		r.Get("/search_podcasts", handlers.SearchPodcasts)
		r.Get("/search_youtube", handlers.SearchYoutube)
		r.Get("/youtube_videos", handlers.YoutubeVideosForChannel)
		r.Get("/people", handlers.GetListedPeople)
		r.Get("/people/search", handlers.GetPeopleBySearch)
		r.Get("/people/posts", handlers.GetListedPosts)
		r.Get("/people/wanteds", handlers.GetListedWanteds)
		r.Get("/people/wanteds/assigned/{pubkey}", handlers.GetPersonAssignedWanteds)
		r.Get("/people/wanteds/header", handlers.GetWantedsHeader)
		r.Get("/people/short", handlers.GetPeopleShortList)
		r.Get("/people/offers", handlers.GetListedOffers)
		r.Get("/admin_pubkeys", handlers.GetAdminPubkeys)
		r.Get("/people/bounty/leaderboard", handlers.GetBountiesLeaderboard)

		r.Get("/ask", db.Ask)
		r.Get("/poll/{challenge}", db.Poll)
		r.Get("/person/{pubkey}", handlers.GetPersonByPubkey)
		r.Get("/person/uuid/{uuid}", handlers.GetPersonByUuid)
		r.Get("/person/uuid/{uuid}/assets", handlers.GetPersonAssetsByUuid)
		r.Get("/person/githubname/{github}", handlers.GetPersonByGithubName)

		r.Get("/github_issue/{owner}/{repo}/{issue}", handlers.GetGithubIssue)
		r.Get("/github_issue/status/open", handlers.GetOpenGithubIssues)
		r.Post("/save", db.PostSave)
		r.Get("/save/{key}", db.PollSave)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)
		r.Post("/channel", handlers.CreateChannel)
		r.Post("/leaderboard/{tribe_uuid}", handlers.CreateLeaderBoard)
		r.Put("/leaderboard/{tribe_uuid}", handlers.UpdateLeaderBoard)
		r.Put("/tribe", handlers.CreateOrEditTribe)
		r.Put("/tribestats", handlers.PutTribeStats)
		r.Delete("/tribe/{uuid}", handlers.DeleteTribe)
		r.Put("/tribeactivity/{uuid}", handlers.PutTribeActivity)
		r.Put("/tribepreview/{uuid}", handlers.SetTribePreview)
		r.Delete("/bot/{uuid}", handlers.DeleteBot)
		r.Post("/person", handlers.CreateOrEditPerson)
		r.Post("/verify/{challenge}", db.Verify)
		r.Post("/badges", handlers.AddOrRemoveBadge)
		r.Put("/bot", handlers.CreateOrEditBot)
		r.Delete("/person/{id}", handlers.DeletePerson)
		r.Delete("/channel/{id}", handlers.DeleteChannel)
		r.Delete("/ticket/{pubKey}/{created}", handlers.DeleteTicketByAdmin)
	})

	r.Group(func(r chi.Router) {
		r.Post("/connectioncodes", handlers.CreateConnectionCode)
		r.Get("/connectioncodes", handlers.GetConnectionCode)
	})

	r.Group(func(r chi.Router) {
		r.Get("/lnauth_login", handlers.ReceiveLnAuthData)
		r.Get("/lnauth", handlers.GetLnurlAuth)
		r.Get("/lnauth_poll", handlers.PollLnurlAuth)
		r.Get("/refresh_jwt", handlers.RefreshToken)
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
