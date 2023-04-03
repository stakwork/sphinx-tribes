package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
)

type database struct {
	db *gorm.DB
}

// DB is the object
var DB database

func initDB() {
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
	db, err := gorm.Open("postgres", dbURL)
	db.LogMode(true)
	if err != nil {
		panic(err)
	}

	DB.db = db

	fmt.Println("db connected")

	// migrate table changes
	db.AutoMigrate(&Person{}, &Channel{}, &LeaderBoard{}, &ConnectionCodes{})

	people := DB.getAllPeople()
	for _, p := range people {
		if p.Uuid == "" {
			DB.addUuidToPerson(p.ID, xid.New().String())
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
func (db database) createOrEditTribe(m Tribe) (Tribe, error) {
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
	if err := db.db.Set("gorm:insert_option", onConflict).Create(&m).Error; err != nil {
		fmt.Println(">>>>>>>>> == ", err)
		return Tribe{}, err
	}
	db.db.Exec(`UPDATE tribes SET tsv =
  	setweight(to_tsvector(name), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C')
	WHERE uuid = '` + m.UUID + "'")
	return m, nil
}

func (db database) createChannel(c Channel) (Channel, error) {

	if c.Created == nil {
		now := time.Now()
		c.Created = &now

	}
	db.db.Create(&c)
	return c, nil

}

// check that update owner_pub_key does in fact throw error
func (db database) createOrEditBot(b Bot) (Bot, error) {
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
	if err := db.db.Set("gorm:insert_option", onConflict).Create(&b).Error; err != nil {
		fmt.Println(err)
		return Bot{}, err
	}
	db.db.Exec(`UPDATE bots SET tsv =
  	setweight(to_tsvector(name), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C')
	WHERE uuid = '` + b.UUID + "'")
	return b, nil
}

// check that update owner_pub_key does in fact throw error
func (db database) createOrEditPerson(m Person) (Person, error) {
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
	if err := db.db.Set("gorm:insert_option", onConflict).Create(&m).Error; err != nil {
		fmt.Println(err)
		return Person{}, err
	}
	db.db.Exec(`UPDATE people SET tsv =
  	setweight(to_tsvector(owner_alias), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C')
	WHERE id = '` + strconv.Itoa(int(m.ID)) + "'")
	return m, nil
}

func (db database) getUnconfirmedTwitter() []Person {
	ms := []Person{}
	db.db.Raw(`SELECT * FROM people where extras -> 'twitter' IS NOT NULL and twitter_confirmed = 'f';`).Find(&ms)
	return ms
}

func (db database) updateTwitterConfirmed(id uint, confirmed bool) {
	if id == 0 {
		return
	}
	db.db.Model(&Person{}).Where("id = ?", id).Updates(map[string]interface{}{
		"twitter_confirmed": confirmed,
	})
}

func (db database) addUuidToPerson(id uint, uuid string) {
	if id == 0 {
		return
	}
	db.db.Model(&Person{}).Where("id = ?", id).Updates(map[string]interface{}{
		"uuid": uuid,
	})
}

func (db database) getUnconfirmedGithub() []Person {
	ms := []Person{}
	db.db.Raw(`SELECT * FROM people where extras -> 'github' IS NOT NULL and github_confirmed = 'f';`).Find(&ms)
	return ms
}

func (db database) updateGithubConfirmed(id uint, confirmed bool) {
	if id == 0 {
		return
	}
	db.db.Model(&Person{}).Where("id = ?", id).Updates(map[string]interface{}{
		"github_confirmed": confirmed,
	})
}

func (db database) updateGithubIssues(id uint, issues map[string]interface{}) {
	db.db.Model(&Person{}).Where("id = ?", id).Updates(map[string]interface{}{
		"github_issues": issues,
	})
}

func (db database) updateTribe(uuid string, u map[string]interface{}) bool {
	if uuid == "" {
		return false
	}
	db.db.Model(&Tribe{}).Where("uuid = ?", uuid).Updates(u)
	return true
}

func (db database) updateChannel(id uint, u map[string]interface{}) bool {
	if id == 0 {
		return false
	}
	db.db.Model(&Channel{}).Where("id= ?", id).Updates(u)
	return true
}

func (db database) updatePerson(id uint, u map[string]interface{}) bool {
	if id == 0 {
		return false
	}
	db.db.Model(&Person{}).Where("id = ?", id).Updates(u)
	return true
}

func (db database) updateTribeUniqueName(uuid string, u string) {
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

func (db database) getOpenGithubIssues(r *http.Request) (int64, error) {
	ms := []GithubOpenIssue{}

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

func (db database) getListedTribes(r *http.Request) []Tribe {
	ms := []Tribe{}
	keys := r.URL.Query()
	tags := keys.Get("tags") // this is a string of tags separated by commas
	offset, limit, sortBy, direction, search := getPaginationParams(r)

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

func (db database) getTribesByOwner(pubkey string) []Tribe {
	ms := []Tribe{}
	db.db.Where("owner_pub_key = ? AND (unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)", pubkey).Find(&ms)
	return ms
}

func (db database) getAllTribesByOwner(pubkey string) []Tribe {
	ms := []Tribe{}
	db.db.Where("owner_pub_key = ? AND (deleted = 'f' OR deleted is null)", pubkey).Find(&ms)
	return ms
}

func (db database) getChannelsByTribe(tribe_uuid string) []Channel {
	ms := []Channel{}
	db.db.Where("tribe_uuid = ? AND (deleted = 'f' OR deleted is null)", tribe_uuid).Find(&ms)
	return ms
}

func (db database) getChannel(id uint) Channel {
	ms := Channel{}
	db.db.Where("id = ?  AND (deleted = 'f' OR deleted is null)", id).Find(&ms)
	return ms
}

func (db database) getListedBots(r *http.Request) []Bot {
	ms := []Bot{}
	offset, limit, sortBy, direction, search := getPaginationParams(r)

	// db.db.Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Find(&ms)
	db.db.Offset(offset).Limit(limit).Order(sortBy+" "+direction).Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Where("LOWER(name) LIKE ?", "%"+search+"%").Find(&ms)

	return ms
}

func (db database) getListedPeople(r *http.Request) []Person {
	ms := []Person{}
	offset, limit, sortBy, direction, search := getPaginationParams(r)

	// if search is empty, returns all
	db.db.Offset(offset).Limit(limit).Order(sortBy+" "+direction+" NULLS LAST").Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Where("LOWER(owner_alias) LIKE ?", "%"+search+"%").Find(&ms)
	return ms
}

func (db database) getAllPeople() []Person {
	ms := []Person{}
	// if search is empty, returns all
	db.db.Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").Find(&ms)
	return ms
}

func (db database) getPeopleBySearch(r *http.Request) []Person {
	ms := []Person{}
	offset, limit, sortBy, direction, search := getPaginationParams(r)

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

func addNotMineToExtrasRawQuery(query string, pubkey string) string {
	return query + ` AND people.owner_pub_key != ` + pubkey + ` `
}

func (db database) getListedPosts(r *http.Request) ([]PeopleExtra, error) {
	ms := []PeopleExtra{}
	// set limit

	offset, limit, sortBy, _, search := getPaginationParams(r)

	rawQuery := makeExtrasListQuery("post")

	// if logged in, dont get mine
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)
	if pubKeyFromAuth != "" {
		rawQuery = addNotMineToExtrasRawQuery(rawQuery, pubKeyFromAuth)
	}

	// sort by newest
	result := db.db.Offset(offset).Limit(limit).Order("arr.item_object->>'"+sortBy+"' DESC").Raw(
		rawQuery, "%"+search+"%").Find(&ms)

	return ms, result.Error
}

func (db database) getListedWanteds(r *http.Request) ([]PeopleExtra, error) {

	ms := []PeopleExtra{}
	// set limit
	offset, limit, sortBy, _, search := getPaginationParams(r)

	rawQuery := makeExtrasListQuery("wanted")

	// 3/1/2022 = 1646172712, we do this to disclude early test tickets
	rawQuery = addNewerThanTimestampToExtrasRawQuery(rawQuery, 1646172712)

	// if logged in, dont get mine
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)
	if pubKeyFromAuth != "" {
		rawQuery = addNotMineToExtrasRawQuery(rawQuery, pubKeyFromAuth)
	}

	// sort by newest
	result := db.db.Offset(offset).Limit(limit).Order("arr.item_object->>'"+sortBy+"' DESC").Raw(
		rawQuery, "%"+search+"%").Find(&ms)

	return ms, result.Error
}

func (db database) getPeopleForNewTicket(languages []interface{}) ([]Person, error) {
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

func (db database) getListedOffers(r *http.Request) ([]PeopleExtra, error) {
	ms := []PeopleExtra{}
	// set limit
	offset, limit, sortBy, _, search := getPaginationParams(r)

	rawQuery := makeExtrasListQuery("offer")

	// if logged in, dont get mine
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)
	if pubKeyFromAuth != "" {
		rawQuery = addNotMineToExtrasRawQuery(rawQuery, pubKeyFromAuth)
	}

	// sort by newest
	result := db.db.Offset(offset).Limit(limit).Order("arr.item_object->>'"+sortBy+"' DESC").Raw(
		rawQuery, "%"+search+"%").Find(&ms)

	return ms, result.Error
}

func (db database) updateBot(uuid string, u map[string]interface{}) bool {
	if uuid == "" {
		return false
	}
	db.db.Model(&Bot{}).Where("uuid = ?", uuid).Updates(u)
	return true
}

func (db database) getAllTribes() []Tribe {
	ms := []Tribe{}
	db.db.Where("(deleted = 'f' OR de leted is null)").Find(&ms)
	return ms
}

func (db database) getTribesTotal() uint64 {
	var count uint64
	db.db.Model(&Tribe{}).Where("deleted = 'false' OR deleted is null").Count(&count)
	return count
}

func (db database) getTribeByIdAndPubkey(uuid string, pubkey string) Tribe {
	m := Tribe{}
	//db.db.Where("uuid = ? AND (deleted = 'f' OR deleted is null) AND owner_pubkey = ?", uuid, pubkey).Find(&m)
	db.db.Where("uuid = ? AND owner_pub_key = ?", uuid, pubkey).Find(&m)
	return m
}

func (db database) getTribe(uuid string) Tribe {
	m := Tribe{}
	db.db.Where("uuid = ? AND (deleted = 'f' OR deleted is null)", uuid).Find(&m)
	return m
}

func (db database) getPerson(id uint) Person {
	m := Person{}
	db.db.Where("id = ? AND (deleted = 'f' OR deleted is null)", id).Find(&m)
	return m
}

func (db database) getPersonByPubkey(pubkey string) Person {
	m := Person{}
	db.db.Where("owner_pub_key = ? AND (deleted = 'f' OR deleted is null)", pubkey).Find(&m)

	return m
}

func (db database) getPersonByUuid(uuid string) Person {
	m := Person{}
	db.db.Where("uuid = ? AND (deleted = 'f' OR deleted is null)", uuid).Find(&m)

	return m
}

func (db database) getPersonByGithubName(github_name string) Person {
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

func (db database) getFirstTribeByFeedURL(feedURL string) Tribe {
	m := Tribe{}
	db.db.Where("feed_url = ? AND (deleted = 'f' OR deleted is null)", feedURL).First(&m)
	return m
}

func (db database) getBot(uuid string) Bot {
	m := Bot{}
	db.db.Where("uuid = ? AND (deleted = 'f' OR deleted is null)", uuid).Find(&m)
	return m
}

func (db database) getTribeByUniqueName(un string) Tribe {
	m := Tribe{}
	db.db.Where("unique_name = ? AND (deleted = 'f' OR deleted is null)", un).Find(&m)
	return m
}

func (db database) getBotsByOwner(pubkey string) []Bot {
	bs := []Bot{}
	db.db.Where("owner_pub_key = ?", pubkey).Find(&bs)
	return bs
}

func (db database) getBotByUniqueName(un string) Bot {
	m := Bot{}
	db.db.Where("unique_name = ? AND (deleted = 'f' OR deleted is null)", un).Find(&m)
	return m
}

func (db database) getPersonByUniqueName(un string) Person {
	m := Person{}
	db.db.Where("unique_name = ? AND (deleted = 'f' OR deleted is null)", un).Find(&m)
	return m
}

func (db database) searchTribes(s string) []Tribe {
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

func (db database) searchBots(s string, limit, offset int) []BotRes {
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

func (db database) searchPeople(s string, limit, offset int) []Person {
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

func (db database) createLeaderBoard(uuid string, leaderboards []LeaderBoard) ([]LeaderBoard, error) {
	m := LeaderBoard{}
	db.db.Where("tribe_uuid = ?", uuid).Delete(&m)
	for _, leaderboard := range leaderboards {
		leaderboard.TribeUuid = uuid
		db.db.Create(leaderboard)
	}
	return leaderboards, nil

}

func (db database) getLeaderBoard(uuid string) []LeaderBoard {
	m := []LeaderBoard{}
	db.db.Where("tribe_uuid = ?", uuid).Find(&m)
	return m
}

func (db database) getLeaderBoardByUuidAndAlias(uuid string, alias string) LeaderBoard {
	m := LeaderBoard{}
	db.db.Where("tribe_uuid = ? and alias = ?", uuid, alias).Find(&m)
	return m
}

func (db database) updateLeaderBoard(uuid string, alias string, u map[string]interface{}) bool {
	if uuid == "" {
		return false
	}
	db.db.Model(&LeaderBoard{}).Where("tribe_uuid = ? and alias = ?", uuid, alias).Updates(u)
	return true
}

func (db database) countDevelopers() uint64 {
	var count uint64
	db.db.Model(&Person{}).Where("deleted = 'f' OR deleted is null").Count(&count)
	return count
}

func (db database) countBounties() uint64 {
	var count struct {
		Sum uint64 `db:"sum"`
	}
	db.db.Raw(`Select sum(jsonb_array_length(extras -> 'wanted')) from people where 
                   people.deleted = 'f' OR people.deleted is null`).Scan(&count)
	return count.Sum
}

func (db database) getPeopleListShort(count uint32) *[]PersonInShort {
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

func (db database) createConnectionCode(c ConnectionCodes) (ConnectionCodes, error) {
	if c.DateCreated == nil {
		now := time.Now()
		c.DateCreated = &now
	}
	db.db.Create(&c)
	return c, nil
}

func (db database) getConnectionCode() ConnectionCodesShort {
	c := ConnectionCodesShort{}

	db.db.Raw(`SELECT connection_string, date_created FROM connectioncodes WHERE is_used =? ORDER BY id DESC LIMIT 1`, false).Find(&c)

	db.db.Model(&ConnectionCodes{}).Where("connection_string = ?", c.ConnectionString).Updates(map[string]interface{}{
		"is_used": true,
	})

	return c
}

func (db database) getLnUser(lnKey string) uint64 {
	var count uint64

	db.db.Model(&Person{}).Where("owner_pub_key = ?", lnKey).Count(&count)

	return count
}

func (db database) createLnUser(lnKey string) (Person, error) {
	now := time.Now()
	p := Person{}

	if db.getLnUser(lnKey) == 0 {
		p.OwnerPubKey = lnKey
		p.OwnerAlias = lnKey
		p.UniqueName, _ = personUniqueNameFromName(p.OwnerAlias)
		p.Created = &now
		p.Tags = pq.StringArray{}
		p.Uuid = xid.New().String()
		p.Extras = map[string]interface{}{}
		p.GithubIssues = map[string]interface{}{}

		db.db.Create(&p)
	}
	return p, nil
}

type Extras struct {
	Owner_pubkey             string `json:"owner_pubkey"`
	Total_bounties_completed uint   `json:"total_bounties_completed"`
	Total_sats_earned        uint   `json:"total_sats_earned"`
}

func (db database) getBountiesLeaderboard() []Extras {
	ms := []Extras{}

	db.db.Raw(`SELECT item->'assignee'->>'owner_pubkey' as owner_pubkey, COUNT(item->'assignee'->>'owner_pubkey') as total_bounties_completed, item->>'price' as total_sats_earned from people  p, LATERAL jsonb_array_elements(p.extras->'wanted') r (item) where extras
	#> '{wanted}' is not null GROUP BY r.item;`).Find(&ms)

	return ms
}
