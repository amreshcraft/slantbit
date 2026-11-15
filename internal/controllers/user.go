package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/amreshcraft/slantbit/internal/models"
	"github.com/amreshcraft/slantbit/internal/repository"
	"github.com/amreshcraft/slantbit/internal/service"
)

type UserController struct {
	service *service.UserService
	logger  *slog.Logger
}

func NewUserController(s *service.UserService, l *slog.Logger) *UserController {
	return &UserController{service: s, logger: l}
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.FetchAllUsers(r.Context())
	if err != nil {
		c.logger.Error("failed to fetch users", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	createdUser, err := c.service.RegisterUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) || errors.Is(err, repository.ErrDuplicatePhone) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.logger.Info("user registered", "user_id", createdUser.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}