package repository

import (
	"context"
	"errors"

	"github.com/vasiliy-maslov/go-subscription-service/internal/model"

	"github.com/google/uuid"
)

// ErrNotFound возвращается, когда подписка не найдена в хранилище.
var ErrNotFound = errors.New("subscription not found")

// SubscriptionRepository определяет методы для работы с хранилищем подписок.
type SubscriptionRepository interface {
	Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, sub model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUserID(ctx context.Context, userID uuid.UUID, serviceName *string) ([]model.Subscription, error)
}
