# Sharing Vision Backend v2

Golang microservice untuk use case **Post Article** — sistem manajemen artikel dengan fitur enterprise.

## Fitur Lengkap

### Backend (Go + Gin)
| Fitur | Status |
|-------|--------|
| ✅ Autentikasi JWT (Register / Login) | Done |
| ✅ CRUD Artikel (Admin) dengan validasi | Done |
| ✅ CRUD Kategori (Admin) — Master Data | Done |
| ✅ Public Article Listing + Pagination | Done |
| ✅ Public Article Detail by ID | Done |
| ✅ Path-based Pagination (`/page/:limit/:offset`) | Done |
| ✅ Public Category List | Done |
| ✅ Search Artikel (Full-Text Search PostgreSQL) | Done |
| ✅ Filter by Kategori | Done |
| ✅ Dashboard Stats (Published / Draft / Thrash) | Done |
| ✅ Audit Log | Done |
| ✅ Rate Limiting (Public & Admin terpisah) | Done |
| ✅ Redis Caching | Done |
| ✅ PostgreSQL (Neon) | Done |
| ✅ CORS (Development & Production) | Done |
| ✅ Unit Tests | Done |
| ✅ K6 Performance Tests | Done |
| ✅ Swagger / OpenAPI Documentation | Done |
| ✅ Postman Collection | Done |

### Frontend (Next.js + Tailwind CSS)
| Fitur | Status |
|-------|--------|
| ✅ Halaman Publik dengan Hero Section | Done |
| ✅ Grid Artikel dengan Loading Skeleton | Done |
| ✅ Pagination dengan Page Numbers | Done |
| ✅ Search Artikel (Debounce 400ms) | Done |
| ✅ Filter Kategori (Dropdown) | Done |
| ✅ Halaman Detail Artikel Publik | Done |
| ✅ Login Page | Done |
| ✅ Dashboard Admin dengan Stat Cards | Done |
| ✅ All Posts dengan Tabs (Published / Drafts / Trashed) | Done |
| ✅ Add New Article Form | Done |
| ✅ Edit Article Form | Done |
| ✅ Category Dropdown dari Master Data | Done |
| ✅ Categories Management (CRUD) | Done |
| ✅ Audit Logs Page | Done |
| ✅ Sidebar Navigation | Done |
| ✅ Toast Notifications | Done |
| ✅ Loading Skeletons & Empty States | Done |
| ✅ Responsive Design (Mobile First) | Done |
| ✅ Security Headers | Done |
| ✅ Build Optimization | Done |

##  Arsitektur

### Backend Folder Structure
```
sharing-vision-backend-v2/
├── cmd/server/main.go          # Entry point & routing
├── internal/
│   ├── auth/                    # JWT, password hashing, migrations
│   ├── cache/                   # Redis client & caching logic
│   ├── config/                  # Environment config loader
│   ├── db/                      # PostgreSQL connection
│   ├── handler/                 # HTTP handlers (article, auth, category)
│   ├── middleware/              # Auth middleware, rate limiter
│   ├── model/                   # Data structures (article, user, category)
│   ├── repository/              # Database queries (article, user, category)
│   └── service/                 # Business logic layer
├── migrations/                  # SQL migration files
├── docs/                        # Swagger JSON, Postman collection
├── tests/
│   ├── unit/                    # Go unit tests
│   └── performance/k6_scripts/  # K6 load test scripts
└── scripts/                     # Utility scripts
```

### Frontend Folder Structure
```
sharing-vision-frontend/
├── app/
│   ├── page.tsx                 # Halaman publik (landing + artikel)
│   ├── layout.tsx               # Root layout + ToastContainer
│   ├── globals.css              # Custom Tailwind classes
│   ├── login/page.tsx           # Halaman login
│   ├── articles/[id]/page.tsx   # Detail artikel publik
│   └── dashboard/
│       ├── layout.tsx           # Dashboard layout + sidebar
│       ├── page.tsx             # Dashboard utama
│       ├── articles/
│       │   ├── page.tsx         # All Posts (tabs + table)
│       │   ├── new/page.tsx     # Add new article
│       │   └── [id]/edit/page.tsx  # Edit article
│       ├── categories/page.tsx  # Manage categories
│       └── audit-logs/page.tsx  # View audit logs
├── components/
│   ├── Navbar.tsx               # Navigation bar
│   ├── Footer.tsx               # Footer
│   └── Toast.tsx                # Toast notification system
└── lib/api.ts                   # API client
```

##  Setup & Instalasi

### Prerequisites
- Go 1.22+
- Node.js 18+
- PostgreSQL (Neon recommended)
- Redis (opsional, untuk caching)

### 1. Backend Setup

```bash
# Clone & masuk ke folder backend
git clone <repo-url>
cd sharing-vision-backend-v2

# Copy & isi .env
cp .env.example .env
```

**Minimal `.env` file:**
```env
# App
APP_PORT=8081
APP_ENV=development

# Database (Neon PostgreSQL)
DB_HOST=<your-neon-host>
DB_PORT=5432
DB_USER=<your-neon-user>
DB_PASSWORD=<your-neon-password>
DB_NAME=sv_portal
DB_SSLMODE=require

# Redis (optional, skip if no Redis)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=<your-secret-key-change-this>
JWT_EXPIRY=24h

# Rate Limit
RATE_LIMIT_PUBLIC_RPS=100
RATE_LIMIT_PUBLIC_BURST=10
RATE_LIMIT_ADMIN_RPS=300
RATE_LIMIT_ADMIN_BURST=30
```

```bash
# Install dependensi & jalankan
go mod tidy
go run ./cmd/server
```

Backend akan berjalan di `http://localhost:8081`.

### 2. Frontend Setup

```bash
# Clone & masuk ke folder frontend
cd sharing-vision-frontend

# Install dependensi
npm install

# Copy env
cp .env.example .env.local
```

**`.env.local`:**
```env
NEXT_PUBLIC_API_URL=http://localhost:8081/api/v1
```

```bash
# Development
npm run dev

# Production build
npm run build
npm start
```

Frontend akan berjalan di `http://localhost:3000`.

## Admin Default

Setelah pertama kali backend dijalankan, sistem akan otomatis membuat admin:

| Email | Password |
|-------|----------|
| `admin@sharingvision.id` | `admin123` |

## API Endpoints

### Public (No Auth Required)
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/auth/login` | Login pengguna |
| POST | `/api/v1/auth/register` | Registrasi pengguna |
| GET | `/api/v1/articles` | List artikel published (query: `limit`, `offset`, `q`, `category`) |
| GET | `/api/v1/articles/page/:limit/:offset` | List artikel dengan path pagination |
| GET | `/api/v1/articles/:id` | Detail artikel |
| GET | `/api/v1/categories` | List kategori |

### Admin (Requires JWT + Role: admin)
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/admin/articles` | Buat artikel baru |
| GET | `/api/v1/admin/articles` | List semua artikel |
| GET | `/api/v1/admin/articles/:id` | Detail artikel |
| PUT | `/api/v1/admin/articles/:id` | Update artikel |
| DELETE | `/api/v1/admin/articles/:id` | Soft-delete (→ thrash) |
| POST | `/api/v1/admin/categories` | Tambah kategori |
| GET | `/api/v1/admin/categories` | List kategori |
| GET | `/api/v1/admin/categories/:id` | Detail kategori |
| PUT | `/api/v1/admin/categories/:id` | Update kategori |
| DELETE | `/api/v1/admin/categories/:id` | Hapus kategori |
| GET | `/api/v1/admin/dashboard` | Statistik dashboard |
| GET | `/api/v1/admin/audit-logs` | Audit log aktivitas |

## Dokumentasi API

### Swagger
Buka file `docs/swagger.json` di [Swagger Editor](https://editor.swagger.io/) atau import ke Postman.

### Postman Collection
Import `docs/postman-collection.json` ke Postman.
- Set variable `base_url` ke URL backend
- Jalankan **Login** untuk auto-set token admin

## Testing

### Unit Tests
```bash
cd sharing-vision-backend-v2
go test ./tests/unit/... -v
```

Test coverage:
- ✅ Auth service (register, login, bootstrap admin, JWT)
- ✅ Article service (create, read, update, delete, count, categories)
- ✅ Category service (create, duplicate, update, delete, list)

### Performance Tests (k6)
```bash
# Public articles
k6 run tests/performance/k6_scripts/public_articles.js

# Admin articles
k6 run tests/performance/k6_scripts/admin_articles.js

# Login
k6 run tests/performance/k6_scripts/login.js

# Categories
k6 run tests/performance/k6_scripts/categories.js
```

## Deployment

### Backend
```bash
# Build binary
go build -o bin/sv-backend ./cmd/server

# Run
./bin/sv-backend
```

### Frontend
```bash
# Build
npm run build

# Start
npm start
```

** Update `.env.local` untuk production:**
```env
NEXT_PUBLIC_API_URL=https://be.antasource.xyz/api/v1
```

### CORS
Untuk production, tambahkan domain di CORS config di `cmd/server/main.go`:
```go
AllowOrigins: []string{
    "http://localhost:3000",
    "https://fdz.antasource.xyz",
    "https://be.antasource.xyz",
},
```

## Teknologi

### Backend
- **Go 1.22+** — Bahasa pemrograman
- **Gin** — HTTP framework
- **PostgreSQL (Neon)** — Database (serverless Postgres)
- **Redis** — Caching
- **JWT (golang-jwt)** — Autentikasi
- **bcrypt** — Password hashing
- **Zap** — Structured logging
- **Testify** — Unit testing
- **k6** — Performance testing

### Frontend
- **Next.js 14** — React framework
- **Tailwind CSS** — Utility-first CSS
- **TypeScript** — Type safety
- **Heroicons** — SVG icons

