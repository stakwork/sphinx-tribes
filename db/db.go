package db

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/rs/xid"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

// DB is the object
var DB database

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	fmt.Printf("db url : %v", dbURL)

	if dbURL == "" {
		rdsHost := os.Getenv("RDS_HOSTNAME")
		rdsPort := os.Getenv("RDS_PORT")
		rdsDbName := os.Getenv("RDS_DB_NAME")
		rdsUsername := os.Getenv("RDS_USERNAME")
		rdsPassword := os.Getenv("RDS_PASSWORD")
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", rdsUsername, rdsPassword, rdsHost, rdsPort, rdsDbName)
	}

	if dbURL == "" {
		panic("DB env vars not found")
	}

	var err error

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB.db = db

	fmt.Println("db connected")

	// migrate table changes
	db.Debug().AutoMigrate(&Person{})
	db.AutoMigrate(&Channel{})
	db.AutoMigrate(&LeaderBoard{})
	db.AutoMigrate(&ConnectionCodes{})

	people := DB.GetAllPeople()
	for _, p := range people {
		if p.Uuid == "" {
			DB.AddUuidToPerson(p.ID, xid.New().String())
		}
	}

}

var updatables = []string{
	"name", "description", "tags", "img",
	"owner_alias", "price_to_join", "price_per_message",
	"escrow_amount", "escrow_millis",
	"unlisted", "private", "deleted",
	"app_url", "bots", "feed_url", "feed_type",
	"owner_route_hint", "updated", "pin",
	"profile_filters",
}
var botupdatables = []string{
	"name", "description", "tags", "img",
	"owner_alias", "price_per_use",
	"unlisted", "deleted",
	"owner_route_hint", "updated",
}
var peopleupdatables = []string{
	"description", "tags", "img",
	"owner_alias",
	"unlisted", "deleted",
	"owner_route_hint",
	"price_to_meet", "updated",
	"extras",
}
var channelupdatables = []string{
	"name", "deleted"}

// check that update owner_pub_key does in fact throw error
func (db database) CreateOrEditTribe(m Tribe) (Tribe, error) {
	if m.OwnerPubKey == "" {
		return Tribe{}, errors.New("no pub key")
	}
	onConflict := "ON CONFLICT (uuid) DO UPDATE SET"
	for i, u := range updatables {
		onConflict = onConflict + fmt.Sprintf(" %s=EXCLUDED.%s", u, u)
		if i < len(updatables)-1 {
			onConflict = onConflict + ","
		}
	}
	if m.Name == "" {
		m.Name = "name"
	}
	if m.Description == "" {
		m.Description = "description"
	}
	if m.Tags == nil {
		m.Tags = []string{}
	}
	if m.Badges == nil {
		m.Badges = []string{}
	}

	if db.db.Model(&m).Where("uuid = ?", m.UUID).Updates(&m).RowsAffected == 0 {
		db.db.Create(&m)
	}

	db.db.Exec(`UPDATE tribes SET tsv =
  	setweight(to_tsvector(name), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C')
	WHERE uuid = '` + m.UUID + "'")
	return m, nil
}

func (db database) CreateChannel(c Channel) (Channel, error) {

	if c.Created == nil {
		now := time.Now()
		c.Created = &now

	}
	db.db.Create(&c)
	return c, nil

}

// check that update owner_pub_key does in fact throw error
func (db database) CreateOrEditBot(b Bot) (Bot, error) {
	if b.OwnerPubKey == "" {
		return Bot{}, errors.New("no pub key")
	}
	if b.UniqueName == "" {
		return Bot{}, errors.New("no unique name")
	}
	onConflict := "ON CONFLICT (uuid) DO UPDATE SET"
	for i, u := range botupdatables {
		onConflict = onConflict + fmt.Sprintf(" %s=EXCLUDED.%s", u, u)
		if i < len(botupdatables)-1 {
			onConflict = onConflict + ","
		}
	}
	if b.Name == "" {
		b.Name = "name"
	}
	if b.Description == "" {
		b.Description = "description"
	}
	if b.Tags == nil {
		b.Tags = []string{}
	}

	if db.db.Model(&b).Where("uuid = ?", b.UUID).Updates(&b).RowsAffected == 0 {
		db.db.Create(&b)
	}

	db.db.Exec(`UPDATE bots SET tsv =
  	setweight(to_tsvector(name), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C')
	WHERE uuid = '` + b.UUID + "'")
	return b, nil
}

// check that update owner_pub_key does in fact throw error
func (db database) CreateOrEditPerson(m Person) (Person, error) {
	if m.OwnerPubKey == "" {
		return Person{}, errors.New("no pub key")
	}
	onConflict := "ON CONFLICT (id) DO UPDATE SET"
	for i, u := range peopleupdatables {
		onConflict = onConflict + fmt.Sprintf(" %s=EXCLUDED.%s", u, u)
		if i < len(peopleupdatables)-1 {
			onConflict = onConflict + ","
		}
	}
	if m.OwnerAlias == "" {
		m.OwnerAlias = "name"
	}
	if m.Description == "" {
		m.Description = "description"
	}
	if m.Tags == nil {
		m.Tags = []string{}
	}
	if m.Extras == nil {
		m.Extras = map[string]interface{}{}
	}
	if m.GithubIssues == nil {
		m.GithubIssues = map[string]interface{}{}
	}

	if db.db.Model(&m).Where("id = ?", m.ID).Updates(&m).RowsAffected == 0 {
		db.db.Create(&m)
	}

	db.db.Exec(`UPDATE people SET tsv =
  	setweight(to_tsvector(owner_alias), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C')
	WHERE id = '` + strconv.Itoa(int(m.ID)) + "'")

	return m, nil
}

func (db database) GetUnconfirmedTwitter() []Person {
	ms := []Person{}
	db.db.Raw(`SELECT * FROM people where extras -> 'twitter' IS NOT NULL and twitter_confirmed = 'f';`).Find(&ms)
	return ms
}

func (db database) UpdateTwitterConfirmed(id uint, confirmed bool) {
	if id == 0 {
		return
	}
	db.db.Model(&Person{}).Where("id = ?", id).Updates(map[string]interface{}{
		"twitter_confirmed": confirmed,
	})
}

func (db database) AddUuidToPerson(id uint, uuid string) {
	if id == 0 {
		return
	}
	db.db.Model(&Person{}).Where("id = ?", id).Updates(map[string]interface{}{
		"uuid": uuid,
	})
}

func (db database) GetUnconfirmedGithub() []Person {
	ms := []Person{}
	db.db.Raw(`SELECT * FROM people where extras -> 'github' IS NOT NULL and github_confirmed = 'f';`).Find(&ms)
	return ms
}

func (db database) UpdateGithubConfirmed(id uint, confirmed bool) {
	if id == 0 {
		return
	}
	db.db.Model(&Person{}).Where("id = ?", id).Updates(map[string]interface{}{
		"github_confirmed": confirmed,
	})
}

func (db database) UpdateGithubIssues(id uint, issues map[string]interface{}) {
	db.db.Model(&Person{}).Where("id = ?", id).Updates(map[string]interface{}{
		"github_issues": issues,
	})
}

func (db database) UpdateTribe(uuid string, u map[string]interface{}) bool {
	if uuid == "" {
		return false
	}
	db.db.Model(&Tribe{}).Where("uuid = ?", uuid).Updates(u)
	return true
}

func (db database) UpdateChannel(id uint, u map[string]interface{}) bool {
	if id == 0 {
		return false
	}
	db.db.Model(&Channel{}).Where("id= ?", id).Updates(u)
	return true
}

func (db database) UpdatePerson(id uint, u map[string]interface{}) bool {
	if id == 0 {
		return false
	}
	db.db.Model(&Person{}).Where("id = ?", id).Updates(u)
	return true
}

func (db database) UpdateTribeUniqueName(uuid string, u string) {
	if uuid == "" {
		return
	}
	// fmt.Println(u)
	db.db.Model(&Tribe{}).Where("uuid = ?", uuid).Update("unique_name", u)
}

type GithubOpenIssue struct {
	Status   string `json:"status"`
	Assignee string `json:"assignee"`
}

type GithubOpenIssueCount struct {
	Count int64 `json:"count"`
}

func (db database) GetOpenGithubIssues(r *http.Request) (int64, error) {
	ms := []GithubOpenIssueCount{}

	// set limit
	result := db.db.Raw(
		`SELECT COUNT(value)
		FROM (
			SELECT * 
			FROM people 
			WHERE github_issues IS NOT NULL 
			AND github_issues != 'null'
			) p,
		jsonb_each(github_issues) t2 
		WHERE value @> '{"status": "open"}' OR value @> '{"status": ""}'`).Find(&ms)

	return result.RowsAffected, result.Error
}

func (db database) GetListedTribes(r *http.Request) []Tribe {
	ms := []Tribe{}
	keys := r.URL.Query()
	tags := keys.Get("tags") // this is a string of tags separated by commas
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)

	thequery := db.db.Offset(offset).Limit(limit).Order(sortBy+" "+direction).Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Where("LOWER(name) LIKE ?", "%"+search+"%")

	if tags != "" {
		// pull out the tags and add them in here
		t := strings.Split(tags, ",")
		for _, s := range t {
			thequery = thequery.Where("'" + s + "'" + " = any (tags)")
		}
	}

	thequery.Find(&ms)
	return ms
}

func (db database) GetTribesByOwner(pubkey string) []Tribe {
	ms := []Tribe{}
	db.db.Where("owner_pub_key = ? AND (unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)", pubkey).Find(&ms)
	return ms
}

func (db database) GetAllTribesByOwner(pubkey string) []Tribe {
	ms := []Tribe{}
	db.db.Where("owner_pub_key = ? AND (deleted = 'f' OR deleted is null)", pubkey).Find(&ms)
	return ms
}

func (db database) GetChannelsByTribe(tribe_uuid string) []Channel {
	ms := []Channel{}
	db.db.Where("tribe_uuid = ? AND (deleted = 'f' OR deleted is null)", tribe_uuid).Find(&ms)
	return ms
}

func (db database) GetChannel(id uint) Channel {
	ms := Channel{}
	db.db.Where("id = ?  AND (deleted = 'f' OR deleted is null)", id).Find(&ms)
	return ms
}

func (db database) GetListedBots(r *http.Request) []Bot {
	ms := []Bot{}
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)

	// db.db.Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Find(&ms)
	db.db.Offset(offset).Limit(limit).Order(sortBy+" "+direction).Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Where("LOWER(name) LIKE ?", "%"+search+"%").Find(&ms)

	return ms
}

func (db database) GetListedPeople(r *http.Request) []Person {
	ms := []Person{}
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)

	// if search is empty, returns all
	db.db.Offset(offset).Limit(limit).Order(sortBy+" "+direction+" NULLS LAST").Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Where("LOWER(owner_alias) LIKE ?", "%"+search+"%").Find(&ms)
	return ms
}

func (db database) GetAllPeople() []Person {
	ms := []Person{}
	// if search is empty, returns all
	db.db.Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Find(&ms)
	return ms
}

func (db database) GetPeopleBySearch(r *http.Request) []Person {
	ms := []Person{}
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)

	// if search is empty, returns all

	// return if like owner_alias, unique_name, or equals pubkey
	db.db.Offset(offset).Limit(limit).Order(sortBy+" "+direction+" NULLS LAST").Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Where("LOWER(owner_alias) LIKE ?", "%"+search+"%").Or("LOWER(unique_name) LIKE ?", "%"+search+"%").Or("LOWER(owner_pub_key) = ?", search).Find(&ms)
	return ms
}

type PeopleExtra struct {
	Body   string `json:"body"`
	Person string `json:"person"`
}

func makeExtrasListQuery(columnName string) string {
	// this is safe because columnName is not provided by the user, its hard-coded in db.go
	return `SELECT 		
	json_build_object('owner_pubkey', owner_pub_key, 'owner_alias', owner_alias, 'img', img, 'unique_name', unique_name, 'id', id, '` + columnName + `', extras->'` + columnName + `', 'github_issues', github_issues) #>> '{}' as person,
	arr.item_object as body
	FROM people,
	jsonb_array_elements(extras->'` + columnName + `') with ordinality 
	arr(item_object, position)
	WHERE people.deleted != true
	AND people.unlisted != true 
	AND LOWER(arr.item_object->>'title') LIKE ?
	AND CASE
			WHEN arr.item_object->>'show' = 'false' THEN false
			ELSE true
		END
	`
}

func makePersonExtrasListQuery(columnName string) string {
	// this is safe because columnName is not provided by the user, its hard-coded in db.go
	return `SELECT 		
	json_build_object('owner_pubkey', owner_pub_key, 'owner_alias', owner_alias, 'img', img, 'unique_name', unique_name, 'id', id, '` + columnName + `', extras->'` + columnName + `', 'github_issues', github_issues) #>> '{}' as person,
	arr.item_object as body
	FROM people,
	jsonb_array_elements(extras->'` + columnName + `') with ordinality 
	arr(item_object, position)
	WHERE arr.item_object->'assignee'->>'owner_pubkey' = ? 
	AND LOWER(arr.item_object->>'title') LIKE ?
	AND CASE
			WHEN arr.item_object->>'show' = 'false' THEN false
			ELSE true
		END`
}

func addNewerThanXDaysToExtrasRawQuery(query string, days int) string {
	secondsInDay := 86400
	newerThan := secondsInDay * days
	t := strconv.Itoa(newerThan)
	return query + ` AND CAST(arr.item_object->>'created' AS INT) > (extract(epoch from now()) - ` + t + `) `
}

func addNewerThanTimestampToExtrasRawQuery(query string, timestamp int) string {
	t := strconv.Itoa(timestamp)
	return query + ` AND CAST(arr.item_object->>'created' AS INT) > ` + t
}

func addOrderToExtrasRawQuery(query string) string {
	return query + `ORDER BY cast(arr.item_object->>'created' as integer) DESC`
}

func addNotMineToExtrasRawQuery(query string, pubkey string) string {
	return query + ` AND people.owner_pub_key != ` + pubkey + ` `
}

func (db database) GetListedPosts(r *http.Request) ([]PeopleExtra, error) {
	ms := []PeopleExtra{}
	// set limit

	offset, limit, sortBy, _, search := utils.GetPaginationParams(r)

	rawQuery := makeExtrasListQuery("post")

	// if logged in, dont get mine
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth != "" {
		rawQuery = addNotMineToExtrasRawQuery(rawQuery, pubKeyFromAuth)
	}

	// sort by newest
	result := db.db.Offset(offset).Limit(limit).Order("arr.item_object->>'"+sortBy+"' DESC").Raw(
		rawQuery, "%"+search+"%").Find(&ms)

	return ms, result.Error
}

func (db database) GetListedWanteds(r *http.Request) ([]PeopleExtra, error) {
	pubkey := chi.URLParam(r, "pubkey")
	ms := []PeopleExtra{}

	var rawQuery string
	var result *gorm.DB
	// set limit
	offset, limit, sortBy, _, search := utils.GetPaginationParams(r)

	if pubkey == "" {
		rawQuery = makeExtrasListQuery("wanted")
	} else {
		rawQuery = makePersonExtrasListQuery("wanted")
	}

	// 3/1/2022 = 1646172712, we do this to disclude early test tickets
	rawQuery = addNewerThanTimestampToExtrasRawQuery(rawQuery, 1646172712)

	// Order the wanted in descending order by created date
	rawQuery = addOrderToExtrasRawQuery(rawQuery)

	// if logged in, dont get mine
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth != "" {
		rawQuery = addNotMineToExtrasRawQuery(rawQuery, pubKeyFromAuth)
	}

	// sort by newest
	if pubkey == "" {
		result = db.db.Offset(offset).Limit(limit).Order("arr.item_object->>'"+sortBy+"' DESC").Raw(
			rawQuery, "%"+search+"%").Find(&ms)
	} else {
		result = db.db.Offset(offset).Limit(limit).Order("arr.item_object->>'"+sortBy+"' DESC").Raw(
			rawQuery, pubkey, "%"+search+"%").Find(&ms)
	}

	return ms, result.Error
}

func (db database) GetPeopleForNewTicket(languages []interface{}) ([]Person, error) {
	ms := []Person{}

	query := "Select owner_pub_key, json_build_object('coding_languages',extras->'coding_languages') as extras from people" +
		" where (deleted != true AND unlisted != true) AND " +
		"extras->'alert' = 'true' AND ("

	for _, lang := range languages {
		l, ok := lang.(map[string]interface{})
		if !ok {
			return ms, errors.New("could not parse coding languages correctly")
		}
		label, ok2 := l["label"].(string)
		if !ok2 {
			return ms, errors.New("could not find label in language")
		}
		query += "extras->'coding_languages' @> '[{\"label\": \"" + label + "\"}]' OR "
	}
	query = query[:len(query)-4]
	query += ");"
	err := db.db.Raw(query).Find(&ms).Error
	return ms, err
}

func (db database) GetListedOffers(r *http.Request) ([]PeopleExtra, error) {
	ms := []PeopleExtra{}
	// set limit
	offset, limit, sortBy, _, search := utils.GetPaginationParams(r)

	rawQuery := makeExtrasListQuery("offer")

	// if logged in, dont get mine
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth != "" {
		rawQuery = addNotMineToExtrasRawQuery(rawQuery, pubKeyFromAuth)
	}

	// sort by newest
	result := db.db.Offset(offset).Limit(limit).Order("arr.item_object->>'"+sortBy+"' DESC").Raw(
		rawQuery, "%"+search+"%").Find(&ms)

	return ms, result.Error
}

func (db database) UpdateBot(uuid string, u map[string]interface{}) bool {
	if uuid == "" {
		return false
	}
	db.db.Model(&Bot{}).Where("uuid = ?", uuid).Updates(u)
	return true
}

func (db database) GetAllTribes() []Tribe {
	ms := []Tribe{}
	db.db.Where("(deleted = 'f' OR de leted is null)").Find(&ms)
	return ms
}

func (db database) GetTribesTotal() int64 {
	var count int64
	db.db.Model(&Tribe{}).Where("deleted = 'false' OR deleted is null").Count(&count)
	return count
}

func (db database) GetTribeByIdAndPubkey(uuid string, pubkey string) Tribe {
	m := Tribe{}
	//db.db.Where("uuid = ? AND (deleted = 'f' OR deleted is null) AND owner_pubkey = ?", uuid, pubkey).Find(&m)
	db.db.Where("uuid = ? AND owner_pub_key = ?", uuid, pubkey).Find(&m)
	return m
}

func (db database) GetTribe(uuid string) Tribe {
	m := Tribe{}
	db.db.Where("uuid = ? AND (deleted = 'f' OR deleted is null)", uuid).Find(&m)
	return m
}

func (db database) GetPerson(id uint) Person {
	m := Person{}
	db.db.Where("id = ? AND (deleted = 'f' OR deleted is null)", id).Find(&m)
	return m
}

func (db database) GetPersonByPubkey(pubkey string) Person {
	m := Person{}
	db.db.Where("owner_pub_key = ? AND (deleted = 'f' OR deleted is null)", pubkey).Find(&m)

	return m
}

func (db database) GetPersonByUuid(uuid string) Person {
	m := Person{}
	db.db.Where("uuid = ? AND (deleted = 'f' OR deleted is null)", uuid).Find(&m)

	return m
}

func (db database) GetPersonByGithubName(github_name string) Person {
	m := Person{}

	db.db.Raw(`SELECT 		
	json_build_object('owner_pubkey', owner_pub_key, 'owner_alias', owner_alias, 'img', img, 'unique_name', unique_name, 'id', id, 'wanted', extras->'wanted', 'github_issues', github_issues) #>> '{}' as person,
	FROM people,
	jsonb_array_elements(extras->'github') with ordinality 
	arr(item_object, position)
	WHERE people.deleted != true
	AND people.unlisted != true
	AND CASE
			WHEN arr.item_object->>'value' = ? THEN true
			ELSE false
		END`, github_name).First(&m)

	return m
}

func (db database) GetFirstTribeByFeedURL(feedURL string) Tribe {
	m := Tribe{}
	db.db.Where("feed_url = ? AND (deleted = 'f' OR deleted is null)", feedURL).First(&m)
	return m
}

func (db database) GetBot(uuid string) Bot {
	m := Bot{}
	db.db.Where("uuid = ? AND (deleted = 'f' OR deleted is null)", uuid).Find(&m)
	return m
}

func (db database) GetTribeByUniqueName(un string) Tribe {
	m := Tribe{}
	db.db.Where("unique_name = ? AND (deleted = 'f' OR deleted is null)", un).Find(&m)
	return m
}

func (db database) GetBotsByOwner(pubkey string) []Bot {
	bs := []Bot{}
	db.db.Where("owner_pub_key = ?", pubkey).Find(&bs)
	return bs
}

func (db database) GetBotByUniqueName(un string) Bot {
	m := Bot{}
	db.db.Where("unique_name = ? AND (deleted = 'f' OR deleted is null)", un).Find(&m)
	return m
}

func (db database) GetPersonByUniqueName(un string) Person {
	m := Person{}
	db.db.Where("unique_name = ? AND (deleted = 'f' OR deleted is null)", un).Find(&m)
	return m
}

func (db database) SearchTribes(s string) []Tribe {
	ms := []Tribe{}
	if s == "" {
		return ms
	}
	// set limit
	db.db.Raw(
		`SELECT uuid, owner_pub_key, name, img, description, ts_rank(tsv, q) as rank
		FROM tribes, to_tsquery(?) q
		WHERE tsv @@ q
		AND (deleted = 'f' OR deleted is null)
		ORDER BY rank DESC LIMIT 100;`, s).Find(&ms)
	return ms
}

func (db database) SearchBots(s string, limit, offset int) []BotRes {
	ms := []BotRes{}
	if s == "" {
		return ms
	}
	// set limit
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)
	db.db.Raw(
		`SELECT uuid, owner_pub_key, name, unique_name, img, description, tags, price_per_use, ts_rank(tsv, q) as rank
		FROM bots, to_tsquery(?) q
		WHERE tsv @@ q
		AND (deleted = 'f' OR deleted is null)
		ORDER BY rank DESC 
		LIMIT ? OFFSET ?;`, s, limitStr, offsetStr).Find(&ms)
	return ms
}

func (db database) SearchPeople(s string, limit, offset int) []Person {
	ms := []Person{}
	if s == "" {
		return ms
	}
	// set limit
	limitStr := strconv.Itoa(limit)
	offsetStr := strconv.Itoa(offset)
	db.db.Raw(
		`SELECT id, owner_pub_key, unique_name, img, description, tags, ts_rank(tsv, q) as rank
		FROM people, to_tsquery(?) q
		WHERE tsv @@ q
		AND (deleted = 'f' OR deleted is null)
		ORDER BY rank DESC 
		LIMIT ? OFFSET ?;`, s, limitStr, offsetStr).Find(&ms)
	return ms
}

func (db database) CreateLeaderBoard(uuid string, leaderboards []LeaderBoard) ([]LeaderBoard, error) {
	m := LeaderBoard{}
	db.db.Where("tribe_uuid = ?", uuid).Delete(&m)
	for _, leaderboard := range leaderboards {
		leaderboard.TribeUuid = uuid
		db.db.Create(leaderboard)
	}
	return leaderboards, nil

}

func (db database) GetLeaderBoard(uuid string) []LeaderBoard {
	m := []LeaderBoard{}
	db.db.Where("tribe_uuid = ?", uuid).Find(&m)
	return m
}

func (db database) GetLeaderBoardByUuidAndAlias(uuid string, alias string) LeaderBoard {
	m := LeaderBoard{}
	db.db.Where("tribe_uuid = ? and alias = ?", uuid, alias).Find(&m)
	return m
}

func (db database) UpdateLeaderBoard(uuid string, alias string, u map[string]interface{}) bool {
	if uuid == "" {
		return false
	}
	db.db.Model(&LeaderBoard{}).Where("tribe_uuid = ? and alias = ?", uuid, alias).Updates(u)
	return true
}

func (db database) CountDevelopers() int64 {
	var count int64
	db.db.Model(&Person{}).Where("deleted = 'f' OR deleted is null").Count(&count)
	return count
}

func (db database) CountBounties() uint64 {
	var count struct {
		Sum uint64 `db:"sum"`
	}
	db.db.Raw(`Select sum(jsonb_array_length(extras -> 'wanted')) from people where 
                   people.deleted = 'f' OR people.deleted is null`).Scan(&count)
	return count.Sum
}

func (db database) GetPeopleListShort(count uint32) *[]PersonInShort {
	p := []PersonInShort{}
	db.db.Raw(
		`SELECT id, owner_pub_key, unique_name, img, uuid, owner_alias
		FROM people
		WHERE
		(deleted = 'f' OR deleted is null)
		ORDER BY random() 
		LIMIT ?;`, count).Find(&p)
	return &p
}

func (db database) CreateConnectionCode(c ConnectionCodes) (ConnectionCodes, error) {
	if c.DateCreated == nil {
		now := time.Now()
		c.DateCreated = &now
	}
	db.db.Create(&c)
	return c, nil
}

func (db database) GetConnectionCode() ConnectionCodesShort {
	c := ConnectionCodesShort{}

	db.db.Raw(`SELECT connection_string, date_created FROM connectioncodes WHERE is_used =? ORDER BY id DESC LIMIT 1`, false).Find(&c)

	db.db.Model(&ConnectionCodes{}).Where("connection_string = ?", c.ConnectionString).Updates(map[string]interface{}{
		"is_used": true,
	})

	return c
}

func (db database) GetLnUser(lnKey string) int64 {
	var count int64

	db.db.Model(&Person{}).Where("owner_pub_key = ?", lnKey).Count(&count)

	return count
}

func (db database) CreateLnUser(lnKey string) (Person, error) {
	now := time.Now()
	p := Person{}

	if db.GetLnUser(lnKey) == 0 {
		p.OwnerPubKey = lnKey
		p.OwnerAlias = lnKey
		p.UniqueName, _ = PersonUniqueNameFromName(p.OwnerAlias)
		p.Created = &now
		p.Tags = pq.StringArray{}
		p.Uuid = xid.New().String()
		p.Extras = make(map[string]interface{})
		p.GithubIssues = make(map[string]interface{})

		db.db.Create(&p)
	}
	return p, nil
}

func PersonUniqueNameFromName(name string) (string, error) {
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
		existing := DB.GetPersonByUniqueName(uniquepath)
		if existing.ID != 0 {
			n = n + 1
		} else {
			path = uniquepath
			break
		}
	}
	return path, nil
}

type Extras struct {
	Owner_pubkey             string `json:"owner_pubkey"`
	Total_bounties_completed uint   `json:"total_bounties_completed"`
	Total_sats_earned        uint   `json:"total_sats_earned"`
}

type LeaderData map[string]interface{}

func (db database) GetBountiesLeaderboard() []LeaderData {
	ms := []Extras{}
	var users = []LeaderData{}

	db.db.Raw(`SELECT item->'assignee'->>'owner_pubkey' as owner_pubkey, COUNT(item->'assignee'->>'owner_pubkey') as total_bounties_completed, item->>'price' as total_sats_earned from people  p, LATERAL jsonb_array_elements(p.extras->'wanted') r (item) where extras
	#> '{wanted}' is not null GROUP BY r.item;`).Find(&ms)

	for _, val := range ms {
		var newLeader = make(map[string]interface{})
		found, index := GetLeaderData(users, val.Owner_pubkey)

		if found == -1 {
			newLeader["owner_pubkey"] = val.Owner_pubkey
			newLeader["total_bounties_completed"] = val.Total_bounties_completed
			newLeader["total_sats_earned"] = val.Total_sats_earned

			users = append(users, newLeader)
		} else {
			total_bounties := users[index]["total_bounties_completed"].(uint)
			total_sats := users[index]["total_sats_earned"].(uint)

			users[index]["total_bounties_completed"] = total_bounties + val.Total_bounties_completed
			users[index]["total_sats_earned"] = total_sats + val.Total_sats_earned
		}
	}
	return users
}

func GetLeaderData(arr []LeaderData, key string) (int, int) {
	found := -1
	index := 0

	for i, v := range arr {
		if v["owner_pubkey"] == key {
			found = 1
			index = i
		}
	}
	return found, index
}
