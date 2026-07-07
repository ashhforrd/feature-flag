# Config Service

HTTP service for creating, updating, listing, and evaluating feature flags.

This service stores flag configuration in PostgreSQL and exposes evaluation APIs that product applications can call at runtime.

## What It Does

- Create and update feature flags
- Store flags in PostgreSQL
- Evaluate flags for a user
- Support targeting rules by user attributes
- Support deterministic percentage rollout
- Return evaluation reasons for debugging

## Requirements

- Go 1.25+
- Docker
- Docker Compose

## Start PostgreSQL

From the repository root:

```sh
docker compose -f infra/docker-compose.yml up -d
```

Check that the container is running:

```sh
docker compose -f infra/docker-compose.yml ps
```

## Run Migration

From the repository root:

```sh
docker compose -f infra/docker-compose.yml exec -T postgres \
  psql -U feature_flags -d feature_flags < services/config-service/migrations/001_create_flags.sql
```

Check the table:

```sh
docker compose -f infra/docker-compose.yml exec postgres \
  psql -U feature_flags -d feature_flags -c "\d flags"
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

By default, the service connects to:

```txt
postgres://feature_flags:feature_flags@localhost:5432/feature_flags?sslmode=disable
```

You can override it with `DATABASE_URL`:

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

## Create Flag

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

## List Flags

```sh
curl http://localhost:8080/flags
```

## Get Flag

```sh
curl http://localhost:8080/flags/new-checkout
```

## Update Flag

```sh
curl -i -X PATCH http://localhost:8080/flags/new-checkout \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": false,
    "rolloutPercentage": 0
  }'
```

## Evaluate Flag

```sh
curl -i -X POST http://localhost:8080/flags/new-checkout/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "id": "user-123",
      "attributes": {
        "country": "ID",
        "email": "alice@example.com"
      }
    },
    "defaultValue": false
  }'
```

Example response when the flag is disabled:

```json
{
  "flagKey": "new-checkout",
  "enabled": false,
  "reason": "FLAG_DISABLED"
}
```

## Targeting Rules

A flag can include targeting rules based on user attributes.

Example:

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

Evaluation responses include a `reason` field to explain the decision.

Common reasons:

- `FLAG_NOT_FOUND`
- `FLAG_DISABLED`
- `MATCHED_RULE`
- `PERCENTAGE_ROLLOUT`
- `DEFAULT_RULE`

## Notes

PostgreSQL is the source of truth for flags. If the server restarts, flags remain available because they are persisted in the database.
