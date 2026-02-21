# YKS Tracker API

TYT & AYT öğrencileri ve eğitmenleri için REST API backend.

## Teknolojiler

- **Go 1.23** + **Gin**
- **PostgreSQL 16** + **GORM**
- **JWT** (access + refresh token)
- **golang-migrate** (migration yönetimi)
- **zap** (structured logging)
- **Docker** + **Docker Compose**

## Hızlı Başlangıç

```bash
# 1. Repo'yu klonla
git clone https://github.com/yourusername/yks-tracker
cd yks-tracker

# 2. .env oluştur
cp .env.example .env
# .env içinde DB_PASSWORD, JWT_ACCESS_SECRET, JWT_REFRESH_SECRET doldur
sh scripts/generate_secret.sh  # secret üretmek için

# 3. DB başlat
make db-start

# 4. Migration uygula
make migrate-up

# 5. Çalıştır
make dev
```

## Docker ile Çalıştırma

```bash
docker compose up --build
# Migration otomatik uygulanır
```

## Migration Komutları

```bash
make migrate-up              # Tüm migration'ları uygula
make migrate-down            # Geri al (tehlikeli!)
make migrate-version         # Mevcut versiyon
make migrate-steps N=-1      # 1 adım geri
make migrate-force V=3       # Dirty state düzelt
make migrate-new NAME=add_x  # Yeni migration dosyası
```

## API Dokümantasyonu

```bash
make swagger                 # Swagger üret
# Sunucu çalışırken: http://localhost:8080/swagger/index.html
```

## Proje Yapısı

```
cmd/           → Çalıştırılabilir binary'ler (api, migrate, seed, worker)
internal/
  config/      → Konfigürasyon + DB bağlantısı
  domain/      → Entity, DTO, hata tipleri (saf Go — framework bağımsız)
  repository/  → Veritabanı erişim katmanı
  service/     → İş mantığı
  handler/     → HTTP katmanı (Gin)
  middleware/  → JWT, CORS, logger, rate-limit
  server/      → Bağımlılık inject + HTTP server
pkg/           → Dışa açık yardımcı paketler
migrations/    → SQL migration dosyaları
test/          → Integration testleri + yardımcılar
```

## Endpoint Özeti

| Method | Path | Açıklama | Rol |
|--------|------|----------|-----|
| POST | /api/v1/auth/register | Kayıt | Public |
| POST | /api/v1/auth/login | Giriş | Public |
| POST | /api/v1/auth/refresh | Token yenile | Public |
| POST | /api/v1/auth/logout | Çıkış | Auth |
| GET | /api/v1/users/me | Profil | Auth |
| GET | /api/v1/subjects | Ders listesi | Auth |
| PATCH | /api/v1/topics/:id/mark | Konu işaretle | Student |
| POST | /api/v1/pomodoros | Pomodoro ekle | Student |
| GET | /api/v1/pomodoros/stats | İstatistik | Student |
| POST | /api/v1/exam-results | Deneme ekle | Student |
| GET | /api/v1/exam-results/stats | Gelişim | Student |
| POST | /api/v1/instructor/students | Öğrenci ekle | Instructor |
| GET | /api/v1/instructor/students | Öğrenci listesi | Instructor |
| GET | /api/v1/instructor/students/:id/progress | Öğrenci ilerlemesi | Instructor |
| GET | /health | Health check | Public |
