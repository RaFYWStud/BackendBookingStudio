# 🎵 Studio Booking Backend API

Backend service untuk sistem booking studio musik. Dibangun dengan **Go + Gin + GORM + PostgreSQL**, menggunakan **JWT (RS256)** untuk autentikasi dan **SMTP** untuk notifikasi email.

---

## 🔧 Tech Stack

-   **Go 1.23** - Programming Language
-   **Gin Web Framework** - HTTP Router
-   **GORM ORM** - Database ORM (PostgreSQL)
-   **JWT (RS256)** - Authentication dengan RSA Keys
-   **SMTP** - Email Notification (Gmail)
-   **Bcrypt** - Password Hashing
-   **PostgreSQL** - Database
-   **Dotenv** - Environment Configuration

---

## 📋 Table of Contents

1. [Prerequisites](#-prerequisites)
2. [Installation & Setup](#-installation--setup)
3. [Environment Variables](#️-environment-variables)
4. [Database Migration](#-database-migration)
5. [API Documentation](#-api-documentation)
    - [Authentication Endpoints](#1-authentication-endpoints)
    - [Studios Endpoints](#2-studios-endpoints)
    - [Bookings Endpoints](#3-bookings-endpoints-customer)
    - [Admin Bookings Endpoints](#4-bookings-admin-endpoints)
6. [Error Handling](#-error-handling)
7. [Email Notifications](#-email-notifications)
8. [Testing Guide](#-testing-guide)

---

## ✅ Prerequisites

-   **Go 1.23+** installed
-   **PostgreSQL 14+** installed and running
-   **Gmail account** with App Password (for SMTP)
-   **Postman** or any API testing tool

---

## 🚀 Installation & Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd BackendBookingStudio
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Generate RSA Keys

```bash
# Generate private key (2048 bit)
openssl genrsa -out private_key.pem 2048

# Generate public key from private key
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

### 4. Setup Environment Variables

Copy `.env.example` to `.env` and fill in your configuration:

```bash
cp .env.example .env
```

### 5. Create PostgreSQL Database

```sql
CREATE DATABASE booking_studio_backend;
```

### 6. Run Database Migration

```bash
# Auto-migrate on first run
go run main.go

# Or manually migrate
go run main.go migrate

# Reset database (drop all tables and re-migrate)
go run main.go reset

# Seed data only
go run main.go seed
```

### 7. Start Server

```bash
# Development (with auto-reload using Air)
air

# Production
go run main.go
```

Server will run on: `http://localhost:8080`

---

## ⚙️ Environment Variables

Create a `.env` file with the following configuration:

```env
# Server Configuration
PORT=8080
IS_PRODUCTION=false
BASE_URL=http://localhost:8080

# Database Configuration
DB_USER=postgres
DB_PASS=your_postgres_password
DB_NAME=booking_studio_backend
DB_HOST=127.0.0.1
DB_PORT=5432
DB_TIME_ZONE=Asia/Jakarta

# CORS Configuration
ALLOW_ORIGIN=*

# JWT Configuration
ACCESS_TOKEN_LIFETIME=3600      # 1 hour in seconds
REFRESH_TOKEN_LIFETIME=86400    # 24 hours in seconds
PRIVATE_KEY=private_key.pem
PUBLIC_KEY=public_key.pem

# SMTP Email Configuration (Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_FROM=your-email@gmail.com
SMTP_PASSWORD=your-app-password  # Gmail App Password
APP_NAME=Studio Booking System
APP_URL=http://localhost:8080

# Admin WhatsApp (for payment confirmation)
ADMIN_WHATSAPP=+6289570608111
ADMIN_WHATSAPP_DISPLAY=0895-7060-8111
ADMIN_NAME=Admin Studio Booking

# Rate Limiting (optional)
RATE_LIMIT_RPS=10    # Requests per second
RATE_LIMIT_BURST=20  # Burst size
```

---

## 📃 Database Migration

### Auto-Migration (Recommended)

Server will automatically run migration on first start if tables don't exist.

### Manual Migration Commands

```bash
# Run migration
go run main.go migrate

# Reset database (drop all tables and re-migrate)
go run main.go reset

# Seed data only
go run main.go seed
```

### Default Seed Data

After migration, you'll get:

**Admin Account:**

-   Email: `admin@studiobooking.com`
-   Password: `admin123`
-   Role: `admin`

**Sample Studios:**

-   Studio Premium A (Jakarta Barat) - Rp 250.000/hour
-   Studio Budget D (Jakarta Utara) - Rp 100.000/hour

---

## 📚 API Documentation

### Base URL

```
http://localhost:8080
```

**⚠️ IMPORTANT: Endpoints do NOT use `/api/` prefix.**

---

## 1. Authentication Endpoints

### 1.1 Register Customer

**Endpoint:** `POST /auth/register`

**Access:** Public

**Request Body:**

```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "password_confirmation": "password123"
}
```

**Validation Rules:**

-   `name`: required, min 2 characters
-   `email`: required, valid email format, unique
-   `password`: required, min 6 characters
-   `password_confirmation`: required, must match password

**Success Response (201 Created):**

```json
{
    "success": true,
    "message": "Registration successful",
    "data": {
        "id": 2,
        "name": "John Doe",
        "email": "john@example.com",
        "role": "customer"
    }
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "password_confirmation": "password123"
  }'
```

**Error Responses:**

```json
// 400 Bad Request - Email already exists
{
    "status": 400,
    "error": "Bad Request",
    "message": "email already registered"
}
```

---

### 1.2 Login

**Endpoint:** `POST /auth/login`

**Access:** Public

**Request Body:**

```json
{
    "email": "admin@studiobooking.com",
    "password": "admin123"
}
```

**Success Response (200 OK):**

```json
{
    "success": true,
    "message": "Login successful",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "name": "Admin",
            "email": "admin@studiobooking.com",
            "role": "admin"
        }
    }
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@studiobooking.com",
    "password": "admin123"
  }'
```

**Frontend Implementation:**

```javascript
// Save token to localStorage
const response = await fetch("http://localhost:8080/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
});

const data = await response.json();
if (data.success) {
    localStorage.setItem("token", data.data.token);
    localStorage.setItem("user", JSON.stringify(data.data.user));
}
```

---

### 1.3 Get Profile

**Endpoint:** `GET /auth/profile`

**Access:** Protected (requires authentication)

**Headers:**

```
Authorization: Bearer <jwt_token>
```

**Success Response (200 OK):**

```json
{
    "success": true,
    "data": {
        "id": 2,
        "name": "John Doe",
        "email": "john@example.com",
        "role": "customer"
    }
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/auth/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 2. Studios Endpoints

### 2.1 Get All Studios (Public)

**Endpoint:** `GET /studios`

**Access:** Public

**Query Parameters:**

| Parameter   | Type    | Required | Description                            | Example                                            |
| ----------- | ------- | -------- | -------------------------------------- | -------------------------------------------------- |
| `location`  | string  | No       | Filter by location                     | `Jakarta`                                          |
| `min_price` | integer | No       | Minimum price per hour                 | `100000`                                           |
| `max_price` | integer | No       | Maximum price per hour                 | `300000`                                           |
| `is_active` | boolean | No       | Filter active studios                  | `true`                                             |
| `search`    | string  | No       | Search by studio name                  | `Premium`                                          |
| `page`      | integer | No       | Page number (default: 1)               | `1`                                                |
| `limit`     | integer | No       | Items per page (default: 10, max: 100) | `10`                                               |
| `sort_by`   | string  | No       | Sort order                             | `price_asc`, `price_desc`, `name_asc`, `name_desc` |

**Example Request:**

```
GET /studios?location=Jakarta&min_price=100000&max_price=300000&page=1&limit=10&sort_by=price_asc
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/studios?location=Jakarta&min_price=100000&max_price=300000&page=1&limit=10"
```

**Success Response (200 OK):**

```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "Studio Premium A",
            "description": "Studio premium dengan peralatan kelas dunia...",
            "location": "Jakarta Barat",
            "price_per_hour": 250000,
            "image_url": "https://images.unsplash.com/photo-1598653222000-6b7b7a552625",
            "facilities": [
                "AC",
                "Professional Drum Set",
                "Multiple Amplifiers",
                "Grand Piano"
            ],
            "operating_hours": "08:00-23:00",
            "is_active": true,
            "created_at": "2025-11-21 10:00:00",
            "updated_at": "2025-11-21 10:00:00"
        }
    ],
    "pagination": {
        "current_page": 1,
        "page_size": 10,
        "total_pages": 1,
        "total_records": 2
    }
}
```

---

### 2.2 Get Studio by ID (Public)

**Endpoint:** `GET /studios/:id`

**Access:** Public

**Example Request:**

```
GET /studios/1
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/studios/1
```

---

### 2.3 Check Availability (Public)

**Endpoint:** `POST /studios/:id/availability`

**Access:** Public

**Request Body:**

```json
{
    "date": "2025-11-25",
    "start_time": "14:00",
    "end_time": "17:00"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/studios/1/availability \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2025-11-25",
    "start_time": "14:00",
    "end_time": "17:00"
  }'
```

---

### 2.4 Create Studio (Admin Only)

**Endpoint:** `POST /studios`

**Access:** Admin Only

**Headers:**

```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "name": "Studio Acoustic B",
    "description": "Studio dengan akustik sempurna untuk rekaman vokal",
    "location": "Jakarta Selatan",
    "price_per_hour": 200000,
    "image_url": "https://images.unsplash.com/photo-1598653222000-6b7b7a552625",
    "facilities": [
        "AC",
        "Soundproof Booth",
        "Professional Microphones",
        "Audio Interface"
    ],
    "operating_hours": "09:00-22:00"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/studios \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Studio Acoustic B",
    "description": "Studio dengan akustik sempurna",
    "location": "Jakarta Selatan",
    "price_per_hour": 200000,
    "image_url": "https://example.com/image.jpg",
    "facilities": ["AC", "Soundproof Booth"],
    "operating_hours": "09:00-22:00"
  }'
```

---

### 2.5 Update Studio - Partial (Admin Only) - **PATCH**

**Endpoint:** `PATCH /studios/:id`

**Access:** Admin Only

**Headers:**

```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body (Send only fields you want to update):**

```json
{
    "price_per_hour": 350000,
    "is_active": false
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/studios/2 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "price_per_hour": 350000,
    "is_active": false
  }'
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Studio updated successfully (partial update)",
  "data": {
    "id": 2,
    "name": "Studio Budget D",
    "price_per_hour": 350000,
    "is_active": false,
    ...
  }
}
```

---

### 2.6 Delete Studio (Admin Only)

**Endpoint:** `DELETE /studios/:id`

**Access:** Admin Only

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/studios/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## 3. Bookings Endpoints (Customer)

### 3.1 Create Booking

**Endpoint:** `POST /bookings`

**Access:** Customer (authenticated)

**Headers:**

```
Authorization: Bearer <customer_token>
Content-Type: application/json
```

**Request Body:**

```json
{
    "studio_id": 1,
    "booking_date": "2025-11-25",
    "start_time": "14:00",
    "end_time": "17:00"
}
```

**⚠️ Note:** `duration_hours` is **auto-calculated** from time difference.

**cURL Example:**

```bash
curl -X POST http://localhost:8080/bookings \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "studio_id": 1,
    "booking_date": "2025-11-25",
    "start_time": "14:00",
    "end_time": "17:00"
  }'
```

**Success Response (201 Created):**

```json
{
    "success": true,
    "message": "Booking berhasil dibuat. Total pembayaran: Rp 750.000.\n\n📱 Silakan hubungi admin untuk pembayaran:\nWhatsApp: 0895-7060-8111",
    "data": {
        "id": 1,
        "user_id": 2,
        "studio_id": 1,
        "booking_date": "2025-11-25",
        "start_time": "14:00",
        "end_time": "17:00",
        "duration_hours": 3,
        "total_price": 750000,
        "status": "pending",
        "created_at": "2025-11-21 15:30:00",
        "updated_at": "2025-11-21 15:30:00",
        "studio": {
            "id": 1,
            "name": "Studio Premium A",
            "location": "Jakarta Barat",
            "price_per_hour": 250000
        }
    }
}
```

---

### 3.2 Get My Bookings

**Endpoint:** `GET /bookings`

**Access:** Customer (authenticated)

**Query Parameters:**

| Parameter   | Type    | Description                                      |
| ----------- | ------- | ------------------------------------------------ |
| `status`    | string  | `pending`, `confirmed`, `completed`, `cancelled` |
| `studio_id` | integer | Filter by studio                                 |
| `page`      | integer | Page number (default: 1)                         |
| `limit`     | integer | Items per page (default: 10)                     |

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/bookings?status=pending&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN"
```

---

### 3.3 Get Booking Detail

**Endpoint:** `GET /bookings/:id`

**Access:** Customer/Admin (authenticated)

**cURL Example:**

```bash
curl -X GET http://localhost:8080/bookings/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 3.4 Cancel Booking

**Endpoint:** `POST /bookings/:id/cancel`

**Access:** Customer (authenticated)

**Request Body:**

```json
{
    "reason": "Ada keperluan mendadak, mohon maaf tidak bisa datang"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/bookings/1/cancel \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Ada keperluan mendadak, mohon maaf tidak bisa datang"
  }'
```

---

## 4. Bookings Admin Endpoints

### 4.1 Get All Bookings (Admin)

**Endpoint:** `GET /bookings/admin`

**Access:** Admin Only

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/bookings/admin?status=pending&studio_id=1" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

### 4.2 Update Booking Status (Admin)

**Endpoint:** `PUT /bookings/admin/:id/status`

**Access:** Admin Only

**Request Body:**

```json
{
    "status": "confirmed",
    "admin_notes": "Pembayaran diterima via BCA tanggal 21 Nov 2024, pukul 15:30. Total Rp 750.000"
}
```

**Available Status Values:**

-   `pending` - Menunggu pembayaran
-   `confirmed` - Sudah dibayar (dikonfirmasi admin)
-   `completed` - Selesai digunakan
-   `cancelled` - Dibatalkan

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/bookings/admin/1/status \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "confirmed",
    "admin_notes": "Pembayaran diterima via BCA tanggal 21 Nov 2024"
  }'
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "message": "Booking status updated to confirmed. Customer has been notified via email.",
  "data": {
    "id": 1,
    "status": "confirmed",
    "admin_notes": "Pembayaran diterima via BCA...",
    ...
  }
}
```

---

## 🚨 Error Handling

All errors follow a consistent format:

### Error Response Structure

```json
{
    "status": 400,
    "error": "Error Type",
    "message": "Detailed error message"
}
```

### Common HTTP Status Codes

| Status Code | Error Type            | Description                             |
| ----------- | --------------------- | --------------------------------------- |
| `400`       | Bad Request           | Invalid input, validation error         |
| `401`       | Unauthorized          | Missing or invalid authentication token |
| `403`       | Forbidden             | Authenticated but no permission         |
| `404`       | Not Found             | Resource not found                      |
| `500`       | Internal Server Error | Server error                            |

---

## 📧 Email Notifications

System automatically sends emails for:

1. **Booking Created** - When customer creates new booking (status: PENDING)
2. **Booking Confirmed** - When admin confirms payment (status: CONFIRMED)
3. **Booking Cancelled** - When booking is cancelled by customer/admin

---

## 🧪 Testing Guide

### Quick Test Flow

```bash
# 1. Register customer
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123","password_confirmation":"password123"}'

# 2. Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 3. Get studios
curl -X GET http://localhost:8080/studios

# 4. Create booking
curl -X POST http://localhost:8080/bookings \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"studio_id":1,"booking_date":"2025-11-25","start_time":"14:00","end_time":"17:00"}'
```

---

## 📊 API Endpoints Summary

| Method       | Endpoint                     | Auth           | Description             |
| ------------ | ---------------------------- | -------------- | ----------------------- |
| **Auth**     |
| POST         | `/auth/register`             | Public         | Register customer       |
| POST         | `/auth/login`                | Public         | Login                   |
| GET          | `/auth/profile`              | Customer/Admin | Get profile             |
| **Studios**  |
| GET          | `/studios`                   | Public         | Get all studios         |
| GET          | `/studios/:id`               | Public         | Get studio by ID        |
| POST         | `/studios/:id/availability`  | Public         | Check availability      |
| POST         | `/studios`                   | Admin          | Create studio           |
| PUT          | `/studios/:id`               | Admin          | Update studio (full)    |
| PATCH        | `/studios/:id`               | Admin          | Update studio (partial) |
| DELETE       | `/studios/:id`               | Admin          | Delete studio           |
| **Bookings** |
| POST         | `/bookings`                  | Customer       | Create booking          |
| GET          | `/bookings`                  | Customer       | Get my bookings         |
| GET          | `/bookings/:id`              | Customer/Admin | Get booking detail      |
| POST         | `/bookings/:id/cancel`       | Customer       | Cancel booking          |
| GET          | `/bookings/admin`            | Admin          | Get all bookings        |
| PUT          | `/bookings/admin/:id/status` | Admin          | Update booking status   |

---

## 🔒 Security Best Practices

1. **Always use HTTPS in production**
2. **Never expose `.env` file**
3. **Rotate JWT keys regularly**
4. **Use strong passwords** (min 8 characters)
5. **Enable rate limiting** in production
6. **Validate all user inputs**

---

## 📄 License

This project is licensed under the MIT License.

---

**Happy Coding! 🎸🎵**
