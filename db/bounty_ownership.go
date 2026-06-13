package db

import (
	"errors"
	"time"
)

func (db database) TransferWorkspaceBountyOwnership(workspaceUUID string, fromOwnerID string, toOwnerID string) (int64, error) {
	if workspaceUUID == "" {
		return 0, errors.New("workspace_uuid is required")
	}
	if fromOwnerID == "" {
		return 0, errors.New("from_owner_id is required")
	}
	if toOwnerID == "" {
		return 0, errors.New("to_owner_id is required")
	}
	if fromOwnerID == toOwnerID {
		return 0, errors.New("from_owner_id and to_owner_id must be different")
	}

	now := time.Now()
	result := db.db.Model(&NewBounty{}).
		Where("workspace_uuid = ? AND owner_id = ? AND paid = ?", workspaceUUID, fromOwnerID, false).
		Updates(map[string]interface{}{
			"owner_id": toOwnerID,
			"updated":  &now,
		})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
