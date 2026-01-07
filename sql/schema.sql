-- schema.sql

-- Habilitar extensión para UUID (solo se ejecuta una vez)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE credentials (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Crear índice para user_id (aunque UNIQUE crea índice implícito)
CREATE INDEX idx_credentials_user_id ON credentials(user_id);

CREATE TABLE recovery_codes (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  code TEXT NOT NULL UNIQUE,
  is_valid BOOLEAN NOT NULL DEFAULT TRUE,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Crear índice para user_id en recovery_codes
CREATE INDEX idx_recovery_codes_user_id ON recovery_codes(user_id);


CREATE TABLE refresh_tokens (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL UNIQUE,
  refresh_token TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Crear índice para user_id en refresh_tokens
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);