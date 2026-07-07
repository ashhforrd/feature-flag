CREATE TABLE IF NOT EXISTS exposure_events (
    id UUID PRIMARY KEY,
    flag_key TEXT NOT NULL,
    user_id TEXT NOT NULL,
    enabled BOOLEAN NOT NULL,
    reason TEXT NOT NULL,
    bucket INTEGER,
    rollout_percentage INTEGER,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_exposure_events_flag_key_created_at
    ON exposure_events (flag_key, created_at DESC);