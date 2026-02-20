-- migrations/002_create_products.up.sql
CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(200)    NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2)  NOT NULL,
    stock       INTEGER         NOT NULL DEFAULT 0,
    user_id     BIGINT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_products_user_id    ON products(user_id);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);
