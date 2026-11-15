package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/amreshcraft/slantbit/internal/config"
	"github.com/amreshcraft/slantbit/internal/controllers"
	"github.com/amreshcraft/slantbit/internal/database"
	"github.com/amreshcraft/slantbit/internal/logger"
	"github.com/amreshcraft/slantbit/internal/repository"
	"github.com/amreshcraft/slantbit/internal/routes"
	"github.com/amreshcraft/slantbit/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	// Init Logger with file support
	log, closeLog := logger.InitLogger(cfg.AppEnv)
	defer closeLog()

	dbPool, err := database.ConnectPostgres(cfg.DBUrl)
	if err != nil {
		log.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	log.Info("Connected to PostgreSQL successfully")

	userRepo := repository.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepo)
	userController := controllers.NewUserController(userService, log)

	// Pass 'log' to SetupRouter
	router := routes.SetupRouter(userController, log)

	addr := ":" + cfg.Port
	log.Info("Server starting", "port", cfg.Port, "env", cfg.AppEnv)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}