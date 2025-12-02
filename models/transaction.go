package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction menyimpan data peminjaman / pengembalian
type Transaction struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	AlatID          primitive.ObjectID `bson:"alat_id" json:"alat_id"`
	Jumlah          int                `bson:"jumlah" json:"jumlah"`
	TanggalPinjam   time.Time          `bson:"tanggal_pinjam" json:"tanggal_pinjam"`
	TanggalKembali  *time.Time         `bson:"tanggal_kembali,omitempty" json:"tanggal_kembali,omitempty"`
	Status          string             `bson:"status" json:"status"` // "PINJAM", "KEMBALI"
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
