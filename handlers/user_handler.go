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

// UserHandler mengelola endpoint admin terkait user
type UserHandler struct{}

// Request untuk update role user
type updateRoleRequest struct {
	Role string `json:"role"`
}

// ListUsers (admin) menampilkan semua user
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal mengambil data user")
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal decode data user")
		return
	}

	// Hapus password hash sebelum dikirim ke client
	for i := range users {
		users[i].PasswordHash = ""
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Data:    users,
	})
}

// UpdateUserRole (admin) mengubah role user (admin / mahasiswa)
func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "ID user tidak valid")
		return
	}

	var req updateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	defer r.Body.Close()

	if req.Role != "admin" && req.Role != "mahasiswa" {
		utils.WriteError(w, http.StatusBadRequest, "Role harus 'admin' atau 'mahasiswa'")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.UserCollection.UpdateByID(ctx, userID, bson.M{
		"$set": bson.M{"role": req.Role},
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal update role user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Message: "Role user berhasil diupdate",
	})
}
