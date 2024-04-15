package db

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/stakwork/sphinx-tribes/utils"
)

func (db database) GetWorkspaces(r *http.Request) []Organization {
	ms := []Organization{}
	offset, limit, sortBy, direction, search := utils.GetPaginationParams(r)

	// return if like owner_alias, unique_name, or equals pubkey
	db.db.Offset(offset).Limit(limit).Order(sortBy+" "+direction+" ").Where("LOWER(name) LIKE ?", "%"+search+"%").Where("deleted != ?", false).Find(&ms)
	return ms
}

func (db database) GetWorkspacesCount() int64 {
	var count int64
	db.db.Model(&Organization{}).Count(&count)
	return count
}

func (db database) GetWorkspaceByUuid(uuid string) Organization {
	ms := Organization{}

	db.db.Model(&Organization{}).Where("uuid = ?", uuid).Find(&ms)

	return ms
}

func (db database) GetWorkspaceByName(name string) Organization {
	ms := Organization{}

	db.db.Model(&Organization{}).Where("name = ?", name).Find(&ms)

	return ms
}

func (db database) CreateOrEditWorkspace(m Organization) (Organization, error) {
	if m.OwnerPubKey == "" {
		return Organization{}, errors.New("no pub key")
	}

	if db.db.Model(&m).Where("uuid = ?", m.Uuid).Updates(&m).RowsAffected == 0 {
		db.db.Create(&m)
	}

	return m, nil
}

func (db database) GetWorkspaceUsers(uuid string) ([]OrganizationUsersData, error) {
	ms := []OrganizationUsersData{}

	err := db.db.Raw(`SELECT org.org_uuid, org.created as user_created, person.* FROM public.organization_users AS org LEFT OUTER JOIN public.people AS person ON org.owner_pub_key = person.owner_pub_key WHERE org.org_uuid = '` + uuid + `' ORDER BY org.created DESC`).Find(&ms).Error

	return ms, err
}

func (db database) GetWorkspaceUsersCount(uuid string) int64 {
	var count int64
	db.db.Model(&OrganizationUsers{}).Where("org_uuid  = ?", uuid).Count(&count)
	return count
}

func (db database) GetWorkspaceBountyCount(uuid string) int64 {
	var count int64
	db.db.Model(&Bounty{}).Where("org_uuid  = ?", uuid).Count(&count)
	return count
}

func (db database) GetWorkspaceUser(pubkey string, org_uuid string) OrganizationUsers {
	ms := OrganizationUsers{}
	db.db.Where("org_uuid = ?", org_uuid).Where("owner_pub_key = ?", pubkey).Find(&ms)
	return ms
}

func (db database) CreateWorkspaceUser(orgUser OrganizationUsers) OrganizationUsers {
	db.db.Create(&orgUser)

	return orgUser
}

func (db database) DeleteWorkspaceUser(orgUser OrganizationUsersData, org string) OrganizationUsersData {
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

func (db database) GetUserCreatedWorkspaces(pubkey string) []Organization {
	ms := []Organization{}
	db.db.Where("owner_pub_key = ?", pubkey).Where("deleted != ?", true).Find(&ms)
	return ms
}

func (db database) GetUserAssignedWorkspaces(pubkey string) []OrganizationUsers {
	ms := []OrganizationUsers{}
	db.db.Where("owner_pub_key = ?", pubkey).Find(&ms)
	return ms
}

func (db database) AddBudgetHistory(budget BudgetHistory) BudgetHistory {
	db.db.Create(&budget)
	return budget
}

func (db database) CreateWorkspaceBudget(budget BountyBudget) BountyBudget {
	db.db.Create(&budget)
	return budget
}

func (db database) UpdateWorkspaceBudget(budget BountyBudget) BountyBudget {
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

func (db database) GetWorkspaceBudget(org_uuid string) BountyBudget {
	ms := BountyBudget{}
	db.db.Where("org_uuid = ?", org_uuid).Find(&ms)

	return ms
}

func (db database) GetWorkspaceStatusBudget(org_uuid string) StatusBudget {

	orgBudget := db.GetWorkspaceBudget(org_uuid)

	var openBudget uint
	db.db.Model(&Bounty{}).Where("assignee = '' ").Where("paid != true").Select("SUM(price)").Row().Scan(&openBudget)

	var assignedBudget uint
	db.db.Model(&Bounty{}).Where("assignee != '' ").Where("paid != true").Select("SUM(price)").Row().Scan(&assignedBudget)

	var completedBudget uint
	db.db.Model(&Bounty{}).Where("completed = true ").Where("paid != true").Select("SUM(price)").Row().Scan(&completedBudget)

	statusBudget := StatusBudget{
		OrgUuid:         org_uuid,
		CurrentBudget:   orgBudget.TotalBudget,
		OpenBudget:      openBudget,
		AssignedBudget:  assignedBudget,
		CompletedBudget: completedBudget,
	}

	return statusBudget
}

func (db database) GetWorkspaceBudgetHistory(org_uuid string) []BudgetHistoryData {
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
		organizationBudget := db.GetWorkspaceBudget(org_uuid)

		if organizationBudget.OrgUuid == "" {
			now := time.Now()
			orgBudget := BountyBudget{
				OrgUuid:     org_uuid,
				TotalBudget: paymentHistory.Amount,
				Created:     &now,
				Updated:     &now,
			}
			db.CreateWorkspaceBudget(orgBudget)
		} else {
			totalBudget := organizationBudget.TotalBudget
			organizationBudget.TotalBudget = totalBudget + paymentHistory.Amount
			db.UpdateWorkspaceBudget(organizationBudget)
		}
	}

	return paymentHistory
}

func (db database) WithdrawBudget(sender_pubkey string, org_uuid string, amount uint) {
	// get organization budget and add payment to total budget
	organizationBudget := db.GetWorkspaceBudget(org_uuid)
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
	organizationBudget := db.GetWorkspaceBudget(payment.OrgUuid)
	totalBudget := organizationBudget.TotalBudget

	// deduct amount if it's a bounty payment
	if payment.PaymentType == "payment" {
		organizationBudget.TotalBudget = totalBudget - payment.Amount
	}

	db.UpdateWorkspaceBudget(organizationBudget)

	return payment
}

func (db database) GetPaymentHistory(org_uuid string, r *http.Request) []PaymentHistory {
	payment := []PaymentHistory{}

	offset, limit, _, _, _ := utils.GetPaginationParams(r)
	limitQuery := ""

	limitQuery = fmt.Sprintf("LIMIT %d  OFFSET %d", limit, offset)

	query := `SELECT * FROM payment_histories WHERE org_uuid = '` + org_uuid + `' AND status = true ORDER BY created DESC`

	db.db.Raw(query + " " + limitQuery).Find(&payment)
	return payment
}

func (db database) GetWorkspaceInvoices(org_uuid string) []InvoiceList {
	ms := []InvoiceList{}
	db.db.Where("org_uuid = ?", org_uuid).Where("status", false).Find(&ms)
	return ms
}

func (db database) GetWorkspaceInvoicesCount(org_uuid string) int64 {
	var count int64
	ms := InvoiceList{}

	db.db.Model(&ms).Where("org_uuid = ?", org_uuid).Where("status", false).Count(&count)
	return count
}

func (db database) ChangeWorkspaceDeleteStatus(org_uuid string, status bool) Organization {
	ms := Organization{}
	db.db.Model(&ms).Where("uuid", org_uuid).Updates(map[string]interface{}{
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

	result := db.db.Model(&Organization{}).Where("uuid = ?", uuid).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db database) DeleteAllUsersFromWorkspace(org string) error {
	if org == "" {
		return errors.New("no org uuid provided")
	}

	// Delete all users associated with the organization
	result := db.db.Where("org_uuid = ?", org).Delete(&OrganizationUsers{})
	if result.Error != nil {
		return result.Error
	}

	// Delete all user roles associated with the organization
	result = db.db.Where("org_uuid = ?", org).Delete(&UserRoles{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
