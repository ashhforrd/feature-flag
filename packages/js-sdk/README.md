# JavaScript SDK

Small JavaScript client for evaluating feature flags from the config service.

The SDK wraps the HTTP evaluation API so application code can ask a simple question:

```js
const enabled = await client.isEnabled("new-checkout", user, false)
```

## Usage

```js
import { FeatureFlagClient } from "./src/index.js"

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

## Fail-Safe Behavior

If the config service is unavailable or returns a non-OK response, the SDK returns the caller-provided default value.

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

## Test

From `packages/js-sdk`:

```sh
npm test
```
