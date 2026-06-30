# Feature Flags

Feature flags let a team change production behavior without redeploying code.

Instead of shipping a feature directly to every user, application code asks the feature flag platform whether a feature should be enabled for the current user:

```js
const enabled = await flags.isEnabled("new-checkout", user);
```

The platform evaluates the flag and returns a decision. That decision should be predictable, explainable, and safe when something goes wrong.

## Why Feature Flags Matter

Feature flags reduce release risk.

They allow teams to:

- Enable a feature for internal users first.
- Roll out gradually to a small percentage of users.
- Target users by attributes like country, plan, or email domain.
- Run experiments between control and treatment variants.
- Disable a feature quickly during an incident.
- Separate deploy from release.

The important mindset is this:

Deploying code should not mean exposing behavior to everyone immediately.

## Core Evaluation Flow

For the first version, flag evaluation should follow this order:

1. If the flag does not exist, return the caller-provided default value.
2. If no default value is provided, return `false`.
3. If the flag exists but is disabled, return `false`.
4. If any targeting rule matches, return `true`.
5. If percentage rollout is configured, use deterministic bucketing.
6. Otherwise return the default value.

This order is intentionally conservative. Missing or disabled flags should not accidentally expose risky behavior.

## Example Flag

```json
{
  "key": "new-checkout",
  "name": "New Checkout",
  "description": "Gradual rollout for the redesigned checkout flow",
  "enabled": true,
  "rolloutPercentage": 10,
  "targetingRules": [
    {
      "attribute": "email",
      "operator": "ends_with",
      "value": "@company.com"
    },
    {
      "attribute": "country",
      "operator": "equals",
      "value": "ID"
    }
  ]
}
```

## Targeting Rules

Targeting rules decide whether a flag should be enabled for a user based on user attributes.

Initial operators:

- `equals`
- `not_equals`
- `contains`
- `ends_with`
- `in`

For the first version, rules use OR semantics. If any rule matches, the flag is enabled.

Missing attributes should not match. For example, if a rule checks `country = ID` and the user has no `country`, the rule should return false.

## Percentage Rollout

Percentage rollout lets us expose a feature to a stable percentage of users.

Do not use random assignment per request. Randomness would cause the same user to move between enabled and disabled states, which creates inconsistent product behavior and corrupts experiment data.

Use deterministic hashing instead:

```txt
bucket = hash(flagKey + ":" + userId) % 100
```

If rollout is 10%, buckets `0..9` are enabled and buckets `10..99` are disabled.

This means:

- The same user gets the same result for the same flag.
- Increasing rollout from 10% to 25% keeps the original 10% enabled.
- Different flags can bucket the same user differently.

## Safe Defaults

The default behavior should be safe:

- Missing flag: return caller default, otherwise `false`.
- Disabled flag: always return `false`.
- Missing user ID for rollout: treat as not included.
- Unknown targeting operator: do not match.
- SDK cannot reach Config Service: return caller default.

This is the difference between a demo app and production thinking. A feature flag platform must behave well during failure.

## Evaluation Reasons

Every evaluation response should explain why the decision was made.

Example reasons:

- `FLAG_NOT_FOUND`
- `FLAG_DISABLED`
- `MATCHED_RULE`
- `PERCENTAGE_ROLLOUT`
- `DEFAULT_RULE`

Reasons are useful for debugging, support, demos, and interviews. They show that the system is explainable, not just functional.

## First Milestone

The first milestone is not the dashboard. It is the evaluator.

The evaluator is the core logic that both the Config Service and SDK will depend on.

Current implementation:

- [packages/shared/src/evaluator.js](/Users/macbookpro/Documents/Codex/2026-06-30/und/packages/shared/src/evaluator.js)

Current tests:

- [packages/shared/test/evaluator.test.js](/Users/macbookpro/Documents/Codex/2026-06-30/und/packages/shared/test/evaluator.test.js)

Run tests:

```sh
npm test
```

## Next Milestone

Wrap the evaluator in the Config Service MVP:

```txt
POST   /flags
GET    /flags
GET    /flags/:key
PATCH  /flags/:key
POST   /flags/:key/evaluate
```

Start with an in-memory repository. Once the API behavior is correct and tested, move persistence to PostgreSQL.

