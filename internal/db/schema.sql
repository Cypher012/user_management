CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email varchar(255) UNIQUE NOT NULL,
  password_hash varchar(255) NOT NULL,
  is_verified boolean NOT NULL DEFAULT false,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  refresh_token_hash varchar(255) UNIQUE NOT NULL,
  device_name varchar(255),
  device_info text,
  ip_address varchar(45),
  expires_at timestamp NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  last_used timestamp,
  revoked_at timestamp
);

CREATE TYPE user_role AS ENUM (
    'user',
    'admin'
);

CREATE TYPE user_gender AS ENUM (
    'unspecified',
    'male',
    'female'
);

CREATE TABLE profiles (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  name varchar(255) NOT NULL,
  phone_number varchar(100),
  role user_role NOT NULL DEFAULT 'user',
  gender user_gender NOT NULL DEFAULT 'unspecified',
  about_me text,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_user_revoked
  ON sessions (user_id, revoked_at);


CREATE TABLE email_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash text NOT NULL UNIQUE,
    type TEXT NOT NULL,
    expires_at timestamp NOT NULL,
    used_at timestamp,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE INDEX idx_email_token_user_type ON email_tokens(user_id, type);
