CREATE TABLE IF NOT EXISTS schema_migrations (
  version BIGINT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS items (
  id TEXT PRIMARY KEY,
  url TEXT NOT NULL,
  title TEXT NOT NULL,
  site TEXT NOT NULL,
  author TEXT NOT NULL DEFAULT '',
  published_at TEXT NOT NULL DEFAULT '',
  saved_reason TEXT NOT NULL,
  summary TEXT NOT NULL,
  read_time_minutes INTEGER NOT NULL,
  recommended_mode TEXT NOT NULL,
  best_sections TEXT[] NOT NULL DEFAULT '{}',
  read_if TEXT[] NOT NULL DEFAULT '{}',
  skip_if TEXT[] NOT NULL DEFAULT '{}',
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS items_created_at_idx ON items (created_at DESC);
CREATE INDEX IF NOT EXISTS items_status_idx ON items (status);
