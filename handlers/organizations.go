package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"net/url"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
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
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Saving the file
	dst, err := os.Create("uploads/" + header.Filename)

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

	challenge := GetMemeChallenge()
	signer := SignChallenge(challenge.Challenge)
	mErr, mToken := GetMemeToken(challenge.Id, signer.Response.Sig)

	if mErr != "" {
		msg := "Could not get meme token"
		fmt.Println(msg, mErr)
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(msg)
	} else {
		memeImgUrl := UploadMemeImage(file, mToken.Token, header.Filename)
		if memeImgUrl == "" {
			msg := "Could not get meme image"
			fmt.Println(msg)
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(msg)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(memeImgUrl)
		}
	}
}

func GetMemeChallenge() db.MemeChallenge {
	memeChallenge := db.MemeChallenge{}

	url := fmt.Sprintf("%s/ask", config.MemeUrl)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	// Unmarshal result
	err = json.Unmarshal(body, &memeChallenge)

	if err != nil {
		log.Printf("Reading Invoice body failed: %s", err)
	}

	return memeChallenge
}

func SignChallenge(challenge string) db.RelaySignerResponse {
	url := fmt.Sprintf("%s/signer/%s", config.RelayUrl, challenge)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	signerResponse := db.RelaySignerResponse{}

	// Unmarshal result
	err = json.Unmarshal(body, &signerResponse)

	if err != nil {
		log.Printf("Reading Challenge body failed: %s", err)
	}

	return signerResponse
}

func GetMemeToken(id string, sig string) (string, db.MemeTokenSuccess) {
	memeUrl := fmt.Sprintf("%s/verify", config.MemeUrl)

	formData := url.Values{
		"id":     {id},
		"sig":    {sig},
		"pubkey": {config.RelayNodeKey},
	}

	res, err := http.PostForm(memeUrl, formData)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return "", db.MemeTokenSuccess{}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if res.StatusCode == 200 {
		tokenSuccess := db.MemeTokenSuccess{}

		// Unmarshal result
		err = json.Unmarshal(body, &tokenSuccess)

		if err != nil {
			log.Printf("Reading token success body failed: %s", err)
		}

		return "", tokenSuccess
	} else {
		var tokenError string

		// Unmarshal result
		err = json.Unmarshal(body, &tokenError)

		if err != nil {
			log.Printf("Reading token error body failed: %s %d", err, res.StatusCode)
		}

		return tokenError, db.MemeTokenSuccess{}
	}
}

func UploadMemeImage(file multipart.File, token string, fileName string) string {
	url := fmt.Sprintf("%s/public", config.MemeUrl)
	filePath := path.Join("./uploads", fileName)
	fileW, _ := os.Open(filePath)
	defer file.Close()

	fileBody := &bytes.Buffer{}
	writer := multipart.NewWriter(fileBody)
	part, _ := writer.CreateFormFile("file", filepath.Base(filePath))
	io.Copy(part, fileW)
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, fileBody)
	req.Header.Set("Authorization", "BEARER "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)

	// Delete image from uploads folder
	DeleteImageFromUploadsFolder(filePath)

	if err != nil {
		fmt.Println("meme request Error ===", err)
		return ""
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err == nil {
		memeSuccess := db.Meme{}
		// Unmarshal result
		err = json.Unmarshal(body, &memeSuccess)
		if err != nil {
			log.Printf("Reading meme error body failed: %s", err)
		} else {
			return config.MemeUrl + "/public/" + memeSuccess.Muid
		}
	}

	return ""
}

func DeleteImageFromUploadsFolder(filePath string) {
	e := os.Remove(filePath)
	if e != nil {
		log.Printf("Could not delete Image %s %s", filePath, e)
	}
}
