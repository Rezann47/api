#!/bin/bash
# PostgreSQL container ilk başladığında çalışır.
# Test DB'si ve gerekli extension'ları oluşturur.
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    -- Extension'ları ana DB'ye ekle
    \c $POSTGRES_DB
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";
    CREATE EXTENSION IF NOT EXISTS "citext";

    -- Test DB'si (CI/integration test için)
    SELECT 'CREATE DATABASE yks_test' WHERE NOT EXISTS (
        SELECT FROM pg_database WHERE datname = 'yks_test'
    )\gexec

    \c yks_test
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";
    CREATE EXTENSION IF NOT EXISTS "citext";
EOSQL

echo "✓ DB init tamamlandı"
