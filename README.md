# E-Voting Backend 2025

Backend service untuk sistem E-Voting (2025). Ditulis dengan Go + Gin + GORM + PostgreSQL, memakai JWT (RSA) untuk autentikasi dan custom migration sederhana.

## 🔧 Tech Stack

- Go 1.23
- Gin Web Framework
- GORM ORM (PostgreSQL)
- JWT (RS256) dengan kunci RSA
- Custom migration ringan (tanpa external tool)
- Dotenv untuk konfigurasi lokal

## 🗂 Struktur Direktori

```
├── main.go
├── config/
│   ├── config.go              # Load & simpan konfigurasi aplikasi
│   ├── database/database.go   # Inisialisasi koneksi GORM + sql.DB
│   ├── middleware/            # Middleware (CORS, auth)
│   └── pkg/                   # Utilitas (token, error wrapper, utils)
├── controller/                # Registrasi controller + error helper
├── contract/                  # Kontrak/abstraksi layer (Service, Repository)
├── repository/                # Implementasi repository (stub)
├── service/                   # Implementasi service (stub)
├── database/                  # Sistem migration, model, and seed
├── private_key.pem / public_key.pem
├── .env.example
└── README.md
```

Alur eksekusi utama:

```
main.go → config.Load() → token.Load() → server.Run()
 server.Run():
   1. Koneksi database
   2. Jalankan migration (Up / optional Down / DownAll via flag)
   3. Inisialisasi repository → service → controller
   4. Start HTTP server (Gin)
```

## ⚙️ Environment Variables

File `.env.example` saat ini memiliki beberapa nama variabel yang **tidak konsisten** dengan kode di `config/config.go` dan paket token.

| Dibaca Kode | Di .env.example | Perlu Disesuaikan | Keterangan |
|-------------|-----------------|------------------|------------|
| `ACCESS_TOKEN_LIFE_TIME`  | `LIFE_TIME_TOKEN`          | Ubah ke `ACCESS_TOKEN_LIFE_TIME` | Detik masa hidup access token |
| `REFRESH_TOKEN_LIFE_TIME` | `REFRESH_LIFE_TIME_TOKEN`  | Ubah ke `REFRESH_TOKEN_LIFE_TIME`| Detik masa hidup refresh token |
| `PRIVATE_KEY`             | (benar)                    | OK | Path file private key |
| `PUBLIC_KEY`              | (benar)                    | OK | Path file public key |

Contoh `.env` yang benar:

```
PORT=8080
IS_PRODUCTION=false

DB_USER=postgres
DB_PASS=your-password
DB_NAME=hme
DB_HOST=127.0.0.1
DB_PORT=5432
DB_TIME_ZONE=Asia/Jakarta

ALLOW_ORIGIN=*
BASE_URL=http://localhost:8080

ACCESS_TOKEN_LIFE_TIME=3600
REFRESH_TOKEN_LIFE_TIME=86400

PRIVATE_KEY=./private_key.pem
PUBLIC_KEY=./public_key.pem
```

## 🔑 Generate Key

1. Generate private key (2048 bit)
```bash
openssl genrsa -out private_key.pem 2048
```

2. Generate public key dari private key

```bash
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

## 📃 Migration
Migrasi model dan seed ke database
```bash
go run main.go migrate
```

Reset semua migrasi
```bash
go run main.go reset
```

## 🏍️ Flow Process

1. Create Repository Functions
Purpose: Handle all database operations (queries, inserts, updates, deletes).
Example: GetUserByUsername, FindOrCreateUser.
Key point: The repository layer only deals with data access, not business rules.

2. Register Repository in Contract & Main File
Purpose: Define the interface in contract so repository implementations are swappable.
Main File: Inject repository into service.
Key Point: Ensures loose coupling and testability.

3. Service Layer (Business Logic)
Purpose: Call repository functions and apply business rules.
Example:
Login() → validates payload, checks password, generates JWT.
Key Point: Service = brain of the app, applies logic beyond just fetching data.

4. Register Service in Contract & Main File
Purpose: Add the service to the contract and inject it into the controller.
Key Point: Maintains clean modular design.

5. Controller Layer (API Layer)
Purpose: Handle HTTP requests & responses.
Example:
POST /auth/login → calls service.Login().
Key Point: Controller = entry point for clients → communicates with service only.

6. Register Controller in the Main File
Purpose: Add routing to related controller
![flow](image.png)