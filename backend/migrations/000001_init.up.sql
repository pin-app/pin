CREATE TABLE IF NOT EXISTS health_checks (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    remote_addr TEXT,
    user_agent TEXT
);