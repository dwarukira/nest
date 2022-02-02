ALTER TABLE tenants DROP COLUMN account_id,
  DROP COLUMN created_by_id,
  DROP CONSTRAINT fk_tenants_accounts,
  DROP CONSTRAINT fk_tenants_created_by;