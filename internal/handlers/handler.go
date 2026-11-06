package handlers

import (
	razorpay "github.com/razorpay/razorpay-go"

	"slantbit/internal/config"
	"slantbit/internal/mailer"
	"slantbit/internal/repository"
	"slantbit/internal/storage"
)

type Handler struct {
	cfg     *config.Config
	store   *repository.Store
	rzp     *razorpay.Client
	mailer  *mailer.Mailer
	storage *storage.Storage
}

func New(cfg *config.Config, store *repository.Store, rzp *razorpay.Client,
	m *mailer.Mailer, st *storage.Storage) *Handler {
	return &Handler{cfg: cfg, store: store, rzp: rzp, mailer: m, storage: st}
}