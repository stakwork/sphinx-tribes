package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
)

type organizationHandler struct {
	db db.Database
}

func NewOrganizationHandler(db db.Database) *organizationHandler {
	return &organizationHandler{db: db}
}

func (oh *organizationHandler) CreateOrEditOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	now := time.Now()

	org := db.Organization{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &org)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if len(org.Name) == 0 || len(org.Name) > 20 {
		fmt.Printf("invalid organization name %s\n", org.Name)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error: organization name must be present and should not exceed 20 character")
		return
	}

	if pubKeyFromAuth != org.OwnerPubKey {
		hasRole := db.UserHasAccess(pubKeyFromAuth, org.Uuid, db.EditOrg)
		if !hasRole {
			fmt.Println(pubKeyFromAuth)
			fmt.Println(org.OwnerPubKey)
			fmt.Println("mismatched pubkey")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Don't have access to Edit Org")
			return
		}
	}

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Validate struct data
	err = db.Validate.Struct(org)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error: did not pass validation test : %s", err)
		json.NewEncoder(w).Encode(msg)
		return
	}

	if org.Github != "" && !strings.Contains(org.Github, "github.com/") {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Error: not a valid github"
		json.NewEncoder(w).Encode(msg)
		return
	}

	existing := oh.db.GetOrganizationByUuid(org.Uuid)
	if existing.ID == 0 { // new!
		if org.ID != 0 { // can't try to "edit" if it does not exist already
			fmt.Println("cant edit non existing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		name := org.Name

		// check if the organization name already exists
		orgName := oh.db.GetOrganizationByName(name)

		if orgName.Name == name {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Organization name already exists")
			return
		} else {
			org.Created = &now
			org.Updated = &now
			org.Uuid = xid.New().String()
			org.Name = name
		}
	} else {
		if org.ID == 0 {
			// can't create that already exists
			fmt.Println("can't create existing organization")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if org.ID != existing.ID { // can't edit someone else's
			fmt.Println("cant edit another organization")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	p, err := oh.db.CreateOrEditOrganization(org)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	orgs := db.DB.GetOrganizations(r)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgs)
}

func GetOrganizationsCount(w http.ResponseWriter, r *http.Request) {
	count := db.DB.GetOrganizationsCount()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}

func GetOrganizationByUuid(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	org := db.DB.GetOrganizationByUuid(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}

func CreateOrganizationUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	now := time.Now()

	orgUser := db.OrganizationUsers{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &orgUser)

	// get orgnanization
	org := db.DB.GetOrganizationByUuid(orgUser.OrgUuid)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check if the user is the organization admin
	if orgUser.OwnerPubKey == org.OwnerPubKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Cannot add organization admin as a user")
		return
	}

	// check if the user tries to add their self
	if pubKeyFromAuth == orgUser.OwnerPubKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Cannot add userself as a user")
		return
	}

	// if not the orgnization admin
	hasRole := db.UserHasAccess(pubKeyFromAuth, orgUser.OrgUuid, db.AddUser)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to add user")
		return
	}

	// check if the user exists on peoples table
	isUser := db.DB.GetPersonByPubkey(orgUser.OwnerPubKey)
	if isUser.OwnerPubKey != orgUser.OwnerPubKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("User doesn't exists in people")
		return
	}

	// check if user already exists
	userExists := db.DB.GetOrganizationUser(orgUser.OwnerPubKey, orgUser.OrgUuid)

	if userExists.ID != 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("User already exists")
		return
	}

	orgUser.Created = &now
	orgUser.Updated = &now

	// create user
	user := db.DB.CreateOrganizationUser(orgUser)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetOrganizationUsers(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	orgUsers, _ := db.DB.GetOrganizationUsers(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgUsers)
}

func GetOrganizationUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	orgUser := db.DB.GetOrganizationUser(pubKeyFromAuth, uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgUser)
}

func GetOrganizationUsersCount(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	count := db.DB.GetOrganizationUsersCount(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}

func DeleteOrganizationUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	orgUser := db.OrganizationUsersData{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &orgUser)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	org := db.DB.GetOrganizationByUuid(orgUser.OrgUuid)

	if orgUser.OwnerPubKey == org.OwnerPubKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Cannot delete organization admin")
		return
	}

	hasRole := db.UserHasAccess(pubKeyFromAuth, orgUser.OrgUuid, db.DeleteUser)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to delete user")
		return
	}

	db.DB.DeleteOrganizationUser(orgUser, orgUser.OrgUuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgUser)
}

func GetBountyRoles(w http.ResponseWriter, r *http.Request) {
	roles := db.DB.GetBountyRoles()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roles)
}

func AddUserRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")
	user := chi.URLParam(r, "user")
	now := time.Now()

	if uuid == "" || user == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("no uuid, or user pubkey")
		return
	}

	roles := []db.UserRoles{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &roles)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if pubKeyFromAuth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("no pubkey from auth")
		return
	}

	// if not the orgnization admin
	hasRole := db.UserHasAccess(pubKeyFromAuth, uuid, db.AddRoles)
	isUser := db.CheckUser(roles, pubKeyFromAuth)

	if isUser {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("cannot add roles for self")
		return
	}

	// check if the user added his pubkey to the route
	if pubKeyFromAuth == user {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("auth pubkey cannot be the same with user's")
		return
	}

	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("user does not have adequate permissions to add roles")
		return
	}

	rolesMap := db.GetRolesMap()
	insertRoles := []db.UserRoles{}
	for _, role := range roles {
		_, ok := rolesMap[role.Role]
		// if any of the roles does not exists return an error
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("not a valid user role")
			return
		}

		// check if the user has the role he his trying to add to another user
		okUser := db.UserHasAccess(pubKeyFromAuth, uuid, role.Role)
		// if the user does not have any of the roles he wants to add return an error
		if !okUser {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("cannot add a role you don't have")
			return
		}

		// add created time for insert
		role.Created = &now
		insertRoles = append(insertRoles, role)
	}

	// check if user already exists
	userExists := db.DB.GetOrganizationUser(user, uuid)

	// if not the organization admin
	if userExists.OwnerPubKey != user || userExists.OrgUuid != uuid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("User does not exists in the organization")
		return
	}

	db.DB.CreateUserRoles(insertRoles, uuid, user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insertRoles)
}

func GetUserRoles(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	user := chi.URLParam(r, "user")

	userRoles := db.DB.GetUserRoles(uuid, user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userRoles)
}

func GetUserOrganizations(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userId")
	userId, _ := utils.ConvertStringToUint(userIdParam)

	if userId == 0 {
		fmt.Println("provide user id")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	user := db.DB.GetPerson(userId)

	// get the organizations created by the user, then get all the organizations
	// the user has been added to, loop through to get the organization
	organizations := GetCreatedOrganizations(user.OwnerPubKey)

	assignedOrganizations := db.DB.GetUserAssignedOrganizations(user.OwnerPubKey)
	for _, value := range assignedOrganizations {
		uuid := value.OrgUuid
		organization := db.DB.GetOrganizationByUuid(uuid)
		bountyCount := db.DB.GetOrganizationBountyCount(uuid)
		hasRole := db.UserHasAccess(user.OwnerPubKey, uuid, db.ViewReport)

		// don't add deleted organizations to the list
		if !organization.Deleted {
			if hasRole {
				budget := db.DB.GetOrganizationBudget(uuid)
				organization.Budget = budget.TotalBudget
			} else {
				organization.Budget = 0
			}
			organization.BountyCount = bountyCount

			organizations = append(organizations, organization)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organizations)
}

func (oh *organizationHandler) GetUserDropdownOrganizations(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userId")
	userId, _ := utils.ConvertStringToUint(userIdParam)

	if userId == 0 {
		fmt.Println("provide user id")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	user := db.DB.GetPerson(userId)

	// get the organizations created by the user, then get all the organizations
	// the user has been added to, loop through to get the organization
	organizations := GetCreatedOrganizations(user.OwnerPubKey)

	assignedOrganizations := db.DB.GetUserAssignedOrganizations(user.OwnerPubKey)
	for _, value := range assignedOrganizations {
		uuid := value.OrgUuid
		organization := db.DB.GetOrganizationByUuid(uuid)
		bountyCount := db.DB.GetOrganizationBountyCount(uuid)
		hasRole := db.UserHasAccess(user.OwnerPubKey, uuid, db.ViewReport)
		hasBountyRoles := oh.db.UserHasManageBountyRoles(user.OwnerPubKey, uuid)

		// don't add deleted organizations to the list
		if !organization.Deleted && hasBountyRoles {
			if hasRole {
				budget := db.DB.GetOrganizationBudget(uuid)
				organization.Budget = budget.TotalBudget
			} else {
				organization.Budget = 0
			}
			organization.BountyCount = bountyCount

			organizations = append(organizations, organization)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organizations)
}

func GetCreatedOrganizations(pubkey string) []db.Organization {
	organizations := db.DB.GetUserCreatedOrganizations(pubkey)
	// add bounty count to the organization
	for index, value := range organizations {
		uuid := value.Uuid
		bountyCount := db.DB.GetOrganizationBountyCount(uuid)
		hasRole := db.UserHasAccess(pubkey, uuid, db.ViewReport)

		if hasRole {
			budget := db.DB.GetOrganizationBudget(uuid)
			organizations[index].Budget = budget.TotalBudget
		} else {
			organizations[index].Budget = 0
		}
		organizations[index].BountyCount = bountyCount
	}
	return organizations
}

func GetOrganizationBounties(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	// get the organization bounties
	organizationBounties := db.DB.GetOrganizationBounties(r, uuid)

	var bountyResponse []db.BountyResponse = GenerateBountyResponse(organizationBounties)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyResponse)
}

func GetOrganizationBudget(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if not the organization admin
	hasRole := db.UserHasAccess(pubKeyFromAuth, uuid, db.ViewReport)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to view budget")
		return
	}

	// get the organization budget
	organizationBudget := db.DB.GetOrganizationBudget(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organizationBudget)
}

func GetOrganizationBudgetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	// if not the organization admin
	hasRole := db.UserHasAccess(pubKeyFromAuth, uuid, db.ViewReport)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to view budget history")
		return
	}

	// get the organization budget
	organizationBudget := db.DB.GetOrganizationBudgetHistory(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organizationBudget)
}

func GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if not the organization admin
	hasRole := db.UserHasAccess(pubKeyFromAuth, uuid, db.ViewReport)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to view payments")
		return
	}

	// get the organization payment history
	paymentHistory := db.DB.GetPaymentHistory(uuid, r)
	paymentHistoryData := []db.PaymentHistoryData{}

	for _, payment := range paymentHistory {
		sender := db.DB.GetPersonByPubkey(payment.SenderPubKey)
		receiver := db.DB.GetPersonByPubkey(payment.ReceiverPubKey)
		paymentData := db.PaymentHistoryData{
			PaymentHistory: payment,
			SenderName:     sender.UniqueName,
			SenderImg:      sender.Img,
			ReceiverName:   receiver.UniqueName,
			ReceiverImg:    receiver.Img,
		}
		paymentHistoryData = append(paymentHistoryData, paymentData)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paymentHistoryData)
}

func PollBudgetInvoices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orgInvoices := db.DB.GetOrganizationInvoices(uuid)

	for _, inv := range orgInvoices {
		invoiceRes, invoiceErr := GetLightningInvoice(inv.PaymentRequest)

		if invoiceErr.Error != "" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(invoiceErr)
			return
		}

		if invoiceRes.Response.Settled {
			if !inv.Status && inv.Type == "BUDGET" {
				db.DB.AddAndUpdateBudget(inv)
				// Update the invoice status
				db.DB.UpdateInvoice(inv.PaymentRequest)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Polled invoices")
}

func GetInvoicesCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	invoiceCount := db.DB.GetOrganizationInvoicesCount(uuid)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceCount)
}

func (oh *organizationHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	organization := oh.db.GetOrganizationByUuid(uuid)

	if pubKeyFromAuth != organization.OwnerPubKey {
		msg := "only org admin can delete an organization"
		fmt.Println(msg)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msg)
		return
	}

	// Update organization to hide and clear certain fields
	if err := oh.db.UpdateOrganizationForDeletion(uuid); err != nil {
		fmt.Println("Error updating organization:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete all users from the organization
	if err := oh.db.DeleteAllUsersFromOrganization(uuid); err != nil {
		fmt.Println("Error removing users from organization:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// soft delete organization
	org := oh.db.ChangeOrganizationDeleteStatus(uuid, true)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}
