package main

import (
	"fmt"
	"log"
	"net/http"

	"SIPAK/config"
	"SIPAK/handlers"
	"SIPAK/middleware"
	"SIPAK/utils"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"github.com/go-chi/chi/v5"
)

func main() {
	// 1. Load konfigurasi dari .env
	config.LoadConfig()

	// 2. Konek ke MongoDB Atlas
	config.ConnectMongo()

	// 3. Setup router Chi
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}).Handler)


	// Tambah middleware API key global untuk semua endpoint /api
	r.Route("/api", func(api chi.Router) {
		// Semua endpoint di bawah /api harus pakai API Key
		api.Use(middleware.APIKeyMiddleware)

		// ==== AUTH (tanpa JWT, tapi wajib API Key) ====
		authHandler := &handlers.AuthHandler{}
		api.Post("/auth/register", authHandler.Register)
		api.Post("/auth/login", authHandler.Login)

		// ==== ENDPOINT YANG BUTUH JWT ====
		api.Group(func(priv chi.Router) {
			// Semua endpoint di group ini butuh JWT
			priv.Use(middleware.AuthMiddleware)

			// ----- Alat -----
			alatHandler := &handlers.AlatHandler{}
			priv.Get("/alat", alatHandler.ListAlat)
			priv.Get("/alat/{id}", alatHandler.GetAlatByID)

			// ----- Peminjaman -----
			pinjamHandler := &handlers.PeminjamanHandler{}
			priv.Post("/peminjaman", pinjamHandler.PinjamAlat)
			priv.Post("/pengembalian/{id}", pinjamHandler.KembalikanAlat)
			priv.Get("/peminjaman/me", pinjamHandler.ListTransaksiUser)
			priv.Get("/riwayat", pinjamHandler.RiwayatSaya)

			// ----- Admin only group -----
			priv.Group(func(admin chi.Router) {
				admin.Use(middleware.AdminOnly)

				// CRUD alat admin
				admin.Post("/admin/alat", alatHandler.CreateAlat)
				admin.Put("/admin/alat/{id}", alatHandler.UpdateAlat)
				admin.Delete("/admin/alat/{id}", alatHandler.DeleteAlat)

				// User management admin
				userHandler := &handlers.UserHandler{}
				admin.Get("/admin/users", userHandler.ListUsers)
				admin.Patch("/admin/users/{id}/role", userHandler.UpdateUserRole)

				// Semua transaksi (admin)
				admin.Get("/admin/peminjaman", pinjamHandler.ListSemuaTransaksi)
				admin.Get("/admin/riwayat", pinjamHandler.RiwayatSemua)
			})
		})
	})

	// Root endpoint sederhana untuk cek status API
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
			Success: true,
			Message: "SIPAK API berjalan 🚀",
		})
	})

	addr := ":" + config.AppConfig.Port
	fmt.Println("Server jalan di", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
