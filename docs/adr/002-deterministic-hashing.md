# ADR: Deterministic Hashing For Rollout

## Context

Percentage rollout must assign the same user to the same result for the same flag across repeated evaluations.

## Decision

Use stable hashing over `flagKey:userId`, then map the hash into bucket `0..99`.

## Alternatives Considered

- `Math.random()` per request.
- Storing explicit rollout assignments for every user.

## Trade-offs

Hashing avoids storage and is stable, but changing the hash algorithm would reshuffle users.

## Consequences

The hash algorithm is part of product behavior and must be covered by tests.

