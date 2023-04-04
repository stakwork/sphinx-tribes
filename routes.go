package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"

	"github.com/stakwork/sphinx-tribes/feeds"
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

func getAdminPubkeys(w http.ResponseWriter, r *http.Request) {
	adminPubKeys := os.Getenv("ADMIN_PUBKEYS")
	admins := strings.Split(adminPubKeys, ",")
	type PubKeysReturn struct {
		Pubkeys []string `json:"pubkeys"`
	}
	pubkeys := PubKeysReturn{}
	if adminPubKeys != "" {
		for _, admin := range admins {
			pubkeys.Pubkeys = append(pubkeys.Pubkeys, admin)
		}
	}
	json.NewEncoder(w).Encode(pubkeys)
	w.WriteHeader(http.StatusOK)
}

func getGenericFeed(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	feed, err := feeds.ParseFeed(url)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tribeUUID := r.URL.Query().Get("uuid")
	tribe := Tribe{}
	if tribeUUID != "" {
		tribe = DB.getTribe(tribeUUID)
	} else {
		tribe = DB.getFirstTribeByFeedURL(url)
	}

	feed.Value = feeds.AddedValue(feed.Value, tribe.OwnerPubKey)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feed)
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

func searchPodcasts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	podcasts, err := searchPodcastIndex(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fs := []feeds.Feed{}
	for _, pod := range podcasts {
		feed, err1 := feeds.PodcastToGeneric(pod.URL, &pod)
		if err1 == nil {
			fs = append(fs, feed)
		}
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(fs)
}

func searchYoutube(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	fs, err := feeds.YoutubeSearch(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(fs)
}

func youtubeVideosForChannel(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("channelId")
	fs, err := feeds.YoutubeVideosForChannel(q)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(fs)
}

func getAllTribes(w http.ResponseWriter, r *http.Request) {
	tribes := DB.getAllTribes()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

func getTotalribes(w http.ResponseWriter, r *http.Request) {
	tribesTotal := DB.getTribesTotal()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribesTotal)
}

func getListedTribes(w http.ResponseWriter, r *http.Request) {
	tribes := DB.getListedTribes(r)
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
func getPeopleBySearch(w http.ResponseWriter, r *http.Request) {
	people := DB.getPeopleBySearch(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}
func getListedPeople(w http.ResponseWriter, r *http.Request) {
	people := DB.getListedPeople(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}
func getListedPosts(w http.ResponseWriter, r *http.Request) {
	people, err := DB.getListedPosts(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(people)
	}
}

func getListedWanteds(w http.ResponseWriter, r *http.Request) {
	people, err := DB.getListedWanteds(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(people)
	}
}

func getPersonAssignedWanteds(w http.ResponseWriter, r *http.Request) {
	people, err := DB.getListedWanteds(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(people)
	}
}

func getWantedsHeader(w http.ResponseWriter, r *http.Request) {
	var ret struct {
		DeveloperCount uint64           `json:"developer_count"`
		BountiesCount  uint64           `json:"bounties_count"`
		People         *[]PersonInShort `json:"people"`
	}
	ret.DeveloperCount = DB.countDevelopers()
	ret.BountiesCount = DB.countBounties()
	ret.People = DB.getPeopleListShort(3)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func getPeopleShortList(w http.ResponseWriter, r *http.Request) {
	var maxCount uint32 = 10000
	people := DB.getPeopleListShort(maxCount)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}

func getListedOffers(w http.ResponseWriter, r *http.Request) {
	people, err := DB.getListedOffers(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(people)
	}
}

func getTribe(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	tribe := DB.getTribe(uuid)

	var theTribe map[string]interface{}
	j, _ := json.Marshal(tribe)
	json.Unmarshal(j, &theTribe)

	theTribe["channels"] = DB.getChannelsByTribe(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theTribe)
}

func getFirstTribeByFeed(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	tribe := DB.getFirstTribeByFeedURL(url)

	if tribe.UUID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var theTribe map[string]interface{}
	j, _ := json.Marshal(tribe)
	json.Unmarshal(j, &theTribe)

	theTribe["channels"] = DB.getChannelsByTribe(tribe.UUID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theTribe)
}

func getTribeByUniqueName(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "un")
	tribe := DB.getTribeByUniqueName(uuid)

	var theTribe map[string]interface{}
	j, _ := json.Marshal(tribe)
	json.Unmarshal(j, &theTribe)

	theTribe["channels"] = DB.getChannelsByTribe(tribe.UUID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theTribe)
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

func setTribePreview(w http.ResponseWriter, r *http.Request) {
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

	preview := r.URL.Query().Get("preview")
	DB.updateTribe(uuid, map[string]interface{}{
		"preview": preview,
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

func deleteChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

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

	existing := DB.getChannel(uint(id))
	existingTribe := DB.getTribe(existing.TribeUUID)
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

	DB.updateChannel(uint(id), map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func createChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	channel := Channel{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &channel)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	//check that the tribe has the same pubKeyFromAuth
	tribe := DB.getTribe(channel.TribeUUID)
	if tribe.OwnerPubKey != pubKeyFromAuth {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	tribeChannels := DB.getChannelsByTribe(channel.TribeUUID)
	for _, tribeChannel := range tribeChannels {
		if tribeChannel.Name == channel.Name {
			fmt.Println("Channel name already in use")
			w.WriteHeader(http.StatusNotAcceptable)
			return

		}
	}

	channel, err = DB.createChannel(channel)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(channel)
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

func createLeaderBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)
	uuid := chi.URLParam(r, "tribe_uuid")

	leaderBoard := []LeaderBoard{}

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

	//from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &leaderBoard)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = DB.createLeaderBoard(uuid, leaderBoard)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func getLeaderBoard(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "tribe_uuid")
	alias := r.URL.Query().Get("alias")

	if alias == "" {
		leaderBoards := DB.getLeaderBoard(uuid)

		var board = []LeaderBoard{}
		for _, leaderboard := range leaderBoards {
			leaderboard.TribeUuid = ""
			board = append(board, leaderboard)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(board)
	} else {
		leaderBoardFromDb := DB.getLeaderBoardByUuidAndAlias(uuid, alias)

		if leaderBoardFromDb.Alias != alias {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(leaderBoardFromDb)
	}
}

func updateLeaderBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)
	uuid := chi.URLParam(r, "tribe_uuid")

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

	//from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	leaderBoard := LeaderBoard{}

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &leaderBoard)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	leaderBoardFromDb := DB.getLeaderBoardByUuidAndAlias(uuid, leaderBoard.Alias)

	if leaderBoardFromDb.Alias != leaderBoard.Alias {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	leaderBoard.TribeUuid = leaderBoardFromDb.TribeUuid

	DB.updateLeaderBoard(leaderBoardFromDb.TribeUuid, leaderBoardFromDb.Alias, map[string]interface{}{
		"spent":      leaderBoard.Spent,
		"earned":     leaderBoard.Earned,
		"reputation": leaderBoard.Reputation,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func createConnectionCode(w http.ResponseWriter, r *http.Request) {
	code := ConnectionCodes{}
	now := time.Now()

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &code)

	code.IsUsed = false
	code.DateCreated = &now

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = DB.createConnectionCode(code)

	if err != nil {
		fmt.Println("=> ERR create connection code", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func getConnectionCode(w http.ResponseWriter, _ *http.Request) {
	connectionCode := DB.getConnectionCode()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(connectionCode)
}

func getLnurlAuth(w http.ResponseWriter, _ *http.Request) {
	encodeData, err := encodeLNURL()
	responseData := make(map[string]string)

	if err != nil {
		responseData["k1"] = ""
		responseData["encode"] = ""

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not generate LNURL AUTH")
	}

	store.SetLnCache(encodeData.k1, LnStore{encodeData.k1, "", false})

	responseData["k1"] = encodeData.k1
	responseData["encode"] = encodeData.encode

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func pollLnurlAuth(w http.ResponseWriter, r *http.Request) {
	k1 := r.URL.Query().Get("k1")
	responseData := make(map[string]interface{})

	res, err := store.GetLnCache(k1)

	if err != nil {
		responseData["k1"] = ""
		responseData["status"] = false

		fmt.Println("=> ERR polling LNURL AUTH", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("LNURL auth data not found")
	}

	tokenString, err := EncodeToken(res.key)

	if err != nil {
		fmt.Println("error creating JWT")
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	person := DB.getPersonByPubkey(res.key)
	user := returnUserMap(person)

	responseData["k1"] = res.k1
	responseData["status"] = res.status
	responseData["jwt"] = tokenString
	responseData["user"] = user

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func receiveLnAuthData(w http.ResponseWriter, r *http.Request) {
	userKey := r.URL.Query().Get("key")
	k1 := r.URL.Query().Get("k1")

	responseMsg := make(map[string]string)

	if userKey != "" {
		// Save in DB if the user does not exists already
		DB.createLnUser(userKey)

		// Set store data to true
		store.SetLnCache(k1, LnStore{k1, userKey, true})

		responseMsg["status"] = "OK"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseMsg)
	}

	responseMsg["status"] = "ERROR"
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(responseMsg)
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-jwt")

	responseData := make(map[string]interface{})
	claims, err := DecodeToken(token)

	if err != nil {
		fmt.Println("Failed to parse JWT")
		http.Error(w, http.StatusText(401), 401)
		return
	}

	pubkey := fmt.Sprint(claims["pubkey"])

	userCount := DB.getLnUser(pubkey)

	if userCount > 0 {
		// Generate a new token
		tokenString, err := EncodeToken(pubkey)

		if err != nil {
			fmt.Println("error creating  refresh JWT")
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		person := DB.getPersonByPubkey(pubkey)
		user := returnUserMap(person)

		responseData["k1"] = ""
		responseData["status"] = true
		responseData["jwt"] = tokenString
		responseData["user"] = user

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseData)
	}
}

func returnUserMap(p Person) map[string]interface{} {
	user := make(map[string]interface{})

	user["id"] = p.ID
	user["created"] = p.Created
	user["owner_pubkey"] = p.OwnerPubKey
	user["owner_alias"] = p.OwnerAlias
	user["contact_key"] = p.OwnerContactKey
	user["img"] = p.Img
	user["description"] = p.Description
	user["tags"] = p.Tags
	user["unique_name"] = p.UniqueName
	user["pubkey"] = p.OwnerPubKey
	user["extras"] = p.Extras
	user["last_login"] = p.LastLogin
	user["price_to_meet"] = p.PriceToMeet
	user["alias"] = p.OwnerAlias
	user["url"] = host

	return user
}
