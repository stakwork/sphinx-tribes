//go:build !mock
// +build !mock

package db

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
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
)

// check that update owner_pub_key does in fact throw error
func (db database) CreateOrEditTribe(m Tribe) (Tribe, error) {
	if m.OwnerPubKey == "" {
		return Tribe{}, errors.New("no pub key")
	}
	onConflict := "ON CONFLICT (uuid) DO UPDATE SET"
	for i, u := range Updatables {
		onConflict = onConflict + fmt.Sprintf(" %s=EXCLUDED.%s", u, u)
		if i < len(Updatables)-1 {
			onConflict = onConflict + ","
		}
	}
	if m.Name == "" {
		m.Name = "name"
	}
	if m.Description == "" {
		m.Description = ""
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

// check that update owner_pub_key does in fact throws an error
func (db database) CreateOrEditBot(b Bot) (Bot, error) {
	if b.OwnerPubKey == "" {
		return Bot{}, errors.New("no pub key")
	}
	if b.UniqueName == "" {
		return Bot{}, errors.New("no unique name")
	}
	onConflict := "ON CONFLICT (uuid) DO UPDATE SET"
	for i, u := range Botupdatables {
		onConflict = onConflict + fmt.Sprintf(" %s=EXCLUDED.%s", u, u)
		if i < len(Botupdatables)-1 {
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

func (db database) DeleteBot() (bool, error) {
	result := db.db.Exec("DELETE FROM bots")
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// check that update owner_pub_key does in fact throws an error
func (db database) CreateOrEditPerson(m Person) (Person, error) {
	if m.OwnerPubKey == "" {
		return Person{}, errors.New("no pub key")
	}
	onConflict := "ON CONFLICT (id) DO UPDATE SET"
	for i, u := range Peopleupdatables {
		onConflict = onConflict + fmt.Sprintf(" %s=EXCLUDED.%s", u, u)
		if i < len(Peopleupdatables)-1 {
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
	if m.PriceToMeet == 0 {
		updatePriceToMeet := make(map[string]interface{})
		updatePriceToMeet["price_to_meet"] = 0

		db.db.Model(&m).Where("id = ?", m.ID).UpdateColumns(&updatePriceToMeet)
	}

	if db.db.Model(&m).Where("owner_pub_key = ?", m.OwnerPubKey).Updates(&m).RowsAffected == 0 {
		db.db.Create(&m)
	}

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

func (db database) GetTribesByAppUrl(aurl string) []Tribe {
	ms := []Tribe{}
	db.db.Where("LOWER(app_url) LIKE ?", "%"+aurl+"%").Find(&ms)
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

	// avoid dereference error, since r can be nil
	var keys url.Values
	if r != nil {
		keys = r.URL.Query()
	}

	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)

	orderQuery := ""
	limitQuery := ""
	searchQuery := ""

	languageQuery := ""

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = "ORDER BY " + sortBy + "" + "DESC"
	}
	if limit > -1 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}
	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(owner_alias) LIKE %[1]s OR LOWER(unique_name) LIKE %[1]s", "'%"+strings.ToLower(search)+"%'")
	}

	if languageLength > 0 {
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					languageQuery = "AND extras->'coding_languages' @> '[{\"label\": \"" + val + "\"}]'"
				} else {
					languageQuery += " OR extras->'coding_languages' @> '[{\"label\": \"" + val + "\"}]'"
				}
			}
		}

	}

	query := "SELECT * FROM people WHERE (unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)"

	allQuery := query + " " + searchQuery + " " + languageQuery + " " + orderQuery + " " + limitQuery

	db.db.Raw(allQuery).Find(&ms)
	return ms
}

func (db database) ListAllPeople(r *http.Request) []Person {
	keys := r.URL.Query()

	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)

	languageQuery := ""

	if languageLength > 0 {
		langs := ""
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					langs = "'" + val + "'"
				} else {
					langs = langs + ", '" + val + "'"
				}
				languageQuery = "AND coding_languages && ARRAY[" + langs + "]"
			}
		}
	}

	ms := []Person{}

	query := "SELECT * from people WHERE (unlisted = 'f' OR unlisted is null AND (deleted = 'f' OR deleted is null)"
	allQuery := query + languageQuery
	db.db.Raw(allQuery).Find(&ms)
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

	// return if like owner_alias, unique_name, or equals pubkey AND owner_pub_key contains "_" (V2 pubkey)
	db.db.Offset(offset).Limit(limit).Order(sortBy+" "+direction+" NULLS LAST").
		Where("(unlisted = 'f' OR unlisted is null) AND (deleted = 'f' OR deleted is null)").
		Where("owner_route_hint LIKE ? AND owner_route_hint NOT LIKE ?", "%_%", "%:%").
		Where("(LOWER(owner_alias) LIKE ? OR LOWER(unique_name) LIKE ? OR LOWER(owner_pub_key) = ?)",
			"%"+search+"%", "%"+search+"%", search).
		Find(&ms)
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

func (db database) GetUserBountiesCount(personKey string, tabType string) int64 {
	var count int64

	query := db.db.Model(&NewBounty{})
	if tabType == "bounties" {
		query.Where("owner_id", personKey)
	} else if tabType == "assigned" {
		query.Where("assignee", personKey)
	}

	query.Count(&count)
	return count
}

func (db database) GetBountiesCount(r *http.Request) int64 {
	keys := r.URL.Query()
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")
	completed := keys.Get("Completed")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")

	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if completed == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND completed = true AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	var statusQuery string
	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	var count int64

	query := "SELECT COUNT(*) FROM bounty WHERE show != false"
	allQuery := query + " " + statusQuery
	db.db.Raw(allQuery).Scan(&count)
	return count
}

func (db database) GetFilterStatusCount() FilterStattuCount {
	var openCount int64
	var assignedCount int64
	var completedCount int64
	var paidCount int64
	var pendingCount int64
	var failedCount int64

	db.db.Model(&Bounty{}).Where("show != false").Where("assignee = ''").Where("paid != true").Count(&openCount)
	db.db.Model(&Bounty{}).Where("show != false").Where("assignee != ''").Where("paid != true").Count(&assignedCount)
	db.db.Model(&Bounty{}).Where("show != false").Where("assignee != ''").Where("completed = true").Where("paid != true").Count(&completedCount)
	db.db.Model(&Bounty{}).Where("show != false").Where("assignee != ''").Where("paid = true").Count(&paidCount)
	db.db.Model(&Bounty{}).Where("show != false").Where("assignee != ''").Where("payment_pending = true").Count(&pendingCount)
	db.db.Model(&Bounty{}).Where("show != false").Where("assignee != ''").Where("payment_failed = true").Count(&failedCount)

	ms := FilterStattuCount{
		Open:      openCount,
		Assigned:  assignedCount,
		Completed: completedCount,
		Paid:      paidCount,
		Pending:   pendingCount,
		Failed:    failedCount,
	}

	return ms
}

func (db database) GetWorkspaceBounties(r *http.Request, workspace_uuid string) []NewBounty {
	keys := r.URL.Query()
	tags := keys.Get("tags") // this is a string of tags separated by commas
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	completed := keys.Get("Completed")
	paid := keys.Get("Paid")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")
	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)

	ms := []NewBounty{}

	orderQuery := ""
	limitQuery := ""
	searchQuery := ""
	languageQuery := ""

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = " ORDER BY " + sortBy + "" + "DESC"
	}
	if limit > 0 {
		limitQuery = fmt.Sprintf("LIMIT %d", limit)
	}
	if offset > 0 {
		limitQuery += fmt.Sprintf(" OFFSET %d", offset)
	}
	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(title) LIKE %s", "'%"+strings.ToLower(search)+"%'")
	}

	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if completed == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND completed = true AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	var statusQuery string
	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if languageLength > 0 {
		langs := ""
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					langs = "'" + val + "'"
				} else {
					langs = langs + ", '" + val + "'"
				}
				languageQuery = "AND coding_languages && ARRAY[" + langs + "]"
			}
		}
	}

	query := `SELECT * FROM bounty WHERE workspace_uuid = '` + workspace_uuid + `'`
	allQuery := query + " " + statusQuery + " " + searchQuery + " " + languageQuery + " " + orderQuery + " " + limitQuery
	theQuery := db.db.Raw(allQuery)

	if tags != "" {
		// pull out the tags and add them in here
		t := strings.Split(tags, ",")
		for _, s := range t {
			theQuery = theQuery.Where("'" + s + "'" + " = any (tags)")
		}
	}

	theQuery.Scan(&ms)

	return ms
}

func (db database) GetWorkspaceBountiesCount(r *http.Request, workspace_uuid string) int64 {
	keys := r.URL.Query()
	tags := keys.Get("tags") // this is a string of tags separated by commas
	search := keys.Get("search")
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	completed := keys.Get("Completed")
	paid := keys.Get("Paid")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")
	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)

	searchQuery := ""
	languageQuery := ""

	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(title) LIKE %s", "'%"+strings.ToLower(search)+"%'")
	}

	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if completed == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND completed = true AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	var statusQuery string
	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if languageLength > 0 {
		langs := ""
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					langs = "'" + val + "'"
				} else {
					langs = langs + ", '" + val + "'"
				}
				languageQuery = "AND coding_languages && ARRAY[" + langs + "]"
			}
		}
	}

	var count int64

	query := `SELECT COUNT(*) FROM bounty WHERE workspace_uuid = '` + workspace_uuid + `'`
	allQuery := query + " " + statusQuery + " " + searchQuery + " " + languageQuery
	theQuery := db.db.Raw(allQuery)

	if tags != "" {
		// pull out the tags and add them in here
		t := strings.Split(tags, ",")
		for _, s := range t {
			theQuery = theQuery.Where("'" + s + "'" + " = any (tags)")
		}
	}

	theQuery.Scan(&count)

	return count
}

func (db database) GetAssignedBounties(r *http.Request) ([]NewBounty, error) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		return []NewBounty{}, nil
	}

	offset, limit, sortBy, direction, _ := utils.GetPaginationParams(r)
	person := db.GetPersonByUuid(uuid)
	pubkey := person.OwnerPubKey
	keys := r.URL.Query()

	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")

	orderQuery := ""
	limitQuery := ""
	var statusQuery string
	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = " ORDER BY created DESC"
	}
	if offset >= 0 && limit > 1 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}

	ms := []NewBounty{}

	query := `SELECT * FROM public.bounty WHERE assignee = '` + pubkey + `' AND show != false`
	allQuery := query + " " + statusQuery + " " + orderQuery + " " + limitQuery
	err := db.db.Raw(allQuery).Find(&ms).Error
	return ms, err
}

func (db database) GetCreatedBounties(r *http.Request) ([]NewBounty, error) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		return []NewBounty{}, nil
	}

	offset, limit, sortBy, direction, _ := utils.GetPaginationParams(r)
	person := db.GetPersonByUuid(uuid)

	fmt.Println("person", person)
	pubkey := person.OwnerPubKey
	keys := r.URL.Query()

	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")

	orderQuery := ""
	limitQuery := ""
	var statusQuery string
	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = "ORDER BY created DESC"
	}

	if offset >= 0 && limit > 1 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}

	ms := []NewBounty{}

	query := `SELECT * FROM public.bounty WHERE owner_id = '` + pubkey + `'`
	allQuery := query + " " + statusQuery + " " + orderQuery + " " + limitQuery

	err := db.db.Raw(allQuery).Find(&ms).Error

	return ms, err
}

func (db database) GetBountyById(id string) ([]NewBounty, error) {
	ms := []NewBounty{}
	err := db.db.Raw(`SELECT * FROM public.bounty WHERE id = '` + id + `'`).Find(&ms).Error
	return ms, err
}

func (db database) GetNextBountyByCreated(r *http.Request) (uint, error) {
	created := chi.URLParam(r, "created")
	keys := r.URL.Query()
	_, _, _, _, search := utils.GetPaginationParams(r)
	var bountyId uint

	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")
	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)

	var languageQuery string
	var statusQuery string
	var searchQuery string
	var statusConditions []string

	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(title) LIKE %s", "'%"+strings.ToLower(search)+"%'")
	}

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if languageLength > 0 {
		langs := ""
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					langs = "'" + val + "'"
				} else {
					langs = langs + ", '" + val + "'"
				}
				languageQuery = "AND coding_languages && ARRAY[" + langs + "]"
			}
		}
	}

	query := `SELECT id FROM public.bounty WHERE created > '` + created + `' AND show = true`
	orderQuery := "ORDER BY created ASC LIMIT 1"

	allQuery := query + " " + searchQuery + " " + statusQuery + " " + languageQuery + " " + orderQuery

	err := db.db.Raw(allQuery).Find(&bountyId).Error
	return bountyId, err
}

func (db database) GetPreviousBountyByCreated(r *http.Request) (uint, error) {
	created := chi.URLParam(r, "created")
	keys := r.URL.Query()
	var bountyId uint
	_, _, _, _, search := utils.GetPaginationParams(r)

	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	completed := keys.Get("Completed")
	paid := keys.Get("Paid")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")
	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)

	var languageQuery string
	var statusQuery string
	var searchQuery string
	var statusConditions []string

	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(title) LIKE %s", "'%"+strings.ToLower(search)+"%'")
	}

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if completed == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND completed = true AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if languageLength > 0 {
		langs := ""
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					langs = "'" + val + "'"
				} else {
					langs = langs + ", '" + val + "'"
				}
				languageQuery = "AND coding_languages && ARRAY[" + langs + "]"
			}
		}
	}

	query := `SELECT id FROM public.bounty WHERE created < '` + created + `' AND show = true`
	orderQuery := "ORDER BY created DESC LIMIT 1"

	allQuery := query + " " + searchQuery + " " + statusQuery + " " + languageQuery + " " + orderQuery

	err := db.db.Raw(allQuery).Find(&bountyId).Error
	return bountyId, err
}

func (db database) GetNextWorkspaceBountyByCreated(r *http.Request) (uint, error) {
	created := chi.URLParam(r, "created")
	uuid := chi.URLParam(r, "uuid")
	keys := r.URL.Query()
	_, _, _, _, search := utils.GetPaginationParams(r)
	var bountyId uint

	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	completed := keys.Get("Completed")
	paid := keys.Get("Paid")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")
	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)

	var languageQuery string
	var statusQuery string
	var searchQuery string
	var statusConditions []string

	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(title) LIKE %s", "'%"+strings.ToLower(search)+"%'")
	}

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if completed == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND completed = true AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if languageLength > 0 {
		langs := ""
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					langs = "'" + val + "'"
				} else {
					langs = langs + ", '" + val + "'"
				}
				languageQuery = "AND coding_languages && ARRAY[" + langs + "]"
			}
		}
	}

	query := `SELECT id FROM public.bounty WHERE workspace_uuid = '` + uuid + `' AND created > '` + created + `' AND show = true`
	orderQuery := "ORDER BY created ASC LIMIT 1"

	allQuery := query + " " + searchQuery + " " + statusQuery + " " + languageQuery + " " + orderQuery

	err := db.db.Raw(allQuery).Find(&bountyId).Error
	return bountyId, err
}

func (db database) GetPreviousWorkspaceBountyByCreated(r *http.Request) (uint, error) {
	created := chi.URLParam(r, "created")
	uuid := chi.URLParam(r, "uuid")
	keys := r.URL.Query()
	_, _, _, _, search := utils.GetPaginationParams(r)
	var bountyId uint

	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	completed := keys.Get("Completed")
	paid := keys.Get("Paid")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")
	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)

	var languageQuery string
	var statusQuery string
	var searchQuery string
	var statusConditions []string

	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(title) LIKE %s", "'%"+strings.ToLower(search)+"%'")
	}

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if completed == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND completed = true AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if languageLength > 0 {
		langs := ""
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					langs = "'" + val + "'"
				} else {
					langs = langs + ", '" + val + "'"
				}
				languageQuery = "AND coding_languages && ARRAY[" + langs + "]"
			}
		}
	}

	query := `SELECT id FROM public.bounty WHERE workspace_uuid = '` + uuid + `' AND created < '` + created + `' AND show = true`
	orderQuery := "ORDER BY created DESC LIMIT 1"

	allQuery := query + " " + searchQuery + " " + statusQuery + " " + languageQuery + " " + orderQuery

	err := db.db.Raw(allQuery).Find(&bountyId).Error
	return bountyId, err
}

func (db database) GetBountyIndexById(id string) int64 {
	var index int64
	db.db.Raw(`SELECT position FROM(SELECT *, row_number() over( ORDER BY id DESC) as position FROM public.bounty) result WHERE id = '` + id + `' OR created = '` + id + `'`).Scan(&index)
	return index
}

func (db database) GetBountyDataByCreated(created string) ([]NewBounty, error) {
	ms := []NewBounty{}
	err := db.db.Raw(`SELECT * FROM public.bounty WHERE created = '` + created + `'`).Find(&ms).Error
	return ms, err
}

func (db database) AddBounty(b Bounty) (Bounty, error) {
	db.db.Create(&b)
	return b, nil
}

func (db database) GetAllBounties(r *http.Request) []NewBounty {
	keys := r.URL.Query()
	tags := keys.Get("tags") // this is a string of tags separated by commas
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	completed := keys.Get("Completed")
	paid := keys.Get("Paid")
	pending := keys.Get("Pending")
	failed := keys.Get("Failed")
	orgUuid := keys.Get("org_uuid")
	workspaceUuid := keys.Get("workspace_uuid")
	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)
	PhaseUuid := keys.Get("phase_uuid")
	PhasePriority := keys.Get("phase_priority")

	if workspaceUuid == "" && orgUuid != "" {
		workspaceUuid = orgUuid
	}

	ms := []NewBounty{}

	orderQuery := ""
	limitQuery := ""
	searchQuery := ""
	workspaceQuery := ""
	languageQuery := ""
	phaseUuidQuery := ""
	phasePriorityQuery := ""

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = "ORDER BY " + sortBy + "" + "DESC"
	}
	if limit != 0 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}
	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(title) LIKE %s", "'%"+strings.ToLower(search)+"%'")
	}

	if PhaseUuid != "" {
		phaseUuidQuery = "AND phase_uuid = '" + PhaseUuid + "'"
	}

	if PhasePriority != "" {
		phasePriorityQuery = "AND phase_priority = '" + PhasePriority + "'"
	}

	var statusConditions []string

	if open == "true" {
		statusConditions = append(statusConditions, "assignee = '' AND paid != true")
	}
	if assingned == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND paid = false")
	}
	if completed == "true" {
		statusConditions = append(statusConditions, "assignee != '' AND completed = true AND paid = false")
	}
	if paid == "true" {
		statusConditions = append(statusConditions, "paid = true")
	}
	if pending == "true" {
		statusConditions = append(statusConditions, "payment_pending = true")
	}
	if failed == "true" {
		statusConditions = append(statusConditions, "payment_failed = true")
	}

	var statusQuery string
	if len(statusConditions) > 0 {
		statusQuery = " AND (" + strings.Join(statusConditions, " OR ") + ")"
	} else {
		statusQuery = ""
	}

	if workspaceUuid != "" {
		workspaceQuery = "AND workspace_uuid = '" + workspaceUuid + "'"
	}
	if languageLength > 0 {
		langs := ""
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					langs = "'" + val + "'"
				} else {
					langs = langs + ", '" + val + "'"
				}
				languageQuery = "AND coding_languages && ARRAY[" + langs + "]"
			}
		}
	}

	query := "SELECT * FROM public.bounty WHERE show != false"

	allQuery := query + " " + statusQuery + " " + searchQuery + " " + workspaceQuery + " " + languageQuery + " " + phaseUuidQuery + " " + phasePriorityQuery + " " + orderQuery + " " + limitQuery

	theQuery := db.db.Raw(allQuery)

	if tags != "" {
		// pull out the tags and add them in here
		t := strings.Split(tags, ",")
		for _, s := range t {
			theQuery = theQuery.Where("'" + s + "'" + " = any (tags)")
		}
	}

	theQuery.Scan(&ms)

	return ms
}

func (db database) CreateOrEditBounty(b NewBounty) (NewBounty, error) {
	if b.OwnerID == "" {
		return NewBounty{}, errors.New("no pub key")
	}

	if db.db.Model(&b).Where("id = ? OR owner_id = ? AND created = ?", b.ID, b.OwnerID, b.Created).Updates(&b).RowsAffected == 0 {
		db.db.Create(&b)
	}
	return b, nil
}

func (db database) UpdateBountyNullColumn(b NewBounty, column string) NewBounty {
	columnMap := make(map[string]interface{})
	columnMap[column] = ""
	db.db.Model(&b).Where("created = ?", b.Created).UpdateColumns(&columnMap)
	return b
}

func (db database) UpdateBountyBoolColumn(b NewBounty, column string) NewBounty {
	columnMap := make(map[string]interface{})
	columnMap[column] = false
	db.db.Model(&b).Select(column).UpdateColumns(columnMap)
	return b
}

func (db database) DeleteBounty(pubkey string, created string) (NewBounty, error) {
	m := NewBounty{}
	db.db.Where("owner_id", pubkey).Where("created", created).Delete(&m)
	return m, nil
}

func (db database) GetBountyByCreated(created uint) (NewBounty, error) {
	b := NewBounty{}
	err := db.db.Where("created", created).Find(&b).Error
	return b, err
}

func (db database) GetBounty(id uint) NewBounty {
	b := NewBounty{}
	db.db.Where("id", id).Find(&b)
	return b
}

func (db database) UpdateBountyPaymentStatuses(bounty NewBounty) (NewBounty, error) {

	bountyUpdates := map[string]interface{}{
		"paid":            bounty.Paid,
		"payment_pending": bounty.PaymentPending,
		"payment_failed":  bounty.PaymentFailed,
		"completed":       bounty.Completed,
		"paid_date":       bounty.PaidDate,
		"completion_date": bounty.CompletionDate,
	}

	db.db.Model(&NewBounty{}).Where("created", bounty.Created).Updates(bountyUpdates)
	return bounty, nil
}

func (db database) UpdateBounty(b NewBounty) (NewBounty, error) {
	db.db.Where("created", b.Created).Updates(&b)
	return b, nil
}

func (db database) UpdateBountyPayment(b NewBounty) (NewBounty, error) {
	db.db.Model(&b).Where("created", b.Created).Updates(map[string]interface{}{
		"paid": b.Paid,
	})
	db.db.Model(&b).Where("created", b.Created).Updates(b)
	return b, nil
}

func (db database) UpdateBountyCompleted(b NewBounty) (NewBounty, error) {
	db.db.Model(&b).Where("created", b.Created).Updates(map[string]interface{}{
		"completed": b.Completed,
	})
	db.db.Model(&b).Where("created", b.Created).Updates(b)
	return b, nil
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
	db.db.Where("(deleted = 'f' OR deleted is null)").Find(&ms)
	return ms
}

func (db database) DeleteTribe() (bool, error) {
	result := db.db.Exec("DELETE FROM tribes")
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
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
	db.db.Where("owner_pub_key = ? AND (deleted = false OR deleted is null)", pubkey).Find(&m)

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
	s = strings.ReplaceAll(s, " ", " & ")
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
	var count uint64
	db.db.Raw(`Select COUNT(*) from bounty`).Scan(&count)
	return count
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

func (db database) CreateConnectionCode(c []ConnectionCodes) ([]ConnectionCodes, error) {
	if len(c) == 0 {
		return nil, fmt.Errorf("no connection codes provided")
	}
	now := time.Now()
	for _, code := range c {
		if code.DateCreated == nil || code.DateCreated.IsZero() {
			code.DateCreated = &now
		}
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
		p.UniqueName, _ = db.PersonUniqueNameFromName(p.OwnerAlias)
		p.Created = &now
		p.Tags = pq.StringArray{}
		p.Uuid = xid.New().String()
		p.Extras = make(map[string]interface{})
		p.GithubIssues = make(map[string]interface{})

		db.db.Create(&p)
	}
	return p, nil
}

func (db database) PersonUniqueNameFromName(name string) (string, error) {
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
		existing := db.GetPersonByUniqueName(uniquepath)
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
	ms := []BountyLeaderboard{}
	var users = []LeaderData{}

	db.db.Raw(`SELECT t1.owner_pubkey, total_bounties_completed, total_sats_earned FROM
(SELECT assignee as owner_pubkey, 
COUNT(assignee) as total_bounties_completed
From bounty 
where paid=true and assignee != '' 
GROUP BY assignee) t1
 Right Join
(SELECT assignee as owner_pubkey,  
SUM(CAST(price as integer)) as total_sats_earned
From bounty
where paid=true and assignee != ''
GROUP BY assignee) t2
ON t1.owner_pubkey = t2.owner_pubkey
ORDER by total_sats_earned DESC`).Find(&ms)

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

func (db database) GetInvoice(payment_request string) NewInvoiceList {
	ms := NewInvoiceList{}
	db.db.Where("payment_request = ?", payment_request).Find(&ms)
	return ms
}

func (db database) UpdateInvoice(payment_request string) NewInvoiceList {
	ms := NewInvoiceList{}
	db.db.Model(&NewInvoiceList{}).Where("payment_request = ?", payment_request).Update("status", true)
	ms.Status = true
	return ms
}

func (db database) AddInvoice(invoice NewInvoiceList) NewInvoiceList {
	db.db.Create(&invoice)
	return invoice
}

func (db database) DeleteInvoice(payment_request string) NewInvoiceList {
	ms := NewInvoiceList{}
	db.db.Model(&NewInvoiceList{}).Where("payment_request = ?", payment_request).Delete(&ms)
	return ms
}

func (db database) AddUserInvoiceData(userData UserInvoiceData) UserInvoiceData {
	db.db.Create(&userData)
	return userData
}

func (db database) ProcessAddInvoice(invoice NewInvoiceList, userData UserInvoiceData) error {
	tx := db.db.Begin()
	var err error

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return err
	}

	if invoice.WorkspaceUuid == "" {
		tx.Rollback()
		return errors.New("cannot create invoice")
	}

	if err = tx.Create(&invoice).Error; err != nil {
		tx.Rollback()
	}

	if err = tx.Create(&userData).Error; err != nil {
		tx.Rollback()
	}

	return tx.Commit().Error
}

func (db database) ProcessBudgetInvoice(paymentHistory NewPaymentHistory, newInvoice NewInvoiceList) error {
	tx := db.db.Begin()
	var err error

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return err
	}

	if paymentHistory.WorkspaceUuid == "" {
		tx.Rollback()
		return errors.New("cannot create invoice")
	}

	if err = tx.Create(&paymentHistory).Error; err != nil {
		tx.Rollback()
	}

	if err = tx.Create(&newInvoice).Error; err != nil {
		tx.Rollback()
	}

	return tx.Commit().Error
}

func (db database) GetUserInvoiceData(payment_request string) UserInvoiceData {
	ms := UserInvoiceData{}
	db.db.Where("payment_request = ?", payment_request).Find(&ms)
	return ms
}

func (db database) DeleteUserInvoiceData(payment_request string) UserInvoiceData {
	ms := UserInvoiceData{}
	db.db.Where("payment_request = ?", payment_request).Delete(&ms)
	return ms
}

func (db database) DeleteAllBounties() error {
	if err := db.db.Exec("DELETE FROM bounty").Error; err != nil {
		return err
	}
	return nil
}

func (db database) GetProofsByBountyID(bountyID uint) []ProofOfWork {
	var proofs []ProofOfWork
	db.db.Where("bounty_id = ?", bountyID).Find(&proofs)
	return proofs
}

func (db database) CreateProof(proof ProofOfWork) error {
	return db.db.Create(&proof).Error
}

func (db database) DeleteProof(proofID string) error {
	return db.db.Delete(&ProofOfWork{}, "id = ?", proofID).Error
}

func (db database) UpdateProofStatus(proofID string, status ProofOfWorkStatus) error {
	return db.db.Model(&ProofOfWork{}).Where("id = ?", proofID).Update("status", status).Error
}

func (db database) IncrementProofCount(bountyID uint) error { // Ensure bountyID is of type uint
	var bounty NewBounty

	if err := db.db.Where("id = ?", bountyID).First(&bounty).Error; err != nil {
		return err
	}

	return db.db.Model(&bounty).
		Updates(map[string]interface{}{
			"proof_of_work_count": bounty.ProofOfWorkCount + 1,
			"updated":             time.Now(),
		}).Error
}
func (db database) DecrementProofCount(bountyID uint) error {
	var bounty NewBounty
	if err := db.db.Where("id = ?", bountyID).First(&bounty).Error; err != nil {
		return err
	}

	newCount := int(math.Max(0, float64(bounty.ProofOfWorkCount-1)))

	return db.db.Model(&bounty).
		Updates(map[string]interface{}{
			"proof_of_work_count": newCount,
			"updated":             time.Now(),
		}).Error
}

func (db database) CreateBountyTiming(bountyID uint) (*BountyTiming, error) {
	timing := &BountyTiming{
		BountyID: bountyID,
	}
	err := db.db.Create(timing).Error
	return timing, err
}

func (db database) GetBountyTiming(bountyID uint) (*BountyTiming, error) {
	var timing BountyTiming
	err := db.db.Where("bounty_id = ?", bountyID).First(&timing).Error
	if err != nil {
		return nil, err
	}
	return &timing, nil
}

func (db database) UpdateBountyTiming(timing *BountyTiming) error {
	return db.db.Save(timing).Error
}

func (db database) StartBountyTiming(bountyID uint) error {
	now := time.Now()
	timing, err := db.GetBountyTiming(bountyID)
	if err != nil {

		timing, err = db.CreateBountyTiming(bountyID)
		if err != nil {
			return fmt.Errorf("failed to create bounty timing: %w", err)
		}
	}

	if timing.FirstAssignedAt == nil {
		timing.FirstAssignedAt = &now
		if err := db.UpdateBountyTiming(timing); err != nil {
			return fmt.Errorf("failed to update bounty timing: %w", err)
		}
	}

	return nil
}

func (db database) CloseBountyTiming(bountyID uint) error {
	timing, err := db.GetBountyTiming(bountyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve bounty timing: %w", err)
	}

	now := time.Now()
	timing.ClosedAt = &now

	if timing.FirstAssignedAt != nil {
		timing.TotalDurationSeconds = int(now.Sub(*timing.FirstAssignedAt).Seconds())
	}

	if err := db.UpdateBountyTiming(timing); err != nil {
		return fmt.Errorf("failed to close bounty timing: %w", err)
	}

	return nil
}

func (db database) UpdateBountyTimingOnProof(bountyID uint) error {
	timing, err := db.GetBountyTiming(bountyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve bounty timing: %w", err)
	}

	now := time.Now()

	if timing.LastPoWAt != nil {
		workTime := int(now.Sub(*timing.LastPoWAt).Seconds())
		timing.TotalWorkTimeSeconds += workTime
	}

	timing.LastPoWAt = &now
	timing.TotalAttempts++

	if err := db.UpdateBountyTiming(timing); err != nil {
		return fmt.Errorf("failed to update bounty timing: %w", err)
	}

	return nil
}

func (db database) GetWorkspaceBountyCardsData(r *http.Request) []NewBounty {
	keys := r.URL.Query()
	_, _, sortBy, direction, search := utils.GetPaginationParams(r)
	workspaceUuid := keys.Get("workspace_uuid")

	orderQuery := ""
	searchQuery := ""
	workspaceQuery := ""
	timeFilterQuery := ""

	timeFilterQuery = `
		AND (
			(NOT paid AND EXTRACT(EPOCH FROM updated::timestamp) > EXTRACT(EPOCH FROM (NOW() - INTERVAL '4 weeks')))
			OR (paid AND EXTRACT(EPOCH FROM updated::timestamp) > EXTRACT(EPOCH FROM (NOW() - INTERVAL '2 weeks')))
			OR updated IS NULL  -- Preserve existing records without updated timestamp
		)`

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = "ORDER BY created DESC"
	}

	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(title) LIKE %s", "'%"+strings.ToLower(search)+"%'")
	}

	if workspaceUuid != "" {
		workspaceQuery = "WHERE workspace_uuid = '" + workspaceUuid + "'"
	}

	query := "SELECT * FROM public.bounty"
	allQuery := query + " " + workspaceQuery + timeFilterQuery + " " + searchQuery + " " + orderQuery

	ms := []NewBounty{}
	db.db.Raw(allQuery).Scan(&ms)

	return ms
}

func (db database) GetFeaturedBountyById(id string) (FeaturedBounty, error) {
	var bounty FeaturedBounty
	err := db.db.Where("bounty_id = ?", id).First(&bounty).Error
	return bounty, err
}

func (db database) GetAllFeaturedBounties() ([]FeaturedBounty, error) {
	var bounties []FeaturedBounty
	err := db.db.Order("added_at DESC").Find(&bounties).Error
	return bounties, err
}

func (db database) CreateFeaturedBounty(bounty FeaturedBounty) error {
	return db.db.Create(&bounty).Error
}

func (db database) UpdateFeaturedBounty(bountyID string, bounty FeaturedBounty) error {
	return db.db.Model(&FeaturedBounty{}).Where("bounty_id = ?", bountyID).Updates(bounty).Error
}

func (db database) DeleteFeaturedBounty(bountyID string) error {
	return db.db.Where("bounty_id = ?", bountyID).Delete(&FeaturedBounty{}).Error
}
