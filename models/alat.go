package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Alat adalah entitas alat kampus yang bisa dipinjam
type Alat struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nama          string             `bson:"nama" json:"nama"`
	Kategori      string             `bson:"kategori" json:"kategori"`
	Deskripsi     string             `bson:"deskripsi" json:"deskripsi"`
	StokTotal     int                `bson:"stok_total" json:"stok_total"`
	StokTersedia  int                `bson:"stok_tersedia" json:"stok_tersedia"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
