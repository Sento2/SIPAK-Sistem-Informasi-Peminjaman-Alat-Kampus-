package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"SIPAK/config"
	"SIPAK/middleware"
	"SIPAK/models"
	"SIPAK/utils"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PeminjamanHandler mengelola peminjaman dan pengembalian
type PeminjamanHandler struct{}

// Request body peminjaman
type peminjamanRequest struct {
	AlatID string `json:"alat_id"`
	Jumlah int    `json:"jumlah"`
}

type RiwayatPeminjamanResponse struct {
	ID			   primitive.ObjectID `json:"id"`
	AlatID 	       primitive.ObjectID `json:"alat_id"`
	NamaAlat       string             `json:"nama_alat"`
	Jumlah         int                `json:"jumlah"`
	TanggalPinjam  time.Time       `json:"tanggal_pinjam"`
	TanggalKembali *time.Time      `json:"tanggal_kembali,omitempty"`
	Status         string            `json:"status"`
}

// PinjamAlat membuat transaksi peminjaman untuk user yg login
func (h *PeminjamanHandler) PinjamAlat(w http.ResponseWriter, r *http.Request) {
	var req peminjamanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	defer r.Body.Close()

	if req.AlatID == "" || req.Jumlah <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "alat_id dan jumlah wajib, jumlah > 0")
		return
	}

	userIDHex := middleware.GetUserIDFromContext(r)
	userObjID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User ID invalid di token")
		return
	}

	alatID, err := primitive.ObjectIDFromHex(req.AlatID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "alat_id tidak valid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ambil alat dulu
	var alat models.Alat
	if err := config.AlatCollection.FindOne(ctx, bson.M{"_id": alatID}).Decode(&alat); err != nil {
		utils.WriteError(w, http.StatusNotFound, "Alat tidak ditemukan")
		return
	}

	// Cek stok tersedia
	if alat.StokTersedia < req.Jumlah {
		utils.WriteError(w, http.StatusBadRequest, "Stok alat tidak mencukupi")
		return
	}

	// Kurangi stok tersedia
	_, err = config.AlatCollection.UpdateByID(ctx, alatID, bson.M{
		"$inc": bson.M{"stok_tersedia": -req.Jumlah},
		"$set": bson.M{"updated_at": time.Now()},
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal mengupdate stok alat")
		return
	}

	now := time.Now()
	trans := models.Transaction{
		ID:             primitive.NewObjectID(),
		UserID:         userObjID,
		AlatID:         alatID,
		Jumlah:         req.Jumlah,
		TanggalPinjam:  now,
		Status:         "PINJAM",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	_, err = config.TransactionCollection.InsertOne(ctx, trans)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal membuat transaksi peminjaman")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.JSONResponse{
		Success: true,
		Message: "Peminjaman berhasil",
		Data:    trans,
	})
}

// KembalikanAlat mengubah status transaksi menjadi KEMBALI dan menambah stok alat
func (h *PeminjamanHandler) KembalikanAlat(w http.ResponseWriter, r *http.Request) {
	transIDParam := chi.URLParam(r, "id")
	transID, err := primitive.ObjectIDFromHex(transIDParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "ID transaksi tidak valid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var trans models.Transaction
	if err := config.TransactionCollection.FindOne(ctx, bson.M{"_id": transID}).Decode(&trans); err != nil {
		utils.WriteError(w, http.StatusNotFound, "Transaksi tidak ditemukan")
		return
	}

	if trans.Status == "KEMBALI" {
		utils.WriteError(w, http.StatusBadRequest, "Transaksi sudah dikembalikan")
		return
	}

	// Tambahkan stok alat kembali
	_, err = config.AlatCollection.UpdateByID(ctx, trans.AlatID, bson.M{
		"$inc": bson.M{"stok_tersedia": trans.Jumlah},
		"$set": bson.M{"updated_at": time.Now()},
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal menambah stok alat")
		return
	}

	now := time.Now()
	update := bson.M{
		"status":          "KEMBALI",
		"tanggal_kembali": &now,
		"updated_at":      now,
	}

	_, err = config.TransactionCollection.UpdateByID(ctx, trans.ID, bson.M{"$set": update})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal update status transaksi")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Message: "Pengembalian berhasil",
	})
}

// ListTransaksiUser menampilkan semua transaksi milik user yg login
func (h *PeminjamanHandler) ListTransaksiUser(w http.ResponseWriter, r *http.Request) {
	userIDHex := middleware.GetUserIDFromContext(r)
	userObjID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User ID invalid di token")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.TransactionCollection.Find(ctx, bson.M{"user_id": userObjID})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal mengambil data transaksi")
		return
	}
	defer cursor.Close(ctx)

	var list []models.Transaction
	if err := cursor.All(ctx, &list); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal decode data transaksi")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Data:    list,
	})
}

// ListSemuaTransaksi (admin) menampilkan semua transaksi
func (h *PeminjamanHandler) ListSemuaTransaksi(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.TransactionCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal mengambil data transaksi")
		return
	}
	defer cursor.Close(ctx)

	var list []models.Transaction
	if err := cursor.All(ctx, &list); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal decode data transaksi")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Data:    list,
	})
}

// RiwayatSaya menampilkan riwayat peminjaman milik user yang sedang login,
// lengkap dengan nama alat (join ke koleksi alat)
func (h *PeminjamanHandler) RiwayatSaya(w http.ResponseWriter, r *http.Request) {
	userIDHex := middleware.GetUserIDFromContext(r)
	userObjID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User ID invalid di token")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"user_id": userObjID}}},
		bson.D{
			{Key: "$lookup", Value: bson.M{
				"from":         "alat",
				"localField":   "alat_id",
				"foreignField": "_id",
				"as":           "alat",
			}},
		},
		bson.D{
			{Key: "$unwind", Value: bson.M{
				"path":                       "$alat",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.M{
				"_id":             1,
				"alat_id":         "$alat_id",
				"jumlah":          1,
				"tanggal_pinjam":  1,
				"tanggal_kembali": 1,
				"status":          1,
				"nama_alat":       "$alat.nama",
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.M{"tanggal_pinjam": -1}},
		},
	}

	cursor, err := config.TransactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal mengambil data riwayat peminjaman")
		return
	}
	defer cursor.Close(ctx)

	var riwayat []RiwayatPeminjamanResponse
	if err := cursor.All(ctx, &riwayat); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal decode data riwayat peminjaman")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Data:    riwayat,
	})
}

// RiwayatSemua menampilkan riwayat semua transaksi (hanya admin)
func (h *PeminjamanHandler) RiwayatSemua(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$lookup", Value: bson.M{
				"from":         "alat",
				"localField":   "alat_id",
				"foreignField": "_id",
				"as":           "alat",
			}},
		},
		bson.D{
			{Key: "$unwind", Value: bson.M{
				"path":                       "$alat",
				"preserveNullAndEmptyArrays": true,
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.M{
				"_id":             1,
				"alat_id":         "$alat_id",
				"jumlah":          1,
				"tanggal_pinjam":  1,
				"tanggal_kembali": 1,
				"status":          1,
				"nama_alat":       "$alat.nama",
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.M{"tanggal_pinjam": -1}},
		},
	}

	cursor, err := config.TransactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal mengambil data riwayat peminjaman")
		return
	}
	defer cursor.Close(ctx)
	var riwayat []RiwayatPeminjamanResponse
	if err := cursor.All(ctx, &riwayat); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal decode data riwayat peminjaman")
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Data:    riwayat,
	})
}

