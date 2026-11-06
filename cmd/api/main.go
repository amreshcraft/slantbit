package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	razorpay "github.com/razorpay/razorpay-go"

	"slantbit/internal/config"
	"slantbit/internal/database"
	"slantbit/internal/handlers"
	"slantbit/internal/mailer"
	"slantbit/internal/repository"
	"slantbit/internal/routes"
	"slantbit/internal/storage"
)

func main() {
	_ = godotenv.Load() // .env na mile to bhi chalega

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	st, err := storage.New(ctx, cfg.AWSRegion, cfg.S3Bucket)
	if err != nil {
		log.Fatal(err)
	}

	store := repository.New(pool)
	rzp := razorpay.NewClient(cfg.RazorpayKeyID, cfg.RazorpayKeySecret)
	mail := mailer.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)

	h := handlers.New(cfg, store, rzp, mail, st)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      routes.New(h, cfg.FrontendURL),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("listening on :" + cfg.Port)
	log.Fatal(srv.ListenAndServe())
}