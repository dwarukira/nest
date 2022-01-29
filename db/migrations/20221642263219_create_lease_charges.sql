CREATE TABLE IF NOT EXISTS lease_charges(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  amount INT,
  dueDate DATE,
  leaseId uuid NOT NULL,
);