package flags

import "strings"

const (
	ReasonFlagNotFound      = "FLAG_NOT_FOUND"
	ReasonFlagDisabled      = "FLAG_DISABLED"
	ReasonMatchedRule       = "MATCHED_RULE"
	ReasonPercentageRollout = "PERCENTAGE_ROLLOUT"
	ReasonDefaultRule       = "DEFAULT_RULE"
)

func Evaluate(flag *Flag, flagKey string, user map[string]any, defaultValue bool) EvaluateFlagResponse {
	if flag == nil {
		return EvaluateFlagResponse{
			FlagKey: flagKey,
			Enabled: defaultValue,
			Reason:  ReasonFlagNotFound,
		}
	}

	if !flag.Enabled {
		return EvaluateFlagResponse{
			FlagKey: flag.Key,
			Enabled: false,
			Reason:  ReasonFlagDisabled,
		}
	}

	if len(flag.TargetingRules) == 0 {
		return EvaluateFlagResponse{
			FlagKey: flag.Key,
			Enabled: true,
			Reason:  ReasonDefaultRule,
		}
	}

	matchedRule := firstMatchedRule(flag.TargetingRules, user)
	if matchedRule != nil {
		return EvaluateFlagResponse{
			FlagKey: flag.Key,
			Enabled: true,
			Reason:  ReasonMatchedRule,
		}
	}

	return EvaluateFlagResponse{
		FlagKey: flag.Key,
		Enabled: defaultValue,
		Reason:  ReasonDefaultRule,
	}
}

func firstMatchedRule(rules []TargetingRule, user map[string]any) *TargetingRule {
	for _, rule := range rules {
		if matchesRule(rule, user) {
			return &rule
		}
	}

	return nil
}

func matchesRule(rule TargetingRule, user map[string]any) bool {
	actual, exists := user[rule.Attribute]
	if !exists || actual == nil {
		return false
	}

	switch rule.Operator {
	case "equals":
		return actual == rule.Value
	case "not_equals":
		return actual != rule.Value
	case "contains":
		actualString, ok := actual.(string)
		if !ok {
			return false
		}

		valueString, ok := rule.Value.(string)
		if !ok {
			return false
		}

		return strings.Contains(actualString, valueString)
	case "ends_with":
		actualString, ok := actual.(string)
		if !ok {
			return false
		}

		valueString, ok := rule.Value.(string)
		if !ok {
			return false
		}

		return strings.HasSuffix(actualString, valueString)
	default:
		return false
	}
}
