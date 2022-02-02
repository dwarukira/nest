ALTER TABLE tenants
ADD COLUMN account_id uuid,
  ADD COLUMN created_by_id uuid,
  ADD CONSTRAINT fk_tenants_accounts FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_tenants_created_by FOREIGN KEY(created_by_id) REFERENCES users(id) ON DELETE
SET NULL;