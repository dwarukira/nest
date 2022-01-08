CREATE TABLE IF NOT EXISTS tickets_status(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  color VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS tickets(
  id uuid primary key,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  images TEXT[],
  ticket_status_id uuid,
  property_id uuid,
  unit_id uuid,
  CONSTRAINT fk_tickets_status FOREIGN KEY(ticket_status_id) REFERENCES tickets_status(id)
);

CREATE INDEX IF NOT EXISTS idx_tickets_title ON tickets(title);
CREATE INDEX IF NOT EXISTS idx_tickets_title ON tickets(description);