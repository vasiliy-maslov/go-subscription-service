package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/vasiliy-maslov/go-subscription-service/internal/config"
	"github.com/vasiliy-maslov/go-subscription-service/internal/repository"
	"github.com/vasiliy-maslov/go-subscription-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Service service.SubscriptionService
}

func New(logger *slog.Logger) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфиг: %w", err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("не удалось пингануть базу данных: %w", err)
	}

	repo := repository.NewSubscriptionRepo(dbpool)
	subService := service.NewSubscriptionService(repo, logger)

	return &App{Service: subService}, nil
}
