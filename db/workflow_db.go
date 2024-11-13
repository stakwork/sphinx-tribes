package db

import (
	"errors"
	"time"
)

func (db database) CreateWorkflowRequest(req *WfRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	now := time.Now()
	req.CreatedAt = now
	req.UpdatedAt = now

	if req.Status == "" {
		req.Status = StatusNew
	}

	result := db.db.Create(req)
	return result.Error
}

func (db database) UpdateWorkflowRequest(req *WfRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	req.UpdatedAt = time.Now()
	result := db.db.Model(&WfRequest{}).Where("request_id = ?", req.RequestID).Updates(req)

	if result.RowsAffected == 0 {
		return errors.New("no workflow request found to update")
	}

	return result.Error
}

func (db database) GetWorkflowRequestByID(requestID string) (*WfRequest, error) {
	if requestID == "" {
		return nil, errors.New("request ID cannot be empty")
	}

	var req WfRequest
	result := db.db.Model(&WfRequest{}).Where("request_id = ?", requestID).First(&req)

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &req, result.Error
}

func (db database) GetWorkflowRequestsByStatus(status WfRequestStatus) ([]WfRequest, error) {
	var requests []WfRequest

	result := db.db.Model(&WfRequest{}).
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&requests)

	if result.Error != nil {
		return nil, result.Error
	}

	return requests, nil
}

func (db database) GetWorkflowRequest(requestID string) (*WfRequest, error) {
	if requestID == "" {
		return nil, errors.New("request ID cannot be empty")
	}

	var req WfRequest
	result := db.db.Model(&WfRequest{}).Where("request_id = ?", requestID).First(&req)

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &req, result.Error
}

func (db database) UpdateWorkflowRequestStatus(requestID string, status WfRequestStatus, responseData JSONB) error {
	if requestID == "" {
		return errors.New("request ID cannot be empty")
	}

	result := db.db.Model(&WfRequest{}).
		Where("request_id = ?", requestID).
		Updates(map[string]interface{}{
			"status":        status,
			"response_data": responseData,
			"updated_at":    time.Now(),
		})

	if result.RowsAffected == 0 {
		return errors.New("no workflow request found to update")
	}

	return result.Error
}

func (db database) GetWorkflowRequestsByWorkflowID(workflowID string) ([]WfRequest, error) {
	if workflowID == "" {
		return nil, errors.New("workflow ID cannot be empty")
	}

	var requests []WfRequest
	result := db.db.Model(&WfRequest{}).
		Where("workflow_id = ?", workflowID).
		Order("created_at DESC").
		Find(&requests)

	return requests, result.Error
}

func (db database) GetPendingWorkflowRequests(limit int) ([]WfRequest, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	var requests []WfRequest
	result := db.db.Model(&WfRequest{}).
		Where("status = ?", StatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&requests)

	return requests, result.Error
}

func (db database) DeleteWorkflowRequest(requestID string) error {
	if requestID == "" {
		return errors.New("request ID cannot be empty")
	}

	result := db.db.Delete(&WfRequest{}, "request_id = ?", requestID)

	if result.RowsAffected == 0 {
		return errors.New("no workflow request found to delete")
	}

	return result.Error
}

func (db database) CreateProcessingMap(pm *WfProcessingMap) error {
	if pm == nil {
		return errors.New("processing map cannot be nil")
	}

	now := time.Now()
	pm.CreatedAt = now
	pm.UpdatedAt = now

	result := db.db.Create(pm)
	return result.Error
}

func (db database) UpdateProcessingMap(pm *WfProcessingMap) error {
	if pm == nil {
		return errors.New("processing map cannot be nil")
	}

	pm.UpdatedAt = time.Now()
	result := db.db.Model(&WfProcessingMap{}).
		Where("id = ?", pm.ID).
		Updates(pm)

	if result.RowsAffected == 0 {
		return errors.New("no processing map found to update")
	}

	return result.Error
}

func (db database) GetProcessingMapByKey(processType, processKey string) (*WfProcessingMap, error) {
	if processType == "" || processKey == "" {
		return nil, errors.New("process type and key cannot be empty")
	}

	var pm WfProcessingMap
	result := db.db.Model(&WfProcessingMap{}).
		Where("type = ? AND process_key = ?", processType, processKey).
		First(&pm)

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &pm, result.Error
}

func (db database) GetProcessingMapsByType(processType string) ([]WfProcessingMap, error) {
	if processType == "" {
		return nil, errors.New("process type cannot be empty")
	}

	var maps []WfProcessingMap
	result := db.db.Model(&WfProcessingMap{}).
		Where("type = ?", processType).
		Order("created_at DESC").
		Find(&maps)

	return maps, result.Error
}

func (db database) DeleteProcessingMap(id uint) error {
	if id == 0 {
		return errors.New("invalid processing map ID")
	}

	result := db.db.Delete(&WfProcessingMap{}, id)

	if result.RowsAffected == 0 {
		return errors.New("no processing map found to delete")
	}

	return result.Error
}
