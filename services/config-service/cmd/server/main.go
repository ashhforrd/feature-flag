package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ashhforrd/feature-flags/services/config-service/internal/flags"
)

func main() {
	repo := flags.NewRepository()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("POST /flags", func(w http.ResponseWriter, r *http.Request) {
		var req flags.CreateFlagRequest

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

		flag := flags.Flag{
			Key:               req.Key,
			Name:              req.Name,
			Description:       req.Description,
			Enabled:           req.Enabled,
			RolloutPercentage: req.RolloutPercentage,
			TargetingRules:    req.TargetingRules,
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		if err := repo.Create(flag); err != nil {
			if errors.Is(err, flags.ErrFlagAlreadyExists) {
				writeError(w, http.StatusConflict, "flag already exists")
				return
			}

			writeError(w, http.StatusInternalServerError, "failed to create flag")
			return
		}

		writeJSON(w, http.StatusCreated, flag)
	})

	mux.HandleFunc("GET /flags", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, repo.List())
	})

	mux.HandleFunc("GET /flags/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		flag, err := repo.GetByKey(key)
		if err != nil {
			if errors.Is(err, flags.ErrFlagNotFound) {
				writeError(w, http.StatusNotFound, "flag not found")
				return
			}

			writeError(w, http.StatusInternalServerError, "failed to get flag")
			return
		}

		writeJSON(w, http.StatusOK, flag)
	})

	mux.HandleFunc("PATCH /flags/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		flag, err := repo.GetByKey(key)
		if err != nil {
			if errors.Is(err, flags.ErrFlagNotFound) {
				writeError(w, http.StatusNotFound, "flag not found")
				return
			}

			writeError(w, http.StatusInternalServerError, "failed to get flag")
			return
		}

		var req flags.UpdateFlagRequest
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
		
		if err := repo.Update(flag); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update flag")
			return
		}

		writeJSON(w, http.StatusOK, flag)
	})

	mux.HandleFunc("POST /flags/{key}/evaluate", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		var req flags.EvaluateFlagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		defaultValue := false
		if req.DefaultValue != nil {
			defaultValue = *req.DefaultValue
		}

		flag, err := repo.GetByKey(key)
		if err != nil {
			if errors.Is(err, flags.ErrFlagNotFound) {
				writeJSON(w, http.StatusOK, flags.EvaluateFlagResponse{
					FlagKey: key,
					Enabled: defaultValue,
					Reason:  "FLAG_NOT_FOUND",
				})
				return
			}

			writeError(w, http.StatusInternalServerError, "failed to get flag")
			return
		}

		if !flag.Enabled {
			writeJSON(w, http.StatusOK, flags.EvaluateFlagResponse{
				FlagKey: flag.Key,
				Enabled: false,
				Reason:  "FLAG_DISABLED",
			})
			return
		}

		writeJSON(w, http.StatusOK, flags.EvaluateFlagResponse{
			FlagKey: flag.Key,
			Enabled: true,
			Reason:  "DEFAULT_RULE",
		})
	})

	log.Println(("config service listening on :8080"))
	log.Fatal(http.ListenAndServe(":8080", mux))
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
