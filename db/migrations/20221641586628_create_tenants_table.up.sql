CREATE TABLE IF NOT EXISTS tenants(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  email VARCHAR(255),
  phone_number VARCHAR(15),
  invite_token TEXT,
  invite_accepted TIMESTAMP,
  invite_sent TIMESTAMP,
  lease_id uuid NOT NULL,
  CONSTRAINT fk_tenant FOREIGN KEY(lease_id) REFERENCES leases(id)
);