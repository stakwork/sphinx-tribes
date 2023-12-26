//go:build mock
// +build mock

package db

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var bountyCache map[string]Bounty
var organizationCache map[string]Organization
var userRolesCache map[string]map[string][]UserRoles

func (db database) GetConnectionCode() ConnectionCodesShort {
	c := ConnectionCodesShort{}
	c.ConnectionString = "test"
	return c
}

func (db database) GetPeopleForNewTicket(languages []interface{}) ([]Person, error) {
	ms := []Person{}

	return ms, nil
}

func (db database) GetAllPeople() []Person {
	ms := []Person{}
	return ms
}

func (db database) AddUuidToPerson(id uint, uuid string) {
	if id == 0 {
		return
	}
}

func (db database) GetOrganizationByUuid(uuid string) Organization {
	return organizationCache[uuid]
}
func (db database) GetUserRoles(uuid string, pubkey string) []UserRoles {
	orgRoles := userRolesCache[uuid]
	if orgRoles == nil {
		return nil
	}
	return orgRoles[pubkey]
}
func (db database) GetPersonByPubkey(pubkey string) Person {
	m := Person{}

	return m
}

func (db database) UpdatePerson(id uint, u map[string]interface{}) bool {
	if id == 0 {
		return false
	}

	return true
}

func (db database) CreateOrEditTribe(m Tribe) (Tribe, error) {

	return m, nil
}

func (db database) CreateChannel(c Channel) (Channel, error) {

	if c.Created == nil {
		now := time.Now()
		c.Created = &now

	}
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

	}

	return m, nil
}

func (db database) GetUnconfirmedTwitter() []Person {
	ms := []Person{}
	return ms
}

func (db database) UpdateTwitterConfirmed(id uint, confirmed bool) {
	if id == 0 {
		return
	}
}

func (db database) GetUnconfirmedGithub() []Person {
	ms := []Person{}
	return ms
}

func (db database) UpdateGithubConfirmed(id uint, confirmed bool) {
	if id == 0 {
		return
	}
}

func (db database) UpdateGithubIssues(id uint, issues map[string]interface{}) {
}

func (db database) UpdateTribe(uuid string, u map[string]interface{}) bool {
	if uuid == "" {
		return false
	}
	return true
}

func (db database) UpdateChannel(id uint, u map[string]interface{}) bool {
	if id == 0 {
		return false
	}
	return true
}

func (db database) UpdateTribeUniqueName(uuid string, u string) {
	if uuid == "" {
		return
	}
}

type GithubOpenIssue struct {
	Status   string `json:"status"`
	Assignee string `json:"assignee"`
}

type GithubOpenIssueCount struct {
	Count int64 `json:"count"`
}

func (db database) GetOpenGithubIssues(r *http.Request) (int64, error) {

	return 0, nil
}

func (db database) GetListedTribes(r *http.Request) []Tribe {
	ms := []Tribe{}
	return ms
}

func (db database) GetTribesByOwner(pubkey string) []Tribe {
	ms := []Tribe{}
	return ms
}

func (db database) GetAllTribesByOwner(pubkey string) []Tribe {
	ms := []Tribe{}
	return ms
}

func (db database) GetTribesByAppUrl(aurl string) []Tribe {
	ms := []Tribe{}
	return ms
}

func (db database) GetChannelsByTribe(tribe_uuid string) []Channel {
	ms := []Channel{}
	return ms
}

func (db database) GetChannel(id uint) Channel {
	ms := Channel{}
	return ms
}

func (db database) GetListedBots(r *http.Request) []Bot {
	ms := []Bot{}

	return ms
}

func (db database) GetListedPeople(r *http.Request) []Person {
	ms := []Person{}
	return ms
}

func (db database) GetPeopleBySearch(r *http.Request) []Person {
	ms := []Person{}
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

	return ms, nil
}

func (db database) GetUserBountiesCount(personKey string, tabType string) int64 {
	return 0
}

func (db database) GetBountiesCount(r *http.Request) int64 {
	return 0
}

func (db database) GetOrganizationBounties(r *http.Request, org_uuid string) []Bounty {
	ms := []Bounty{}

	return ms
}

func (db database) GetAssignedBounties(pubkey string) ([]Bounty, error) {
	ms := []Bounty{}
	return ms, nil
}

func (db database) GetCreatedBounties(pubkey string) ([]Bounty, error) {
	ms := []Bounty{}
	return ms, nil
}

func (db database) GetBountyById(id string) ([]Bounty, error) {
	var bounty Bounty
	var present bool
	if bounty, present = bountyCache[id]; !present {
		return nil, nil
	}
	return []Bounty{bounty}, nil
}

func (db database) GetBountyIndexById(id string) int64 {
	var index int64
	return index
}

func (db database) GetBountyDataByCreated(created string) ([]Bounty, error) {
	ms := []Bounty{}
	return ms, nil
}

func (db database) AddBounty(b Bounty) (Bounty, error) {
	return b, nil
}

func (db database) GetAllBounties(r *http.Request) []Bounty {
	ms := []Bounty{}

	return ms
}

func (db database) CreateOrEditBounty(b Bounty) (Bounty, error) {
	if b.OwnerID == "" {
		return Bounty{}, errors.New("no pub key")
	}
	if bountyCache == nil {
		bountyCache = make(map[string]Bounty)
	}
	bountyCache[strconv.Itoa(int(b.ID))] = b
	return b, nil
}

func (db database) UpdateBountyNullColumn(b Bounty, column string) Bounty {
	columnMap := make(map[string]interface{})
	columnMap[column] = ""
	return b
}

func (db database) UpdateBountyBoolColumn(b Bounty, column string) Bounty {
	columnMap := make(map[string]interface{})
	columnMap[column] = false
	return b
}

func (db database) DeleteBounty(pubkey string, created string) (Bounty, error) {
	m := Bounty{}
	return m, nil
}

func (db database) GetBountyByCreated(created uint) (Bounty, error) {
	b := Bounty{}
	return b, nil
}

func (db database) GetBounty(id uint) Bounty {
	b := Bounty{}
	return b
}

func (db database) UpdateBounty(b Bounty) (Bounty, error) {
	return b, nil
}

func (db database) UpdateBountyPayment(b Bounty) (Bounty, error) {
	return b, nil
}

func (db database) GetListedOffers(r *http.Request) ([]PeopleExtra, error) {
	ms := []PeopleExtra{}
	return ms, nil
}

func (db database) UpdateBot(uuid string, u map[string]interface{}) bool {
	if uuid == "" {
		return false
	}
	return true
}

func (db database) GetAllTribes() []Tribe {
	ms := []Tribe{}
	return ms
}

func (db database) GetTribesTotal() int64 {
	var count int64
	return count
}

func (db database) GetTribeByIdAndPubkey(uuid string, pubkey string) Tribe {
	m := Tribe{}
	return m
}

func (db database) GetTribe(uuid string) Tribe {
	m := Tribe{}
	return m
}

func (db database) GetPerson(id uint) Person {
	m := Person{}
	return m
}

func (db database) GetPersonByUuid(uuid string) Person {
	m := Person{}

	return m
}

func (db database) GetPersonByGithubName(github_name string) Person {
	m := Person{}

	return m
}

func (db database) GetFirstTribeByFeedURL(feedURL string) Tribe {
	m := Tribe{}
	return m
}

func (db database) GetBot(uuid string) Bot {
	m := Bot{}
	return m
}

func (db database) GetTribeByUniqueName(un string) Tribe {
	m := Tribe{}
	return m
}

func (db database) GetBotsByOwner(pubkey string) []Bot {
	bs := []Bot{}
	return bs
}

func (db database) GetBotByUniqueName(un string) Bot {
	m := Bot{}
	return m
}

func (db database) GetPersonByUniqueName(un string) Person {
	m := Person{}
	return m
}

func (db database) SearchTribes(s string) []Tribe {
	ms := []Tribe{}
	if s == "" {
		return ms
	}
	// set limit
	return ms
}

func (db database) SearchBots(s string, limit, offset int) []BotRes {
	ms := []BotRes{}
	if s == "" {
		return ms
	}
	// set limit
	return ms
}

func (db database) SearchPeople(s string, limit, offset int) []Person {
	ms := []Person{}
	if s == "" {
		return ms
	}
	return ms
}

func (db database) CreateLeaderBoard(uuid string, leaderboards []LeaderBoard) ([]LeaderBoard, error) {
	return leaderboards, nil

}

func (db database) GetLeaderBoard(uuid string) []LeaderBoard {
	m := []LeaderBoard{}
	return m
}

func (db database) GetLeaderBoardByUuidAndAlias(uuid string, alias string) LeaderBoard {
	m := LeaderBoard{}
	return m
}

func (db database) UpdateLeaderBoard(uuid string, alias string, u map[string]interface{}) bool {
	return true
}

func (db database) CountDevelopers() int64 {
	var count int64
	return count
}

func (db database) CountBounties() uint64 {
	var count uint64
	return count
}

func (db database) GetPeopleListShort(count uint32) *[]PersonInShort {
	p := []PersonInShort{}
	return &p
}

func (db database) CreateConnectionCode(c ConnectionCodes) (ConnectionCodes, error) {
	if c.DateCreated == nil {
		now := time.Now()
		c.DateCreated = &now
	}
	if c.ID == 1 {
		errorMsg := errors.New("Error Reaching DB")
		return c, errorMsg
	}
	return c, nil
}

func (db database) GetLnUser(lnKey string) int64 {
	var count int64

	return count
}

func (db database) CreateLnUser(lnKey string) (Person, error) {
	p := Person{}
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
	var users = []LeaderData{}

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

	return ms
}

func (db database) GetOrganizationsCount() int64 {
	var count int64
	return count
}

func (db database) GetOrganizationByName(name string) Organization {
	ms := Organization{}

	return ms
}

func (db database) CreateOrEditOrganization(m Organization) (Organization, error) {
	if m.OwnerPubKey == "" {
		return Organization{}, errors.New("no pub key")
	}
	if organizationCache == nil {
		organizationCache = make(map[string]Organization)
	}
	organizationCache[m.Uuid] = m
	return m, nil
}

func (db database) GetOrganizationUsers(uuid string) ([]OrganizationUsersData, error) {
	ms := []OrganizationUsersData{}

	return ms, nil
}

func (db database) GetOrganizationUsersCount(uuid string) int64 {
	return 0
}

func (db database) GetOrganizationBountyCount(uuid string) int64 {
	var count int64
	return count
}

func (db database) GetOrganizationUser(pubkey string, org_uuid string) OrganizationUsers {
	ms := OrganizationUsers{}
	return ms
}

func (db database) CreateOrganizationUser(orgUser OrganizationUsers) OrganizationUsers {

	return orgUser
}

func (db database) DeleteOrganizationUser(orgUser OrganizationUsersData, org string) OrganizationUsersData {
	return orgUser
}

func (db database) GetBountyRoles() []BountyRoles {
	ms := []BountyRoles{}
	return ms
}

func (db database) CreateUserRoles(roles []UserRoles, uuid string, pubkey string) []UserRoles {
	// delete roles and create new ones
	if userRolesCache == nil {
		userRolesCache = make(map[string]map[string][]UserRoles)
	}
	_, orgRolePresent := userRolesCache[uuid]
	if !orgRolePresent {
		userRolesCache[uuid] = make(map[string][]UserRoles)
	}
	userRolesCache[uuid][pubkey] = roles
	return roles
}

func (db database) GetUserCreatedOrganizations(pubkey string) []Organization {
	ms := []Organization{}
	return ms
}

func (db database) GetUserAssignedOrganizations(pubkey string) []OrganizationUsers {
	ms := []OrganizationUsers{}
	return ms
}

func (db database) AddBudgetHistory(budget BudgetHistory) BudgetHistory {
	return budget
}

func (db database) CreateOrganizationBudget(budget BountyBudget) BountyBudget {
	return budget
}

func (db database) UpdateOrganizationBudget(budget BountyBudget) BountyBudget {
	return budget
}

func (db database) GetPaymentHistoryByCreated(created *time.Time, org_uuid string) PaymentHistory {
	ms := PaymentHistory{}
	return ms
}

func (db database) GetOrganizationBudget(org_uuid string) BountyBudget {
	ms := BountyBudget{}
	return ms
}

func (db database) GetOrganizationBudgetHistory(org_uuid string) []BudgetHistoryData {
	budgetHistory := []BudgetHistoryData{}

	return budgetHistory
}

func (db database) AddAndUpdateBudget(invoice InvoiceList) PaymentHistory {
	paymentHistory := PaymentHistory{}

	return paymentHistory
}

func (db database) WithdrawBudget(sender_pubkey string, org_uuid string, amount uint) {
	// get organization budget and add payment to total budget

}

func (db database) AddPaymentHistory(payment PaymentHistory) PaymentHistory {

	return payment
}

func (db database) GetPaymentHistory(org_uuid string, r *http.Request) []PaymentHistory {
	payment := []PaymentHistory{}

	return payment
}

func (db database) GetInvoice(payment_request string) InvoiceList {
	ms := InvoiceList{}
	return ms
}

func (db database) GetOrganizationInvoices(org_uuid string) []InvoiceList {
	ms := []InvoiceList{}
	return ms
}

func (db database) GetOrganizationInvoicesCount(org_uuid string) int64 {

	return 0
}

func (db database) UpdateInvoice(payment_request string) InvoiceList {
	ms := InvoiceList{}
	ms.Status = true
	return ms
}

func (db database) AddInvoice(invoice InvoiceList) InvoiceList {
	return invoice
}

func (db database) AddUserInvoiceData(userData UserInvoiceData) UserInvoiceData {
	return userData
}

func (db database) GetUserInvoiceData(payment_request string) UserInvoiceData {
	ms := UserInvoiceData{}
	return ms
}

func (db database) DeleteUserInvoiceData(payment_request string) UserInvoiceData {
	ms := UserInvoiceData{}
	return ms
}

func (db database) ChangeOrganizationDeleteStatus(org_uuid string, status bool) Organization {
	ms := Organization{}
	return ms
}

func (db database) GetFilterStatusCount() FilterStattuCount {
	ms := FilterStattuCount{}

	return ms
}
