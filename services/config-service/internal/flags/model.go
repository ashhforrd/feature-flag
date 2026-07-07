package flags

import "time"

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
	FlagKey           string `json:"flagKey"`
	Enabled           bool   `json:"enabled"`
	Reason            string `json:"reason"`
	Bucket            *int   `json:"bucket,omitempty"`
	RolloutPercentage *int   `json:"rolloutPercentage,omitempty"`
}

type ExposureEvent struct {
	FlagKey           string
	UserID            string
	Enabled           bool
	Reason            string
	Bucket            *int
	RolloutPercentage *int
	CreatedAt         time.Time
}

type ExposureSummary struct {
	FlagKey  string `json:"flagKey"`
	Total    int    `json:"total"`
	Enabled  int    `json:"enabled"`
	Disabled int    `json:"disabled"`
}
