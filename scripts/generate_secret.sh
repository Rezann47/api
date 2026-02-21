#!/bin/sh
# JWT secret üretir
# Kullanım: sh scripts/generate_secret.sh

echo "JWT_ACCESS_SECRET=$(openssl rand -hex 32)"
echo "JWT_REFRESH_SECRET=$(openssl rand -hex 32)"
