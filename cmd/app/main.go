package main

import (
	"log/slog"
	"os"

	"github.com/vasiliy-maslov/go-subscription-service/internal/app"
	"github.com/vasiliy-maslov/go-subscription-service/internal/handler/http"
)

// @title Subscription Service API
// @version 1.0
// @description API Server for Subscription Management Application

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger) // Устанавливаем его как логгер по умолчанию

	application, err := app.New(logger)
	if err != nil {
		logger.Error("Не удалось запустить приложение", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Приложение успешно инициализировано.")

	handler := http.NewHandler(application.Service, logger)

	logger.Info("Запускаем HTTP-сервер", slog.String("port", "8080"))

	if err := handler.InitRoutes().Run(":8080"); err != nil {
		logger.Error("Не удалось запустить HTTP-сервер", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
