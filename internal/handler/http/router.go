package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/vasiliy-maslov/go-subscription-service/docs"
)

// InitRoutes инициализирует роуты и связывает их с методами обработчика.
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("/", h.CreateSubscription)
			subscriptions.GET("/:id", h.GetSubscriptionByID)
			subscriptions.PUT("/:id", h.UpdateSubscription)
			subscriptions.DELETE("/:id", h.DeleteSubscription)
			subscriptions.GET("/total_cost", h.CalculateTotalCost)
		}
	}

	return router
}
