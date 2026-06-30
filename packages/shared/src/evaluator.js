export const EvaluationReason = Object.freeze({
  FLAG_NOT_FOUND: "FLAG_NOT_FOUND",
  FLAG_DISABLED: "FLAG_DISABLED",
  MATCHED_RULE: "MATCHED_RULE",
  PERCENTAGE_ROLLOUT: "PERCENTAGE_ROLLOUT",
  DEFAULT_RULE: "DEFAULT_RULE"
});

export function evaluateFlag(flag, user, options = {}) {
  const defaultValue = options.defaultValue ?? false;

  if (!flag) {
    return {
      flagKey: options.flagKey,
      enabled: defaultValue,
      reason: EvaluationReason.FLAG_NOT_FOUND
    };
  }

  if (!flag.enabled) {
    return {
      flagKey: flag.key,
      enabled: false,
      reason: EvaluationReason.FLAG_DISABLED
    };
  }

  const matchedRule = firstMatchedRule(flag.targetingRules ?? [], user ?? {});
  if (matchedRule) {
    return {
      flagKey: flag.key,
      enabled: true,
      reason: EvaluationReason.MATCHED_RULE,
      matchedRule
    };
  }

  const rolloutPercentage = normalizeRolloutPercentage(flag.rolloutPercentage ?? 0);
  if (rolloutPercentage > 0) {
    const bucket = bucketUser(flag.key, user?.id);
    return {
      flagKey: flag.key,
      enabled: bucket < rolloutPercentage,
      reason: EvaluationReason.PERCENTAGE_ROLLOUT,
      bucket,
      rolloutPercentage
    };
  }

  return {
    flagKey: flag.key,
    enabled: defaultValue,
    reason: EvaluationReason.DEFAULT_RULE
  };
}

export function firstMatchedRule(rules, user) {
  return rules.find((rule) => matchesRule(rule, user));
}

export function matchesRule(rule, user) {
  const actual = user?.[rule.attribute];
  if (actual === undefined || actual === null) {
    return false;
  }

  switch (rule.operator) {
    case "equals":
      return actual === rule.value;
    case "not_equals":
      return actual !== rule.value;
    case "contains":
      return String(actual).includes(String(rule.value));
    case "ends_with":
      return String(actual).endsWith(String(rule.value));
    case "in":
      return Array.isArray(rule.value) && rule.value.includes(actual);
    default:
      return false;
  }
}

export function bucketUser(flagKey, userId) {
  if (!userId) {
    return 99;
  }

  return stableHash(`${flagKey}:${userId}`) % 100;
}

export function stableHash(input) {
  let hash = 2166136261;

  for (let index = 0; index < input.length; index += 1) {
    hash ^= input.charCodeAt(index);
    hash = Math.imul(hash, 16777619);
  }

  return hash >>> 0;
}

function normalizeRolloutPercentage(value) {
  if (!Number.isFinite(value)) {
    return 0;
  }

  return Math.max(0, Math.min(100, Math.floor(value)));
}

