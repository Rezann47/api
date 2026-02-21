-include .env
export

MIGRATE     = go run ./cmd/migrate
DC          = docker compose
DC_PROD     = docker compose -f docker-compose.yml -f docker-compose.prod.yml

.PHONY: help \
        dev build test lint swagger \
        db-start db-stop db-logs db-shell \
        docker-build docker-up docker-down docker-logs docker-shell \
        docker-dev docker-prod \
        migrate-up migrate-down migrate-version migrate-steps migrate-force migrate-new \
        dc-migrate-up dc-migrate-down dc-migrate-version \
        pgadmin clean

help: ## Tüm komutları listele
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-28s\033[0m %s\n",$$1,$$2}'

# ════════════════════════════════════════════════════════
# LOCAL DEV (Docker olmadan)
# ════════════════════════════════════════════════════════
dev: ## API'yi local çalıştır (gerekli: .env + PostgreSQL)
	go run ./cmd/api

build: ## Binary'leri oluştur → bin/
	@mkdir -p bin
	CGO_ENABLED=0 go build -o bin/api     ./cmd/api
	CGO_ENABLED=0 go build -o bin/migrate ./cmd/migrate
	@echo "✓ bin/api  bin/migrate hazır"

test: ## Testleri çalıştır
	go test ./... -v -race -count=1

lint: ## Linter (golangci-lint gerekli)
	golangci-lint run ./...

swagger: ## Swagger dokümanı üret (swag gerekli)
	swag init -g cmd/api/main.go -o docs
	@echo "✓ docs/ güncellendi → http://localhost:8080/swagger/index.html"

# ════════════════════════════════════════════════════════
# DATABASE (sadece PostgreSQL)
# ════════════════════════════════════════════════════════
db-start: ## Sadece PostgreSQL container'ı başlat
	$(DC) up -d db
	@echo "⏳ DB hazır olana kadar bekleniyor..."
	@until $(DC) exec db pg_isready -U $${DB_USER:-postgres} > /dev/null 2>&1; do \
		printf '.'; sleep 1; done
	@echo "\n✓ PostgreSQL hazır → localhost:$${DB_PORT:-5432}"

db-stop: ## PostgreSQL'i durdur
	$(DC) stop db

db-logs: ## DB loglarını takip et
	$(DC) logs -f db

db-shell: ## psql shell aç
	$(DC) exec db psql -U $${DB_USER:-postgres} -d $${DB_NAME:-yks_tracker}

# ════════════════════════════════════════════════════════
# DOCKER COMPOSE — Geliştirme
# ════════════════════════════════════════════════════════
docker-build: ## Image'ları derle (cache kullan)
	$(DC) build --parallel

docker-up: ## DB + API başlat (arka planda)
	$(DC) up -d --build
	@echo "✓ API → http://localhost:$${SERVER_PORT:-8080}"
	@echo "✓ Swagger → http://localhost:$${SERVER_PORT:-8080}/swagger/index.html"

docker-down: ## Tüm container'ları durdur ve kaldır
	$(DC) down

docker-logs: ## Tüm logları canlı takip et
	$(DC) logs -f

docker-shell: ## API container içine shell aç
	$(DC) exec api sh

docker-dev: ## Hot-reload ile geliştirme modu başlat
	$(DC) --profile dev up --build api-dev db
	@echo "✓ Hot-reload aktif → http://localhost:$${SERVER_PORT:-8080}"

# ════════════════════════════════════════════════════════
# DOCKER COMPOSE — Production
# ════════════════════════════════════════════════════════
docker-prod: ## Production konfigürasyonuyla başlat
	$(DC_PROD) up -d --build
	@echo "✓ Production modu aktif"

docker-prod-down: ## Production'ı durdur
	$(DC_PROD) down

# ════════════════════════════════════════════════════════
# MİGRATION — Local (go run)
# ════════════════════════════════════════════════════════
migrate-up: ## Tüm migration'ları uygula (local DB)
	$(MIGRATE) up

migrate-down: ## TÜM migration'ları geri al (local DB)
	$(MIGRATE) down

migrate-version: ## Mevcut versiyon (local DB)
	$(MIGRATE) version

migrate-steps: ## n adım — örn: make migrate-steps N=-1
	$(MIGRATE) steps $(N)

migrate-force: ## Dirty fix — örn: make migrate-force V=3
	$(MIGRATE) force $(V)

migrate-new: ## Yeni migration — örn: make migrate-new NAME=add_goals
	$(MIGRATE) new $(NAME)

# ════════════════════════════════════════════════════════
# MİGRATION — Docker container içinde
# ════════════════════════════════════════════════════════
dc-migrate-up: ## Container'da migration uygula
	$(DC) run --rm --profile tools migrate up

dc-migrate-down: ## Container'da migration geri al
	$(DC) run --rm --profile tools migrate down

dc-migrate-version: ## Container'da versiyon göster
	$(DC) run --rm --profile tools migrate version

# ════════════════════════════════════════════════════════
# ARAÇLAR
# ════════════════════════════════════════════════════════
pgadmin: ## pgAdmin arayüzünü başlat → http://localhost:5050
	$(DC) --profile tools up -d pgadmin db
	@echo "✓ pgAdmin → http://localhost:5050"
	@echo "  Email   : $${PGADMIN_EMAIL:-admin@yks.local}"
	@echo "  Password: $${PGADMIN_PASSWORD:-admin}"

clean: ## Tüm container, volume ve image'ları temizle
	$(DC) down -v --rmi local
	docker volume prune -f
	rm -rf bin/ tmp/
	@echo "✓ Temizlik tamamlandı"
