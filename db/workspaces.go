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

	query.Order("name ASC").Find(&ms)
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

func (db database) DeleteWorkspace() (bool, error) {
	result := db.db.Exec("DELETE FROM workspaces")
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
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

func (db database) DeleteWorkspaceUser(orgUser WorkspaceUsersData, workspace_uuid string) WorkspaceUsersData {
	db.db.Where("owner_pub_key = ?", orgUser.OwnerPubKey).Where("workspace_uuid = ?", workspace_uuid).Delete(&WorkspaceUsers{})
	db.db.Where("owner_pub_key = ?", orgUser.OwnerPubKey).Where("workspace_uuid = ?", workspace_uuid).Delete(&UserRoles{})
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

func (db database) DeleteWorkspaceBudget() error {
	err := db.db.Unscoped().Where("1 = 1").Delete(&NewBountyBudget{}).Error
	return err
}

func (db database) GetWorkspaceStatusBudget(workspace_uuid string) StatusBudget {
	workspaceBudget := db.GetWorkspaceBudget(workspace_uuid)

	var openBudget uint
	db.db.Model(&NewBounty{}).Where("workspace_uuid = ?", workspace_uuid).Where("assignee = '' ").Where("paid != true").Select("SUM(price)").Row().Scan(&openBudget)

	var openCount int64
	db.db.Model(&NewBounty{}).Where("workspace_uuid = ?", workspace_uuid).Where("assignee = '' ").Where("paid != true").Count(&openCount)

	var openDifference int = int(workspaceBudget.TotalBudget - openBudget)

	var assignedBudget uint
	db.db.Model(&NewBounty{}).Where("workspace_uuid = ?", workspace_uuid).Where("assignee != '' ").Where("paid != true").Select("SUM(price)").Row().Scan(&assignedBudget)

	var assignedCount int64
	db.db.Model(&NewBounty{}).Where("workspace_uuid = ?", workspace_uuid).Where("assignee != '' ").Where("paid != true").Count(&assignedCount)

	var assignedDifference int = int(workspaceBudget.TotalBudget - assignedBudget)

	var completedBudget uint
	db.db.Model(&NewBounty{}).Where("workspace_uuid = ?", workspace_uuid).Where("completed = true ").Where("paid != true").Select("SUM(price)").Row().Scan(&completedBudget)

	var completedCount int64
	db.db.Model(&NewBounty{}).Where("workspace_uuid = ?", workspace_uuid).Where("completed = true ").Where("paid != true").Count(&completedCount)

	var completedDifference int = int(workspaceBudget.TotalBudget - completedBudget)

	statusBudget := StatusBudget{
		OrgUuid:             workspace_uuid,
		WorkspaceUuid:       workspace_uuid,
		CurrentBudget:       workspaceBudget.TotalBudget,
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

func (db database) ProcessUpdateBudget(non_tx_invoice NewInvoiceList) error {
	// Start db transaction
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

	created := non_tx_invoice.Created
	workspace_uuid := non_tx_invoice.WorkspaceUuid

	invoice := NewInvoiceList{}
	tx.Where("payment_request = ?", non_tx_invoice.PaymentRequest).Find(&invoice)

	if invoice.Status {
		tx.Rollback()
		return errors.New("cannot process already paid invoice")
	}

	if workspace_uuid == "" {
		return errors.New("cannot Create a Workspace Without a Workspace uuid")
	}

	// Get payment history and update budget
	paymentHistory := NewPaymentHistory{}
	tx.Model(&NewPaymentHistory{}).Where("created = ?", created).Where("workspace_uuid = ? ", workspace_uuid).Find(&paymentHistory)

	if paymentHistory.WorkspaceUuid != "" && paymentHistory.Amount != 0 {
		paymentHistory.Status = true

		// Update payment history
		if err = tx.Where("created = ?", created).Where("workspace_uuid = ? ", workspace_uuid).Updates(paymentHistory).Error; err != nil {
			tx.Rollback()
		}

		// get Workspace budget and add payment to total budget
		workspaceBudget := NewBountyBudget{}
		tx.Model(&NewBountyBudget{}).Where("workspace_uuid = ?", workspace_uuid).Find(&workspaceBudget)

		if workspaceBudget.WorkspaceUuid == "" {
			now := time.Now()
			workBudget := NewBountyBudget{
				WorkspaceUuid: workspace_uuid,
				TotalBudget:   paymentHistory.Amount,
				Created:       &now,
				Updated:       &now,
			}

			if err = tx.Create(&workBudget).Error; err != nil {
				tx.Rollback()
			}
		} else {
			totalBudget := workspaceBudget.TotalBudget
			workspaceBudget.TotalBudget = totalBudget + paymentHistory.Amount

			if err = tx.Model(&NewBountyBudget{}).Where("workspace_uuid = ?", workspaceBudget.WorkspaceUuid).Updates(map[string]interface{}{
				"total_budget": workspaceBudget.TotalBudget,
			}).Error; err != nil {
				tx.Rollback()
			}
		}

		// update invoice
		if err = tx.Model(&NewInvoiceList{}).Where("payment_request = ?", invoice.PaymentRequest).Update("status", true).Error; err != nil {
			tx.Rollback()
		}
	}

	return tx.Commit().Error
}

func (db database) AddAndUpdateBudget(invoice NewInvoiceList) NewPaymentHistory {
	// Start db transaction
	tx := db.db.Begin()

	created := invoice.Created
	workspace_uuid := invoice.WorkspaceUuid

	paymentHistory := NewPaymentHistory{}
	tx.Model(&NewPaymentHistory{}).Where("created = ?", created).Where("workspace_uuid = ? ", workspace_uuid).Find(&paymentHistory)

	if paymentHistory.WorkspaceUuid != "" && paymentHistory.Amount != 0 {
		paymentHistory.Status = true
		db.db.Where("created = ?", created).Where("workspace_uuid = ? ", workspace_uuid).Updates(paymentHistory)

		// get Workspace budget and add payment to total budget
		workspaceBudget := NewBountyBudget{}
		tx.Model(&NewBountyBudget{}).Where("workspace_uuid = ?", workspace_uuid).Find(&workspaceBudget)

		if workspaceBudget.WorkspaceUuid == "" {
			now := time.Now()
			workBudget := NewBountyBudget{
				WorkspaceUuid: workspace_uuid,
				TotalBudget:   paymentHistory.Amount,
				Created:       &now,
				Updated:       &now,
			}

			if err := tx.Create(&workBudget).Error; err != nil {
				tx.Rollback()
			}
		} else {
			totalBudget := workspaceBudget.TotalBudget
			workspaceBudget.TotalBudget = totalBudget + paymentHistory.Amount

			if err := tx.Model(&NewBountyBudget{}).Where("workspace_uuid = ?", workspaceBudget.WorkspaceUuid).Updates(map[string]interface{}{
				"total_budget": workspaceBudget.TotalBudget,
			}).Error; err != nil {
				tx.Rollback()
			}
		}
	} else {
		tx.Rollback()
	}

	tx.Commit()

	return paymentHistory
}

func (db database) WithdrawBudget(sender_pubkey string, workspace_uuid string, amount uint) {
	tx := db.db.Begin()
	var err error

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return
	}

	// get Workspace budget and add payment to total budget
	WorkspaceBudget := db.GetWorkspaceBudget(workspace_uuid)
	totalBudget := WorkspaceBudget.TotalBudget

	newBudget := totalBudget - amount

	if err = tx.Model(&NewBountyBudget{}).Where("workspace_uuid = ?", workspace_uuid).Updates(map[string]interface{}{
		"total_budget": newBudget,
	}).Error; err != nil {
		tx.Rollback()
	}

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

	if err = tx.Create(&budgetHistory).Error; err != nil {
		tx.Rollback()
	}
	tx.Commit()
}

func (db database) AddPaymentHistory(payment NewPaymentHistory) NewPaymentHistory {
	db.db.Create(&payment)

	return payment
}

func (db database) ProcessBountyPayment(payment NewPaymentHistory, bounty NewBounty) error {
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

	// add to payment history
	if err = tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return err
	}

	if payment.PaymentStatus != PaymentFailed {
		// get Workspace budget and subtract payment from total budget
		WorkspaceBudget := db.GetWorkspaceBudget(payment.WorkspaceUuid)
		totalBudget := WorkspaceBudget.TotalBudget

		// update budget
		WorkspaceBudget.TotalBudget = totalBudget - payment.Amount
		if err = tx.Model(&NewBountyBudget{}).Where("workspace_uuid = ?", payment.WorkspaceUuid).Updates(map[string]interface{}{
			"total_budget": WorkspaceBudget.TotalBudget,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// updatge bounty status
		if err = tx.Where("created", bounty.Created).Updates(&bounty).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
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

func (db database) GetPendingPaymentHistory() []NewPaymentHistory {
	paymentHistories := []NewPaymentHistory{}

	query := `SELECT * FROM payment_histories WHERE payment_status = '` + PaymentPending + `' AND status = true AND payment_type = 'payment' ORDER BY created DESC`

	db.db.Raw(query).Find(&paymentHistories)
	return paymentHistories
}

func (db database) GetPaymentByBountyId(bountyId uint) NewPaymentHistory {
	paymentHistories := NewPaymentHistory{}

	query := fmt.Sprintf("SELECT * FROM payment_histories WHERE bounty_id = %d AND status = true ORDER BY created DESC", bountyId)

	db.db.Raw(query).Find(&paymentHistories)

	return paymentHistories
}

func (db database) SetPaymentAsComplete(tag string) bool {
	db.db.Model(NewPaymentHistory{}).Where("tag = ?", tag).Update("payment_status", PaymentComplete)
	return true
}

func (db database) SetPaymentStatusByBountyId(bountyId uint, tagResult V2TagRes) bool {
	mapResult := map[string]string{}

	mapResult["payment_status"] = tagResult.Status
	mapResult["error"] = tagResult.Error
	mapResult["tag"] = tagResult.Tag

	db.db.Model(NewPaymentHistory{}).Where("bounty_id = ?", bountyId).Updates(mapResult)
	return true
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

func (db database) ProcessDeleteWorkspace(workspace_uuid string) error {
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

	updates := map[string]interface{}{
		"website":     "",
		"github":      "",
		"description": "",
		"show":        false,
	}

	// Update workspace
	if err = tx.Model(&Workspace{}).Where("uuid = ?", workspace_uuid).Updates(updates).Error; err != nil {
		tx.Rollback()
	}

	// Delete all users associated with the Workspace
	if err = tx.Where("workspace_uuid = ?", workspace_uuid).Delete(&WorkspaceUsers{}).Error; err != nil {
		tx.Rollback()
	}

	// Delete all user roles associated with the Workspace
	if err = tx.Where("workspace_uuid = ?", workspace_uuid).Delete(&WorkspaceUserRoles{}).Error; err != nil {
		tx.Rollback()
	}

	// Change delete status to true
	if err = tx.Model(&Workspace{}).Where("uuid", workspace_uuid).Updates(map[string]interface{}{
		"deleted": true,
	}).Error; err != nil {
		tx.Rollback()
	}

	return tx.Commit().Error
}

func (db database) DeleteAllUsersFromWorkspace(workspace_uuid string) error {
	if workspace_uuid == "" {
		return errors.New("no workspoace uuid provided")
	}

	// Delete all users associated with the Workspace
	result := db.db.Where("workspace_uuid = ?", workspace_uuid).Delete(&WorkspaceUsers{})
	if result.Error != nil {
		return result.Error
	}

	// Delete all user roles associated with the Workspace
	result = db.db.Where("workspace_uuid = ?", workspace_uuid).Delete(&WorkspaceUserRoles{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db database) GetLastWithdrawal(workspace_uuid string) NewPaymentHistory {
	p := NewPaymentHistory{}
	db.db.Model(&NewPaymentHistory{}).Where("workspace_uuid", workspace_uuid).Where("payment_type", "withdraw").Order("created DESC").Limit(1).Find(&p)
	return p
}

func (db database) GetSumOfDeposits(workspace_uuid string) uint {
	var depositAmount uint
	db.db.Model(&NewPaymentHistory{}).Where("workspace_uuid = ?", workspace_uuid).Where("status = ?", true).Where("payment_type = ?", "deposit").Select("SUM(amount)").Row().Scan(&depositAmount)

	return depositAmount
}

func (db database) GetSumOfWithdrawal(workspace_uuid string) uint {
	var depositAmount uint
	db.db.Model(&NewPaymentHistory{}).Where("workspace_uuid = ?", workspace_uuid).Where("status = ?", true).Where("payment_type = ?", "withdraw").Select("SUM(amount)").Row().Scan(&depositAmount)

	return depositAmount
}

func (db database) GetFeaturePhasesBountiesCount(bountyType string, phaseUuid string) int64 {
	var count int64

	query := db.db.Model(&NewBounty{})
	if bountyType == "open" {
		query.Where("phase_uuid", phaseUuid).Where("assignee = '' ")
	} else if bountyType == "assigned" {
		query.Where("phase_uuid", phaseUuid).Where("assignee != '' ").Where("paid != true").Where("completed != true")
	} else if bountyType == "completed" {
		query.Where("phase_uuid", phaseUuid).Where("completed = true")
	} else if bountyType == "paid" {
		query.Where("phase_uuid", phaseUuid).Where("paid = true")
	}

	query.Count(&count)
	return count
}
