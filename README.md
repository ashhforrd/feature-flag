# Feature Flag Platform

A portfolio-focused Feature Flag + Experimentation Platform.

The goal is not just to build a CRUD app. The goal is to practice production engineering: deterministic behavior, safe defaults, failure handling, observability, load testing, and clear trade-off documentation.

## North Star

How can a team safely change production behavior without redeploying code?

## Demo Story

An e-commerce team releases a risky `new-checkout` flow by:

1. Creating a feature flag.
2. Enabling it for internal users.
3. Rolling it out to 1%, 10%, 25%, 50%, then 100%.
4. Tracking exposure, purchase, and checkout error events.
5. Comparing conversion rate between old and new checkout.
6. Disabling the flag quickly if errors spike.

## Repository Layout

```txt
apps/
  dashboard/
  demo-ecommerce-app/
services/
  config-service/
packages/
  shared/
docs/
  adr/
infra/
tests/
```

## Current Milestone

Phase 1 starts with the core feature evaluation rules before HTTP, database, or UI:

- Missing flag returns caller default, otherwise `false`.
- Disabled flag returns disabled.
- Targeting rules use OR semantics.
- Missing user attributes do not match.
- Percentage rollout uses deterministic hashing.

Run the current tests:

```sh
npm test
```

