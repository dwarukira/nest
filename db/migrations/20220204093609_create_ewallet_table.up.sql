CREATE TABLE IF NOT EXISTS wallet(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  user_id uuid NOT NULL,
  balance INT NOT NULL,
  CONSTRAINT fk_wallet_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS wallet_transaction(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  wallet_id uuid NOT NULL,
  amount decimal NOT NULL,
  description TEXT,
  CONSTRAINT fk_wallet_transaction_wallet FOREIGN KEY(wallet_id) REFERENCES wallet(id)
);
