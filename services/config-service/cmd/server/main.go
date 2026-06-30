package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Flag struct {
	Key               string          `json:"key"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Enabled           bool            `json:"enabled"`
	RolloutPercentage int             `json:"rolloutPercentage"`
	TargetingRules    []TargetingRule `json:"targetingRules"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}

type TargetingRule struct {
	Attribute string `json:"attribute"`
	Operator  string `json:"operator"`
	Value     any    `json:"value"`
}

type CreateFlagRequest struct {
	Key               string          `json:"key"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Enabled           bool            `json:"enabled"`
	RolloutPercentage int             `json:"rolloutPercentage"`
	TargetingRules    []TargetingRule `json:"targetingRules"`
}

type UpdateFlagRequest struct {
	Name              *string          `json:"name"`
	Description       *string          `json:"description"`
	Enabled           *bool            `json:"enabled"`
	RolloutPercentage *int             `json:"rolloutPercentage"`
	TargetingRules    *[]TargetingRule `json:"targetingRules"`
}

type EvaluateFlagRequest struct {
	User         map[string]any `json:"user"`
	DefaultValue *bool          `json:"defaultValue"`
}

type EvaluateFlagResponse struct {
	FlagKey string `json:"flagKey"`
	Enabled bool   `json:"enabled"`
	Reason  string `json:"reason"`
}

var flags = map[string]Flag{}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("POST /flags", func(w http.ResponseWriter, r *http.Request) {
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

		if _, exists := flags[req.Key]; exists {
			writeError(w, http.StatusConflict, "flag already exists")
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

		flags[req.Key] = flag

		writeJSON(w, http.StatusCreated, flag)
	})

	mux.HandleFunc("GET /flags", func(w http.ResponseWriter, r *http.Request) {
		result := make([]Flag, 0, len(flags))

		for _, flag := range flags {
			result = append(result, flag)
		}

		writeJSON(w, http.StatusOK, result)
	})

	mux.HandleFunc("GET /flags/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		flag, exists := flags[key]
		if !exists {
			writeError(w, http.StatusNotFound, "flag not found")
			return
		}

		writeJSON(w, http.StatusOK, flag)
	})

	mux.HandleFunc("PATCH /flags/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		flag, exists := flags[key]
		if !exists {
			writeError(w, http.StatusNotFound, "flag not found")
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
		flags[key] = flag

		writeJSON(w, http.StatusOK, flag)
	})

	mux.HandleFunc("POST /flags/{key}/evaluate", func(w http.ResponseWriter, r *http.Request) {
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

		flag, exists := flags[key]
		if !exists {
			writeJSON(w, http.StatusOK, EvaluateFlagResponse{
				FlagKey: key,
				Enabled: defaultValue,
				Reason:  "FLAG_DISABLED",
			})
			return
		}

		writeJSON(w, http.StatusOK, EvaluateFlagResponse{
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
