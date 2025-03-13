package db

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/logger"
	"gorm.io/gorm"
)

func (db database) CreateSkill(skill *Skill) (*Skill, error) {
	if skill.ID == uuid.Nil {
		skill.ID = uuid.New()
	}

	if skill.Name == "" {
		return nil, errors.New("skill name is required")
	}

	if skill.OwnerPubkey == "" {
		return nil, errors.New("owner pubkey is required")
	}

	if skill.Status == "" {
		skill.Status = DraftSkillStatus
	}

	if skill.ChargeModel == "" {
		skill.ChargeModel = FreeChargeModel
	}

	if err := db.db.Create(skill).Error; err != nil {
		logger.Log.Error("failed to create skill", "error", err)
		return nil, fmt.Errorf("failed to create skill: %w", err)
	}

	return skill, nil
}

func (db database) GetAllSkills() ([]Skill, error) {
	var skills []Skill
	if err := db.db.Find(&skills).Error; err != nil {
		logger.Log.Error("failed to get all skills", "error", err)
		return nil, fmt.Errorf("failed to get all skills: %w", err)
	}
	return skills, nil
}

func (db database) GetSkillByID(id uuid.UUID) (*Skill, error) {
	var skill Skill
	if err := db.db.Where("id = ?", id).First(&skill).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("skill not found with ID: %s", id)
		}
		logger.Log.Error("failed to get skill by ID", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}
	return &skill, nil
}

func (db database) UpdateSkillByID(skill *Skill) (*Skill, error) {
	if skill.ID == uuid.Nil {
		return nil, errors.New("skill ID is required")
	}

	existingSkill, err := db.GetSkillByID(skill.ID)
	if err != nil {
		return nil, err
	}

	if err := db.db.Model(&existingSkill).Updates(skill).Error; err != nil {
		logger.Log.Error("failed to update skill", "error", err, "id", skill.ID)
		return nil, fmt.Errorf("failed to update skill: %w", err)
	}

	updatedSkill, err := db.GetSkillByID(skill.ID)
	if err != nil {
		return nil, err
	}

	return updatedSkill, nil
}

func (db database) DeleteSkillByID(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("skill ID is required")
	}

	_, err := db.GetSkillByID(id)
	if err != nil {
		return err
	}

	if err := db.db.Delete(&Skill{ID: id}).Error; err != nil {
		logger.Log.Error("failed to delete skill", "error", err, "id", id)
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	return nil
}

func (db database) CreateSkillInstall(install *SkillInstall) (*SkillInstall, error) {
	if install.ID == uuid.Nil {
		install.ID = uuid.New()
	}

	if install.SkillID == uuid.Nil {
		return nil, errors.New("skill ID is required")
	}

	if install.Client == "" {
		return nil, errors.New("client type is required")
	}

	_, err := db.GetSkillByID(install.SkillID)
	if err != nil {
		return nil, fmt.Errorf("invalid skill ID: %w", err)
	}

	if err := db.db.Create(install).Error; err != nil {
		logger.Log.Error("failed to create skill installation", "error", err)
		return nil, fmt.Errorf("failed to create skill installation: %w", err)
	}

	return install, nil
}

func (db database) GetSkillInstallBySkillsID(skillID uuid.UUID) ([]SkillInstall, error) {
	var installs []SkillInstall
	if err := db.db.Where("skill_id = ?", skillID).Find(&installs).Error; err != nil {
		logger.Log.Error("failed to get skill installations", "error", err, "skill_id", skillID)
		return nil, fmt.Errorf("failed to get skill installations: %w", err)
	}
	return installs, nil
}

func (db database) GetSkillInstallByID(id uuid.UUID) (*SkillInstall, error) {
	var install SkillInstall
	if err := db.db.Where("id = ?", id).First(&install).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("skill installation not found with ID: %s", id)
		}
		logger.Log.Error("failed to get skill installation by ID", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get skill installation: %w", err)
	}
	return &install, nil
}

func (db database) UpdateSkillInstall(install *SkillInstall) (*SkillInstall, error) {
	if install.ID == uuid.Nil {
		return nil, errors.New("skill installation ID is required")
	}

	existingInstall, err := db.GetSkillInstallByID(install.ID)
	if err != nil {
		return nil, err
	}

	if err := db.db.Model(&existingInstall).Updates(install).Error; err != nil {
		logger.Log.Error("failed to update skill installation", "error", err, "id", install.ID)
		return nil, fmt.Errorf("failed to update skill installation: %w", err)
	}

	updatedInstall, err := db.GetSkillInstallByID(install.ID)
	if err != nil {
		return nil, err
	}

	return updatedInstall, nil
}

func (db database) DeleteSkillInstall(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("skill installation ID is required")
	}

	_, err := db.GetSkillInstallByID(id)
	if err != nil {
		return err
	}

	if err := db.db.Delete(&SkillInstall{ID: id}).Error; err != nil {
		logger.Log.Error("failed to delete skill installation", "error", err, "id", id)
		return fmt.Errorf("failed to delete skill installation: %w", err)
	}

	return nil
} 