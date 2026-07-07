# JavaScript SDK

Small JavaScript client for evaluating feature flags and recording conversions through the config service.

The SDK lets application code ask:

```js
const enabled = await client.isEnabled("new-checkout", user, false)
```

and record outcomes:

```js
await client.recordConversion("new-checkout", user.id, "checkout_completed")
```

## Usage

```js
import { FeatureFlagClient } from "./src/index.js"

const client = new FeatureFlagClient({
  baseUrl: "http://localhost:8080"
})

const user = {
  id: "user-123",
  country: "ID",
  email: "alice@example.com"
}

const enabled = await client.isEnabled("new-checkout", user, false)

if (enabled) {
  console.log("Show new checkout")
} else {
  console.log("Show old checkout")
}
```

## Evaluate With Reason

Use `evaluate` when the caller needs the full decision response.

```js
const result = await client.evaluate("new-checkout", user, false)

console.log(result)
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

## Record Conversion

Use `recordConversion` after a user completes the business action being measured.

```js
const recorded = await client.recordConversion(
  "new-checkout",
  "user-123",
  "checkout_completed"
)
```

Returns:

```js
true
```

if the config service records the event successfully, otherwise:

```js
false
```

## Fail-Safe Behavior

If the config service is unavailable or returns a non-OK response, `evaluate` returns the caller-provided default value.

```js
const enabled = await client.isEnabled("new-checkout", user, false)
```

If the request fails, `enabled` will be `false`.

The full `evaluate` response will look like this:

```json
{
  "flagKey": "new-checkout",
  "enabled": false,
  "reason": "SDK_REQUEST_FAILED"
}
```

Conversion recording is also fail-safe: it returns `false` instead of throwing when the request fails.

## API

### `new FeatureFlagClient(options)`

Creates a client.

```js
const client = new FeatureFlagClient({
  baseUrl: "http://localhost:8080"
})
```

Options:

- `baseUrl`: base URL of the config service

### `client.isEnabled(flagKey, user, defaultValue)`

Returns only the boolean flag decision.

```js
const enabled = await client.isEnabled("new-checkout", user, false)
```

### `client.evaluate(flagKey, user, defaultValue)`

Returns the full evaluation response from the config service.

```js
const result = await client.evaluate("new-checkout", user, false)
```

### `client.recordConversion(flagKey, userId, eventName)`

Records a conversion event.

```js
const recorded = await client.recordConversion(
  "new-checkout",
  "user-123",
  "checkout_completed"
)
```

## Test

From `packages/js-sdk`:

```sh
npm test
```
