CREATE TABLE IF NOT EXISTS users(
    id uuid primary key,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    first_name text,
    last_name text,
    email text UNIQUE,
    phone VARCHAR(15) NOT NULL UNIQUE,
    email_verified timestamp,
    phone_verified timestamp,
    avatar text,
    password text
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);