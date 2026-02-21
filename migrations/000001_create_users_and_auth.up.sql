-- ============================================================
-- 000001: Users & Auth
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";  -- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "citext";    -- case-insensitive email

-- ─── ENUM TYPES ───────────────────────────────────────────
CREATE TYPE user_role AS ENUM ('student', 'instructor');

-- ─── USERS ────────────────────────────────────────────────
CREATE TABLE users (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(100) NOT NULL,
    email          CITEXT      NOT NULL,
    password_hash  VARCHAR(255) NOT NULL,
    role           user_role   NOT NULL DEFAULT 'student',
    student_code   VARCHAR(20)  UNIQUE,           -- NULL for instructors
    is_premium     BOOLEAN      NOT NULL DEFAULT FALSE,
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ                    -- soft delete

    -- CITEXT handles case-insensitive email uniqueness automatically
    -- We'll add partial index below for active users
);

-- Partial unique index: email unique among non-deleted users
CREATE UNIQUE INDEX idx_users_email_active
    ON users (email) WHERE deleted_at IS NULL;

-- Partial unique index: student_code unique among non-deleted
CREATE UNIQUE INDEX idx_users_student_code_active
    ON users (student_code) WHERE student_code IS NOT NULL AND deleted_at IS NULL;

-- Role filter (eğitmenler / öğrenciler ayrı listelenirken)
CREATE INDEX idx_users_role ON users (role) WHERE deleted_at IS NULL;

-- ─── REFRESH TOKENS ───────────────────────────────────────
-- Her token satırı bir cihaz/session'ı temsil eder
CREATE TABLE refresh_tokens (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash  CHAR(64)    NOT NULL UNIQUE,   -- SHA-256 hex of token
    user_agent  VARCHAR(500),                  -- cihaz bilgisi
    ip_address  INET,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id  ON refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_expires  ON refresh_tokens (expires_at)
    WHERE revoked_at IS NULL;   -- sadece aktif tokenları tara

-- ─── AUTO-UPDATE updated_at trigger (tüm tablolarda kullanılacak) ──
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
