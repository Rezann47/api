# ════════════════════════════════════════════════════════
#  YKS Tracker — Multi-stage Dockerfile
#  Hedefler:
#    docker build --target dev   -t yks-tracker:dev .   (hot-reload)
#    docker build --target prod  -t yks-tracker:prod .  (minimal image)
# ════════════════════════════════════════════════════════

# ─── 1. Bağımlılık indirme (cache katmanı) ───────────────
FROM golang:1.23-alpine AS deps

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ─── 2. Builder ──────────────────────────────────────────
FROM deps AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY . .

# -ldflags: binary boyutunu küçültür (debug sembollerini kaldırır)
# -trimpath: binary içindeki mutlak yol bilgisini kaldırır (güvenlik)
RUN go build -ldflags="-w -s" -trimpath -o /out/api     ./cmd/api  && \
    go build -ldflags="-w -s" -trimpath -o /out/migrate ./cmd/migrate

# ─── 3. Development hedefi (hot-reload ile) ──────────────
FROM golang:1.23-alpine AS dev

RUN apk add --no-cache tzdata ca-certificates curl && \
    go install github.com/air-verse/air@latest

WORKDIR /app

# Bağımlılıkları önceden indir (volume mount değiştiğinde tekrar indirilmesin)
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kod volume olarak mount edilecek, burada COPY yok
COPY .air.toml ./

EXPOSE 8080

# air: dosya değişikliklerini izleyip otomatik rebuild yapar
CMD ["air", "-c", ".air.toml"]

# ─── 4. Production hedefi (minimal distroless/alpine) ────
FROM alpine:3.20 AS prod

# Timezone ve TLS için gerekli; başka hiçbir şey ekleme
RUN apk add --no-cache tzdata ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Sadece binary ve migrations kopyala
COPY --from=builder /out/api       ./api
COPY --from=builder /out/migrate   ./migrate
COPY migrations                    ./migrations

# root yerine kısıtlı kullanıcı
USER appuser

EXPOSE 8080

# Liveness probe için health check
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./api"]
