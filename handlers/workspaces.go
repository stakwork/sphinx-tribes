package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/utils"
	"gorm.io/gorm"
)

type workspaceHandler struct {
	db                             db.Database
	generateBountyHandler          func(bounties []db.NewBounty) []db.BountyResponse
	getLightningInvoice            func(payment_request string) (db.InvoiceResult, db.InvoiceError)
	userHasAccess                  func(pubKeyFromAuth string, uuid string, role string) bool
	configUserHasAccess            func(pubKeyFromAuth string, uuid string, role string) bool
	configUserHasManageBountyRoles func(pubKeyFromAuth string, uuid string) bool
	userHasManageBountyRoles       func(pubKeyFromAuth string, uuid string) bool
	getAllUserWorkspaces           func(pubKeyFromAuth string) []db.Workspace
}

func NewWorkspaceHandler(database db.Database) *workspaceHandler {
	bHandler := NewBountyHandler(http.DefaultClient, database)
	dbConf := db.NewDatabaseConfig(&gorm.DB{})
	configHandler := db.NewConfigHandler(database)
	return &workspaceHandler{
		db:                             database,
		generateBountyHandler:          bHandler.GenerateBountyResponse,
		getLightningInvoice:            bHandler.GetLightningInvoice,
		userHasAccess:                  dbConf.UserHasAccess,
		configUserHasAccess:            configHandler.UserHasAccess,
		configUserHasManageBountyRoles: configHandler.UserHasManageBountyRoles,
		userHasManageBountyRoles:       dbConf.UserHasManageBountyRoles,
		getAllUserWorkspaces:           GetAllUserWorkspaces,
	}
}

// CreateOrEditWorkspace godoc
//
//	@Summary		Create or Edit Workspace
//	@Description	Create or edit a workspace
//	@Tags			Workspaces
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace	body		db.Workspace	true	"Workspace"
//	@Success		200			{object}	db.Workspace
//	@Router			/workspace [post]
func (oh *workspaceHandler) CreateOrEditWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	now := time.Now()

	workspace := db.Workspace{}
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()
	err := json.Unmarshal(body, &workspace)

	if err != nil {
		logger.Log.Error("[workspaces] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	workspace.Name = strings.TrimSpace(workspace.Name)

	if len(workspace.Name) == 0 || len(workspace.Name) > 20 {
		logger.Log.Info("[workspaces] invalid workspace name %s", workspace.Name)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error: workspace name must be present and should not exceed 20 character")
		return
	}

	if len(workspace.Description) > 120 {
		logger.Log.Info("[workspaces] invalid workspace name %s", workspace.Description)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error: workspace description should not exceed 120 character")
		return
	}

	if pubKeyFromAuth != workspace.OwnerPubKey {
		hasRole := db.UserHasAccess(pubKeyFromAuth, workspace.Uuid, db.EditOrg)
		if !hasRole {
			logger.Log.Info("[workspaces] mismatched pubkey")
			logger.Log.Info("[workspaces] Auth pubkey: %s", pubKeyFromAuth)
			logger.Log.Info("[workspaces] OwnerPubKey: %s", workspace.OwnerPubKey)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Don't have access to Edit workspace")
			return
		}
	}

	// Validate struct data
	err = db.Validate.Struct(workspace)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error: did not pass validation test : %s", err)
		json.NewEncoder(w).Encode(msg)
		return
	}

	if workspace.Github != "" && !strings.Contains(workspace.Github, "github.com/") {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Error: not a valid github"
		json.NewEncoder(w).Encode(msg)
		return
	}

	existing := oh.db.GetWorkspaceByUuid(workspace.Uuid)
	if existing.ID == 0 { // new!
		if workspace.ID != 0 { // can't try to "edit" if it does not exist already
			logger.Log.Info("[workspaces] cant edit non existing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		name := workspace.Name

		// check if the workspace name already exists
		workspaceSameName := oh.db.GetWorkspaceByName(name)
		if workspaceSameName.Name == name {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode("Workspace name already exists - " + name)
			return
		}

		workspace.Created = &now
		workspace.Updated = &now
		if len(workspace.Uuid) == 0 {
			workspace.Uuid = xid.New().String()
		}
	} else {
		workspace.Updated = &now
		workspace.Created = existing.Created
	}

	p, err := oh.db.CreateOrEditWorkspace(workspace)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// GetWorkspaces godoc
//
//	@Summary		Get Workspaces
//	@Description	Get all workspaces
//	@Tags			Workspaces
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	db.Workspace
//	@Router			/workspace [get]
func GetWorkspaces(w http.ResponseWriter, r *http.Request) {
	orgs := db.DB.GetWorkspaces(r)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgs)
}

// GetWorkspacesCount godoc
//
//	@Summary		Get Workspaces Count
//	@Description	Get the count of all workspaces
//	@Tags			Workspaces
//	@Accept			json
//	@Produce		json
//	@Success		200	{int}	int
//	@Router			/workspaces/count [get]
func GetWorkspacesCount(w http.ResponseWriter, r *http.Request) {
	count := db.DB.GetWorkspacesCount()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}

// GetWorkspaceByUuid godoc
//
//	@Summary		Get Workspace by UUID
//	@Description	Get a workspace by its UUID
//	@Tags			Workspaces
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string	true	"Workspace UUID"
//	@Success		200		{object}	db.Workspace
//	@Router			/workspaces/{uuid} [get]
func GetWorkspaceByUuid(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	workspace := db.DB.GetWorkspaceByUuid(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspace)
}

// CreateWorkspaceUser godoc
//
//	@Summary		Create Workspace User
//	@Description	Create a user for a workspace
//	@Tags			Workspace -  Users
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspaceUser	body		db.WorkspaceUsers	true	"Workspace User"
//	@Success		200				{object}	db.WorkspaceUsers
//	@Router			/workspaces/users/{uuid} [post]
func (oh *workspaceHandler) CreateWorkspaceUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	now := time.Now()

	workspaceUser := db.WorkspaceUsers{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		logger.Log.Error("[body] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = json.Unmarshal(body, &workspaceUser)

	if workspaceUser.WorkspaceUuid == "" && workspaceUser.OrgUuid != "" {
		workspaceUser.WorkspaceUuid = workspaceUser.OrgUuid
	}

	// get orgnanization
	workspace := oh.db.GetWorkspaceByUuid(workspaceUser.WorkspaceUuid)

	if err != nil {
		logger.Log.Error("[workspaces] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check if the user is the workspace admin
	if workspaceUser.OwnerPubKey == workspace.OwnerPubKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Cannot add workspace admin as a user")
		return
	}

	// check if the user tries to add their self
	if pubKeyFromAuth == workspaceUser.OwnerPubKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Cannot add userself as a user")
		return
	}

	// if not the orgnization admin
	hasRole := oh.userHasAccess(pubKeyFromAuth, workspaceUser.WorkspaceUuid, db.AddUser)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to add user")
		return
	}

	// check if the user exists on peoples table
	isUser := oh.db.GetPersonByPubkey(workspaceUser.OwnerPubKey)
	if isUser.OwnerPubKey != workspaceUser.OwnerPubKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("User doesn't exists in people")
		return
	}

	// check if user already exists
	userExists := oh.db.GetWorkspaceUser(workspaceUser.OwnerPubKey, workspaceUser.WorkspaceUuid)

	if userExists.ID != 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("User already exists")
		return
	}

	workspaceUser.Created = &now
	workspaceUser.Updated = &now

	// create user
	user := oh.db.CreateWorkspaceUser(workspaceUser)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetWorkspaceUsers godoc
//
//	@Summary		Get Workspace Users
//	@Description	Get users of a workspace by its UUID
//	@Tags			Workspace -  Users
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Success		200		{array}	db.WorkspaceUsers
//	@Router			/workspaces/users/{uuid} [get]
func GetWorkspaceUsers(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	workspaceUsers, _ := db.DB.GetWorkspaceUsers(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceUsers)
}

// GetWorkspaceUser godoc
//
//	@Summary		Get Workspace User
//	@Description	Get a user of a workspace by its UUID
//	@Tags			Workspace -  Users
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string	true	"Workspace UUID"
//	@Success		200		{object}	db.WorkspaceUsers
//	@Router			/workspaces/foruser/{uuid} [get]
func GetWorkspaceUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	workspaceUser := db.DB.GetWorkspaceUser(pubKeyFromAuth, uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceUser)
}

// GetWorkspaceUsersCount godoc
//
//	@Summary		Get Workspace Users Count
//	@Description	Get the count of users in a workspace
//	@Tags			Workspace -  Users
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Success		200		{int}	int
//	@Router			/workspaces/users/{uuid}/count [get]
func GetWorkspaceUsersCount(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	count := db.DB.GetWorkspaceUsersCount(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}

// DeleteWorkspaceUser godoc
//
//	@Summary		Delete Workspace User
//	@Description	Delete a user from a workspace
//	@Tags			Workspace -  Users
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid			path		string					true	"Workspace UUID"
//	@Param			workspaceUser	body		db.WorkspaceUsersData	true	"Workspace User Data"
//	@Success		200				{object}	db.WorkspaceUsersData
//	@Router			/workspaces/users/{uuid} [delete]
func DeleteWorkspaceUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	workspaceUser := db.WorkspaceUsersData{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		logger.Log.Error("[body] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = json.Unmarshal(body, &workspaceUser)

	if workspaceUser.WorkspaceUuid == "" && workspaceUser.OrgUuid != "" {
		workspaceUser.WorkspaceUuid = workspaceUser.OrgUuid
	}

	if err != nil {
		logger.Log.Error("[workspaces] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspace := db.DB.GetWorkspaceByUuid(workspaceUser.WorkspaceUuid)

	if workspaceUser.OwnerPubKey == workspace.OwnerPubKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Cannot delete workspace admin")
		return
	}

	hasRole := db.UserHasAccess(pubKeyFromAuth, workspaceUser.WorkspaceUuid, db.DeleteUser)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to delete user")
		return
	}

	db.DB.DeleteWorkspaceUser(workspaceUser, workspaceUser.WorkspaceUuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceUser)
}

// GetBountyRoles godoc
//
//	@Summary		Get Bounty Roles
//	@Description	Get all bounty roles
//	@Tags			Workspace -  Bounties
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{array}	db.BountyRoles
//	@Router			/workspaces/bounty/roles [get]
func GetBountyRoles(w http.ResponseWriter, r *http.Request) {
	roles := db.DB.GetBountyRoles()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roles)
}

// AddUserRoles godoc
//
//	@Summary		Add User Roles
//	@Description	Add roles to a user in a workspace
//	@Tags			Workspace -  Users
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path	string					true	"Workspace UUID"
//	@Param			user	path	string					true	"User PubKey"
//	@Param			roles	body	db.WorkspaceUserRoles	true	"Workspace User Roles"
//	@Success		200		{array}	db.WorkspaceUserRoles
//	@Router			/workspaces/users/role/{uuid}/{user} [post]
func (oh *workspaceHandler) AddUserRoles(w http.ResponseWriter, r *http.Request) {
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

	roles := []db.WorkspaceUserRoles{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		logger.Log.Error("[body] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = json.Unmarshal(body, &roles)

	if err != nil {
		logger.Log.Error("[workspaces]: %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if pubKeyFromAuth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("no pubkey from auth")
		return
	}

	// if not the orgnization admin
	hasRole := oh.userHasAccess(pubKeyFromAuth, uuid, db.AddRoles)
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
	insertRoles := []db.WorkspaceUserRoles{}
	for _, role := range roles {

		if role.WorkspaceUuid == "" && role.OrgUuid != "" {
			role.WorkspaceUuid = role.OrgUuid
		}

		_, ok := rolesMap[role.Role]
		// if any of the roles does not exists return an error
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("not a valid user role")
			return
		}

		// check if the user has the role he his trying to add to another user
		okUser := oh.userHasAccess(pubKeyFromAuth, uuid, role.Role)
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
	userExists := oh.db.GetWorkspaceUser(user, uuid)

	// if not the workspace admin
	if userExists.OwnerPubKey != user || userExists.WorkspaceUuid != uuid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("User does not exists in the workspace")
		return
	}

	oh.db.CreateUserRoles(insertRoles, uuid, user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insertRoles)
}

// GetUserRoles godoc
//
//	@Summary		Get User Roles
//	@Description	Get roles of a user in a workspace
//	@Tags			Workspace -  Users
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Param			user	path	string	true	"User PubKey"
//	@Success		200		{array}	db.WorkspaceUserRoles
//	@Router			/workspaces/users/role/{uuid}/{user} [get]
func (oh *workspaceHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	user := chi.URLParam(r, "user")

	userRoles := oh.db.GetUserRoles(uuid, user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userRoles)
}

// GetUserWorkspaces godoc
//
//	@Summary		Get User Workspaces
//	@Description	Get workspaces of a user by their ID
//	@Tags			Workspace -  Users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	uint	true	"User ID"
//	@Success		200		{array}	db.Workspace
//	@Router			/workspaces/user/{userId} [get]
func GetUserWorkspaces(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userId")
	userId, _ := utils.ConvertStringToUint(userIdParam)

	if userId == 0 {
		logger.Log.Info("[workspaces] provide user id")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	user := db.DB.GetPerson(userId)
	// get the user workspaces
	workspaces := GetAllUserWorkspaces(user.OwnerPubKey)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaces)
}

// GetUserDropdownWorkspaces godoc
//
//	@Summary		Get User Dropdown Workspaces
//	@Description	Get dropdown workspaces of a user by their ID
//	@Tags			Workspace -  Users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	uint	true	"User ID"
//	@Success		200		{array}	db.Workspace
//	@Router			/workspaces/user/dropdown/{userId} [get]
func (oh *workspaceHandler) GetUserDropdownWorkspaces(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userId")
	userId, _ := utils.ConvertStringToUint(userIdParam)

	if userId == 0 {
		logger.Log.Info("[workspaces] provide user id")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	user := oh.db.GetPerson(userId)

	// get the workspaces created by the user, then get all the workspaces
	// the user has been added to, loop through to get the workspace
	workspaces := oh.GetCreatedWorkspaces(user.OwnerPubKey)
	assignedWorkspaces := oh.db.GetUserAssignedWorkspaces(user.OwnerPubKey)
	for _, value := range assignedWorkspaces {
		uuid := value.WorkspaceUuid
		workspace := oh.db.GetWorkspaceByUuid(uuid)
		bountyCount := oh.db.GetWorkspaceBountyCount(uuid)
		hasRole := oh.configUserHasAccess(user.OwnerPubKey, uuid, db.ViewReport)
		hasBountyRoles := oh.configUserHasManageBountyRoles(user.OwnerPubKey, uuid)

		alreadyAdded := false

		if workspace.OwnerPubKey == user.OwnerPubKey {
			alreadyAdded = true
		}

		// don't add deleted workspaces to the list
		if !workspace.Deleted && hasBountyRoles {

			// check if workspace has already been added to the list
			for _, existingWorkspace := range workspaces {
				if existingWorkspace.Uuid == workspace.Uuid {
					alreadyAdded = true
				}
			}

			if hasRole && !alreadyAdded {
				budget := oh.db.GetWorkspaceBudget(uuid)
				workspace.Budget = budget.TotalBudget
			} else {
				workspace.Budget = 0
			}

			workspace.BountyCount = bountyCount

			if !alreadyAdded {
				workspaces = append(workspaces, workspace)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaces)
}

func GetCreatedWorkspaces(pubkey string) []db.Workspace {
	workspaces := db.DB.GetUserCreatedWorkspaces(pubkey)
	// add bounty count to the workspace
	for index, value := range workspaces {
		uuid := value.Uuid
		bountyCount := db.DB.GetWorkspaceBountyCount(uuid)
		hasRole := db.UserHasAccess(pubkey, uuid, db.ViewReport)

		if hasRole {
			budget := db.DB.GetWorkspaceBudget(uuid)
			workspaces[index].Budget = budget.TotalBudget
		} else {
			workspaces[index].Budget = 0
		}
		workspaces[index].BountyCount = bountyCount
	}
	return workspaces
}

func (oh *workspaceHandler) GetCreatedWorkspaces(pubkey string) []db.Workspace {
	workspaces := oh.db.GetUserCreatedWorkspaces(pubkey)
	// add bounty count to the workspace
	for index, value := range workspaces {
		uuid := value.Uuid
		bountyCount := oh.db.GetWorkspaceBountyCount(uuid)
		hasRole := oh.configUserHasAccess(pubkey, uuid, db.ViewReport)

		if hasRole {
			budget := oh.db.GetWorkspaceBudget(uuid)
			workspaces[index].Budget = budget.TotalBudget
		} else {
			workspaces[index].Budget = 0
		}
		workspaces[index].BountyCount = bountyCount
	}
	return workspaces
}

// GetWorkspaceBounties godoc
//
//	@Summary		Get Workspace Bounties
//	@Description	Get bounties of a workspace by its UUID
//	@Tags			Workspace -  Bounties
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Success		200		{array}	db.BountyResponse
//	@Router			/workspaces/bounties/{uuid} [get]
func (oh *workspaceHandler) GetWorkspaceBounties(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	// get the workspace bounties
	workspaceBounties := oh.db.GetWorkspaceBounties(r, uuid)

	var bountyResponse []db.BountyResponse = oh.generateBountyHandler(workspaceBounties)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyResponse)
}

// GetWorkspaceBountiesCount godoc
//
//	@Summary		Get Workspace Bounties Count
//	@Description	Get the count of bounties in a workspace
//	@Tags			Workspace -  Bounties
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Success		200		{int}	int
//	@Router			/workspaces/bounties/{uuid}/count [get]
func (oh *workspaceHandler) GetWorkspaceBountiesCount(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	workspaceBountiesCount := oh.db.GetWorkspaceBountiesCount(r, uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceBountiesCount)
}

// GetWorkspaceBudget godoc
//
//	@Summary		Get Workspace Budget
//	@Description	Get the budget of a workspace by its UUID
//	@Tags			Workspace -  Payments
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string	true	"Workspace UUID"
//	@Success		200		{object}	db.StatusBudget
//	@Router			/workspaces/budget/{uuid} [get]
func (oh *workspaceHandler) GetWorkspaceBudget(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if not the workspace admin
	hasRole := oh.userHasAccess(pubKeyFromAuth, uuid, db.ViewReport)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to view budget")
		return
	}

	// get the workspace budget
	workspaceBudget := oh.db.GetWorkspaceStatusBudget(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceBudget)
}

// GetWorkspaceBudgetHistory godoc
//
//	@Summary		Get Workspace Budget History
//	@Description	Get the budget history of a workspace by its UUID
//	@Tags			Workspace -  Payments
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Success		200		{array}	db.BudgetHistoryData
//	@Router			/workspaces/budget/history/{uuid} [get]
func (oh *workspaceHandler) GetWorkspaceBudgetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	// if not the workspace admin
	hasRole := oh.userHasAccess(pubKeyFromAuth, uuid, db.ViewReport)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to view budget history")
		return
	}

	// get the workspace budget
	workspaceBudget := oh.db.GetWorkspaceBudgetHistory(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceBudget)
}

// GetPaymentHistory godoc
//
//	@Summary		Get Payment History
//	@Description	Get the payment history of a workspace by its UUID
//	@Tags			Workspace -  Payments
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Success		200		{array}	db.PaymentHistoryData
//	@Router			/workspaces/payments/{uuid} [get]
func GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if not the workspace admin
	hasRole := db.UserHasAccess(pubKeyFromAuth, uuid, db.ViewReport)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Don't have access to view payments")
		return
	}

	// get the workspace payment history
	paymentHistory := db.DB.GetPaymentHistory(uuid, r)
	paymentHistoryData := []db.PaymentHistoryData{}

	for _, payment := range paymentHistory {
		sender := db.DB.GetPersonByPubkey(payment.SenderPubKey)
		receiver := db.DB.GetPersonByPubkey(payment.ReceiverPubKey)
		paymentData := db.PaymentHistoryData{
			NewPaymentHistory: payment,
			SenderName:        sender.UniqueName,
			SenderImg:         sender.Img,
			ReceiverName:      receiver.UniqueName,
			ReceiverImg:       receiver.Img,
		}
		paymentHistoryData = append(paymentHistoryData, paymentData)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paymentHistoryData)
}

// UpdateWorkspacePendingPayments godoc
//
//	@Summary		Update Workspace Pending Payments
//	@Description	Update pending payments of a workspace by its UUID
//	@Tags			Workspace -  Payments
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace_uuid	path		string	true	"Workspace UUID"
//	@Success		200				{string}	string	"Updated Payments Successfully"
//	@Router			/workspaces/{workspace_uuid}/payments [put]
func UpdateWorkspacePendingPayments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	workspace_uuid := chi.URLParam(r, "workspace_uuid")

	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	paymentsHistory := db.DB.GetWorkspacePendingPayments(workspace_uuid)

	for _, payment := range paymentsHistory {
		tag := payment.Tag
		tagResult := GetInvoiceStatusByTag(tag)

		if tagResult.Status == db.PaymentComplete {
			db.DB.SetPaymentAsComplete(tag)

			bounty := db.DB.GetBounty(payment.ID)

			if bounty.ID > 0 {
				now := time.Now()

				bounty.Paid = true
				bounty.PaymentPending = false
				bounty.PaymentFailed = false
				bounty.PaidDate = &now
				bounty.Completed = true
				bounty.CompletionDate = &now

				db.DB.UpdateBounty(bounty)
			}
		} else if tagResult.Status == db.PaymentFailed {
			// Handle failed payments
			bounty := db.DB.GetBounty(payment.ID)

			if bounty.ID > 0 {
				db.DB.SetPaymentStatusByBountyId(bounty.ID, tagResult)

				bounty.Paid = false
				bounty.PaymentPending = false
				bounty.PaymentFailed = true

				db.DB.UpdateBounty(bounty)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated Payments Successfully")
}

// PollBudgetInvoices godoc
//
//	@Summary		Poll Budget Invoices
//	@Description	Poll budget invoices of a workspace by its UUID
//	@Tags			Workspace -  Payments
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string	true	"Workspace UUID"
//	@Success		200		{string}	string	"Polled invoices"
//	@Router			/workspaces/poll/invoices/{uuid} [get]
func (oh *workspaceHandler) PollBudgetInvoices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workInvoices := oh.db.GetWorkspaceInvoices(uuid)

	for _, inv := range workInvoices {
		invoiceRes, invoiceErr := oh.getLightningInvoice(inv.PaymentRequest)

		if invoiceErr.Error != "" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(invoiceErr)
			return
		}

		if invoiceRes.Response.Settled {
			if !inv.Status && inv.Type == "BUDGET" {
				oh.db.ProcessUpdateBudget(inv)
			}
		} else {
			// Cheeck if time has expired
			isInvoiceExpired := utils.GetInvoiceExpired(inv.PaymentRequest)
			// If the invoice has expired and it is not paid delete from the DB
			if isInvoiceExpired {
				oh.db.DeleteInvoice(inv.PaymentRequest)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Polled invoices")
}

// PollUserWorkspacesBudget godoc
//
//	@Summary		Poll User Workspaces Budget
//	@Description	Poll budget invoices of all workspaces of a user
//	@Tags			Workspace -  Payments
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{string}	string	"Polled user workspace invoices"
//	@Router			/workspaces/poll/user/invoices [get]
func (oh *workspaceHandler) PollUserWorkspacesBudget(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get the user workspaces
	workspaces := oh.getAllUserWorkspaces(pubKeyFromAuth)
	// loop through the worksppaces and get each workspace invoice
	for _, space := range workspaces {

		// get all workspace invoice
		workInvoices := oh.db.GetWorkspaceInvoices(space.Uuid)

		for _, inv := range workInvoices {
			invoiceRes, invoiceErr := oh.getLightningInvoice(inv.PaymentRequest)

			if invoiceErr.Error != "" {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(invoiceErr)
				return
			}

			if invoiceRes.Response.Settled {
				if !inv.Status && inv.Type == "BUDGET" {
					oh.db.ProcessUpdateBudget(inv)
				}
			} else {
				// Cheeck if time has expired
				isInvoiceExpired := utils.GetInvoiceExpired(inv.PaymentRequest)
				// If the invoice has expired and it is not paid delete from the DB
				if isInvoiceExpired {
					oh.db.DeleteInvoice(inv.PaymentRequest)
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Polled user workspace invoices")
}

// GetInvoicesCount godoc
//
//	@Summary		Get Invoices Count
//	@Description	Get the count of invoices in a workspace by its UUID
//	@Tags			Workspace -  Payments
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Success		200		{int}	int
//	@Router			/workspaces/invoices/count/{uuid} [get]
func GetInvoicesCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	invoiceCount := db.DB.GetWorkspaceInvoicesCount(uuid)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceCount)
}

// GetAllUserInvoicesCount godoc
//
//	@Summary		Get All User Invoices Count
//	@Description	Get the count of all invoices of a user
//	@Tags			Workspace -  Payments
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{int}	int
//	@Router			/workspaces/user/invoices/count [get]
func GetAllUserInvoicesCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	allCount := int64(0)
	workspaces := GetAllUserWorkspaces(pubKeyFromAuth)
	for _, space := range workspaces {
		invoiceCount := db.DB.GetWorkspaceInvoicesCount(space.Uuid)
		allCount += invoiceCount
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allCount)
}

// DeleteWorkspace godoc
//
//	@Summary		Delete Workspace
//	@Description	Delete a workspace by its UUID
//	@Tags			Workspaces
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	db.Workspace
//	@Router			/workspaces/delete/{uuid} [delete]
func (oh *workspaceHandler) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "uuid")

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspace := oh.db.GetWorkspaceByUuid(uuid)
	if pubKeyFromAuth != workspace.OwnerPubKey {
		msg := "only workspace admin can delete an workspace"
		logger.Log.Info("[workspaces] %s", msg)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msg)
		return
	}

	// Soft delete Workspace and delete user data
	if err := oh.db.ProcessDeleteWorkspace(uuid); err != nil {
		msg := "Error removing users from workspace"
		logger.Log.Error("%s: %v", msg, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspace)
}

// UpdateWorkspace godoc
//
//	@Summary		Update Workspace
//	@Description	Update a workspace
//	@Tags			Workspaces
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace	body		db.Workspace	true	"Workspace"
//	@Success		200			{object}	db.Workspace
//	@Router			/workspaces/mission [post]
//	@Router			/workspaces/tactics [post]
//	@Router			/workspaces/schematicurl [post]
func (oh *workspaceHandler) UpdateWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspace := db.Workspace{}
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()
	err := json.Unmarshal(body, &workspace)

	if err != nil {
		logger.Log.Error("[workspaces] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if pubKeyFromAuth != workspace.OwnerPubKey {
		hasRole := db.UserHasAccess(pubKeyFromAuth, workspace.Uuid, db.EditOrg)
		if !hasRole {
			logger.Log.Info("[workspaces] mismatched pubkey")
			logger.Log.Info("Auth Pubkey: %s", pubKeyFromAuth)
			logger.Log.Info("OwnerPubKey: %s", workspace.OwnerPubKey)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Don't have access to Edit workspace")
			return
		}
	}

	// Validate struct data
	err = db.Validate.Struct(workspace)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error: did not pass validation test : %s", err)
		json.NewEncoder(w).Encode(msg)
		return
	}

	p, err := oh.db.CreateOrEditWorkspace(workspace)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// CreateOrEditWorkspaceRepository godoc
//
//	@Summary		Create or Edit Workspace Repository
//	@Description	Create or edit a repository for a workspace
//	@Tags			Workspace -  Repositories
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspaceRepo	body		db.WorkspaceRepositories	true	"Workspace Repository"
//	@Success		200				{object}	db.WorkspaceRepositories
//	@Router			/workspaces/repositories [post]
func (oh *workspaceHandler) CreateOrEditWorkspaceRepository(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspaceRepo := db.WorkspaceRepositories{}
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()
	err := json.Unmarshal(body, &workspaceRepo)

	if err != nil {
		logger.Log.Error("[workspaces] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if len(workspaceRepo.Uuid) == 0 {
		workspaceRepo.Uuid = xid.New().String()
		workspaceRepo.CreatedBy = pubKeyFromAuth
	}

	workspaceRepo.UpdatedBy = pubKeyFromAuth

	if workspaceRepo.Uuid == "" {
		workspaceRepo.Uuid = xid.New().String()
	} else {
		workspaceRepo.UpdatedBy = pubKeyFromAuth
	}

	// Validate struct data
	err = db.Validate.Struct(workspaceRepo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error: did not pass validation test : %s", err)
		json.NewEncoder(w).Encode(msg)
		return
	}

	// Check if workspace exists
	workpace := oh.db.GetWorkspaceByUuid(workspaceRepo.WorkspaceUuid)
	if workpace.Uuid != workspaceRepo.WorkspaceUuid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Workspace does not exists")
		return
	}

	p, err := oh.db.CreateOrEditWorkspaceRepository(workspaceRepo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// GetWorkspaceRepositorByWorkspaceUuid godoc
//
//	@Summary		Get Workspace Repositories
//	@Description	Get repositories of a workspace by its UUID
//	@Tags			Workspace -  Repositories
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"Workspace UUID"
//	@Success		200		{array}	db.WorkspaceRepositories
//	@Router			/workspaces/repositories/{uuid} [get]
func (oh *workspaceHandler) GetWorkspaceRepositorByWorkspaceUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	workspaceFeatures := oh.db.GetWorkspaceRepositorByWorkspaceUuid(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceFeatures)
}

// GetWorkspaceRepoByWorkspaceUuidAndRepoUuid godoc
//
//	@Summary		Get Workspace Repository by UUID
//	@Description	Get a repository of a workspace by its UUID
//	@Tags			Workspace -  Repositories
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace_uuid	path		string	true	"Workspace UUID"
//	@Param			uuid			path		string	true	"Repository UUID"
//	@Success		200				{object}	db.WorkspaceRepositories
//	@Router			/workspaces/repository/{uuid} [get]
func (oh *workspaceHandler) GetWorkspaceRepoByWorkspaceUuidAndRepoUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspace_uuid := chi.URLParam(r, "workspace_uuid")
	uuid := chi.URLParam(r, "uuid")
	WorkspaceRepository, err := oh.db.GetWorkspaceRepoByWorkspaceUuidAndRepoUuid(workspace_uuid, uuid)
	if err != nil {
		logger.Log.Error("[workspaces] workspace repository not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Repository not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(WorkspaceRepository)
}

// DeleteWorkspaceRepository godoc
//
//	@Summary		Delete Workspace Repository
//	@Description	Delete a repository from a workspace by its UUID
//	@Tags			Workspace -  Repositories
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace_uuid	path		string	true	"Workspace UUID"
//	@Param			uuid			path		string	true	"Repository UUID"
//	@Success		200				{string}	string	"Repository deleted successfully"
//	@Router			/workspaces/repository/{uuid} [delete]
func (oh *workspaceHandler) DeleteWorkspaceRepository(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspace_uuid := chi.URLParam(r, "workspace_uuid")
	uuid := chi.URLParam(r, "uuid")

	oh.db.DeleteWorkspaceRepository(workspace_uuid, uuid)

	w.WriteHeader(http.StatusOK)
}

func isValidUUID(uuid string) bool {

	regexPattern := `^[a-zA-Z0-9\-]+$`
	rgx := regexp.MustCompile(regexPattern)
	return rgx.MatchString(uuid) && len(uuid) > 0
}

// GetFeaturesByWorkspaceUuid godoc
//
//	@Summary		Get Features by Workspace UUID
//	@Description	Get features of a workspace by its UUID
//	@Tags			Workspaces
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace_uuid	path	string	true	"Workspace UUID"
//	@Success		200				{array}	db.WorkspaceFeatures
//	@Router			/workspaces/{workspace_uuid}/features [get]
func (oh *workspaceHandler) GetFeaturesByWorkspaceUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "workspace_uuid")

	if uuid == "" {
		logger.Log.Info("workspace_uuid parameter is missing")
		http.Error(w, "Missing workspace_uuid parameter", http.StatusBadRequest)
		return
	}

	if !isValidUUID(uuid) {
		logger.Log.Info("invalid UUID format or contains special characters")
		http.Error(w, "Invalid UUID format or contains special characters", http.StatusBadRequest)
		return
	}

	workspaceFeatures := oh.db.GetFeaturesByWorkspaceUuid(uuid, r)

	for i, feature := range workspaceFeatures {
		phases := oh.db.GetPhasesByFeatureUuid(feature.Uuid)
		var totalCompleted, totalAssigned, totalOpen int

		for _, phase := range phases {
			completed := oh.db.GetFeaturePhasesBountiesCount("completed", phase.Uuid)
			assigned := oh.db.GetFeaturePhasesBountiesCount("assigned", phase.Uuid)
			open := oh.db.GetFeaturePhasesBountiesCount("open", phase.Uuid)

			totalAssigned += int(assigned)
			totalCompleted += int(completed)
			totalOpen += int(open)
		}

		workspaceFeatures[i].BountiesCountCompleted = totalCompleted
		workspaceFeatures[i].BountiesCountAssigned = totalAssigned
		workspaceFeatures[i].BountiesCountOpen = totalOpen
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceFeatures)
}

// GetLastWithdrawal godoc
//
//	@Summary		Get Last Withdrawal
//	@Description	Get the last withdrawal of a workspace by its UUID
//	@Tags			Workspace -  Payments
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace_uuid	path	string	true	"Workspace UUID"
//	@Success		200				{int}	int		"Hours since last withdrawal"
//	@Router			/workspaces/{workspace_uuid}/lastwithdrawal [get]
func (oh *workspaceHandler) GetLastWithdrawal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspace_uuid := chi.URLParam(r, "workspace_uuid")
	lastWithdrawal := oh.db.GetLastWithdrawal(workspace_uuid)

	log.Println("This workspaces last withdrawal is", workspace_uuid, lastWithdrawal)

	hoursDiff := int64(1)

	if lastWithdrawal.ID > 0 {
		now := time.Now()
		withdrawCreated := lastWithdrawal.Created
		withdrawTime := utils.ConvertTimeToTimestamp(withdrawCreated.String())

		hoursDiff = utils.GetHoursDifference(int64(withdrawTime), &now)
		log.Println("This workspaces last withdrawal hours difference is", hoursDiff)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hoursDiff)
}

func GetAllUserWorkspaces(pubkey string) []db.Workspace {
	// get the workspaces created by the user, then get all the workspaces
	// the user has been added to, loop through to get the workspace
	workspaces := GetCreatedWorkspaces(pubkey)
	assignedWorkspaces := db.DB.GetUserAssignedWorkspaces(pubkey)
	for _, value := range assignedWorkspaces {
		uuid := value.WorkspaceUuid
		workspace := db.DB.GetWorkspaceByUuid(uuid)
		bountyCount := db.DB.GetWorkspaceBountyCount(uuid)
		hasRole := db.UserHasAccess(pubkey, uuid, db.ViewReport)

		// don't add workspace to the list if user is the owner of the workspace
		alreadyAdded := false
		if workspace.OwnerPubKey == pubkey {
			alreadyAdded = true
		}

		// don't add deleted workspaces to the list
		if !workspace.Deleted {
			if hasRole {
				budget := db.DB.GetWorkspaceBudget(uuid)
				workspace.Budget = budget.TotalBudget
			} else {
				workspace.Budget = 0
			}
			workspace.BountyCount = bountyCount

			// check if workspace has already been added to the list
			for _, existingWorkspace := range workspaces {
				if existingWorkspace.Uuid == workspace.Uuid {
					alreadyAdded = true
				}
			}

			if !alreadyAdded {
				workspaces = append(workspaces, workspace)
			}
		}
	}

	return workspaces
}

// CreateOrEditWorkspaceCodeGraph godoc
//
//	@Summary		Create or Edit Workspace Code Graph
//	@Description	Create or edit a code graph for a workspace
//	@Tags			Workspace -  Code Graph
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			codeGraph	body		db.WorkspaceCodeGraph	true	"Workspace Code Graph"
//	@Success		200			{object}	db.WorkspaceCodeGraph
//	@Router			/workspaces/codegraph [post]
func (oh *workspaceHandler) CreateOrEditWorkspaceCodeGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	codeGraph := db.WorkspaceCodeGraph{}
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()
	err := json.Unmarshal(body, &codeGraph)

	if err != nil {
		logger.Log.Error("[workspaces] %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if len(codeGraph.Uuid) == 0 {
		codeGraph.Uuid = xid.New().String()
		codeGraph.CreatedBy = pubKeyFromAuth
	}
	codeGraph.UpdatedBy = pubKeyFromAuth

	err = db.Validate.Struct(codeGraph)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error: did not pass validation test : %s", err)
		json.NewEncoder(w).Encode(msg)
		return
	}

	workspace := oh.db.GetWorkspaceByUuid(codeGraph.WorkspaceUuid)
	if workspace.Uuid != codeGraph.WorkspaceUuid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Workspace does not exist")
		return
	}

	p, err := oh.db.CreateOrEditCodeGraph(codeGraph)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// GetWorkspaceCodeGraphByUUID godoc
//
//	@Summary		Get Workspace Code Graph by UUID
//	@Description	Get a code graph of a workspace by its UUID
//	@Tags			Workspace -  Code Graph
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string	true	"Code Graph UUID"
//	@Success		200		{object}	db.WorkspaceCodeGraph
//	@Router			/workspaces/codegraph/{uuid} [get]
func (oh *workspaceHandler) GetWorkspaceCodeGraphByUUID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	codeGraph, err := oh.db.GetCodeGraphByUUID(uuid)
	if err != nil {
		logger.Log.Error("[workspaces] code graph not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Code graph not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codeGraph)
}

// GetCodeGraphByWorkspaceUuid godoc
//
//	@Summary		Get Code Graph by Workspace UUID
//	@Description	Get code graphs of a workspace by its UUID
//	@Tags			Workspace -  Code Graph
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace_uuid	path	string	true	"Workspace UUID"
//	@Success		200				{array}	db.WorkspaceCodeGraph
//	@Router			/workspaces/{workspace_uuid}/codegraph [get]
func (oh *workspaceHandler) GetCodeGraphByWorkspaceUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspace_uuid := chi.URLParam(r, "workspace_uuid")
	codeGraph, err := oh.db.GetCodeGraphByWorkspaceUuid(workspace_uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get code graphs"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codeGraph)
}

// DeleteWorkspaceCodeGraph godoc
//
//	@Summary		Delete Workspace Code Graph
//	@Description	Delete a code graph from a workspace by its UUID
//	@Tags			Workspace -  Code Graph
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace_uuid	path		string	true	"Workspace UUID"
//	@Param			uuid			path		string	true	"Code Graph UUID"
//	@Success		200				{string}	string	"Code graph deleted successfully"
//	@Router			/workspaces/{workspace_uuid}/codegraph/{uuid} [delete]
func (oh *workspaceHandler) DeleteWorkspaceCodeGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspace_uuid := chi.URLParam(r, "workspace_uuid")
	uuid := chi.URLParam(r, "uuid")

	err := oh.db.DeleteCodeGraph(workspace_uuid, uuid)
	if err != nil {
		if err.Error() == "code graph not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Code graph not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete code graph"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Code graph deleted successfully"})
}

func (oh *workspaceHandler) RefreshCodeGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[workspaces] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspaceUuid := chi.URLParam(r, "workspace_uuid")

	codeGraph, err := oh.db.GetCodeGraphByWorkspaceUuid(workspaceUuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get code graph"})
		return
	}

	repos := oh.db.GetWorkspaceRepositorByWorkspaceUuid(workspaceUuid)

	repoURLs := []string{}
	for _, repo := range repos {
		repoURLs = append(repoURLs, repo.Url)
	}
	reposStr := strings.Join(repoURLs, ",")

	requestURL := fmt.Sprintf("%s/git/sync?source_link=%s", codeGraph.Url, reposStr)

	response, err := http.Get(requestURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to sync repositories"})
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read response"})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetWorkspaceEnvVars proxies env var fetch to 3rd party
func (oh *workspaceHandler) GetWorkspaceEnvVars(w http.ResponseWriter, r *http.Request) {
	workspaceUUID := chi.URLParam(r, "workspace_uuid")
	codespaces, err := oh.db.GetCodeSpaceMapByWorkspace(workspaceUUID)
	if err != nil || len(codespaces) == 0 || codespaces[0].CodeSpaceURL == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "codespaceURL not found for workspace"})
		return
	}
	codespaceURL := codespaces[0].CodeSpaceURL
	url := "https://workspaces.sphinx.chat/api/pools/" + codespaceURL

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+codespaces[0].PoolAPIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to contact 3rd party service"})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}
	var result struct {
		Config struct {
			EnvVars []map[string]interface{} `json:"env_vars"`
		} `json:"config"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to decode response"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result.Config.EnvVars)
}

// UpdateWorkspaceEnvVars proxies env var update to 3rd party
func (oh *workspaceHandler) UpdateWorkspaceEnvVars(w http.ResponseWriter, r *http.Request) {
	workspaceUUID := chi.URLParam(r, "workspace_uuid")
	codespaces, err := oh.db.GetCodeSpaceMapByWorkspace(workspaceUUID)
	if err != nil || len(codespaces) == 0 || codespaces[0].CodeSpaceURL == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "codespaceURL not found for workspace"})
		return
	}
	codespaceURL := codespaces[0].CodeSpaceURL
	url := "https://workspaces.sphinx.chat/api/pools/" + codespaceURL

	var body struct {
		EnvVars []map[string]string `json:"env_vars"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}
	b, _ := json.Marshal(map[string]interface{}{"env_vars": body.EnvVars})
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(b)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+codespaces[0].PoolAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to contact 3rd party service"})
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
