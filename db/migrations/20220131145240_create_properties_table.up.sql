CREATE TABLE IF NOT EXISTS properties(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  owner_id uuid NOT NULL,
  CONSTRAINT fk_user FOREIGN KEY(owner_id) REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_properties_name ON properties(name);
