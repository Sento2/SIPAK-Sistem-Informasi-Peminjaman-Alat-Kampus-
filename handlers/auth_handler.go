package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"SIPAK/config"
	"SIPAK/models"
	"SIPAK/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler menampung method untuk auth (login, register)
type AuthHandler struct{}

// Request body untuk register
type registerRequest struct {
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Request body untuk login
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register membuat user baru (default role: mahasiswa)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	defer r.Body.Close()

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" || req.Nama == "" {
		utils.WriteError(w, http.StatusBadRequest, "Nama, email, dan password wajib diisi")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Cek apakah email sudah digunakan
	var existing models.User
	err := config.UserCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existing)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, "Email sudah terdaftar")
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal hash password")
		return
	}

	user := models.User{
		ID:           primitive.NewObjectID(),
		Nama:         req.Nama,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "mahasiswa",
		CreatedAt:    time.Now(),
	}

	_, err = config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal menyimpan user")
		return
	}

	// Jangan kembalikan password hash
	user.PasswordHash = ""

	utils.WriteJSON(w, http.StatusCreated, utils.JSONResponse{
		Success: true,
		Message: "Registrasi berhasil",
		Data:    user,
	})
}

// Login memverifikasi user dan mengembalikan JWT
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	defer r.Body.Close()

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Email dan password wajib diisi")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := config.UserCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Email atau password salah")
		return
	}

	// Cocokkan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Email atau password salah")
		return
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID.Hex(), user.Role)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Gagal membuat token")
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Success: true,
		Message: "Login berhasil",
		Data: map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":    user.ID.Hex(),
				"nama":  user.Nama,
				"email": user.Email,
				"role":  user.Role,
			},
		},
	})
}
