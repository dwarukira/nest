CREATE TABLE IF NOT EXISTS units(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  default_rent INT DEFAULT 0,
  property_id uuid NOT NULL,
  CONSTRAINT fk_property FOREIGN KEY(property_id) REFERENCES properties(id)
);
CREATE INDEX IF NOT EXISTS idx_units_name ON units(name);