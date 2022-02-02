CREATE TABLE IF NOT EXISTS lease_charge_types(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  name VARCHAR(255) NOT NULL,
  description TEXT
);
DROP TYPE IF EXISTS lease_charge_types_enum;
CREATE TYPE lease_charge_types_enum AS ENUM(
  'RENT',
  'SECURITY_DEPOSIT',
  'OTHER_DEPOSIT',
  'FEE',
  'OTHER'
);
CREATE TABLE IF NOT EXISTS lease_charges(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  lease_id uuid,
  charge_type lease_charge_types_enum,
  lease_charge_type_id uuid,
  amount int NOT NULL,
  name VARCHAR(255) NOT NULL,
  due_date TIMESTAMP,
  description TEXT,
  CONSTRAINT fk_lease_charges_leases FOREIGN KEY(lease_id) REFERENCES leases(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_lease_charges_lease_id ON lease_charges(name);