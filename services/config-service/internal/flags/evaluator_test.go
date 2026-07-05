package flags

import "testing"

func TestEvaluateMissingFlagReturnsDefaultValue(t *testing.T) {
	result := Evaluate(nil, "new-checkout", map[string]any{}, true)

	if !result.Enabled {
		t.Fatalf("expected enabled true")
	}

	if result.Reason != ReasonFlagNotFound {
		t.Fatalf("expected reason %s, got %s", ReasonFlagNotFound, result.Reason)
	}
}

func TestEvaluateDisabledFlagReturnsFalse(t *testing.T) {
	flag := Flag{
		Key:     "new-checkout",
		Enabled: false,
	}

	result := Evaluate(&flag, "new-checkout", map[string]any{}, true)

	if result.Enabled {
		t.Fatalf("expected enabled false")
	}

	if result.Reason != ReasonFlagDisabled {
		t.Fatalf("expected reason %s, got %s", ReasonFlagDisabled, result.Reason)
	}
}

func TestEvaluateEnabledFlagReturnsDefaultRule(t *testing.T) {
	flag := Flag{
		Key:     "new-checkout",
		Enabled: true,
	}

	result := Evaluate(&flag, "new-checkout", map[string]any{}, false)

	if !result.Enabled {
		t.Fatalf("expected enabled true")
	}

	if result.Reason != ReasonDefaultRule {
		t.Fatalf("expected reason %s, got %s", ReasonDefaultRule, result.Reason)
	}
}

func TestEvaluateMatchedCountryRuleReturnsTrue(t *testing.T) {
	flag := Flag{
		Key:     "new-checkout",
		Enabled: true,
		TargetingRules: []TargetingRule{
			{
				Attribute: "country",
				Operator:  "equals",
				Value:     "ID",
			},
		},
	}

	result := Evaluate(&flag, "new-checkout", map[string]any{
		"id":      "user_123",
		"country": "ID",
	}, false)

	if !result.Enabled {
		t.Fatalf("expected enabled true")
	}

	if result.Reason != ReasonMatchedRule {
		t.Fatalf("expected reason %s, got %s", ReasonMatchedRule, result.Reason)
	}
}

func TestEvaluateEmailEndsWithRuleReturnsTrue(t *testing.T) {
	flag := Flag{
		Key:     "new-checkout",
		Enabled: true,
		TargetingRules: []TargetingRule{
			{
				Attribute: "email",
				Operator:  "ends_with",
				Value:     "@company.com",
			},
		},
	}

	result := Evaluate(&flag, "new-checkout", map[string]any{
		"id":    "user_123",
		"email": "alice@company.com",
	}, false)

	if !result.Enabled {
		t.Fatalf("expected enabled true")
	}

	if result.Reason != ReasonMatchedRule {
		t.Fatalf("expected reason %s, got %s", ReasonMatchedRule, result.Reason)
	}
}

func TestEvaluateNotEqualsRuleReturnsTrue(t *testing.T) {
	flag := Flag{
		Key:     "new-checkout",
		Enabled: true,
		TargetingRules: []TargetingRule{
			{
				Attribute: "country",
				Operator:  "not_equals",
				Value:     "SG",
			},
		},
	}

	result := Evaluate(&flag, "new-checkout", map[string]any{
		"id":      "user_123",
		"country": "ID",
	}, false)

	if !result.Enabled {
		t.Fatalf("expected enabled true")
	}

	if result.Reason != ReasonMatchedRule {
		t.Fatalf("expected reason %s, got %s", ReasonMatchedRule, result.Reason)
	}
}

func TestEvaluateContainsRuleReturnsTrue(t *testing.T) {
	flag := Flag{
		Key:     "beta-dashboard",
		Enabled: true,
		TargetingRules: []TargetingRule{
			{
				Attribute: "email",
				Operator:  "contains",
				Value:     "alice",
			},
		},
	}

	result := Evaluate(&flag, "beta-dashboard", map[string]any{
		"id":    "user_123",
		"email": "alice@example.com",
	}, false)

	if !result.Enabled {
		t.Fatalf("expected enabled true")
	}

	if result.Reason != ReasonMatchedRule {
		t.Fatalf("expected reason %s, got %s", ReasonMatchedRule, result.Reason)
	}
}

func TestEvaluateMissingUserAttributeDoesNotMatchRule(t *testing.T) {
	flag := Flag{
		Key:     "new-checkout",
		Enabled: true,
		TargetingRules: []TargetingRule{
			{
				Attribute: "country",
				Operator:  "equals",
				Value:     "ID",
			},
		},
	}

	result := Evaluate(&flag, "new-checkout", map[string]any{
		"id": "user_123",
	}, false)

	if result.Enabled {
		t.Fatalf("expected enabled false")
	}

	if result.Reason != ReasonDefaultRule {
		t.Fatalf("expected reason %s, got %s", ReasonDefaultRule, result.Reason)
	}
}

func TestEvaluatePercentageRolloutIsDeterministic(t *testing.T) {
	flag := Flag{
		Key:               "new-checkout",
		Enabled:           true,
		RolloutPercentage: 10,
	}

	user := map[string]any{
		"id": "user_123",
	}

	first := Evaluate(&flag, "new-checkout", user, false)
	second := Evaluate(&flag, "new-checkout", user, false)

	if first.Enabled != second.Enabled {
		t.Fatalf("expected deterministic rollout result")
	}

	if first.Reason != second.Reason {
		t.Fatalf("expected deterministic rollout reason")
	}
}
func TestEvaluateHundredPercentRolloutEnablesUser(t *testing.T) {
	flag := Flag{
		Key:               "new-checkout",
		Enabled:           true,
		RolloutPercentage: 100,
	}

	result := Evaluate(&flag, "new-checkout", map[string]any{
		"id": "user_123",
	}, false)

	if !result.Enabled {
		t.Fatalf("expected enabled true")
	}

	if result.Reason != ReasonPercentageRollout {
		t.Fatalf("expected reason %s, got %s", ReasonPercentageRollout, result.Reason)
	}
}
