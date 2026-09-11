-- 020: TTL config + extraction metadata (change 014, step 01).
-- Additive: new tables only; no existing schema is touched.
-- The filename sorts after 019_session_relay.sql so it applies last.

CREATE TABLE IF NOT EXISTS ttl_config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS extraction_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT,
    content_hash TEXT,
    extracted_at TIMESTAMP,
    model TEXT
);
