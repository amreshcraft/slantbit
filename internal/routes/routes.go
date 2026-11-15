package routes

import (
	"log/slog"
	"net/http"

	"github.com/amreshcraft/slantbit/internal/controllers"
	middlewares "github.com/amreshcraft/slantbit/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(userCtrl *controllers.UserController, log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Built-in Chi Middlewares
	r.Use(middleware.RequestID)
    r.Use(middlewares.RequestLogger(log))   // HTTP request logs
	r.Use(middleware.Recoverer) // Panic recovery (server crash nahi hoga)

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userCtrl.GetUsers)
			r.Post("/", userCtrl.CreateUser)
		})
	})

	return r
}