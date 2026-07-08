# Feature Flags JavaScript SDK

JavaScript SDK for evaluating feature flags and recording conversion events against the Feature Flag Platform config service.

Use this package inside an application when you want to ask the platform whether a feature should be enabled for a specific user.

## Install

```sh
npm install @ashhforrd/feature-flags-js
```

For local development inside this monorepo, the demo app imports the SDK directly from `packages/js-sdk/src/index.js`.

## Quick Start

```js
import { FeatureFlagClient } from "@ashhforrd/feature-flags-js"

const client = new FeatureFlagClient({
  baseUrl: "http://localhost:8080"
})

const user = {
  id: "user-123",
  attributes: {
    country: "ID",
    email: "alice@example.com"
  }
}

const enabled = await client.isEnabled("new-checkout", user, false)

if (enabled) {
  // Show the new feature.
} else {
  // Show the existing experience.
}
```

## Evaluate a Flag

Use `evaluate` when the app needs the full decision response, including the reason and rollout metadata.

```js
const result = await client.evaluate("new-checkout", user, false)

console.log(result.enabled)
console.log(result.reason)
```

Example response:

```json
{
  "flagKey": "new-checkout",
  "enabled": true,
  "reason": "PERCENTAGE_ROLLOUT",
  "bucket": 42,
  "rolloutPercentage": 50
}
```

## Boolean Helper

Use `isEnabled` when the app only needs `true` or `false`.

```js
const enabled = await client.isEnabled("new-checkout", user, false)
```

## Record a Conversion

Use `recordConversion` after the user completes the business action being measured.

```js
const recorded = await client.recordConversion(
  "new-checkout",
  "user-123",
  "checkout_completed"
)
```

It returns `true` when the config service records the event successfully. It returns `false` if the request fails.

## Fail-Safe Behavior

Feature flag SDKs should not crash product experiences when the flag service is unavailable.

If `evaluate` fails because of a network error or non-OK response, it returns the caller-provided default value:

```js
const result = await client.evaluate("new-checkout", user, false)
```

Fallback response:

```json
{
  "flagKey": "new-checkout",
  "enabled": false,
  "reason": "SDK_REQUEST_FAILED"
}
```

`recordConversion` is also fail-safe. It returns `false` instead of throwing.

## API Reference

### `new FeatureFlagClient(options)`

Creates a client.

```js
const client = new FeatureFlagClient({
  baseUrl: "http://localhost:8080"
})
```

Options:

- `baseUrl`: config service base URL. Required.

### `client.evaluate(flagKey, user, defaultValue)`

Returns the full evaluation response from the config service.

Parameters:

- `flagKey`: unique flag key, for example `new-checkout`.
- `user`: user object with `id` and optional `attributes`.
- `defaultValue`: value to return when the flag service cannot make a decision.

### `client.isEnabled(flagKey, user, defaultValue)`

Returns only the boolean decision.

### `client.recordConversion(flagKey, userId, eventName)`

Records a conversion event for experiment results.

Parameters:

- `flagKey`: unique flag key.
- `userId`: user identifier.
- `eventName`: conversion event name, for example `checkout_completed`.

## Backend Requirement

This SDK expects the config service to expose these endpoints:

- `POST /flags/{key}/evaluate`
- `POST /flags/{key}/conversions`

## Development

Run tests from `packages/js-sdk`:

```sh
npm test
```

Check the package contents before publishing:

```sh
npm pack --dry-run
```
