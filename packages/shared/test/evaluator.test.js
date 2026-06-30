import assert from "node:assert/strict";
import test from "node:test";

import {
  EvaluationReason,
  bucketUser,
  evaluateFlag,
  matchesRule
} from "../src/evaluator.js";

test("missing flag returns caller default value", () => {
  const result = evaluateFlag(null, { id: "user_123" }, {
    flagKey: "new-checkout",
    defaultValue: true
  });

  assert.equal(result.enabled, true);
  assert.equal(result.reason, EvaluationReason.FLAG_NOT_FOUND);
});

test("missing flag defaults to false when caller default is not provided", () => {
  const result = evaluateFlag(null, { id: "user_123" }, {
    flagKey: "new-checkout"
  });

  assert.equal(result.enabled, false);
  assert.equal(result.reason, EvaluationReason.FLAG_NOT_FOUND);
});

test("disabled flag returns false even when targeting rule would match", () => {
  const result = evaluateFlag({
    key: "new-checkout",
    enabled: false,
    targetingRules: [
      { attribute: "country", operator: "equals", value: "ID" }
    ]
  }, {
    id: "user_123",
    country: "ID"
  });

  assert.equal(result.enabled, false);
  assert.equal(result.reason, EvaluationReason.FLAG_DISABLED);
});

test("targeting rules use OR semantics", () => {
  const result = evaluateFlag({
    key: "new-checkout",
    enabled: true,
    targetingRules: [
      { attribute: "country", operator: "equals", value: "SG" },
      { attribute: "plan", operator: "equals", value: "premium" }
    ]
  }, {
    id: "user_123",
    country: "ID",
    plan: "premium"
  });

  assert.equal(result.enabled, true);
  assert.equal(result.reason, EvaluationReason.MATCHED_RULE);
  assert.deepEqual(result.matchedRule, {
    attribute: "plan",
    operator: "equals",
    value: "premium"
  });
});

test("missing user attribute does not match a rule", () => {
  assert.equal(
    matchesRule({ attribute: "country", operator: "equals", value: "ID" }, {}),
    false
  );
});

test("supports first-version targeting operators", () => {
  const user = {
    email: "alice@company.com",
    country: "ID",
    plan: "premium"
  };

  assert.equal(matchesRule({ attribute: "country", operator: "equals", value: "ID" }, user), true);
  assert.equal(matchesRule({ attribute: "country", operator: "not_equals", value: "SG" }, user), true);
  assert.equal(matchesRule({ attribute: "email", operator: "contains", value: "alice" }, user), true);
  assert.equal(matchesRule({ attribute: "email", operator: "ends_with", value: "@company.com" }, user), true);
  assert.equal(matchesRule({ attribute: "plan", operator: "in", value: ["premium", "enterprise"] }, user), true);
});

test("percentage rollout is deterministic for the same flag and user", () => {
  const first = bucketUser("new-checkout", "user_123");
  const second = bucketUser("new-checkout", "user_123");

  assert.equal(first, second);
});

test("percentage rollout keeps enabled users when percentage increases", () => {
  const flag = {
    key: "new-checkout",
    enabled: true,
    rolloutPercentage: 10
  };
  const user = { id: "user_123" };
  const resultAt10 = evaluateFlag(flag, user);
  const resultAt25 = evaluateFlag({ ...flag, rolloutPercentage: 25 }, user);

  if (resultAt10.enabled) {
    assert.equal(resultAt25.enabled, true);
  }
});

test("zero rollout falls back to caller default", () => {
  const result = evaluateFlag({
    key: "new-checkout",
    enabled: true,
    rolloutPercentage: 0
  }, {
    id: "user_123"
  }, {
    defaultValue: false
  });

  assert.equal(result.enabled, false);
  assert.equal(result.reason, EvaluationReason.DEFAULT_RULE);
});


test("100 percent rollout enables any user with an id", () => {
  const result = evaluateFlag({
    key: "new-checkout",
    enabled: true,
    rolloutPercentage: 100
  }, {
    id: "user_123"
  });

  assert.equal(result.enabled, true)
  assert.equal(result.reason, EvaluationReason.PERCENTAGE_ROLLOUT)
} )