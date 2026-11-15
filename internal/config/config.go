package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	AppEnv string
	DBUrl  string
	RazorpayKeyID     string
	RazorpayKeySecret string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Port:   getEnv("PORT", "8080"),
		AppEnv: getEnv("APP_ENV", "development"),
		DBUrl:  getEnv("DATABASE_URL", ""),
		RazorpayKeyID:     getEnv("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret: getEnv("RAZORPAY_KEY_SECRET", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}