# Demo E-commerce App

Minimal checkout demo that uses the JavaScript SDK to evaluate the `new-checkout` feature flag and record checkout conversions.

The app shows one of two checkout experiences:

- Classic checkout when `new-checkout` is disabled
- One-page checkout when `new-checkout` is enabled for the selected user

Clicking the checkout button records a `checkout_completed` conversion event.

## Requirements

- PostgreSQL running through Docker Compose
- Config service running on `http://localhost:8080`
- `new-checkout` flag created in the config service

## Run

From `apps/demo-ecommerce-app`:

```sh
npm install
npm run dev
```

Open:

```txt
http://127.0.0.1:5173
```

By default the app calls `/api`, which Vite proxies to:

```txt
http://localhost:8080
```

Override the config service URL with:

```sh
VITE_CONFIG_SERVICE_URL="http://localhost:8080" npm run dev
```

If you override the URL directly, the config service must allow browser CORS requests.

## Demo Flow

1. Start Postgres.
2. Run config service migrations.
3. Start the config service.
4. Create or update the `new-checkout` flag.
5. Start this app.
6. Switch customer profiles to trigger different evaluations.
7. Click the checkout button to record conversions.
8. Check experiment results from the config service.

## Toggle Classic Checkout

```sh
curl -X PATCH http://localhost:8080/flags/new-checkout \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": false,
    "rolloutPercentage": 0
  }'
```

Refresh the app. It should show classic checkout.

## Toggle New Checkout

```sh
curl -X PATCH http://localhost:8080/flags/new-checkout \
  -H "Content-Type: application/json" \
  -d '{
    "enabled": true,
    "rolloutPercentage": 100
  }'
```

Refresh the app. It should show one-page checkout.

## Check Results

```sh
curl http://localhost:8080/flags/new-checkout/results
```

The enabled and disabled groups show exposures, conversions, and conversion rates.

## Build

```sh
npm run build
```
