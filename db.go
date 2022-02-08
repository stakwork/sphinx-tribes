package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type database struct {
	db *gorm.DB
}

// DB is the object
var DB database

func initDB() {
	dbURL := os.Getenv("DATABASE_URL")
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
	db.AutoMigrate(&Person{})

	// data := map[string]string{
	// 	"assignee": "Evanfeenstra",
	// 	"status":   "open",
	// }
	// DB.updateGithubIssues(1, map[string]interface{}{
	// 	"stakwork/sphinx-relay/229": data,
	// })
}

var updatables = []string{
	"name", "description", "tags", "img",
	"owner_alias", "price_to_join", "price_per_message",
	"escrow_amount", "escrow_millis",
	"unlisted", "private", "deleted",
	"app_url", "bots", "feed_url", "feed_type",
	"owner_route_hint", "updated",
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
	if err := db.db.Set("gorm:insert_option", onConflict).Create(&m).Error; err != nil {
		fmt.Println(err)
		return Tribe{}, err
	}
	db.db.Exec(`UPDATE tribes SET tsv =
  	setweight(to_tsvector(name), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C')
	WHERE uuid = '` + m.UUID + "'")
	return m, nil
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

func (db database) getOpenGithubIssues(r *http.Request) ([]GithubOpenIssue, error) {
	ms := []GithubOpenIssue{}
	// set limit
	result := db.db.Raw(
		`SELECT value
		FROM (
			SELECT * 
			FROM people 
			WHERE github_issues IS NOT NULL 
			AND github_issues != 'null'
			) p,
		jsonb_each(github_issues) t2 
		WHERE value @> '{"status": "open"}' OR value @> '{"status": ""}'`).Find(&ms)

	return ms, result.Error
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

type PeopleExtra struct {
	Body   string `json:"body"`
	Person string `json:"person"`
}

func (db database) getListedPosts(r *http.Request) ([]PeopleExtra, error) {
	ms := []PeopleExtra{}
	// set limit

	offset, limit, sortBy, direction, search := getPaginationParams(r)

	result := db.db.Offset(offset).Limit(limit).Order(sortBy + " " + direction).Raw(
		`SELECT 
		json_build_object('owner_pubkey', owner_pub_key, 'owner_alias', owner_alias, 'img', img, 'unique_name', unique_name, 'id', id, 'post', extras->'post') #>> '{}' as person,
		to_json(jsonb_array_elements(extras->'post'::text)) #>> '{}' as body 
		FROM people
		WHERE (unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)
		AND extras IS NOT NULL
		AND extras != '{}'::jsonb
		AND extras->'post' IS NOT NULL
		AND extras->'post' != '[]'::jsonb
		AND EXISTS (SELECT *
			FROM jsonb_array_elements_Text(extras -> 'post') as x(title)
			WHERE LOWER(x.title) LIKE '%` + search + `%')`).Find(&ms)

	return ms, result.Error
}

func (db database) getListedWanteds(r *http.Request) ([]PeopleExtra, error) {
	ms := []PeopleExtra{}
	// set limit
	offset, limit, sortBy, direction, search := getPaginationParams(r)

	result := db.db.Offset(offset).Limit(limit).Order(sortBy + " " + direction).Raw(
		`SELECT 
		json_build_object('owner_pubkey', owner_pub_key, 'owner_alias', owner_alias, 'img', img, 'unique_name', unique_name, 'id', id, 'wanted', extras->'wanted', 'github_issues', github_issues) #>> '{}' as person,
		to_json(jsonb_array_elements(extras->'wanted'::text)) #>> '{}' as body 
		FROM people
		WHERE (unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)
		AND extras IS NOT NULL
		AND extras != '{}'::jsonb
		AND extras->'wanted' IS NOT NULL
		AND extras->'wanted' != '[]'::jsonb
		AND EXISTS (SELECT *
			FROM jsonb_array_elements_Text(extras -> 'wanted') as x(title)
			WHERE LOWER(x.title) LIKE '%` + search + `%')`).Find(&ms)

	return ms, result.Error
}

func (db database) getListedOffers(r *http.Request) ([]PeopleExtra, error) {
	ms := []PeopleExtra{}
	// set limit
	offset, limit, sortBy, direction, search := getPaginationParams(r)

	result := db.db.Offset(offset).Limit(limit).Order(sortBy + " " + direction).Raw(
		`SELECT 
		json_build_object('owner_pubkey', owner_pub_key, 'owner_alias', owner_alias, 'img', img, 'unique_name', unique_name, 'id', id, 'offer', extras->'offer') #>> '{}' as person,
		to_json(jsonb_array_elements(extras->'offer'::text)) #>> '{}' as body 
		FROM people
		WHERE (unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)
		AND extras IS NOT NULL
		AND extras != '{}'::jsonb
		AND extras->'offer' IS NOT NULL
		AND extras->'offer' != '[]'::jsonb
		AND EXISTS (SELECT *
			FROM jsonb_array_elements_Text(extras -> 'offer') as x(title)
			WHERE LOWER(x.title) LIKE '%` + search + `%')`).Find(&ms)

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
	db.db.Where("(deleted = 'f' OR deleted is null)").Find(&ms)
	return ms
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
		FROM tribes, to_tsquery('` + s + `') q
		WHERE tsv @@ q
		AND (deleted = 'f' OR deleted is null)
		ORDER BY rank DESC LIMIT 100;`).Find(&ms)
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
		FROM bots, to_tsquery('` + s + `') q
		WHERE tsv @@ q
		AND (deleted = 'f' OR deleted is null)
		ORDER BY rank DESC 
		LIMIT ` + limitStr + ` OFFSET ` + offsetStr + `;`).Find(&ms)
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
		FROM people, to_tsquery('` + s + `') q
		WHERE tsv @@ q
		AND (deleted = 'f' OR deleted is null)
		ORDER BY rank DESC 
		LIMIT ` + limitStr + ` OFFSET ` + offsetStr + `;`).Find(&ms)
	return ms
}
