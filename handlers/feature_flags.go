package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
)

type FeatureFlagHandler struct {
	db db.Database
}

type FeatureFlagResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type CreateFeatureFlagRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	Endpoints   []string `json:"endpoints"`
}

type UpdateFeatureFlagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type AddFeatureFlagEndpointRequest struct {
	Endpoints []string `json:"endpoints"`
}

type UpdateFeatureFlagEndpointRequest struct {
	Path string `json:"new_endpoint_path"`
}

func NewFeatureFlagHandler(database db.Database) *FeatureFlagHandler {
	return &FeatureFlagHandler{
		db: database,
	}
}

// GetFeatureFlags godoc
//
//	@Summary		Get all feature flags
//	@Description	Get a list of all feature flags
//	@Tags			Feature Flag
//	@Success		200	{object}	FeatureFlagResponse
//	@Router			/feature_flags [get]
func (fh *FeatureFlagHandler) GetFeatureFlags(w http.ResponseWriter, r *http.Request) {
	flags, err := fh.db.GetFeatureFlags()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Failed to fetch feature flags",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FeatureFlagResponse{
		Success: true,
		Message: "Feature flags fetched successfully",
		Data:    flags,
	})
}

// CreateFeatureFlag godoc
//
//	@Summary		Create a new feature flag
//	@Description	Create a new feature flag with specified details
//	@Tags			Feature Flag
//	@Param			feature_flag	body		CreateFeatureFlagRequest	true	"Feature flag details"
//	@Success		201				{object}	FeatureFlagResponse
//	@Router			/feature_flags [post]
func (fh *FeatureFlagHandler) CreateFeatureFlag(w http.ResponseWriter, r *http.Request) {
	var request CreateFeatureFlagRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	flag := &db.FeatureFlag{
		UUID:        uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		Enabled:     request.Enabled,
	}

	createdFlag, err := fh.db.AddFeatureFlag(flag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Failed to create feature flag",
		})
		return
	}

	var endpoints []db.Endpoint
	for _, path := range request.Endpoints {
		endpoint := &db.Endpoint{
			UUID:            uuid.New(),
			Path:            path,
			FeatureFlagUUID: createdFlag.UUID,
		}
		createdEndpoint, err := fh.db.AddEndpoint(endpoint)
		if err != nil {
			_ = fh.db.DeleteFeatureFlag(createdFlag.UUID)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(FeatureFlagResponse{
				Success: false,
				Message: "Failed to create endpoint",
			})
			return
		}
		endpoints = append(endpoints, createdEndpoint)
	}

	createdFlag.Endpoints = endpoints

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(FeatureFlagResponse{
		Success: true,
		Message: "Feature flag created successfully",
		Data:    createdFlag,
	})
}

// UpdateFeatureFlag godoc
//
//	@Summary		Update an existing feature flag
//	@Description	Update the details of an existing feature flag
//	@Tags			Feature Flag
//	@Param			id				path		string						true	"Feature flag ID"
//	@Param			feature_flag	body		UpdateFeatureFlagRequest	true	"Updated feature flag details"
//	@Success		200				{object}	FeatureFlagResponse
//	@Router			/feature_flags/{id} [put]
func (fh *FeatureFlagHandler) UpdateFeatureFlag(w http.ResponseWriter, r *http.Request) {
	flagID := chi.URLParam(r, "id")
	if flagID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Feature flag ID is required",
		})
		return
	}

	flagUUID, err := uuid.Parse(flagID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid feature flag ID",
		})
		return
	}

	var request UpdateFeatureFlagRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	flag := &db.FeatureFlag{
		UUID:        flagUUID,
		Name:        request.Name,
		Description: request.Description,
		Enabled:     request.Enabled,
	}

	updatedFlag, err := fh.db.UpdateFeatureFlag(flag)
	if err != nil {
		if err.Error() == "feature flag not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(FeatureFlagResponse{
				Success: false,
				Message: "Feature flag not found",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Failed to update feature flag",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FeatureFlagResponse{
		Success: true,
		Message: "Feature flag updated successfully",
		Data:    updatedFlag,
	})
}

// DeleteFeatureFlag godoc
//
//	@Summary		Delete a feature flag
//	@Description	Delete a feature flag by ID
//	@Tags			Feature Flag
//	@Param			id	path		string	true	"Feature flag ID"
//	@Success		200	{object}	FeatureFlagResponse
//	@Router			/feature_flags/{id} [delete]
func (fh *FeatureFlagHandler) DeleteFeatureFlag(w http.ResponseWriter, r *http.Request) {
	flagID := chi.URLParam(r, "id")
	if flagID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Feature flag ID is required",
		})
		return
	}

	flagUUID, err := uuid.Parse(flagID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid feature flag ID",
		})
		return
	}

	if err := fh.db.DeleteFeatureFlag(flagUUID); err != nil {
		if err.Error() == "feature flag not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(FeatureFlagResponse{
				Success: false,
				Message: "Feature flag not found",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Failed to delete feature flag",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FeatureFlagResponse{
		Success: true,
		Message: "Feature flag deleted successfully",
	})
}

// AddFeatureFlagEndpoint godoc
//
//	@Summary		Add endpoints to a feature flag
//	@Description	Add new endpoints to an existing feature flag
//	@Tags			Feature Flag
//	@Param			feature_flag_id	path		string							true	"Feature flag ID"
//	@Param			endpoints		body		AddFeatureFlagEndpointRequest	true	"Endpoints to add"
//	@Success		201				{object}	FeatureFlagResponse
//	@Router			/feature_flags/{feature_flag_id}/endpoints [post]
func (fh *FeatureFlagHandler) AddFeatureFlagEndpoint(w http.ResponseWriter, r *http.Request) {
	flagID := chi.URLParam(r, "feature_flag_id")
	if flagID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Feature flag ID is required",
		})
		return
	}

	flagUUID, err := uuid.Parse(flagID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid feature flag ID",
		})
		return
	}

	var request AddFeatureFlagEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	for _, path := range request.Endpoints {
		_, err := fh.db.GetEndpointByPath(path)
		if err == nil {
			continue
		}

		endpoint := &db.Endpoint{
			UUID:            uuid.New(),
			Path:            path,
			FeatureFlagUUID: flagUUID,
		}

		_, err = fh.db.AddEndpoint(endpoint)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(FeatureFlagResponse{
				Success: false,
				Message: "Failed to create endpoint",
			})
			return
		}
	}

	updatedFlag, err := fh.db.GetFeatureFlagByUUID(flagUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Failed to fetch updated feature flag",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(FeatureFlagResponse{
		Success: true,
		Message: "Endpoints processed successfully",
		Data:    updatedFlag,
	})
}

// UpdateFeatureFlagEndpoint godoc
//
//	@Summary		Update an endpoint of a feature flag
//	@Description	Update the details of an endpoint of an existing feature flag
//	@Tags			Feature Flag
//	@Param			feature_flag_id	path		string								true	"Feature flag ID"
//	@Param			endpoint_id		path		string								true	"Endpoint ID"
//	@Param			endpoint		body		UpdateFeatureFlagEndpointRequest	true	"Updated endpoint details"
//	@Success		200				{object}	FeatureFlagResponse
//	@Router			/feature_flags/{feature_flag_id}/endpoints/{endpoint_id} [put]
func (fh *FeatureFlagHandler) UpdateFeatureFlagEndpoint(w http.ResponseWriter, r *http.Request) {
	flagID := chi.URLParam(r, "feature_flag_id")
	endpointID := chi.URLParam(r, "endpoint_id")
	if flagID == "" || endpointID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Feature flag ID and endpoint ID are required",
		})
		return
	}

	flagUUID, err := uuid.Parse(flagID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid feature flag ID",
		})
		return
	}

	_, err = fh.db.GetFeatureFlagByUUID(flagUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Feature flag not found",
		})
		return
	}

	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid endpoint ID",
		})
		return
	}

	var request UpdateFeatureFlagEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.Path == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Path is required",
		})
		return
	}

	endpoint := &db.Endpoint{
		UUID:            endpointUUID,
		Path:            request.Path,
		FeatureFlagUUID: flagUUID,
	}

	_, err = fh.db.UpdateEndpoint(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Failed to update endpoint",
		})
		return
	}

	updatedFlag, err := fh.db.GetFeatureFlagByUUID(flagUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Failed to fetch updated feature flag",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FeatureFlagResponse{
		Success: true,
		Message: "Endpoint updated successfully",
		Data:    updatedFlag,
	})
}

// DeleteFeatureFlagEndpoint godoc
//
//	@Summary		Delete an endpoint of a feature flag
//	@Description	Delete an endpoint of a feature flag by ID
//	@Tags			Feature Flag
//	@Param			feature_flag_id	path		string	true	"Feature flag ID"
//	@Param			endpoint_id		path		string	true	"Endpoint ID"
//	@Success		200				{object}	FeatureFlagResponse
//	@Router			/feature_flags/{feature_flag_id}/endpoints/{endpoint_id} [delete]
func (fh *FeatureFlagHandler) DeleteFeatureFlagEndpoint(w http.ResponseWriter, r *http.Request) {
	flagID := chi.URLParam(r, "feature_flag_id")
	endpointID := chi.URLParam(r, "endpoint_id")

	if flagID == "" || endpointID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Feature flag ID and endpoint ID are required",
		})
		return
	}

	flagUUID, err := uuid.Parse(flagID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid feature flag ID",
		})
		return
	}

	_, err = fh.db.GetFeatureFlagByUUID(flagUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Feature flag not found",
		})
		return
	}

	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Invalid endpoint ID",
		})
		return
	}

	endpoint, err := fh.db.GetEndpointByUUID(endpointUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Endpoint not found",
		})
		return
	}

	if endpoint.FeatureFlagUUID != flagUUID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Endpoint does not belong to this feature flag",
		})
		return
	}

	if err := fh.db.DeleteEndpoint(endpointUUID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FeatureFlagResponse{
			Success: false,
			Message: "Failed to delete endpoint",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FeatureFlagResponse{
		Success: true,
		Message: "Endpoint deleted successfully",
	})
}
