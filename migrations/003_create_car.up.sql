-- migrations/001_create_users.up.sql
CREATE TABLE IF NOT EXISTS cars (
    id         BIGSERIAL PRIMARY KEY,
    make       VARCHAR(100)  NOT NULL,
    model      VARCHAR(150)  NOT NULL ,
    year      INT           NOT NULL,

    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ   -- Soft delete için NULL olabilir
);

CREATE INDEX IF NOT EXISTS idx_cars_make      ON cars(make);
CREATE INDEX IF NOT EXISTS idx_cars_model     ON cars(model);
CREATE INDEX IF NOT EXISTS idx_cars_year      ON cars(year);
CREATE INDEX IF NOT EXISTS idx_cars_deleted_at ON cars(deleted_at);