package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/amreshcraft/slantbit/internal/models"
	"github.com/amreshcraft/slantbit/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FetchAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}




func (s *UserService) FetchUserByEmail(ctx context.Context, email string) (models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return models.User{}, errors.New("email cannot be empty")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}

	// Security: API response se password hash strip kar dein
	user.Password = ""
	return user, nil
}

func (s *UserService) RegisterUser(ctx context.Context, user models.User) (models.User, error) {
	// 1. Basic Validation
	if user.Name == "" || user.Email == "" || user.Phone == "" || user.Password == "" {
		return models.User{}, errors.New("name, phone, email, and password are required")
	}

	// 2. Hash Password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedBytes)

	// 3. Save to DB
	return s.repo.Create(ctx, user)
}