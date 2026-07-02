package flags

const (
	ReasonFlagNotFound      = "FLAG_NOT_FOUND"
	ReasonFlagDisabled      = "FLAG_DISABLED"
	ReasonMatchedRule       = "MATCHED_RULE"
	ReasonPercentageRollout = "PERCENTAGE_ROLLOUT"
	ReasonDefaultRule       = "DEFAULT_RULE"
)

func Evaluate(flag *Flag, flagKey string, defaultValue bool) EvaluateFlagResponse {
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

	return EvaluateFlagResponse{
		FlagKey: flag.Key,
		Enabled: true,
		Reason:  ReasonDefaultRule,
	}
}
