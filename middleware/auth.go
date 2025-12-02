package middleware

import (
	"context"
	"net/http"
	"strings"

	"SIPAK/config"
	"SIPAK/utils"
)

// contextKey dipakai untuk menyimpan data user di context request
type contextKey string

const (
	ContextUserID contextKey = "userID"
	ContextRole   contextKey = "role"
)

// APIKeyMiddleware memeriksa header X-API-Key
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" || apiKey != config.AppConfig.APIKey {
			utils.WriteError(w, http.StatusUnauthorized, "API key invalid atau tidak ada")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware memeriksa JWT di header Authorization: Bearer <token>
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, "Authorization header tidak ditemukan")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.WriteError(w, http.StatusUnauthorized, "Format Authorization salah (harus Bearer token)")
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "Token invalid atau kadaluarsa")
			return
		}

		// Simpan userID & role ke context supaya bisa dipakai di handler
		ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextRole, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly middleware yang memastikan role = admin
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(ContextRole).(string)
		if !ok || role != "admin" {
			utils.WriteError(w, http.StatusForbidden, "Hanya admin yang bisa mengakses endpoint ini")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext helper untuk ambil userID di handler
func GetUserIDFromContext(r *http.Request) string {
	id, _ := r.Context().Value(ContextUserID).(string)
	return id
}

// GetRoleFromContext helper untuk ambil role di handler
func GetRoleFromContext(r *http.Request) string {
	role, _ := r.Context().Value(ContextRole).(string)
	return role
}
