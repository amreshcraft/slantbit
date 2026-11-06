package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"slantbit/internal/handlers"
)

func New(h *handlers.Handler, frontendURL string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{frontendURL},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	r.Get("/api/courses", h.ListCourses)
	r.Post("/api/orders", h.CreateOrder)
	r.Post("/api/v1/webhooks/razorpay", h.RazorpayWebhook)

	r.Get("/api/access/{token}", h.GetAccess)
	r.Post("/api/access/{token}/reviews", h.CreateReview)

	return r
}