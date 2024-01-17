//go:build !mock
// +build !mock

package db

import (
	"errors"
	"fmt"
	"net/http"
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

func (db database) GetUserBountiesCount(personKey string, tabType string) int64 {
	var count int64

	query := db.db.Model(&Bounty{})
	if tabType == "wanted" {
		query.Where("owner_id", personKey)
	} else if tabType == "usertickets" {
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
	openQuery := ""
	assignedQuery := ""
	paidQuery := ""

	if open != "" && open == "true" {
		openQuery = "AND assignee = '' AND paid != true"
		assignedQuery = ""
	}
	if assingned != "" && assingned == "true" {
		if open != "" && open == "true" {
			assignedQuery = "OR assignee != '' AND paid = false"
		} else {
			assignedQuery = "AND assignee != '' AND paid = false"
		}
	}
	if paid != "" && paid == "true" {
		if open != "" && open == "true" || assingned != "" && assingned == "true" {
			paidQuery = "OR paid = true"
		} else if open != "" && open == "true" && assingned == "" && assingned != "true" {
			assignedQuery = ""
		} else {
			paidQuery = "AND paid = true"
		}
	}

	var count int64

	query := "SELECT COUNT(*) FROM bounty WHERE show != false"
	allQuery := query + " " + openQuery + " " + assignedQuery + " " + paidQuery
	db.db.Raw(allQuery).Scan(&count)
	return count
}

func (db database) GetFilterStatusCount() FilterStattuCount {
	var openCount int64
	var assignedCount int64
	var paidCount int64

	db.db.Model(&Bounty{}).Where("show != false").Where("assignee = ''").Where("paid != true").Count(&openCount)
	db.db.Model(&Bounty{}).Where("show != false").Where("assignee != ''").Where("paid != true").Count(&assignedCount)
	db.db.Model(&Bounty{}).Where("show != false").Where("assignee != ''").Where("paid = true").Count(&paidCount)

	ms := FilterStattuCount{
		Open:     openCount,
		Assigned: assignedCount,
		Paid:     paidCount,
	}

	return ms
}

func (db database) GetOrganizationBounties(r *http.Request, org_uuid string) []Bounty {
	keys := r.URL.Query()
	tags := keys.Get("tags") // this is a string of tags separated by commas
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")
	ms := []Bounty{}

	orderQuery := ""
	limitQuery := ""
	searchQuery := ""
	openQuery := ""
	assignedQuery := ""
	paidQuery := ""

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = " ORDER BY " + sortBy + "" + "DESC"
	}
	if offset != 0 && limit != 0 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}
	if search != "" {
		searchQuery = fmt.Sprintf("WHERE LOWER(title) LIKE %s", "'%"+search+"%'")
	}
	if open != "" && open == "true" {
		openQuery = "AND assignee = '' AND paid != true"
		assignedQuery = ""
	}
	if assingned != "" && assingned == "true" {
		if open != "" && open == "true" {
			assignedQuery = "OR assignee != '' AND paid = false"
		} else {
			assignedQuery = "AND assignee != '' AND paid = false"
		}
	}
	if paid != "" && paid == "true" {
		if open != "" && open == "true" || assingned != "" && assingned == "true" {
			paidQuery = "OR paid = true"
		} else if open != "" && open == "true" && assingned == "" && assingned != "true" {
			assignedQuery = ""
		} else {
			paidQuery = "AND paid = true"
		}
	}

	query := `SELECT * FROM bounty WHERE org_uuid = '` + org_uuid + `'`
	allQuery := query + " " + openQuery + " " + assignedQuery + " " + paidQuery + " " + searchQuery + " " + orderQuery + " " + limitQuery
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

func (db database) GetAssignedBounties(r *http.Request) ([]Bounty, error) {
	offset, limit, sortBy, direction, _ := utils.GetPaginationParams(r)
	pubkey := chi.URLParam(r, "pubkey")

	orderQuery := ""
	limitQuery := ""

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + "DESC"
	} else {
		orderQuery = " ORDER BY paid DESC"
	}
	if offset != 0 && limit != 0 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}

	ms := []Bounty{}

	query := `SELECT * FROM public.bounty WHERE assignee = '` + pubkey + `' AND show != false`
	allQuery := query + " " + orderQuery + " " + limitQuery
	err := db.db.Raw(allQuery).Find(&ms).Error
	return ms, err
}

func (db database) GetCreatedBounties(r *http.Request) ([]Bounty, error) {
	offset, limit, sortBy, direction, _ := utils.GetPaginationParams(r)
	pubkey := chi.URLParam(r, "pubkey")

	orderQuery := ""
	limitQuery := ""

	fmt.Println("Sort BY", sortBy, limit)

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + "ASC"
	} else {
		orderQuery = "ORDER BY paid DESC"
	}

	if offset != 0 && limit != 0 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}

	ms := []Bounty{}

	query := `SELECT * FROM public.bounty WHERE owner_id = '` + pubkey + `'`
	allQuery := query + " " + orderQuery + " " + limitQuery

	err := db.db.Raw(allQuery).Find(&ms).Error
	return ms, err
}

func (db database) GetBountyById(id string) ([]Bounty, error) {
	ms := []Bounty{}
	err := db.db.Raw(`SELECT * FROM public.bounty WHERE id = '` + id + `'`).Find(&ms).Error
	return ms, err
}

func (db database) GetBountyIndexById(id string) int64 {
	var index int64
	db.db.Raw(`SELECT position FROM(SELECT *, row_number() over( ORDER BY id DESC) as position FROM public.bounty) result WHERE id = '` + id + `' OR created = '` + id + `'`).Scan(&index)
	return index
}

func (db database) GetBountyDataByCreated(created string) ([]Bounty, error) {
	ms := []Bounty{}
	err := db.db.Raw(`SELECT * FROM public.bounty WHERE created = '` + created + `'`).Find(&ms).Error
	return ms, err
}

func (db database) AddBounty(b Bounty) (Bounty, error) {
	db.db.Create(&b)
	return b, nil
}

func (db database) GetAllBounties(r *http.Request) []Bounty {
	keys := r.URL.Query()
	tags := keys.Get("tags") // this is a string of tags separated by commas
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)
	open := keys.Get("Open")
	assingned := keys.Get("Assigned")
	paid := keys.Get("Paid")
	orgUuid := keys.Get("org_uuid")
	languages := keys.Get("languages")
	languageArray := strings.Split(languages, ",")
	languageLength := len(languageArray)

	ms := []Bounty{}

	orderQuery := ""
	limitQuery := ""
	searchQuery := ""
	openQuery := ""
	assignedQuery := ""
	paidQuery := ""
	orgQuery := ""
	languageQuery := ""

	if sortBy != "" && direction != "" {
		orderQuery = "ORDER BY " + sortBy + " " + direction
	} else {
		orderQuery = "ORDER BY " + sortBy + "" + "DESC"
	}
	if limit != 0 {
		limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)
	}
	if search != "" {
		searchQuery = fmt.Sprintf("AND LOWER(title) LIKE %s", "'%"+search+"%'")
	}
	if open != "" && open == "true" {
		openQuery = "AND assignee = '' AND paid != true"
		assignedQuery = ""
	}
	if assingned != "" && assingned == "true" {
		if open != "" && open == "true" {
			assignedQuery = "OR assignee != '' AND paid != true"
		} else {
			assignedQuery = "AND assignee != '' AND paid != true"
		}
	}
	if paid != "" && paid == "true" {
		if open != "" && open == "true" || assingned != "" && assingned == "true" {
			paidQuery = "OR paid = true"
		} else if open != "" && open == "true" && assingned == "" && assingned != "true" {
			assignedQuery = ""
		} else {
			paidQuery = "AND paid = true"
		}
	}
	if orgUuid != "" {
		orgQuery = "AND org_uuid = '" + orgUuid + "'"
	}
	if languageLength > 0 {
		for i, val := range languageArray {
			if val != "" {
				if i == 0 {
					languageQuery = "AND coding_languages && ARRAY['" + val + "']"
				} else {
					query := "OR coding_languages && ARRAY['" + val + "']"
					languageQuery = languageQuery + " " + query
				}
			}
		}
	}

	query := "SELECT * FROM public.bounty WHERE show != false"

	allQuery := query + " " + openQuery + " " + assignedQuery + " " + paidQuery + " " + searchQuery + " " + orgQuery + " " + languageQuery + " " + orderQuery + " " + limitQuery

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

func (db database) CreateOrEditBounty(b Bounty) (Bounty, error) {
	if b.OwnerID == "" {
		return Bounty{}, errors.New("no pub key")
	}

	if db.db.Model(&b).Where("id = ? OR owner_id = ? AND created = ?", b.ID, b.OwnerID, b.Created).Updates(&b).RowsAffected == 0 {
		db.db.Create(&b)
	}
	return b, nil
}

func (db database) UpdateBountyNullColumn(b Bounty, column string) Bounty {
	columnMap := make(map[string]interface{})
	columnMap[column] = ""
	db.db.Model(&b).Where("created = ?", b.Created).UpdateColumns(&columnMap)
	return b
}

func (db database) UpdateBountyBoolColumn(b Bounty, column string) Bounty {
	columnMap := make(map[string]interface{})
	columnMap[column] = false
	db.db.Model(&b).Select(column).UpdateColumns(columnMap)
	return b
}

func (db database) DeleteBounty(pubkey string, created string) (Bounty, error) {
	m := Bounty{}
	db.db.Where("owner_id", pubkey).Where("created", created).Delete(&m)
	return m, nil
}

func (db database) GetBountyByCreated(created uint) (Bounty, error) {
	b := Bounty{}
	err := db.db.Where("created", created).Find(&b).Error
	return b, err
}

func (db database) GetBounty(id uint) Bounty {
	b := Bounty{}
	db.db.Where("id", id).Find(&b)
	return b
}

func (db database) UpdateBounty(b Bounty) (Bounty, error) {
	db.db.Where("created", b.Created).Updates(&b)
	return b, nil
}

func (db database) UpdateBountyPayment(b Bounty) (Bounty, error) {
	db.db.Model(&b).Where("created", b.Created).Updates(map[string]interface{}{
		"paid": b.Paid,
	})
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

func (db database) GetOrganizations(r *http.Request) []Organization {
	ms := []Organization{}
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)

	// return if like owner_alias, unique_name, or equals pubkey
	db.db.Offset(offset).Limit(limit).Order(sortBy+" "+direction+" ").Where("LOWER(name) LIKE ?", "%"+search+"%").Where("deleted != ?", false).Find(&ms)
	return ms
}

func (db database) GetOrganizationsCount() int64 {
	var count int64
	db.db.Model(&Organization{}).Count(&count)
	return count
}

func (db database) GetOrganizationByUuid(uuid string) Organization {
	ms := Organization{}

	db.db.Model(&Organization{}).Where("uuid = ?", uuid).Find(&ms)

	return ms
}

func (db database) GetOrganizationByName(name string) Organization {
	ms := Organization{}

	db.db.Model(&Organization{}).Where("name = ?", name).Find(&ms)

	return ms
}

func (db database) CreateOrEditOrganization(m Organization) (Organization, error) {
	if m.OwnerPubKey == "" {
		return Organization{}, errors.New("no pub key")
	}

	if db.db.Model(&m).Where("uuid = ?", m.Uuid).Updates(&m).RowsAffected == 0 {
		db.db.Create(&m)
	}

	return m, nil
}

func (db database) GetOrganizationUsers(uuid string) ([]OrganizationUsersData, error) {
	ms := []OrganizationUsersData{}

	err := db.db.Raw(`SELECT org.org_uuid, org.created as user_created, person.* FROM public.organization_users AS org LEFT OUTER JOIN public.people AS person ON org.owner_pub_key = person.owner_pub_key WHERE org.org_uuid = '` + uuid + `' ORDER BY org.created DESC`).Find(&ms).Error

	return ms, err
}

func (db database) GetOrganizationUsersCount(uuid string) int64 {
	var count int64
	db.db.Model(&OrganizationUsers{}).Where("org_uuid  = ?", uuid).Count(&count)
	return count
}

func (db database) GetOrganizationBountyCount(uuid string) int64 {
	var count int64
	db.db.Model(&Bounty{}).Where("org_uuid  = ?", uuid).Count(&count)
	return count
}

func (db database) GetOrganizationUser(pubkey string, org_uuid string) OrganizationUsers {
	ms := OrganizationUsers{}
	db.db.Where("org_uuid = ?", org_uuid).Where("owner_pub_key = ?", pubkey).Find(&ms)
	return ms
}

func (db database) CreateOrganizationUser(orgUser OrganizationUsers) OrganizationUsers {
	db.db.Create(&orgUser)

	return orgUser
}

func (db database) DeleteOrganizationUser(orgUser OrganizationUsersData, org string) OrganizationUsersData {
	db.db.Where("owner_pub_key = ?", orgUser.OwnerPubKey).Where("org_uuid = ?", org).Delete(&OrganizationUsers{})
	db.db.Where("owner_pub_key = ?", orgUser.OwnerPubKey).Where("org_uuid = ?", org).Delete(&UserRoles{})
	return orgUser
}

func (db database) GetBountyRoles() []BountyRoles {
	ms := []BountyRoles{}
	db.db.Find(&ms)
	return ms
}

func (db database) CreateUserRoles(roles []UserRoles, uuid string, pubkey string) []UserRoles {
	// delete roles and create new ones
	db.db.Where("org_uuid = ?", uuid).Where("owner_pub_key = ?", pubkey).Delete(&UserRoles{})
	db.db.Create(&roles)

	return roles
}

func (db database) GetUserRoles(uuid string, pubkey string) []UserRoles {
	ms := []UserRoles{}
	db.db.Where("org_uuid = ?", uuid).Where("owner_pub_key = ?", pubkey).Find(&ms)
	return ms
}

func (db database) GetUserCreatedOrganizations(pubkey string) []Organization {
	ms := []Organization{}
	db.db.Where("owner_pub_key = ?", pubkey).Where("deleted != ?", true).Find(&ms)
	return ms
}

func (db database) GetUserAssignedOrganizations(pubkey string) []OrganizationUsers {
	ms := []OrganizationUsers{}
	db.db.Where("owner_pub_key = ?", pubkey).Find(&ms)
	return ms
}

func (db database) AddBudgetHistory(budget BudgetHistory) BudgetHistory {
	db.db.Create(&budget)
	return budget
}

func (db database) CreateOrganizationBudget(budget BountyBudget) BountyBudget {
	db.db.Create(&budget)
	return budget
}

func (db database) UpdateOrganizationBudget(budget BountyBudget) BountyBudget {
	db.db.Model(&BountyBudget{}).Where("org_uuid = ?", budget.OrgUuid).Updates(map[string]interface{}{
		"total_budget": budget.TotalBudget,
	})
	return budget
}

func (db database) GetPaymentHistoryByCreated(created *time.Time, org_uuid string) PaymentHistory {
	ms := PaymentHistory{}
	db.db.Where("created = ?", created).Where("org_uuid = ? ", org_uuid).Find(&ms)
	return ms
}

func (db database) GetOrganizationBudget(org_uuid string) BountyBudget {
	ms := BountyBudget{}
	db.db.Where("org_uuid = ?", org_uuid).Find(&ms)
	return ms
}

func (db database) GetOrganizationBudgetHistory(org_uuid string) []BudgetHistoryData {
	budgetHistory := []BudgetHistoryData{}

	db.db.Raw(`SELECT budget.id, budget.org_uuid, budget.amount, budget.created, budget.updated, budget.payment_type, budget.status, budget.sender_pub_key, sender.unique_name AS sender_name FROM public.budget_histories AS budget LEFT OUTER JOIN public.people AS sender ON budget.sender_pub_key = sender.owner_pub_key WHERE budget.org_uuid = '` + org_uuid + `' ORDER BY budget.created DESC`).Find(&budgetHistory)
	return budgetHistory
}

func (db database) AddAndUpdateBudget(invoice InvoiceList) PaymentHistory {
	created := invoice.Created
	org_uuid := invoice.OrgUuid

	paymentHistory := db.GetPaymentHistoryByCreated(created, org_uuid)

	if paymentHistory.OrgUuid != "" && paymentHistory.Amount != 0 {
		paymentHistory.Status = true
		db.db.Where("created = ?", created).Where("org_uuid = ? ", org_uuid).Updates(paymentHistory)

		// get organization budget and add payment to total budget
		organizationBudget := db.GetOrganizationBudget(org_uuid)

		if organizationBudget.OrgUuid == "" {
			now := time.Now()
			orgBudget := BountyBudget{
				OrgUuid:     org_uuid,
				TotalBudget: paymentHistory.Amount,
				Created:     &now,
				Updated:     &now,
			}
			db.CreateOrganizationBudget(orgBudget)
		} else {
			totalBudget := organizationBudget.TotalBudget
			organizationBudget.TotalBudget = totalBudget + paymentHistory.Amount
			db.UpdateOrganizationBudget(organizationBudget)
		}
	}

	return paymentHistory
}

func (db database) WithdrawBudget(sender_pubkey string, org_uuid string, amount uint) {
	// get organization budget and add payment to total budget
	organizationBudget := db.GetOrganizationBudget(org_uuid)
	totalBudget := organizationBudget.TotalBudget

	newBudget := totalBudget - amount
	db.db.Model(&BountyBudget{}).Where("org_uuid = ?", org_uuid).Updates(map[string]interface{}{
		"total_budget": newBudget,
	})

	now := time.Now()

	budgetHistory := PaymentHistory{
		OrgUuid:        org_uuid,
		Amount:         amount,
		Status:         true,
		PaymentType:    "withdraw",
		Created:        &now,
		Updated:        &now,
		SenderPubKey:   sender_pubkey,
		ReceiverPubKey: "",
		BountyId:       0,
	}
	db.AddPaymentHistory(budgetHistory)
}

func (db database) AddPaymentHistory(payment PaymentHistory) PaymentHistory {
	db.db.Create(&payment)

	// get organization budget and subtract payment from total budget
	organizationBudget := db.GetOrganizationBudget(payment.OrgUuid)
	totalBudget := organizationBudget.TotalBudget

	// deduct amount if it's a bounty payment
	if payment.PaymentType == "payment" {
		organizationBudget.TotalBudget = totalBudget - payment.Amount
	}

	db.UpdateOrganizationBudget(organizationBudget)

	return payment
}

func (db database) GetPaymentHistory(org_uuid string, r *http.Request) []PaymentHistory {
	payment := []PaymentHistory{}

	offset, limit, _, _, _ := utils.GetPaginationParams(r)
	limitQuery := ""

	limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)

	query := `SELECT * FROM payment_histories WHERE org_uuid = '` + org_uuid + `' ORDER BY created DESC`

	db.db.Raw(query + " " + limitQuery).Find(&payment)
	return payment
}

func (db database) GetInvoice(payment_request string) InvoiceList {
	ms := InvoiceList{}
	db.db.Where("payment_request = ?", payment_request).Find(&ms)
	return ms
}

func (db database) GetOrganizationInvoices(org_uuid string) []InvoiceList {
	ms := []InvoiceList{}
	db.db.Where("org_uuid = ?", org_uuid).Where("status", false).Find(&ms)
	return ms
}

func (db database) GetOrganizationInvoicesCount(org_uuid string) int64 {
	var count int64
	ms := InvoiceList{}

	db.db.Model(&ms).Where("org_uuid = ?", org_uuid).Where("status", false).Count(&count)
	return count
}

func (db database) UpdateInvoice(payment_request string) InvoiceList {
	ms := InvoiceList{}
	db.db.Model(&InvoiceList{}).Where("payment_request = ?", payment_request).Update("status", true)
	ms.Status = true
	return ms
}

func (db database) AddInvoice(invoice InvoiceList) InvoiceList {
	db.db.Create(&invoice)
	return invoice
}

func (db database) AddUserInvoiceData(userData UserInvoiceData) UserInvoiceData {
	db.db.Create(&userData)
	return userData
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

func (db database) ChangeOrganizationDeleteStatus(org_uuid string, status bool) Organization {
	ms := Organization{}
	db.db.Model(&ms).Where("uuid", org_uuid).Updates(map[string]interface{}{
		"deleted": status,
	})
	return ms
}
