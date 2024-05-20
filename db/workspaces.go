package db

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/stakwork/sphinx-tribes/utils"
)

func (db database) GetWorkspaces(r *http.Request) []Workspace {
	ms := []Workspace{}
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)

	query := db.db.Model(&ms).Where("LOWER(name) LIKE ?", "%"+search+"%").Where("deleted != ?", true)

	if limit > 1 {
		query.Offset(offset).Limit(limit).Order(sortBy + " " + direction + " ")
	}

	query.Find(&ms)
	return ms
}

func (db database) GetWorkspacesCount() int64 {
	var count int64
	db.db.Model(&Workspace{}).Count(&count)
	return count
}

func (db database) GetWorkspaceByUuid(uuid string) Workspace {
	ms := Workspace{}

	db.db.Model(&Workspace{}).Where("uuid = ?", uuid).Find(&ms)

	return ms
}

func (db database) GetWorkspaceByName(name string) Workspace {
	ms := Workspace{}

	db.db.Model(&Workspace{}).Where("name = ?", name).Find(&ms)

	return ms
}

func (db database) CreateOrEditWorkspace(m Workspace) (Workspace, error) {
	if m.OwnerPubKey == "" {
		return Workspace{}, errors.New("no pub key")
	}

	if db.db.Model(&m).Where("uuid = ?", m.Uuid).Updates(&m).RowsAffected == 0 {
		db.db.Create(&m)
	}

	return m, nil
}

func (db database) CreateOrEditWorkspaceRepository(m WorkspaceRepositories) (WorkspaceRepositories, error) {
	m.Name = strings.TrimSpace(m.Name)
	m.Url = strings.TrimSpace(m.Url)

	now := time.Now()
	m.Updated = &now

	if db.db.Model(&m).Where("uuid = ?", m.Uuid).Updates(&m).RowsAffected == 0 {
		m.Created = &now
		db.db.Create(&m)
	}

	db.db.Model(&WorkspaceRepositories{}).Where("uuid = ?", m.Uuid).Find(&m)

	return m, nil
}

func (db database) GetWorkspaceRepositorByWorkspaceUuid(uuid string) []WorkspaceRepositories {
	ms := []WorkspaceRepositories{}

	db.db.Model(&WorkspaceRepositories{}).Where("workspace_uuid = ?", uuid).Order("Created").Find(&ms)

	return ms
}

func (db database) GetWorkspaceRepoByWorkspaceUuidAndRepoUuid(workspace_uuid string, uuid string) (WorkspaceRepositories, error) {
	var ms WorkspaceRepositories

	result := db.db.Model(&WorkspaceRepositories{}).Where("workspace_uuid = ?", workspace_uuid).Where("uuid = ?", uuid).Find(&ms)
	if result.RowsAffected == 0 {
		return ms, fmt.Errorf("workspace repository not found")
	}

	return ms, nil
}

func (db database) DeleteWorkspaceRepository(workspace_uuid string, uuid string) bool {
	db.db.Where("workspace_uuid = ?", workspace_uuid).Where("uuid = ?", uuid).Delete(&WorkspaceRepositories{})
	return true
}

func (db database) GetWorkspaceUsers(uuid string) ([]WorkspaceUsersData, error) {
	ms := []WorkspaceUsersData{}

	err := db.db.Raw(`SELECT org.workspace_uuid, org.created as user_created, person.* FROM public.Workspace_users AS org LEFT OUTER JOIN public.people AS person ON org.owner_pub_key = person.owner_pub_key WHERE org.workspace_uuid = '` + uuid + `' ORDER BY org.created DESC`).Find(&ms).Error

	return ms, err
}

func (db database) GetWorkspaceUsersCount(uuid string) int64 {
	var count int64
	db.db.Model(&WorkspaceUsers{}).Where("workspace_uuid  = ?", uuid).Count(&count)
	return count
}

func (db database) GetWorkspaceBountyCount(uuid string) int64 {
	var count int64
	db.db.Model(&Bounty{}).Where("workspace_uuid  = ?", uuid).Count(&count)
	return count
}

func (db database) GetWorkspaceUser(pubkey string, workspace_uuid string) WorkspaceUsers {
	ms := WorkspaceUsers{}
	db.db.Where("workspace_uuid = ?", workspace_uuid).Where("owner_pub_key = ?", pubkey).Find(&ms)
	return ms
}

func (db database) CreateWorkspaceUser(orgUser WorkspaceUsers) WorkspaceUsers {
	db.db.Create(&orgUser)

	return orgUser
}

func (db database) DeleteWorkspaceUser(orgUser WorkspaceUsersData, org string) WorkspaceUsersData {
	db.db.Where("owner_pub_key = ?", orgUser.OwnerPubKey).Where("workspace_uuid = ?", org).Delete(&WorkspaceUsers{})
	db.db.Where("owner_pub_key = ?", orgUser.OwnerPubKey).Where("workspace_uuid = ?", org).Delete(&UserRoles{})
	return orgUser
}

func (db database) GetBountyRoles() []BountyRoles {
	ms := []BountyRoles{}
	db.db.Find(&ms)
	return ms
}

func (db database) CreateUserRoles(roles []WorkspaceUserRoles, uuid string, pubkey string) []WorkspaceUserRoles {
	// delete roles and create new ones
	db.db.Where("workspace_uuid = ?", uuid).Where("owner_pub_key = ?", pubkey).Delete(&WorkspaceUserRoles{})
	db.db.Create(&roles)

	return roles
}

func (db database) GetUserRoles(uuid string, pubkey string) []WorkspaceUserRoles {
	ms := []WorkspaceUserRoles{}
	db.db.Where("workspace_uuid = ?", uuid).Where("owner_pub_key = ?", pubkey).Find(&ms)
	return ms
}

func (db database) GetUserCreatedWorkspaces(pubkey string) []Workspace {
	ms := []Workspace{}
	db.db.Where("owner_pub_key = ?", pubkey).Where("deleted != ?", true).Find(&ms)
	return ms
}

func (db database) GetUserAssignedWorkspaces(pubkey string) []WorkspaceUsers {
	ms := []WorkspaceUsers{}
	db.db.Where("owner_pub_key = ?", pubkey).Find(&ms)
	return ms
}

func (db database) AddBudgetHistory(budget BudgetHistory) BudgetHistory {
	db.db.Create(&budget)
	return budget
}

func (db database) CreateWorkspaceBudget(budget NewBountyBudget) NewBountyBudget {
	db.db.Create(&budget)
	return budget
}

func (db database) UpdateWorkspaceBudget(budget NewBountyBudget) NewBountyBudget {
	db.db.Model(&NewBountyBudget{}).Where("workspace_uuid = ?", budget.WorkspaceUuid).Updates(map[string]interface{}{
		"total_budget": budget.TotalBudget,
	})
	return budget
}

func (db database) GetPaymentHistoryByCreated(created *time.Time, workspace_uuid string) NewPaymentHistory {
	ms := NewPaymentHistory{}
	db.db.Model(&NewPaymentHistory{}).Where("created = ?", created).Where("workspace_uuid = ? ", workspace_uuid).Find(&ms)
	return ms
}

func (db database) GetWorkspaceBudget(workspace_uuid string) NewBountyBudget {
	ms := NewBountyBudget{}
	db.db.Model(&NewBountyBudget{}).Where("workspace_uuid = ?", workspace_uuid).Find(&ms)

	return ms
}

func (db database) GetWorkspaceStatusBudget(workspace_uuid string) StatusBudget {
	orgBudget := db.GetWorkspaceBudget(workspace_uuid)

	var openBudget uint
	db.db.Model(&NewBounty{}).Where("assignee = '' ").Where("paid != true").Select("SUM(price)").Row().Scan(&openBudget)

	var openCount int64
	db.db.Model(&NewBounty{}).Where("assignee = '' ").Where("paid != true").Count(&openCount)

	var openDifference int = int(orgBudget.TotalBudget - openBudget)

	var assignedBudget uint
	db.db.Model(&NewBounty{}).Where("assignee != '' ").Where("paid != true").Select("SUM(price)").Row().Scan(&assignedBudget)

	var assignedCount int64
	db.db.Model(&NewBounty{}).Where("assignee != '' ").Where("paid != true").Count(&assignedCount)

	var assignedDifference int = int(orgBudget.TotalBudget - assignedBudget)

	var completedBudget uint
	db.db.Model(&NewBounty{}).Where("completed = true ").Where("paid != true").Select("SUM(price)").Row().Scan(&completedBudget)

	var completedCount int64
	db.db.Model(&NewBounty{}).Where("completed = true ").Where("paid != true").Count(&completedCount)

	var completedDifference int = int(orgBudget.TotalBudget - completedBudget)

	statusBudget := StatusBudget{
		OrgUuid:             workspace_uuid,
		WorkspaceUuid:       workspace_uuid,
		CurrentBudget:       orgBudget.TotalBudget,
		OpenBudget:          openBudget,
		OpenCount:           openCount,
		OpenDifference:      openDifference,
		AssignedBudget:      assignedBudget,
		AssignedCount:       assignedCount,
		AssignedDifference:  assignedDifference,
		CompletedBudget:     completedBudget,
		CompletedCount:      completedCount,
		CompletedDifference: completedDifference,
	}

	return statusBudget
}

func (db database) GetWorkspaceBudgetHistory(workspace_uuid string) []BudgetHistoryData {
	budgetHistory := []BudgetHistoryData{}

	db.db.Raw(`SELECT budget.id, budget.workspace_uuid, budget.amount, budget.created, budget.updated, budget.payment_type, budget.status, budget.sender_pub_key, sender.unique_name AS sender_name FROM public.budget_histories AS budget LEFT OUTER JOIN public.people AS sender ON budget.sender_pub_key = sender.owner_pub_key WHERE budget.workspace_uuid = '` + workspace_uuid + `' ORDER BY budget.created DESC`).Find(&budgetHistory)
	return budgetHistory
}

func (db database) AddAndUpdateBudget(invoice NewInvoiceList) NewPaymentHistory {
	created := invoice.Created
	workspace_uuid := invoice.WorkspaceUuid

	paymentHistory := db.GetPaymentHistoryByCreated(created, workspace_uuid)

	if paymentHistory.WorkspaceUuid != "" && paymentHistory.Amount != 0 {
		paymentHistory.Status = true
		db.db.Where("created = ?", created).Where("workspace_uuid = ? ", workspace_uuid).Updates(paymentHistory)

		// get Workspace budget and add payment to total budget
		WorkspaceBudget := db.GetWorkspaceBudget(workspace_uuid)

		if WorkspaceBudget.WorkspaceUuid == "" {
			now := time.Now()
			workBudget := NewBountyBudget{
				WorkspaceUuid: workspace_uuid,
				TotalBudget:   paymentHistory.Amount,
				Created:       &now,
				Updated:       &now,
			}
			db.CreateWorkspaceBudget(workBudget)
		} else {
			totalBudget := WorkspaceBudget.TotalBudget
			WorkspaceBudget.TotalBudget = totalBudget + paymentHistory.Amount
			db.UpdateWorkspaceBudget(WorkspaceBudget)
		}
	}

	return paymentHistory
}

func (db database) WithdrawBudget(sender_pubkey string, workspace_uuid string, amount uint) {
	// get Workspace budget and add payment to total budget
	WorkspaceBudget := db.GetWorkspaceBudget(workspace_uuid)
	totalBudget := WorkspaceBudget.TotalBudget

	newBudget := totalBudget - amount
	db.db.Model(&NewBountyBudget{}).Where("workspace_uuid = ?", workspace_uuid).Updates(map[string]interface{}{
		"total_budget": newBudget,
	})

	now := time.Now()

	budgetHistory := NewPaymentHistory{
		WorkspaceUuid:  workspace_uuid,
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

func (db database) AddPaymentHistory(payment NewPaymentHistory) NewPaymentHistory {
	db.db.Create(&payment)

	// get Workspace budget and subtract payment from total budget
	WorkspaceBudget := db.GetWorkspaceBudget(payment.WorkspaceUuid)
	totalBudget := WorkspaceBudget.TotalBudget

	// deduct amount if it's a bounty payment
	if payment.PaymentType == "payment" {
		WorkspaceBudget.TotalBudget = totalBudget - payment.Amount
	}

	db.UpdateWorkspaceBudget(WorkspaceBudget)

	return payment
}

func (db database) GetPaymentHistory(workspace_uuid string, r *http.Request) []NewPaymentHistory {
	payment := []NewPaymentHistory{}

	offset, limit, _, _, _ := utils.GetPaginationParams(r)
	limitQuery := ""

	limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)

	query := `SELECT * FROM payment_histories WHERE workspace_uuid = '` + workspace_uuid + `' AND status = true ORDER BY created DESC`

	db.db.Raw(query + " " + limitQuery).Find(&payment)
	return payment
}

func (db database) GetWorkspaceInvoices(workspace_uuid string) []NewInvoiceList {
	ms := []NewInvoiceList{}
	db.db.Where("workspace_uuid = ?", workspace_uuid).Where("status", false).Find(&ms)
	return ms
}

func (db database) GetWorkspaceInvoicesCount(workspace_uuid string) int64 {
	var count int64
	ms := NewInvoiceList{}

	db.db.Model(&ms).Where("workspace_uuid = ?", workspace_uuid).Where("status", false).Count(&count)
	return count
}

func (db database) ChangeWorkspaceDeleteStatus(workspace_uuid string, status bool) Workspace {
	ms := Workspace{}
	db.db.Model(&ms).Where("uuid", workspace_uuid).Updates(map[string]interface{}{
		"deleted": status,
	})
	return ms
}

func (db database) UpdateWorkspaceForDeletion(uuid string) error {
	updates := map[string]interface{}{
		"website":     "",
		"github":      "",
		"description": "",
		"show":        false,
	}

	result := db.db.Model(&Workspace{}).Where("uuid = ?", uuid).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db database) DeleteAllUsersFromWorkspace(org string) error {
	if org == "" {
		return errors.New("no org uuid provided")
	}

	// Delete all users associated with the Workspace
	result := db.db.Where("workspace_uuid = ?", org).Delete(&WorkspaceUsers{})
	if result.Error != nil {
		return result.Error
	}

	// Delete all user roles associated with the Workspace
	result = db.db.Where("workspace_uuid = ?", org).Delete(&UserRoles{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
