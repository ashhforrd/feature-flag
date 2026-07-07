package flags

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /flags", h.createFlag)
	mux.HandleFunc("GET /flags", h.listFlags)
	mux.HandleFunc("GET /flags/{key}", h.getFlag)
	mux.HandleFunc("PATCH /flags/{key}", h.updateFlag)
	mux.HandleFunc("POST /flags/{key}/evaluate", h.evaluateFlag)
}

func (h *Handler) createFlag(w http.ResponseWriter, r *http.Request) {
	var req CreateFlagRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Key == "" {
		writeError(w, http.StatusBadRequest, "key is required")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	if req.RolloutPercentage < 0 || req.RolloutPercentage > 100 {
		writeError(w, http.StatusBadRequest, "rolloutPercentage must be between 0 and 100")
		return
	}

	now := time.Now().UTC()

	flag := Flag{
		Key:               req.Key,
		Name:              req.Name,
		Description:       req.Description,
		Enabled:           req.Enabled,
		RolloutPercentage: req.RolloutPercentage,
		TargetingRules:    req.TargetingRules,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := h.repo.Create(flag); err != nil {
		if errors.Is(err, ErrFlagAlreadyExists) {
			writeError(w, http.StatusConflict, "flag already exists")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to create flag")
		return
	}

	writeJSON(w, http.StatusCreated, flag)
}

func (h *Handler) listFlags(w http.ResponseWriter, r *http.Request) {
	flags, err := h.repo.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list flags")
		return
	}

	writeJSON(w, http.StatusOK, flags)
}

func (h *Handler) getFlag(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	flag, err := h.repo.GetByKey(key)
	if err != nil {
		if errors.Is(err, ErrFlagNotFound) {
			writeError(w, http.StatusNotFound, "flag not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to get flag")
		return
	}

	writeJSON(w, http.StatusOK, flag)
}

func (h *Handler) updateFlag(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	flag, err := h.repo.GetByKey(key)
	if err != nil {
		if errors.Is(err, ErrFlagNotFound) {
			writeError(w, http.StatusNotFound, "flag not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to get flag")
		return
	}

	var req UpdateFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.RolloutPercentage != nil {
		if *req.RolloutPercentage < 0 || *req.RolloutPercentage > 100 {
			writeError(w, http.StatusBadRequest, "rolloutPercentage must be between 0 and 100")
			return
		}
		flag.RolloutPercentage = *req.RolloutPercentage
	}

	if req.Name != nil {
		flag.Name = *req.Name
	}

	if req.Description != nil {
		flag.Description = *req.Description
	}

	if req.Enabled != nil {
		flag.Enabled = *req.Enabled
	}

	if req.TargetingRules != nil {
		flag.TargetingRules = *req.TargetingRules
	}

	flag.UpdatedAt = time.Now().UTC()

	if err := h.repo.Update(flag); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update flag")
		return
	}

	writeJSON(w, http.StatusOK, flag)
}

func (h *Handler) evaluateFlag(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	var req EvaluateFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	defaultValue := false
	if req.DefaultValue != nil {
		defaultValue = *req.DefaultValue
	}

	flag, err := h.repo.GetByKey(key)
	if err != nil {
		if errors.Is(err, ErrFlagNotFound) {
			result := Evaluate(nil, key, req.User, defaultValue)
			writeJSON(w, http.StatusOK, result)
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to get flag")
		return
	}

	result := Evaluate(&flag, key, req.User, defaultValue)
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
