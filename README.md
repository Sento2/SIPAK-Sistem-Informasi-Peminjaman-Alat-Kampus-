# 🎓 SIPAK — Sistem Informasi Peminjaman Alat Kampus

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"/>
  <img src="https://img.shields.io/badge/MongoDB-Atlas-47A248?style=for-the-badge&logo=mongodb&logoColor=white" alt="MongoDB"/>
  <img src="https://img.shields.io/badge/Chi_Router-v5-blue?style=for-the-badge" alt="Chi Router"/>
  <img src="https://img.shields.io/badge/JWT-Auth-orange?style=for-the-badge&logo=jsonwebtokens&logoColor=white" alt="JWT"/>
</p>

**SIPAK** adalah REST API berbasis **Golang** untuk mengelola sistem peminjaman alat di kampus. Aplikasi ini dikembangkan sebagai **Tugas Besar Rekayasa Perangkat Lunak (RPL)**.

---

## ✨ Fitur Utama

| Fitur                 | Deskripsi                         |
| --------------------- | --------------------------------- |
| 🔐 **Autentikasi**    | Register & Login dengan JWT Token |
| 👥 **Multi-Role**     | Akun Mahasiswa & Admin            |
| 📦 **Manajemen Alat** | CRUD alat kampus (Admin only)     |
| 🔄 **Peminjaman**     | Pinjam & kembalikan alat          |
| 📊 **Riwayat**        | Lacak transaksi peminjaman        |
| 🔒 **Keamanan**       | API Key + JWT Authentication      |
| 🌐 **CORS**           | Support cross-origin requests     |

---

## 🧱 Tech Stack

| Kategori             | Teknologi                                             |
| -------------------- | ----------------------------------------------------- |
| **Backend**          | Go 1.25 + [Chi Router](https://github.com/go-chi/chi) |
| **Database**         | MongoDB Atlas                                         |
| **Authentication**   | JWT (JSON Web Token) + API Key                        |
| **Password Hashing** | bcrypt (`golang.org/x/crypto`)                        |

### 📚 Library yang Digunakan

```
github.com/go-chi/chi/v5       → HTTP Router
go.mongodb.org/mongo-driver    → MongoDB Driver
github.com/golang-jwt/jwt/v5   → JWT Authentication
github.com/joho/godotenv       → Environment Variables
github.com/rs/cors             → CORS Middleware
golang.org/x/crypto            → Password Hashing (bcrypt)
```

---

## 📁 Struktur Project

```
SIPAK/
├── 📄 go.mod                  # Go module definition
├── 📄 main.go                 # Entry point aplikasi
├── 📄 .env                    # Environment variables (jangan di-commit!)
├── 📁 config/
│   └── config.go              # Konfigurasi & koneksi MongoDB
├── 📁 models/
│   ├── user.go                # Model User (Mahasiswa/Admin)
│   ├── alat.go                # Model Alat Kampus
│   └── transaction.go         # Model Transaksi Peminjaman
├── 📁 handlers/
│   ├── auth_handler.go        # Handler Login & Register
│   ├── alat_handler.go        # Handler CRUD Alat
│   ├── peminjaman_handler.go  # Handler Peminjaman & Pengembalian
│   └── user_handler.go        # Handler Manajemen User (Admin)
├── 📁 middleware/
│   └── auth.go                # Middleware API Key, JWT, AdminOnly
└── 📁 utils/
    ├── jwt.go                 # Helper generate & validate JWT
    └── response.go            # Helper JSON response
```

---

## ⚙️ Konfigurasi Environment

Buat file `.env` di root project:

```env
# MongoDB Configuration
MONGO_URI=mongodb+srv://user:password@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
DB_NAME=sipak_db

# Security
JWT_SECRET=your_super_secret_jwt_key
API_KEY=your_super_secret_api_key

# Server
PORT=8080
```

| Variable     | Deskripsi                        |
| ------------ | -------------------------------- |
| `MONGO_URI`  | Connection string MongoDB Atlas  |
| `DB_NAME`    | Nama database yang digunakan     |
| `JWT_SECRET` | Secret key untuk signing JWT     |
| `API_KEY`    | API Key untuk header `X-API-Key` |
| `PORT`       | Port server (default: 8080)      |

---

## 🚀 Cara Menjalankan

### Prerequisites

- Go 1.21+ terinstall
- Akun MongoDB Atlas (atau MongoDB lokal)

### Langkah-langkah

```bash
# 1. Clone repository
git clone https://github.com/Sento2/SIPAK-Sistem-Informasi-Peminjaman-Alat-Kampus-.git
cd SIPAK-Sistem-Informasi-Peminjaman-Alat-Kampus-

# 2. Install dependencies
go mod tidy

# 3. Buat file .env (sesuaikan dengan konfigurasi Anda)
cp .env.example .env

# 4. Jalankan server
go run .
```

Jika berhasil:

```
✅ Koneksi MongoDB berhasil
Server jalan di :8080
```

---

## 📖 API Documentation

### 📌 Header Wajib

| Header          | Deskripsi            | Required                     |
| --------------- | -------------------- | ---------------------------- |
| `X-API-Key`     | API Key dari .env    | ✅ Semua endpoint `/api/*`   |
| `Authorization` | `Bearer <JWT_TOKEN>` | ✅ Endpoint yang butuh login |

---

### 🔓 Authentication Endpoints

#### Register User

```http
POST /api/auth/register
```

```json
{
  "nama": "John Doe",
  "email": "john@mail.com",
  "password": "123456",
  "nim": "F55124001",
  "jurusan": "Teknik Informatika"
}
```

#### Login User

```http
POST /api/auth/login
```

```json
{
  "email": "john@mail.com",
  "password": "123456"
}
```

---

### 📦 Alat Endpoints

#### List Semua Alat

```http
GET /api/alat
```

#### Detail Alat by ID

```http
GET /api/alat/{id}
```

#### Tambah Alat (Admin Only)

```http
POST /api/admin/alat
```

```json
{
  "nama": "Proyektor Epson",
  "kategori": "Elektronik",
  "deskripsi": "Proyektor ruang kelas",
  "stok_total": 5
}
```

#### Update Alat (Admin Only)

```http
PUT /api/admin/alat/{id}
```

#### Hapus Alat (Admin Only)

```http
DELETE /api/admin/alat/{id}
```

---

### 🔄 Peminjaman Endpoints

#### Pinjam Alat

```http
POST /api/peminjaman
```

```json
{
  "alat_id": "67a35021ea8a689c444a92d0",
  "jumlah": 2
}
```

#### Kembalikan Alat

```http
POST /api/pengembalian/{transaction_id}
```

#### Riwayat Peminjaman Saya

```http
GET /api/riwayat
```

#### List Transaksi Saya (Aktif)

```http
GET /api/peminjaman/me
```

---

### 👑 Admin Endpoints

#### List Semua User

```http
GET /api/admin/users
```

#### Update Role User

```http
PATCH /api/admin/users/{id}/role
```

```json
{
  "role": "admin"
}
```

#### List Semua Transaksi

```http
GET /api/admin/peminjaman
```

#### Riwayat Semua Transaksi

```http
GET /api/admin/riwayat
```

---

### 🏠 Status Server

```http
GET /
```

Response:

```json
{
  "success": true,
  "message": "SIPAK API berjalan 🚀"
}
```

---

## 📊 Database Schema

### User Collection

| Field           | Type     | Description           |
| --------------- | -------- | --------------------- |
| `_id`           | ObjectID | Primary key           |
| `nama`          | string   | Nama lengkap          |
| `email`         | string   | Email (unique)        |
| `password_hash` | string   | Password ter-hash     |
| `role`          | string   | `admin` / `mahasiswa` |
| `nim`           | string   | NIM mahasiswa         |
| `jurusan`       | string   | Jurusan               |
| `created_at`    | datetime | Waktu registrasi      |

### Alat Collection

| Field           | Type     | Description        |
| --------------- | -------- | ------------------ |
| `_id`           | ObjectID | Primary key        |
| `nama`          | string   | Nama alat          |
| `kategori`      | string   | Kategori alat      |
| `deskripsi`     | string   | Deskripsi alat     |
| `stok_total`    | int      | Total stok         |
| `stok_tersedia` | int      | Stok yang tersedia |
| `created_at`    | datetime | Waktu dibuat       |
| `updated_at`    | datetime | Waktu update       |

### Transaction Collection

| Field             | Type     | Description                |
| ----------------- | -------- | -------------------------- |
| `_id`             | ObjectID | Primary key                |
| `user_id`         | ObjectID | FK ke User                 |
| `alat_id`         | ObjectID | FK ke Alat                 |
| `jumlah`          | int      | Jumlah dipinjam            |
| `tanggal_pinjam`  | datetime | Tanggal pinjam             |
| `tanggal_kembali` | datetime | Tanggal kembali (nullable) |
| `status`          | string   | `PINJAM` / `KEMBALI`       |

---

## 🔒 Security Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│  API Key    │────▶│    JWT      │
│  (Request)  │     │ Middleware  │     │ Middleware  │
└─────────────┘     └─────────────┘     └─────────────┘
                           │                   │
                           ▼                   ▼
                    ┌─────────────┐     ┌─────────────┐
                    │   Reject    │     │  Handler    │
                    │   (401)     │     │  (Success)  │
                    └─────────────┘     └─────────────┘
```

---

## 📋 API Response Format

### Success Response

```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": { ... }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Pesan error"
}
```

---

## 👨‍💻 Tim Pengembang

**Kelompok 4 - Rekayasa Perangkat Luna (RPL)**

Ketua Kelompok : Mia Islamia F5512114

Teknik Informatika - Universitas Tadulako

FrontEnd:

Moh Reza Dwi Syahputra F55124085

BackEnd :

Moh Magribi Ramadhan F55124104

Moh Fiqih F55124108


## 📝 License

Project ini dikembangkan untuk keperluan akademik (Tugas Besar Rekayasa Perangkat Lunak).

---

<p align="center">
  Made with ❤️ using Go
</p>
