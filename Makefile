# Makefile

.PHONY: run build test migrate-up migrate-down docker-up docker-down tidy

## Uygulamayı çalıştır
run:
	go run ./cmd/main.go

## Binary derle
build:
	go build -o bin/app ./cmd/main.go

## Bağımlılıkları indir
tidy:
	go mod tidy

## Test çalıştır
test:
	go test ./... -v -cover

## PostgreSQL'i Docker ile başlat
docker-up:
	docker run --name go-crud-db \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=secret \
		-e POSTGRES_DB=go_crud_db \
		-p 5432:5432 \
		-d postgres:15-alpine

## Docker container'ı durdur
docker-down:
	docker stop go-crud-db && docker rm go-crud-db

## Manuel migration (GORM AutoMigrate yerine)
migrate-up:
	psql $(DATABASE_URL) -f migrations/001_create_users.up.sql
	psql $(DATABASE_URL) -f migrations/002_create_products.up.sql

## Rollback
migrate-down:
	psql $(DATABASE_URL) -f migrations/rollback.sql

## Linter
lint:
	golangci-lint run ./...
