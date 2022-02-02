CREATE TABLE IF NOT EXISTS accounts(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  name VARCHAR(255) NOT NULL,
  description TEXT
);


CREATE TABLE IF NOT EXISTS memberships(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  account_id uuid NOT NULL,
  user_id uuid NOT NULL,
  CONSTRAINT fk_membership_account FOREIGN KEY(account_id) REFERENCES accounts(id),
  CONSTRAINT fk_membership_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_accounts_name ON accounts(name);