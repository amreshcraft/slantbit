package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/amreshcraft/slantbit/internal/models")

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrDuplicatePhone = errors.New("phone number already exists")
	ErrUserNotFound   = errors.New("user not found")
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts user with hashed password & auto-generates UUID
func (r *UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	query := `
		INSERT INTO users (name, phone, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, phone, email
	`

	var created models.User
	err := r.db.QueryRow(ctx, query, user.Name, user.Phone, user.Email, user.Password).Scan(
		&created.ID,
		&created.Name,
		&created.Phone,
		&created.Email,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		// 23505 = unique_violation in PostgreSQL
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "phone") {
				return models.User{}, ErrDuplicatePhone
			}
			return models.User{}, ErrDuplicateEmail
		}
		return models.User{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return created, nil
}

// GetAll fetches all users (password column select nahi karenge)
func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, phone, email FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// GetByEmail fetches user including password hash (Login ke liye kaam aayega)
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query := `SELECT id, name, phone, email, password FROM users WHERE email = $1`

	var u models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Phone,
		&u.Email,
		&u.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return u, nil
}