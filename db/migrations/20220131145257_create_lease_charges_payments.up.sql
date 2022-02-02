CREATE TABLE IF NOT EXISTS lease_charges_payments(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  payment_date TIMESTAMP,
  amount INT,
  lease_charge_id uuid NOT NULL,
  CONSTRAINT fk_lease_charges_payments_lease_charges FOREIGN KEY(lease_charge_id) REFERENCES lease_charges(id) ON DELETE CASCADE
)