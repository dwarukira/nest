DROP TYPE IF EXISTS lease_status;
CREATE TYPE lease_status AS ENUM('DRAFT', 'ACTIVE', 'INACTIVE');
CREATE TABLE IF NOT EXISTS leases(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  lease_number VARCHAR(255) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE,
  monthly_rent INT,
  security_deposit INT,
  unit_id uuid NOT NULL,
  status lease_status,
  rent_due_day_of_month INT,
  CONSTRAINT fk_lease FOREIGN KEY(unit_id) REFERENCES units(id)
);