package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
)

func CreateOrEditOrganization(w http.ResponseWriter, r *http.Request) {
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

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pubKeyFromAuth != org.OwnerPubKey {
		fmt.Println(pubKeyFromAuth)
		fmt.Println(org.OwnerPubKey)
		fmt.Println("mismatched pubkey")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := db.DB.GetOrganizationByUuid(org.Uuid)
	if existing.ID == 0 { // new!
		if org.ID != 0 { // cant try to "edit" if not exists already
			fmt.Println("cant edit non existing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		name := org.Name

		// check if the organization name already exists
		orgName := db.DB.GetOrganizationByName(name)

		if orgName.Name == name {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Organization name alreday exists")
			return
		} else {
			org.Created = &now
			org.Updated = &now
			org.Uuid = xid.New().String()
			org.Name = name
		}
	} else {
		if org.ID == 0 {
			// cant create that already exists
			fmt.Println("can't create existing organization")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if org.ID != existing.ID { // cant edit someone else's
			fmt.Println("cant edit another organization")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	p, err := db.DB.CreateOrEditOrganization(org)
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

	db.DB.DeleteOrganizationUser(orgUser)

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
	userRoles := db.DB.GetUserRoles(uuid, pubKeyFromAuth)
	isUser := db.CheckUser(userRoles, pubKeyFromAuth)

	if isUser {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("cannot add roles for self")
		return
	}

	// check if the user added his pubkey to the route
	if pubKeyFromAuth == user {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("cannot add roles for self")
		return
	}

	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("user cannot add roles")
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

		// add created time for insert
		role.Created = &now
		insertRoles = append(insertRoles, role)
	}

	// check if user already exists
	userExists := db.DB.GetOrganizationUser(user, uuid)

	// if not the orgnization admin
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
	organizations := db.DB.GetUserCreatedOrganizations(user.OwnerPubKey)
	// add bounty count to the organization
	for index, value := range organizations {
		uuid := value.Uuid
		bountyCount := db.DB.GetOrganizationBountyCount(uuid)
		hasRole := db.UserHasAccess(user.OwnerPubKey, uuid, db.ViewReport)

		if hasRole {
			budget := db.DB.GetOrganizationBudget(uuid)
			organizations[index].Budget = budget.TotalBudget
		} else {
			organizations[index].Budget = 0
		}
		organizations[index].BountyCount = bountyCount
	}

	assignedOrganizations := db.DB.GetUserAssignedOrganizations(user.OwnerPubKey)
	for _, value := range assignedOrganizations {
		uuid := value.OrgUuid
		organization := db.DB.GetOrganizationByUuid(uuid)
		bountyCount := db.DB.GetOrganizationBountyCount(uuid)
		hasRole := db.UserHasAccess(user.OwnerPubKey, uuid, db.ViewReport)

		if hasRole {
			budget := db.DB.GetOrganizationBudget(uuid)
			organization.Budget = budget.TotalBudget
		} else {
			organization.Budget = 0
		}
		organization.BountyCount = bountyCount

		organizations = append(organizations, organization)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(organizations)
}

func GetOrganizationBounties(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	// get the organization bounties
	organizationBounties := db.DB.GetOrganizationBounties(r, uuid)

	var bountyResponse []db.BountyResponse = generateBountyResponse(organizationBounties)
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

	// if not the orgnization admin
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

	// if not the orgnization admin
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
	keys := r.URL.Query()
	page := keys.Get("page")
	limit := keys.Get("limit")

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if not the orgnization admin
	hasRole := db.UserHasAccess(pubKeyFromAuth, uuid, db.ViewReport)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to view payments")
		return
	}

	// get the organization payment history
	paymentHistory := db.DB.GetPaymentHistory(uuid, page, limit)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paymentHistory)
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

func MemeImageUpload(w http.ResponseWriter, r *http.Request) {
	// Parsing uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Println("FIle  ===", header.Filename)

	// Saving the file
	dst, err := os.Create("uploads/" + header.Filename)
	fmt.Println("FILE +++", dst)
	if err != nil {
		http.Error(w, "Unable to save file 1", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Unable to save file 2", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File uploaded successfully"))

	// buf := bytes.NewBuffer(nil)
	// if _, err := io.Copy(buf, file); err != nil {
	// 	fmt.Println("Error ===", err)
	// }

	// err = os.WriteFile("/uploads/"+header.Filename, buf.Bytes(), 0644)
	// if err != nil {
	// 	fmt.Println("WRITE FILE ERROR", err)
	// }
	// TODO

	// Upload the file
	// Send to meme server ask endpoint to get a challenge
	// GET /ask
	// Send to RELAY to sign the chaallenge
	// GET http://localhost:3001/signer/${r.challenge}

	// POST TO MEME SERVER TO GET TO GET TOKEN
	// POST /verify"
	// form: { id: r.id, sig: r2.response.sig, pubkey: node.pubkey },
	// ID FROM CHALLENGE
	// SIG FROM RELAY
	// NODE PUBKEY fROM LND NODE

	// WHEN UT RETURNS SEND THE IMAGE TO MEME SERVER PUBLIC
}
