package http

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/vasiliy-maslov/go-subscription-service/internal/model"
	"github.com/vasiliy-maslov/go-subscription-service/internal/repository"
	"github.com/vasiliy-maslov/go-subscription-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateResponse struct {
	ID uuid.UUID `json:"id"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}

// Handler - это слой, который связывает HTTP-запросы с бизнес-логикой.
type Handler struct {
	service service.SubscriptionService
	logger  *slog.Logger
}

// NewHandler создает новый экземпляр обработчика.
func NewHandler(s service.SubscriptionService, logger *slog.Logger) *Handler {
	return &Handler{
		service: s,
		logger:  logger,
	}
}

// GetSubscriptionByID godoc
// @Summary Get a subscription by ID
// @Description Retrieves full details of a subscription by its UUID.
// @Tags subscriptions
// @Produce  json
// @Param   id path string true "Subscription UUID" Format(uuid)
// @Success 200 {object} model.Subscription
// @Failure 400 {object} ErrorResponse "Invalid UUID format"
// @Failure 404 {object} ErrorResponse "Subscription not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscriptionByID(c *gin.Context) {
	const op = "handler.GetSubscriptionByID"
	idStr := c.Param("id")
	log := h.logger.With(slog.String("op", op), slog.String("id", idStr))

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Warn("Некорректный формат UUID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
		return
	}

	log.Info("Запрос на получение подписки")

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Error("Сервис вернул ошибку при получении", slog.String("error", err.Error()))
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Adds a new subscription to the database based on the provided data.
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param   subscription body model.Subscription true "Subscription data to create. ID, CreatedAt, UpdatedAt will be ignored."
// @Success 201 {object} CreateResponse
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	const op = "handler.CreateSubscription"
	log := h.logger.With(slog.String("op", op))

	var input model.Subscription
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn("Не удалось прочитать тело запроса", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("Запрос на создание подписки", slog.Any("input", input))

	createdID, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		log.Error("Сервис вернул ошибку при создании", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": createdID})
}

// UpdateSubscription godoc
// @Summary Update an existing subscription
// @Description Updates the details of an existing subscription by its UUID.
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param   id path string true "Subscription UUID" Format(uuid)
// @Param   subscription body model.Subscription true "New subscription data. All fields must be provided."
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse "Invalid request body or UUID format"
// @Failure 404 {object} ErrorResponse "Subscription not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	const op = "handler.UpdateSubscription"
	idStr := c.Param("id")
	log := h.logger.With(slog.String("op", op), slog.String("id", idStr))

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Warn("Некорректный формат UUID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
		return
	}

	var input model.Subscription
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn("Не удалось прочитать тело запроса", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("Запрос на обновление подписки", slog.Any("input", input))

	err = h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		log.Error("Сервис вернул ошибку при обновлении", slog.String("error", err.Error()))
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Deletes a subscription by its UUID.
// @Tags subscriptions
// @Produce  json
// @Param   id path string true "Subscription UUID" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Invalid UUID format"
// @Failure 404 {object} ErrorResponse "Subscription not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	const op = "handler.DeleteSubscription"
	idStr := c.Param("id")
	log := h.logger.With(slog.String("op", op), slog.String("id", idStr))

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Warn("Некорректный формат UUID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
		return
	}

	log.Info("Запрос на удаление подписки")

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		log.Error("Сервис вернул ошибку при удалении", slog.String("error", err.Error()))
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CalculateTotalCost godoc
// @Summary Calculate total subscription cost
// @Description Calculates the total cost of subscriptions for a user over a specified period.
// @Tags subscriptions
// @Produce  json
// @Param   user_id query string true "User UUID" Format(uuid)
// @Param   start_period query string true "Start period in YYYY-MM format" Example("2024-01")
// @Param   end_period query string true "End period in YYYY-MM format" Example("2024-12")
// @Param   service_name query string false "Optional: filter by service name"
// @Success 200 {object} TotalCostResponse
// @Failure 400 {object} ErrorResponse "Missing or invalid query parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /subscriptions/total_cost [get]
func (h *Handler) CalculateTotalCost(c *gin.Context) {
	const op = "handler.CalculateTotalCost"
	log := h.logger.With(slog.String("op", op))

	userIDStr, ok := c.GetQuery("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	startPeriodStr, ok := c.GetQuery("start_period")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_period is required"})
		return
	}
	startPeriod, err := time.Parse("2006-01", startPeriodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_period format, use YYYY-MM"})
		return
	}

	endPeriodStr, ok := c.GetQuery("end_period")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_period is required"})
		return
	}
	endPeriod, err := time.Parse("2006-01", endPeriodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_period format, use YYYY-MM"})
		return
	}

	log.Info("Запрос на расчет стоимости",
		slog.String("user_id", userIDStr),
		slog.String("start_period", startPeriodStr),
		slog.String("end_period", endPeriodStr),
	)

	// Делаем конец периода последним днем месяца
	endPeriod = endPeriod.AddDate(0, 1, -1)

	var serviceName *string
	if name, exists := c.GetQuery("service_name"); exists {
		serviceName = &name
	}

	totalCost, err := h.service.CalculateTotalCost(c.Request.Context(), userID, serviceName, startPeriod, endPeriod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_cost": totalCost})
}
