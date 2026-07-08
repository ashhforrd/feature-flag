# Feature Flag Platform

A Feature Flag + Experimentation Platform.

The project demonstrates how a team can safely change production behavior without redeploying code: create a flag, evaluate it for a user, roll it out gradually, track exposures, record conversions, and compare experiment outcomes.

## What It Includes

- Config service written in Go
- PostgreSQL persistence for flags and analytics events
- Deterministic percentage rollout
- Attribute-based targeting rules
- JavaScript SDK for application developers, published as `@ashhforrd/feature-flags-js`
- Demo ecommerce checkout app using React + Vite
- Exposure tracking
- Conversion tracking
- Experiment result endpoint for conversion rate comparison

## Demo Story

An ecommerce team releases a risky `new-checkout` flow by:

1. Creating the `new-checkout` feature flag.
2. Enabling it for selected users or rollout percentages.
3. Showing classic checkout or one-page checkout in the demo app.
4. Recording exposure events when users are evaluated.
5. Recording conversion events when users complete checkout.
6. Comparing conversion rates between disabled and enabled variants.

## Repository Layout

```txt
apps/
  demo-ecommerce-app/     React demo app using the published JS SDK
  dashboard/              Dashboard for managing flags and viewing experiment results
services/
  config-service/         Go HTTP API and PostgreSQL repositories
packages/
  js-sdk/                 JavaScript SDK for evaluating flags
  shared/                 Earlier shared evaluator experiments
docs/                     Architecture and design notes
infra/                    Docker Compose infrastructure
```

## Architecture

```txt
Demo Ecommerce App
        |
        | uses
        v
JavaScript SDK
        |
        | HTTP
        v
Config Service
        |
        | SQL
        v
PostgreSQL
```

The config service is the source of truth for flags and analytics events. Product apps call the SDK, and the SDK talks to the config service.

## JavaScript SDK

The SDK is published on npm:

```sh
npm install @ashhforrd/feature-flags-js
```

Example usage:

```js
import { FeatureFlagClient } from "@ashhforrd/feature-flags-js"

const client = new FeatureFlagClient({
  baseUrl: "http://localhost:8080"
})

const enabled = await client.isEnabled("new-checkout", user, false)
```

## Quick Start

### 1. Start PostgreSQL

From the repository root:

```sh
docker compose -f infra/docker-compose.yml up -d
```

### 2. Run Migrations

```sh
docker compose -f infra/docker-compose.yml exec -T postgres \
  psql -U feature_flags -d feature_flags < services/config-service/migrations/001_create_flags.sql

docker compose -f infra/docker-compose.yml exec -T postgres \
  psql -U feature_flags -d feature_flags < services/config-service/migrations/002_create_exposure_events.sql

docker compose -f infra/docker-compose.yml exec -T postgres \
  psql -U feature_flags -d feature_flags < services/config-service/migrations/003_create_conversion_events.sql
```

### 3. Run Config Service

```sh
cd services/config-service
go test ./...
go run ./cmd/server
```

Health check:

```sh
curl http://localhost:8080/health
```

### 4. Create Demo Flag

```sh
curl -i -X POST http://localhost:8080/flags \
  -H "Content-Type: application/json" \
  -d '{
    "key": "new-checkout",
    "name": "New Checkout",
    "description": "Gradual rollout for redesigned checkout",
    "enabled": true,
    "rolloutPercentage": 100,
    "targetingRules": []
  }'
```

### 5. Run Demo App

In another terminal:

```sh
cd apps/demo-ecommerce-app
npm install
npm run dev
```

Open:

```txt
http://127.0.0.1:5173
```

Click the checkout button to record a conversion.

### 6. Check Experiment Results

```sh
curl http://localhost:8080/flags/new-checkout/results
```

Example response:

```json
{
  "flagKey": "new-checkout",
  "enabled": {
    "exposures": 3,
    "conversions": 2,
    "conversionRate": 0.6666666666666666
  },
  "disabled": {
    "exposures": 2,
    "conversions": 1,
    "conversionRate": 0.5
  }
}
```

## Useful Endpoints

```txt
GET    /health
POST   /flags
GET    /flags
GET    /flags/{key}
PATCH  /flags/{key}
POST   /flags/{key}/evaluate
GET    /flags/{key}/exposures
POST   /flags/{key}/conversions
GET    /flags/{key}/results
```

## Current Status

Done:

- Core evaluator
- Config service API
- PostgreSQL persistence
- JavaScript SDK
- Demo ecommerce app
- Exposure tracking
- Conversion tracking
- Experiment result calculation

Next:

- More complete analytics views
- Authentication and project/environment scoping
- Operational hardening and load testing
