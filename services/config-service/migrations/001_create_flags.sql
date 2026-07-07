CREATE TABLE IF NOT EXISTS flags (
    id UUID PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    enabled BOOLEAN NOT NULL DEFAULT false,
    rollout_percentage INTEGER NOT NULL DEFAULT 0,
    targeting_rules JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT flags_key_not_empty CHECK (length(trim(key)) > 0),
    CONSTRAINT flags_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT flags_rollout_percentage_range CHECK (
        rollout_percentage >= 0 AND rollout_percentage <= 100
    )
);