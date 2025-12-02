package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User merepresentasikan akun dalam sistem (admin / mahasiswa)
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nama         string             `bson:"nama" json:"nama"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Role         string             `bson:"role" json:"role"` // "admin" atau "mahasiswa"
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	NIM		     string             `bson:"nim,omitempty" json:"nim,omitempty"`
	Jurusan	 	 string             `bson:"jurusan,omitempty" json:"jurusan,omitempty"`
}
