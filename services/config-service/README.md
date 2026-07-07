# Config Service

Go HTTP service for managing feature flags, evaluating them for users, and recording basic experimentation analytics.

The service stores flag configuration, exposure events, and conversion events in PostgreSQL.

## What It Does

- Create, list, get, and update feature flags
- Evaluate flags for a user
- Support targeting rules by user attributes
- Support deterministic percentage rollout
- Record exposure events during evaluation
- Record conversion events from product apps
- Return exposure summaries and experiment results

## Requirements

- Go 1.25+
- Docker
- Docker Compose

## Start PostgreSQL

From the repository root:

```sh
docker compose -f infra/docker-compose.yml up -d
```

Check the container:

```sh
docker compose -f infra/docker-compose.yml ps
```

## Run Migrations

From the repository root:

```sh
docker compose -f infra/docker-compose.yml exec -T postgres \
  psql -U feature_flags -d feature_flags < services/config-service/migrations/001_create_flags.sql

docker compose -f infra/docker-compose.yml exec -T postgres \
  psql -U feature_flags -d feature_flags < services/config-service/migrations/002_create_exposure_events.sql

docker compose -f infra/docker-compose.yml exec -T postgres \
  psql -U feature_flags -d feature_flags < services/config-service/migrations/003_create_conversion_events.sql
```

Check tables:

```sh
docker compose -f infra/docker-compose.yml exec postgres \
  psql -U feature_flags -d feature_flags -c "\dt"
```

## Run Tests

From `services/config-service`:

```sh
go test ./...
```

## Run Server

From `services/config-service`:

```sh
go run ./cmd/server
```

Default database URL:

```txt
postgres://feature_flags:feature_flags@localhost:5432/feature_flags?sslmode=disable
```

Override it with:

```sh
DATABASE_URL="postgres://feature_flags:feature_flags@localhost:5432/feature_flags?sslmode=disable" \
  go run ./cmd/server
```

## Health Check

```sh
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

## Flags

### Create Flag

```sh
curl -i -X POST http://localhost:8080/flags \
  -H "Content-Type: application/json" \
  -d '{
    "key": "new-checkout",
    "name": "New Checkout",
    "description": "Gradual rollout for redesigned checkout",
    "enabled": true,
    "rolloutPercentage": 50,
    "targetingRules": []
  }'
```

### List Flags

```sh
curl http://localhost:8080/flags
```

### Get Flag

```sh
curl http://localhost:8080/flags/new-checkout
```

### Update Flag

```sh
curl -i -X PATCH http://localhost:8080/flags/new-checkout \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": true,
    "rolloutPercentage": 100
  }'
```

## Evaluation

```sh
curl -i -X POST http://localhost:8080/flags/new-checkout/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "id": "user-123",
      "country": "ID",
      "email": "alice@example.com"
    },
    "defaultValue": false
  }'
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

Evaluation records an exposure event when the request includes a user id.

## Analytics

### Exposure Summary

```sh
curl http://localhost:8080/flags/new-checkout/exposures
```

Example response:

```json
{
  "flagKey": "new-checkout",
  "total": 12,
  "enabled": 7,
  "disabled": 5
}
```

### Record Conversion

```sh
curl -i -X POST http://localhost:8080/flags/new-checkout/conversions \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user-123",
    "eventName": "checkout_completed"
  }'
```

Expected response:

```json
{"status":"recorded"}
```

### Experiment Results

```sh
curl http://localhost:8080/flags/new-checkout/results
```

Example response:

```json
{
  "flagKey": "new-checkout",
  "enabled": {
    "exposures": 10,
    "conversions": 3,
    "conversionRate": 0.3
  },
  "disabled": {
    "exposures": 8,
    "conversions": 1,
    "conversionRate": 0.125
  }
}
```

## Targeting Rules

A flag can include targeting rules based on user attributes.

```json
{
  "attribute": "country",
  "operator": "equals",
  "value": "ID"
}
```

Supported operators:

- `equals`
- `not_equals`
- `contains`
- `ends_with`

Rules use OR semantics: if any rule matches, the flag is enabled for the user.

## Evaluation Reasons

- `FLAG_NOT_FOUND`
- `FLAG_DISABLED`
- `MATCHED_RULE`
- `PERCENTAGE_ROLLOUT`
- `DEFAULT_RULE`

## Notes

PostgreSQL is the source of truth. If the server restarts, flags and analytics events remain available.
