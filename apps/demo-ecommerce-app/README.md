# Demo E-commerce App

Minimal checkout demo that uses the JavaScript SDK to evaluate the `new-checkout` feature flag.

The app shows one of two checkout experiences:

- Classic checkout when `new-checkout` is disabled
- One-page checkout when `new-checkout` is enabled for the selected user

## Requirements

- Config service running on `http://localhost:8080`
- `new-checkout` flag created in the config service

## Run

```sh
npm install
npm run dev
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

1. Start Postgres and the config service.
2. Create or update the `new-checkout` flag.
3. Start this app.
4. Switch customer profiles to see how targeting and rollout affect the checkout UI.
