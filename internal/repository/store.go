package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"slantbit/internal/models"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyReviewed = errors.New("already reviewed")

type Store struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) ListCourses(ctx context.Context) ([]models.Course, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, title, slug, description, price_paise
		 FROM courses WHERE active ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := []models.Course{}
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.PricePaise); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}

func (s *Store) GetCourse(ctx context.Context, id int64) (*models.Course, error) {
	var c models.Course
	err := s.db.QueryRow(ctx,
		`SELECT id, title, slug, description, price_paise, whatsapp_link, pdf_s3_key
		 FROM courses WHERE id=$1 AND active`, id).
		Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.PricePaise, &c.WhatsappLink, &c.PdfKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &c, err
}

func (s *Store) CreateOrder(ctx context.Context, courseID int64, name, email, rzpOrderID string, amount int64) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO orders (course_id, buyer_name, buyer_email, razorpay_order_id, amount_paise)
		 VALUES ($1,$2,$3,$4,$5)`, courseID, name, email, rzpOrderID, amount)
	return err
}

// CompletePurchase idempotent hai: dobara call hone par (nil, nil) return karta hai
func (s *Store) CompletePurchase(ctx context.Context, rzpOrderID, paymentID string) (*models.PurchaseResult, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var orderID, courseID int64
	var name, email string
	err = tx.QueryRow(ctx,
		`UPDATE orders SET status='PAID', razorpay_payment_id=$2
		 WHERE razorpay_order_id=$1 AND status='CREATED'
		 RETURNING id, course_id, buyer_name, buyer_email`, rzpOrderID, paymentID).
		Scan(&orderID, &courseID, &name, &email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil // already processed ya order exist nahi karta
	}
	if err != nil {
		return nil, err
	}

	token, err := newToken()
	if err != nil {
		return nil, err
	}

	if _, err = tx.Exec(ctx,
		`INSERT INTO enrollments (order_id, course_id, buyer_name, buyer_email, access_token)
		 VALUES ($1,$2,$3,$4,$5)`, orderID, courseID, name, email, token); err != nil {
		return nil, err
	}

	var title string
	if err = tx.QueryRow(ctx, `SELECT title FROM courses WHERE id=$1`, courseID).Scan(&title); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.PurchaseResult{
		BuyerName: name, BuyerEmail: email, CourseTitle: title, AccessToken: token,
	}, nil
}

func (s *Store) GetAccessByToken(ctx context.Context, token string) (*models.Access, error) {
	var a models.Access
	err := s.db.QueryRow(ctx,
		`SELECT e.id, c.id, c.title, c.whatsapp_link, c.pdf_s3_key
		 FROM enrollments e JOIN courses c ON c.id = e.course_id
		 WHERE e.access_token=$1`, token).
		Scan(&a.EnrollmentID, &a.CourseID, &a.CourseTitle, &a.WhatsappLink, &a.PdfKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &a, err
}

func (s *Store) CreateReview(ctx context.Context, enrollmentID, courseID int64, rating int, comment string) error {
	tag, err := s.db.Exec(ctx,
		`INSERT INTO reviews (enrollment_id, course_id, rating, comment)
		 VALUES ($1,$2,$3,$4)
		 ON CONFLICT (enrollment_id) DO NOTHING`, enrollmentID, courseID, rating, comment)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrAlreadyReviewed
	}
	return nil
}

func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}