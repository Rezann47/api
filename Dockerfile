# ── Build Stage ──────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Bağımlılıkları önce kopyala (Docker cache optimizasyonu)
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Binary derle (CGO kapalı = static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

# ── Runtime Stage (minimal image) ────────────────────────────
FROM alpine:3.19

# Güvenlik: root olmayan kullanıcı
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

# Sertifika (HTTPS istekler için)
RUN apk --no-cache add ca-certificates

USER appuser

EXPOSE 8080

CMD ["./app"]
