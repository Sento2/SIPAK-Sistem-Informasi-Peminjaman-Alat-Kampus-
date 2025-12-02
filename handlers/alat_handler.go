package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"SIPAK/config"
	"SIPAK/models"
	"SIPAK/utils"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AlatHandler mengelola CRUD alat
type AlatHandler struct{}

// Request body untuk membuat/mengupdate alat
type alatRequest struct {
	Nama      string `json:"nama"`
	Kategori  string `json:"kategori"`
	Deskripsi string `json:"deskripsi"`
	StokTotal int    `json:"stok_total"`
}

// CreateAlat (admin) menambah alat baru
func (h *AlatHandler) CreateAlat(w http.ResponseWriter, r *http.Request) {
	var req alatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	defer r.Body.Close()

	if req.Nama == "" || req.StokTotal <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Nama dan stok_total wajib, stok_total > 0")
		return
	}

	now := time.Now()
	alat := models.Alat{
		ID:           primitive.NewObjectID(),
		Nama:         req.Nama,
		Kategori:     req.Kategori,
		Deskripsi:    req.Deskripsi,
		StokTotal:    req.StokTotal,
		StokTersedia: req.StokTotal,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.AlatCollection.InsertOne(ctx, alat)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal menyimpan alat")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.JSONResponse{
		Success: true,
		Message: "Alat berhasil ditambahkan",
		Data:    alat,
	})
}

// ListAlat menampilkan daftar alat (public: mahasiswa & admin)
func (h *AlatHandler) ListAlat(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.AlatCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal mengambil data alat")
		return
	}
	defer cursor.Close(ctx)

	var alatList []models.Alat
	if err := cursor.All(ctx, &alatList); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal decode data alat")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Data:    alatList,
	})
}

// GetAlatByID mengambil detail alat
func (h *AlatHandler) GetAlatByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var alat models.Alat
	err = config.AlatCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&alat)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Alat tidak ditemukan")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Data:    alat,
	})
}

// UpdateAlat (admin) mengubah data alat
func (h *AlatHandler) UpdateAlat(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req alatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	defer r.Body.Close()

	update := bson.M{
		"nama":       req.Nama,
		"kategori":   req.Kategori,
		"deskripsi":  req.Deskripsi,
		"updated_at": time.Now(),
	}

	if req.StokTotal > 0 {
		// NOTE: ini sederhana, tidak menghitung stok_tersedia dari transaksi
		update["stok_total"] = req.StokTotal
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.AlatCollection.UpdateByID(ctx, objID, bson.M{"$set": update})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal mengupdate alat")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Message: "Alat berhasil diupdate",
	})
}

// DeleteAlat (admin) menghapus alat
func (h *AlatHandler) DeleteAlat(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.AlatCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal menghapus alat")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Message: "Alat berhasil dihapus",
	})
}
