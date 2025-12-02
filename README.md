# 🎓 SIPAK — Sistem Informasi Peminjaman Alat Kampus

SIPAK adalah aplikasi **REST API** berbasis **Golang** yang digunakan untuk mengelola peminjaman alat di kampus.  
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
Buat file .env di root project:
  -`MONGO_URI=mongodb+srv://user:password@cluster0.xxxxxx.mongodb.net/?retryWrites=true&w=majority`
  -`DB_NAME=sipak_db`
  -`JWT_SECRET=supersecretjwt`
  -`API_KEY=supersecretapikey`

  -`PORT=8080`

  -`JWT_SECRET=supersecretjwt`
  -`API_KEY=supersecretapikey`

Catatan:

MONGO_URI → URI dari MongoDB Atlas

DB_NAME → nama database yang akan dipakai

JWT_SECRET → secret key untuk JWT

API_KEY → API key yang harus dikirim via header X-API-Key

PORT → port server (default 8080 kalau kosong)

🚀 Cara Menjalankan Project
1. Clone repo & masuk ke folder project:
   

PORT=8080

