package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (db database) CreateCodeSpaceMap(codeSpace CodeSpaceMap) (CodeSpaceMap, error) {
	if codeSpace.ID == uuid.Nil {
		codeSpace.ID = uuid.New()
	}
	
	var existingMap CodeSpaceMap
	result := db.db.Where("workspace_id = ? AND user_pubkey = ?", codeSpace.WorkspaceID, codeSpace.UserPubkey).First(&existingMap)
	
	if result.RowsAffected > 0 {
		now := time.Now()
		existingMap.CodeSpaceURL = codeSpace.CodeSpaceURL
		existingMap.Username = codeSpace.Username
		existingMap.GithubPat = codeSpace.GithubPat
		existingMap.BaseBranch = codeSpace.BaseBranch
		existingMap.UpdatedAt = now
		
		db.db.Save(&existingMap)
		return existingMap, nil
	}
	
	now := time.Now()
	codeSpace.CreatedAt = now
	codeSpace.UpdatedAt = now
	
	if err := db.db.Create(&codeSpace).Error; err != nil {
		return CodeSpaceMap{}, err
	}
	
	return codeSpace, nil
}

func (db database) GetCodeSpaceMaps() ([]CodeSpaceMap, error) {
	var codespaces []CodeSpaceMap
	if err := db.db.Find(&codespaces).Error; err != nil {
		return nil, err
	}
	return codespaces, nil
}

func (db database) GetCodeSpaceMapByWorkspace(workspaceID string) ([]CodeSpaceMap, error) {
	var codespaces []CodeSpaceMap
	result := db.db.Where("workspace_id = ?", workspaceID).Find(&codespaces)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return codespaces, nil
}

func (db database) GetCodeSpaceMapByUser(userPubkey string) ([]CodeSpaceMap, error) {
	var codespaces []CodeSpaceMap
	result := db.db.Where("user_pubkey = ?", userPubkey).Find(&codespaces)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return codespaces, nil
}

func (db database) GetCodeSpaceMapByURL(codeSpaceURL string) ([]CodeSpaceMap, error) {
	var codespaces []CodeSpaceMap
	result := db.db.Where("code_space_url = ?", codeSpaceURL).Find(&codespaces)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return codespaces, nil
}

func (db database) GetCodeSpaceMapByWorkspaceAndUser(workspaceID, userPubkey string) (CodeSpaceMap, error) {
	var codespace CodeSpaceMap
	result := db.db.Where("workspace_id = ? AND user_pubkey = ?", workspaceID, userPubkey).First(&codespace)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return CodeSpaceMap{}, errors.New("codespace mapping not found")
		}
		return CodeSpaceMap{}, result.Error
	}
	
	return codespace, nil
}

func (db database) GetCodeSpaceMapByID(id uuid.UUID) (CodeSpaceMap, error) {
	var codespace CodeSpaceMap
	result := db.db.Where("id = ?", id).First(&codespace)
	
	if result.RowsAffected == 0 {
		return CodeSpaceMap{}, errors.New("codespace mapping not found")
	}
	
	return codespace, nil
}

func (db database) UpdateCodeSpaceMap(id uuid.UUID, updates map[string]interface{}) (CodeSpaceMap, error) {
	var codespace CodeSpaceMap
	result := db.db.Where("id = ?", id).First(&codespace)
	
	if result.RowsAffected == 0 {
		return CodeSpaceMap{}, errors.New("codespace mapping not found")
	}
	
	updates["updated_at"] = time.Now()
	if err := db.db.Model(&codespace).Updates(updates).Error; err != nil {
		return CodeSpaceMap{}, err
	}
	
	if err := db.db.Where("id = ?", id).First(&codespace).Error; err != nil {
		return CodeSpaceMap{}, err
	}
	
	return codespace, nil
}

func (db database) DeleteCodeSpaceMap(id uuid.UUID) error {
	var codespace CodeSpaceMap
	result := db.db.Where("id = ?", id).First(&codespace)
	
	if result.RowsAffected == 0 {
		return errors.New("codespace mapping not found")
	}
	
	if err := db.db.Delete(&codespace).Error; err != nil {
		return err
	}
	
	return nil
} 