# ADR: SDK Local Cache

## Context

Feature flag checks may run inside production request paths. Calling the Config Service on every request would add latency and make client applications depend on Config Service availability.

## Decision

The SDK should fetch configuration periodically, cache it locally, and evaluate flags locally when possible.

## Alternatives Considered

- Server-side evaluation on every request.
- Hybrid evaluation where the SDK can call the service if a flag is missing locally.

## Trade-offs

Local evaluation is fast and resilient, but config can be stale until the next refresh.

## Consequences

The evaluator must be deterministic and portable so the service and SDK can share behavior.

