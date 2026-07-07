CREATE TABLE IF NOT EXISTS conversion_events (
    id UUID PRIMARY KEY,
    flag_key TEXT NOT NULL,
    user_id TEXT NOT NULL,
    event_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_conversion_events_flag_key_created_at
    ON conversion_events (flag_key, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_conversion_events_flag_key_user_id
    ON conversion_events (flag_key, user_id);