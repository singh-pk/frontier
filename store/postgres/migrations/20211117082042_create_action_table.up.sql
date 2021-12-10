CREATE TABLE IF NOT EXISTS actions (
  id  varchar PRIMARY KEY,
  name varchar UNIQUE NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  deleted_at timestamptz
);