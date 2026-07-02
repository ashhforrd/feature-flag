package flags

import "testing"

func TestEvaluateMissingFlagReturnsDefaultValue(t *testing.T) {
	result := Evaluate(nil, "new-checkout", true)

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

	result := Evaluate(&flag, "new-checkout", true)

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

	result := Evaluate(&flag, "new-checkout", false)

	if !result.Enabled {
		t.Fatalf("expected enabled true")
	}

	if result.Reason != ReasonDefaultRule {
		t.Fatalf("expected reason %s, got %s", ReasonDefaultRule, result.Reason)
	}
}