package main

import (
	"fmt"
	"user-service/internal/app"
	"user-service/internal/config"
	"user-service/internal/infrastructure/logger"
	"user-service/pkg/jwt"
)

func main() {
	cfg := config.LoadConfig()

	DBDSN := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	log := logger.NewSlogLogger()

	jwt.Init(cfg.JWTSecret)

	app := app.New(log, DBDSN, cfg.TokenTTL)

	app.Run(cfg.AppHost, cfg.AppPort)
}
