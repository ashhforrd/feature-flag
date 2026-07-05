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