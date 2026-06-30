# System Design

## Initial Architecture

```txt
Demo App -> TypeScript SDK -> Config Service -> PostgreSQL
```

The first implementation should keep the system simple. The Config Service owns flag metadata and can evaluate a flag for a user. The SDK will later cache config and evaluate locally to avoid adding latency or availability risk to client applications.

## Later Architecture

```txt
Demo App -> TypeScript SDK -> Config Service -> PostgreSQL
Demo App / SDK -> Event Ingestion API -> Queue -> Analytics Worker -> Analytics Storage -> Dashboard
```

## Core Design Pressure

Feature flags sit directly in production request paths. The system must prefer predictable, explainable behavior over cleverness.

Important failure questions:

- What if the Config Service is down?
- What if the SDK cache is stale?
- What if a flag does not exist?
- What if a user is missing an attribute?
- What if rollout changes from 10% to 25%?
- What if event ingestion receives duplicates?

