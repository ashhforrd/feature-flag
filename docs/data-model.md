# Data Model

## Initial Flag Model

```txt
flags
- id
- key
- name
- description
- enabled
- rollout_percentage
- targeting_rules
- created_at
- updated_at
```

## Audit Log Model

```txt
audit_logs
- id
- actor_id
- action
- resource_type
- resource_id
- before
- after
- reason
- created_at
```

Audit logs should be append-only. Production flag changes should require a reason.

