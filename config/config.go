package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config menampung semua konfigurasi utama aplikasi
type Config struct {
	MongoURI  string
	DBName    string
	JWTSecret string
	APIKey    string
	Port      string
}

// AppConfig adalah variabel global untuk menyimpan konfigurasi
var AppConfig Config

// MongoClient adalah client global MongoDB
var MongoClient *mongo.Client

// Koleksi global agar mudah dipakai di handler
var (
	UserCollection         *mongo.Collection
	AlatCollection         *mongo.Collection
	TransactionCollection  *mongo.Collection
)

// LoadConfig membaca file .env lalu isi AppConfig
func LoadConfig() {
	// Muat file .env (jika ada)
	_ = godotenv.Load()

	AppConfig = Config{
		MongoURI:  os.Getenv("MONGO_URI"),
		DBName:    os.Getenv("DB_NAME"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		APIKey:    os.Getenv("API_KEY"),
		Port:      os.Getenv("PORT"),
	}

	if AppConfig.Port == "" {
		AppConfig.Port = "8080"
	}

	// Validasi sederhana
	if AppConfig.MongoURI == "" || AppConfig.DBName == "" {
		log.Fatal("MONGO_URI atau DB_NAME belum di-set di .env")
	}
	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET belum di-set di .env")
	}
	if AppConfig.APIKey == "" {
		log.Fatal("API_KEY belum di-set di .env")
	}
}

// ConnectMongo menghubungkan aplikasi ke MongoDB Atlas
func ConnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(AppConfig.MongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Gagal konek MongoDB: %v", err)
	}

	// Tes koneksi
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Gagal ping MongoDB: %v", err)
	}

	MongoClient = client
	db := client.Database(AppConfig.DBName)

	// Inisialisasi koleksi
	UserCollection = db.Collection("users")
	AlatCollection = db.Collection("alat")
	TransactionCollection = db.Collection("transactions")

	fmt.Println("✅ Koneksi MongoDB berhasil")
}
