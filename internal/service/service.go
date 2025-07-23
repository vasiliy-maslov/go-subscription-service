package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/vasiliy-maslov/go-subscription-service/internal/model"
	"github.com/vasiliy-maslov/go-subscription-service/internal/repository"

	"github.com/google/uuid"
)

// SubscriptionService определяет интерфейс для бизнес-логики работы с подписками.
type SubscriptionService interface {
	Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, sub model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateTotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, startPeriod, endPeriod time.Time) (int, error)
}

type subscriptionService struct {
	repo   repository.SubscriptionRepository
	logger *slog.Logger
}

// NewSubscriptionService создает новый экземпляр сервиса.
func NewSubscriptionService(repo repository.SubscriptionRepository, logger *slog.Logger) SubscriptionService {
	return &subscriptionService{
		repo:   repo,
		logger: logger,
	}
}

func (s *subscriptionService) Create(ctx context.Context, sub model.Subscription) (uuid.UUID, error) {
	const op = "service.Create"
	log := s.logger.With(
		slog.String("op", op),
		slog.String("user_id", sub.UserID.String()),
	)

	log.Info("Создание подписки")

	id, err := s.repo.Create(ctx, sub)
	if err != nil {
		log.Error("Не удалось создать подписку в репозитории", slog.String("error", err.Error()))
		return uuid.Nil, err
	}

	log.Info("Подписка успешно создана", slog.String("subscription_id", id.String()))
	return id, nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (model.Subscription, error) {
	const op = "service.GetByID"
	log := s.logger.With(
		slog.String("op", op),
		slog.String("subscription_id", id.String()),
	)

	log.Info("Получение подписки")

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("Не удалось получить подписку из репозитория", slog.String("error", err.Error()))
		return model.Subscription{}, err
	}

	log.Info("Подписка успешно получена")
	return sub, nil
}

func (s *subscriptionService) Update(ctx context.Context, id uuid.UUID, sub model.Subscription) error {
	const op = "service.Update"
	log := s.logger.With(
		slog.String("op", op),
		slog.String("subscription_id", id.String()),
	)

	log.Info("Обновление подписки")

	err := s.repo.Update(ctx, id, sub)
	if err != nil {
		log.Error("Не удалось обновить подписку в репозитории", slog.String("error", err.Error()))
		return err
	}

	log.Info("Подписка успешно обновлена")
	return nil
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "service.Delete"
	log := s.logger.With(
		slog.String("op", op),
		slog.String("subscription_id", id.String()),
	)

	log.Info("Удаление подписки")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Error("Не удалось удалить подписку в репозитории", slog.String("error", err.Error()))
		return err
	}

	log.Info("Подписка успешно удалена")
	return nil
}

// CalculateTotalCost вычисляет суммарную стоимость подписок за период.
func (s *subscriptionService) CalculateTotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, startPeriod, endPeriod time.Time) (int, error) {
	const op = "service.CalculateTotalCost"
	log := s.logger.With(
		slog.String("op", op),
		slog.String("user_id", userID.String()),
	)

	log.Info("Начат расчет суммарной стоимости")

	// 1. Получаем все релевантные подписки из базы
	subscriptions, err := s.repo.ListByUserID(ctx, userID, serviceName)
	if err != nil {
		log.Error("Не удалось получить список подписок", slog.String("error", err.Error()))
		return 0, err
	}

	totalCost := 0

	// 2. Итерируемся по каждой подписке и считаем ее вклад
	for _, sub := range subscriptions {
		// Определяем фактический конец подписки. Если его нет, считаем, что она активна до конца нашего периода.
		subEnd := endPeriod
		if sub.EndDate != nil && sub.EndDate.Before(endPeriod) {
			subEnd = *sub.EndDate
		}

		// Находим период пересечения [sub.StartDate, subEnd] и [startPeriod, endPeriod]
		overlapStart := maxTime(sub.StartDate, startPeriod)
		overlapEnd := minTime(subEnd, endPeriod)

		// Если нет пересечения, пропускаем подписку
		if overlapStart.After(overlapEnd) {
			continue
		}

		// Считаем количество полных месяцев в периоде пересечения.
		months := 0
		currentMonth := time.Date(overlapStart.Year(), overlapStart.Month(), 1, 0, 0, 0, 0, time.UTC)
		for !currentMonth.After(overlapEnd) {
			months++
			currentMonth = currentMonth.AddDate(0, 1, 0)
		}

		totalCost += sub.Price * months
	}

	log.Info("Расчет успешно завершен", slog.Int("total_cost", totalCost))
	return totalCost, nil
}

// Вспомогательные функции для работы с датами (можно вынести в отдельный пакет utils)
func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
