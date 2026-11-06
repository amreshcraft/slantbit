package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	FrontendURL string

	RazorpayKeyID         string
	RazorpayKeySecret     string
	RazorpayWebhookSecret string

	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	SMTPFrom string

	AWSRegion string
	S3Bucket  string
}

func Load() (*Config, error) {
	c := &Config{
		Port:                  getenv("PORT", "8080"),
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		FrontendURL:           getenv("FRONTEND_URL", "http://localhost:3000"),
		RazorpayKeyID:         os.Getenv("RAZORPAY_KEY_ID"),
		RazorpayKeySecret:     os.Getenv("RAZORPAY_KEY_SECRET"),
		RazorpayWebhookSecret: os.Getenv("RAZORPAY_WEBHOOK_SECRET"),
		SMTPHost:              os.Getenv("SMTP_HOST"),
		SMTPPort:              getenv("SMTP_PORT", "587"),
		SMTPUser:              os.Getenv("SMTP_USER"),
		SMTPPass:              os.Getenv("SMTP_PASS"),
		SMTPFrom:              os.Getenv("SMTP_FROM"),
		AWSRegion:             getenv("AWS_REGION", "ap-south-1"),
		S3Bucket:              os.Getenv("S3_BUCKET"),
	}

	required := map[string]string{
		"DATABASE_URL":            c.DatabaseURL,
		"RAZORPAY_KEY_ID":         c.RazorpayKeyID,
		"RAZORPAY_KEY_SECRET":     c.RazorpayKeySecret,
		"RAZORPAY_WEBHOOK_SECRET": c.RazorpayWebhookSecret,
		"S3_BUCKET":               c.S3Bucket,
	}
	for k, v := range required {
		if v == "" {
			return nil, fmt.Errorf("missing env var: %s", k)
		}
	}
	return c, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}