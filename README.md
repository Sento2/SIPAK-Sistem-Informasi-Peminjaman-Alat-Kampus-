# 🎓 SIPAK — Sistem Informasi Peminjaman Alat Kampus

SIPAK adalah aplikasi berbasis **Golang** yang digunakan untuk mengelola peminjaman alat di kampus.  
Fitur utama:

- Manajemen akun **mahasiswa** & **admin**
- Daftar alat (CRUD alat oleh admin)
- Peminjaman & pengembalian alat
- Manajemen role user (admin/mahasiswa)
- Keamanan dengan **JWT** dan **API Key**
- Database menggunakan **MongoDB Atlas**

---

## 🧱 Tech Stack

- **Backend**: Go + [Chi Router](https://github.com/go-chi/chi)
- **Database**: MongoDB Atlas
- **Auth**:
  - JWT (JSON Web Token)
  - API Key (header `X-API-Key`)
- **Library utama**:
  - `github.com/go-chi/chi/v5`
  - `go.mongodb.org/mongo-driver`
  - `github.com/golang-jwt/jwt/v5`
  - `golang.org/x/crypto`
  - `github.com/joho/godotenv`

---

## 📁 Struktur Folder

```bash
sipak/
├── go.mod
├── main.go
├── .env               # konfigurasi environment (jangan di-commit)
├── config/
│   └── config.go      # koneksi MongoDB & konfigurasi global
├── models/
│   ├── user.go        # model User
│   ├── alat.go        # model Alat
│   └── transaction.go # model Transaction (peminjaman)
├── utils/
│   ├── jwt.go         # helper JWT
│   └── response.go    # helper response JSON
├── middleware/
│   └── auth.go        # middleware API Key, JWT, AdminOnly
└── handlers/
    ├── auth_handler.go        # login & register
    ├── alat_handler.go        # CRUD alat
    ├── peminjaman_handler.go  # peminjaman & pengembalian
    └── user_handler.go        # manajemen user (admin)
```
⚙️ Konfigurasi Environment
- **Buat file .env di root project**:
```
  MONGO_URL=mongodb+srv://user:password@cluster0.xxxxxx.mongodb.net/?retryWrites=true&w=majority
  DB_NAME=sipak_db
  JWT_SECRET=supersecretjwt
  API_KEY=supersecretapikey

  PORT=8080
  
  JWT_SECRET=supersecretjwt
  API_KEY=supersecretapikey
```
Catatan:

MONGO_URI → URI dari MongoDB Atlas

DB_NAME → nama database yang akan dipakai

JWT_SECRET → secret key untuk JWT

API_KEY → API key yang harus dikirim via header X-API-Key

PORT → port server (default 8080 kalau kosong)

## 🚀 Cara Menjalankan Project
1. Clone repo & masuk ke folder project:
   git clone https://github.com/Sento2/SIPAK-Sistem-Informasi-Peminjaman-Alat-Kampus-
   cd SIPAK-Sistem-Informasi-Peminjaman-Alat-Kampus-
2. Install Dependency dan Driver Mongo DB
   go mod tidy
   go get go.mongodb.org/mongo-driver/v2/mongo
3. Pastikan .env sudah dibuat dengan benar.
4. Jalankan server:
   go run .
5. Jika berhasil, akan muncul log:
   ✅ Koneksi MongoDB berhasil
   Server jalan di :8080

# 🔐 Alur Penggunaan API — SIPAK (Sistem Informasi Peminjaman Alat Kampus)

Dokumentasi ini menjelaskan seluruh endpoint yang tersedia dalam API SIPAK, lengkap dengan:
- URL Endpoint  
- Header wajib  
- Contoh request  
- Contoh body JSON  

Semua endpoint berada pada prefix:
```
http://127.0.0.1:3000/api
```

Untuk server hosting, sesuaikan domain Anda.

---

## 📌 Header Wajib

| Header | Nilai |
|--------|--------|
| `x-api-key` | API_KEY di .env |
| `Authorization` | Bearer `<JWT_TOKEN>` *(hanya endpoint tertentu)* |

---

# 🟣 AUTHENTICATION

## 📝 Register User
**POST** `/api/auth/register`

### 📤 Request Body
```json
{
  "nama": "Kelompok 8",
  "email": "kelompok8@mail.com",
  "password": "123456",
  "nim": "A11.2024.00123",
  "jurusan": "Teknik Informatika"
}
```

---

## 🔑 Login User
**POST** `/api/auth/login`

### 📤 Request Body
```json
{
  "email": "kelompok8@mail.com",
  "password": "123456"
}
```

### 📥 Response (token)
```json
{
  "token": "JWT_TOKEN_HERE"
}
```

---

# 🟦 ALAT (Mahasiswa & Admin)

## 📄 List Alat
**GET** `/api/alat`

---

## 🔍 Detail Alat
**GET** `/api/alat/{id}`

Contoh:
```
GET /api/alat/67a35021ea8a689c444a92d0
```

---

# 🟥 ALAT – ADMIN ONLY

## ➕ Tambah Alat
**POST** `/api/admin/alat`
```json
{
  "nama": "Proyektor Epson",
  "kategori": "Elektronik",
  "deskripsi": "Proyektor ruang kelas",
  "stok_total": 5
}
```

---

## ✏️ Update Alat
**PUT** `/api/admin/alat/{id}`
```json
{
  "nama": "Proyektor Epson X200",
  "kategori": "Elektronik",
  "deskripsi": "Update spesifikasi",
  "stok_total": 8
}
```

---

## 🗑️ Hapus Alat
**DELETE** `/api/admin/alat/{id}`

---

# 🟩 PEMINJAMAN (User Login)

## 📦 Pinjam Alat
**POST** `/api/peminjaman`
```json
{
  "alat_id": "67a35021ea8a689c444a92d0",
  "jumlah": 2
}
```

---

## 📤 Kembalikan Alat
**POST** `/api/pengembalian/{id}`
```
POST /api/pengembalian/67a3790abbb123abc902f11
```

---

## 📚 Riwayat Peminjaman Saya
**GET** `/api/peminjaman/me`

---

# 🟥 PEMINJAMAN – ADMIN ONLY

## 📋 Semua Transaksi
**GET** `/api/admin/peminjaman`

---

# 🟧 USER MANAGEMENT – ADMIN ONLY

## 👥 List User
**GET** `/api/admin/users`

---

## 🔄 Update Role User
**PATCH** `/api/admin/users/{id}/role`

```json
{
  "role": "admin"
}
```
Atau:
```json
{
  "role": "mahasiswa"
}
```

---

# 🟢 STATUS SERVER
**GET** `/`

Response:
```json
{
  "success": true,
  "message": "SIPAK API berjalan 🚀"
}
```