package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
)

type skillHandler struct {
	db db.Database
}

func NewSkillHandler(db db.Database) *skillHandler {
	return &skillHandler{
		db: db,
	}
}


func (sh *skillHandler) CreateSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[skill] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("failed to read request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error reading request body"})
		return
	}
	defer r.Body.Close()

	var skill db.Skill
	if err := json.Unmarshal(body, &skill); err != nil {
		logger.Log.Error("failed to unmarshal skill data", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid skill data format"})
		return
	}

	skill.OwnerPubkey = pubKeyFromAuth

	createdSkill, err := sh.db.CreateSkill(&skill)
	if err != nil {
		logger.Log.Error("failed to create skill", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSkill)
}

func (sh *skillHandler) GetAllSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := sh.db.GetAllSkills()
	if err != nil {
		logger.Log.Error("failed to get all skills", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(skills)
}

func (sh *skillHandler) GetSkillByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Skill ID is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid skill ID format"})
		return
	}

	skill, err := sh.db.GetSkillByID(id)
	if err != nil {
		logger.Log.Error("failed to get skill by ID", "error", err, "id", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(skill)
}

func (sh *skillHandler) UpdateSkillByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[skill] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Skill ID is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid skill ID format"})
		return
	}

	existingSkill, err := sh.db.GetSkillByID(id)
	if err != nil {
		logger.Log.Error("failed to get skill by ID", "error", err, "id", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if existingSkill.OwnerPubkey != pubKeyFromAuth {
		logger.Log.Info("[skill] unauthorized update attempt", "pubkey", pubKeyFromAuth, "owner", existingSkill.OwnerPubkey)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to update this skill"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("failed to read request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error reading request body"})
		return
	}
	defer r.Body.Close()

	var updatedSkill db.Skill
	if err := json.Unmarshal(body, &updatedSkill); err != nil {
		logger.Log.Error("failed to unmarshal skill data", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid skill data format"})
		return
	}

	updatedSkill.ID = id
	updatedSkill.OwnerPubkey = pubKeyFromAuth

	result, err := sh.db.UpdateSkillByID(&updatedSkill)
	if err != nil {
		logger.Log.Error("failed to update skill", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (sh *skillHandler) DeleteSkillByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[skill] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Skill ID is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid skill ID format"})
		return
	}

	existingSkill, err := sh.db.GetSkillByID(id)
	if err != nil {
		logger.Log.Error("failed to get skill by ID", "error", err, "id", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if existingSkill.OwnerPubkey != pubKeyFromAuth {
		logger.Log.Info("[skill] unauthorized delete attempt", "pubkey", pubKeyFromAuth, "owner", existingSkill.OwnerPubkey)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to delete this skill"})
		return
	}

	if err := sh.db.DeleteSkillByID(id); err != nil {
		logger.Log.Error("failed to delete skill", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "true",
		"message": "Skill successfully deleted",
	})
}

func (sh *skillHandler) CreateSkillInstall(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[skill] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Skill ID is required"})
		return
	}

	skillID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid skill ID format"})
		return
	}

	existingSkill, err := sh.db.GetSkillByID(skillID)
	if err != nil {
		logger.Log.Error("failed to get skill by ID", "error", err, "id", skillID)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if existingSkill.OwnerPubkey != pubKeyFromAuth {
		logger.Log.Info("[skill] unauthorized install creation attempt", "pubkey", pubKeyFromAuth, "owner", existingSkill.OwnerPubkey)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to create installations for this skill"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("failed to read request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error reading request body"})
		return
	}
	defer r.Body.Close()

	var install db.SkillInstall
	if err := json.Unmarshal(body, &install); err != nil {
		logger.Log.Error("failed to unmarshal installation data", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid installation data format"})
		return
	}

	install.SkillID = skillID

	createdInstall, err := sh.db.CreateSkillInstall(&install)
	if err != nil {
		logger.Log.Error("failed to create skill installation", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdInstall)
}

func (sh *skillHandler) GetAllSkillInstallsBySkillID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Skill ID is required"})
		return
	}

	skillID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid skill ID format"})
		return
	}

	_, err = sh.db.GetSkillByID(skillID)
	if err != nil {
		logger.Log.Error("failed to get skill by ID", "error", err, "id", skillID)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	installs, err := sh.db.GetSkillInstallBySkillsID(skillID)
	if err != nil {
		logger.Log.Error("failed to get skill installations", "error", err, "skill_id", skillID)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(installs)
}

func (sh *skillHandler) DeleteSkillInstallByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[skill] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Installation ID is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid installation ID format"})
		return
	}

	install, err := sh.db.GetSkillInstallByID(id)
	if err != nil {
		logger.Log.Error("failed to get skill installation by ID", "error", err, "id", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	skill, err := sh.db.GetSkillByID(install.SkillID)
	if err != nil {
		logger.Log.Error("failed to get skill by ID", "error", err, "id", install.SkillID)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if skill.OwnerPubkey != pubKeyFromAuth {
		logger.Log.Info("[skill] unauthorized delete attempt", "pubkey", pubKeyFromAuth, "owner", skill.OwnerPubkey)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to delete installations for this skill"})
		return
	}

	if err := sh.db.DeleteSkillInstallByID(id); err != nil {
		logger.Log.Error("failed to delete skill installation", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "true",
		"message": "Skill installation successfully deleted",
	})
}

func (sh *skillHandler) UpdateSkillInstallByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[skill] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Installation ID is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid installation ID format"})
		return
	}

	existingInstall, err := sh.db.GetSkillInstallByID(id)
	if err != nil {
		logger.Log.Error("failed to get skill installation by ID", "error", err, "id", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	skill, err := sh.db.GetSkillByID(existingInstall.SkillID)
	if err != nil {
		logger.Log.Error("failed to get skill by ID", "error", err, "id", existingInstall.SkillID)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if skill.OwnerPubkey != pubKeyFromAuth {
		logger.Log.Info("[skill] unauthorized update attempt", "pubkey", pubKeyFromAuth, "owner", skill.OwnerPubkey)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to update installations for this skill"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("failed to read request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error reading request body"})
		return
	}
	defer r.Body.Close()

	var updatedInstall db.SkillInstall
	if err := json.Unmarshal(body, &updatedInstall); err != nil {
		logger.Log.Error("failed to unmarshal installation data", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid installation data format"})
		return
	}

	updatedInstall.ID = id
	updatedInstall.SkillID = existingInstall.SkillID

	result, err := sh.db.UpdateSkillInstallByID(&updatedInstall)
	if err != nil {
		logger.Log.Error("failed to update skill installation", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (sh *skillHandler) GetSkillInstallByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Installation ID is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid installation ID format"})
		return
	}

	install, err := sh.db.GetSkillInstallByID(id)
	if err != nil {
		logger.Log.Error("failed to get skill installation by ID", "error", err, "id", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(install)
}