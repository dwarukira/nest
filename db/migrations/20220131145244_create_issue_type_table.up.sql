CREATE TABLE IF NOT EXISTS issue_types(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  name VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL
);