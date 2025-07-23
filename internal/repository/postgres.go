package repository

import (
	"context"
	"errors"

	"github.com/vasiliy-maslov/go-subscription-service/internal/model" // Проверь имя модуля

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ SubscriptionRepository = (*SubscriptionRepo)(nil)

type SubscriptionRepo struct {
	db *pgxpool.Pool
}

// NewSubscriptionRepo создает новый экземпляр репозитория.
func NewSubscriptionRepo(db *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

// Create создает новую запись о подписке в базе данных.
func (r *SubscriptionRepo) Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error) {
	sub.ID = uuid.New()

	query := `
		INSERT INTO subscriptions (id, user_id, service_name, price, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

	_, err := r.db.Exec(ctx, query,
		sub.ID, sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate)

	if err != nil {
		return uuid.Nil, err
	}

	return sub.ID, nil
}

// GetByID получает подписку по ее ID.
func (r *SubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (model.Subscription, error) {
	var sub model.Subscription

	query := `
		SELECT id, user_id, service_name, price, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscription{}, ErrNotFound
		}
		return model.Subscription{}, err
	}

	return sub, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, id uuid.UUID, sub model.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4, updated_at = NOW()
		WHERE id = $5`

	res, err := r.db.Exec(ctx, query, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete удаляет подписку по ID.
func (r *SubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *SubscriptionRepo) ListByUserID(ctx context.Context, userID uuid.UUID, serviceName *string) ([]model.Subscription, error) {
	query := `
		SELECT id, user_id, service_name, price, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1`

	args := []any{userID} // Начинаем собирать аргументы для запроса

	if serviceName != nil {
		query += " AND service_name = $2"
		args = append(args, *serviceName)
	}

	query += " ORDER BY start_date DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		if err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.ServiceName, &sub.Price,
			&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}
