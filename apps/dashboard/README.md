# Dashboard

Minimal operations dashboard for the feature flag platform.

The dashboard lets you:

- List feature flags
- Select a flag
- Toggle enabled state
- Change rollout percentage
- View exposure counts
- View enabled vs disabled conversion rates

Numbers are animated with `react-countup`.

## Requirements

- Config service running on `http://localhost:8080`
- PostgreSQL migrations applied
- At least one flag created

## Run

```sh
npm install
npm run dev
```

Open:

```txt
http://127.0.0.1:5174
```

By default the dashboard calls `/api`, which Vite proxies to:

```txt
http://localhost:8080
```

Override the config service URL with:

```sh
VITE_CONFIG_SERVICE_URL="http://localhost:8080" npm run dev
```

## Build

```sh
npm run build
```
