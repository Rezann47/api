# Go Gin CRUD — Production-Ready API

PostgreSQL + GORM + JWT + Katmanlı Mimari

## 📁 Proje Yapısı

```
go-gin-crud/
├── cmd/
│   └── main.go              # Giriş noktası, DI bağlantısı, server başlatma
├── config/
│   ├── config.go            # Env değişkenleri okuma
│   └── database.go          # GORM + PostgreSQL bağlantı pool
├── internal/
│   ├── model/
│   │   └── model.go         # GORM modelleri (User, Product)
│   ├── dto/
│   │   └── dto.go           # Request/Response struct'ları, validasyon
│   ├── repository/
│   │   └── repository.go    # Veritabanı işlemleri (CRUD + sayfalama)
│   ├── service/
│   │   └── service.go       # Business logic, bcrypt, JWT
│   ├── handler/
│   │   └── handler.go       # HTTP handler'lar, route mantığı
│   └── middleware/
│       └── middleware.go    # JWT auth, CORS, logger, rate limiter
├── migrations/
│   ├── 001_create_users.up.sql
│   └── 002_create_products.up.sql
├── docs/
│   └── api_test.http        # VS Code REST Client test dosyası
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── Makefile
```
handler -> service -> repository -> model

## 🚀 Hızlı Başlangıç

### 1. PostgreSQL başlat
```bash
make docker-up
# veya: docker-compose up db -d
```

### 2. .env dosyası oluştur
```bash
cp .env.example .env
# .env içini düzenle
```

### 3. Uygulamayı çalıştır
```bash
make tidy   # bağımlılıkları indir
make run    # uygulamayı başlat
```

### Docker Compose ile (tümü birden)
```bash
docker-compose up --build
```

## 🔗 API Endpoint'leri

| Method | Endpoint                  | Auth  | Açıklama               |
|--------|---------------------------|-------|------------------------|
| GET    | /health                   | -     | Sağlık kontrolü        |
| POST   | /api/v1/auth/register     | -     | Kayıt ol               |
| POST   | /api/v1/auth/login        | -     | Giriş yap, JWT al      |
| GET    | /api/v1/me                | JWT   | Kendi profili           |
| GET    | /api/v1/users             | Admin | Tüm kullanıcılar        |
| GET    | /api/v1/users/:id         | Admin | Tek kullanıcı           |
| PUT    | /api/v1/users/:id         | Admin | Güncelle               |
| DELETE | /api/v1/users/:id         | Admin | Sil (soft delete)       |
| POST   | /api/v1/products          | JWT   | Ürün oluştur           |
| GET    | /api/v1/products          | JWT   | Ürünleri listele        |
| GET    | /api/v1/products/:id      | JWT   | Tek ürün               |
| PUT    | /api/v1/products/:id      | JWT   | Ürün güncelle          |
| DELETE | /api/v1/products/:id      | JWT   | Ürün sil               |

## 🏗️ Mimari: Katmanlı Yapı

```
HTTP Request
     ↓
  Middleware (JWT, CORS, RateLimit, Logger)
     ↓
  Handler   — HTTP parse/validate/respond
     ↓
  Service   — Business logic, bcrypt, JWT
     ↓
  Repository — SQL/GORM operasyonları
     ↓
  PostgreSQL
```

### Her Katmanın Sorumluluğu

**Handler**: İstek alır, DTO'ya dönüştürür, service'i çağırır, response döner.  
**Service**: İş kuralları, şifre hashleme, token üretme, yetki kontrolleri.  
**Repository**: SADECE veritabanı işlemleri. Hiçbir iş mantığı içermez.  
**Model**: Tablo yapısı. GORM tag'leri burada tanımlanır.  
**DTO**: Dış dünyaya açılan veri kapısı. `binding:"required"` ile validasyon.

## 🔒 Güvenlik

- **Şifre**: bcrypt cost=12 ile hashlenir, asla plain text saklanmaz
- **JWT**: HS256, expire süreli, her istekte doğrulanır
- **Soft Delete**: Kayıtlar fiziksel olarak silinmez, `deleted_at` dolar
- **Role Based**: user / admin rolü, handler seviyesinde kontrol
- **Rate Limiting**: IP başına dakika başı istek sınırı
- **Input Validation**: Gin binding tag'leri ile otomatik

## 📊 Sayfalama Örneği

```bash
GET /api/v1/products?page=2&limit=10&search=laptop&sort=price&order=desc
```

Response:
```json
{
  "success": true,
  "data": {
    "data": [...],
    "total": 47,
    "page": 2,
    "limit": 10,
    "total_pages": 5
  }
}
```
