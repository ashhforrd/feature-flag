# API Design

## Phase 1 Endpoints

```txt
POST   /flags
GET    /flags
GET    /flags/:key
PATCH  /flags/:key
POST   /flags/:key/evaluate
```

## Evaluation Request

```json
{
  "user": {
    "id": "user_123",
    "email": "alice@example.com",
    "country": "ID",
    "plan": "premium"
  },
  "defaultValue": false
}
```

## Evaluation Response

```json
{
  "flagKey": "new-checkout",
  "enabled": true,
  "reason": "PERCENTAGE_ROLLOUT",
  "bucket": 7,
  "rolloutPercentage": 10
}
```

## Safety Rule

Missing flags return the caller-provided default value. If no default is provided, return `false`.

